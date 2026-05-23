package tests

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"meshcore-australia-live-map/backend/internal/api"
	"meshcore-australia-live-map/backend/internal/live"
	"meshcore-australia-live-map/backend/internal/resolve"
)

func TestPublicLiveStateStripsSensitiveFieldsAndInvalidCoordinates(t *testing.T) {
	lat := -33.8688
	lng := 151.2093
	invalidLat := 0.0
	invalidLng := 0.0
	publicKey := "AA00000000000000000000000000000000000000000000000000000000000000"
	packetHash := "secret-packet-hash"
	pathHex := "AABBCC"
	summary := "private packet text"

	state := live.State{
		ServerTime: 1747665456000,
		Nodes: []live.Node{
			{
				NodeID:           "node-a",
				PublicKey:        publicKey,
				Name:             publicKey[:8],
				Role:             "repeater",
				Latitude:         &lat,
				Longitude:        &lng,
				LastSeen:         1747665456000,
				FirstSeen:        1747660000000,
				IATAsHeardIn:     []string{"SYD"},
				ObservationCount: 3,
			},
			{
				NodeID:    "node-invalid",
				PublicKey: "BB00000000000000000000000000000000000000000000000000000000000000",
				Name:      "Invalid",
				Role:      "companion",
				Latitude:  &invalidLat,
				Longitude: &invalidLng,
			},
		},
		Observers: []live.Observer{
			{
				PublicKey: publicKey,
				IATA:      "SYD",
				Name:      publicKey[:8],
				Latitude:  &lat,
				Longitude: &lng,
			},
		},
		RecentPackets: []live.PacketObservation{
			{
				ID:                12,
				PacketHash:        packetHash,
				PayloadTypeName:   "PLAIN_TEXT",
				RouteTypeName:     "FLOOD",
				ObserverPublicKey: publicKey,
				IATA:              "SYD",
				HeardAt:           1747665456000,
				HopCount:          1,
				PathHex:           pathHex,
				ResolutionReason:  "debug reason",
				Summary:           summary,
			},
			{
				ID:                13,
				PacketHash:        "observer-only-hash",
				PayloadTypeName:   "ADVERT",
				RouteTypeName:     "FLOOD",
				ObserverPublicKey: publicKey,
				IATA:              "SYD",
				HeardAt:           1747665457000,
				HopCount:          0,
				ResolutionStatus:  resolve.StatusNoPath,
			},
		},
		RecentEdgeEvents: []live.EdgeEvent{
			{
				ID:              1,
				PacketHash:      packetHash,
				PayloadTypeName: "PLAIN_TEXT",
				HeardAt:         1747665456000,
				RenderReason:    "resolved_path_high_confidence",
				Segments: []live.EdgeSegment{
					{
						From:       live.EdgeEndpoint{NodeID: "node-a", Name: publicKey[:8], Lat: -33.8688, Lng: 151.2093},
						To:         live.EdgeEndpoint{NodeID: publicKey, Name: publicKey[:8], Lat: -37.8136, Lng: 144.9631},
						DistanceKM: 94,
					},
				},
			},
		},
	}

	publicState := live.BuildPublicLiveState(state, live.PublicStats{Packets: 1, MQTTConnected: true})
	if len(publicState.Nodes) != 1 {
		t.Fatalf("public nodes = %d, want 1", len(publicState.Nodes))
	}
	if publicState.Nodes[0].Label == publicKey[:8] {
		t.Fatalf("public label leaked a public key prefix")
	}
	if len(publicState.Routes) != 1 {
		t.Fatalf("public routes = %d, want 1", len(publicState.Routes))
	}
	if publicState.Routes[0].To.NodeID == publicKey {
		t.Fatalf("public route endpoint leaked a public key")
	}
	if publicState.RecentActivity[0].AnimationState != live.PublicAnimationRoute || publicState.RecentActivity[0].ResolutionBucket != live.PublicBucketRouted {
		t.Fatalf("routed activity metadata = %#v, want route/routed", publicState.RecentActivity[0])
	}
	observerActivity := publicState.RecentActivity[1]
	if observerActivity.AnimationState != live.PublicAnimationObserver || observerActivity.ResolutionBucket != live.PublicBucketObserverOnly {
		t.Fatalf("observer activity metadata = %#v, want observer/observer_only", observerActivity)
	}
	if observerActivity.ObserverLocation == nil {
		t.Fatalf("observer activity missing sanitized observer location")
	}
	if observerActivity.ObserverLocation.Label == publicKey[:8] {
		t.Fatalf("observer location leaked public key label")
	}

	body, err := json.Marshal(publicState)
	if err != nil {
		t.Fatal(err)
	}
	raw := string(body)
	for _, forbidden := range []string{
		"publicKey",
		"packetHash",
		"observerPublicKey",
		"pathHex",
		"summary",
		"resolutionReason",
		publicKey,
		publicKey[:8],
		packetHash,
		pathHex,
		summary,
	} {
		if strings.Contains(raw, forbidden) {
			t.Fatalf("public JSON contains forbidden value %q: %s", forbidden, raw)
		}
	}
}

func TestPublicResolutionBuckets(t *testing.T) {
	tests := []struct {
		name   string
		packet live.PacketObservation
		route  bool
		want   string
	}{
		{name: "routed", route: true, want: live.PublicBucketRouted},
		{name: "observer only", packet: live.PacketObservation{ResolutionStatus: resolve.StatusNoPath}, want: live.PublicBucketObserverOnly},
		{name: "missing location", packet: live.PacketObservation{ResolutionStatus: resolve.StatusMissingCoords}, want: live.PublicBucketMissingLoc},
		{name: "rf gated", packet: live.PacketObservation{ResolutionStatus: resolve.StatusMissingRF}, want: live.PublicBucketRFGated},
		{name: "distance gated", packet: live.PacketObservation{ResolutionStatus: resolve.StatusDistanceGate}, want: live.PublicBucketDistanceGated},
		{name: "invalid", packet: live.PacketObservation{InvalidForMap: true}, want: live.PublicBucketNotMapSafe},
		{name: "ambiguous", packet: live.PacketObservation{ResolutionStatus: resolve.StatusAmbiguous}, want: live.PublicBucketUnresolved},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := live.PublicResolutionBucket(tt.packet, tt.route); got != tt.want {
				t.Fatalf("PublicResolutionBucket() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestPublicRouteAggregationUsesStableIDsCountsAndBuckets(t *testing.T) {
	edge := func(packetHash string, heardAt int64) live.EdgeEvent {
		return live.EdgeEvent{
			PacketHash:      packetHash,
			PayloadTypeName: "ADVERT",
			HeardAt:         heardAt,
			Segments: []live.EdgeSegment{
				{
					From:       live.EdgeEndpoint{NodeID: "node-a", Name: "A", Lat: -33.8688, Lng: 151.2093},
					To:         live.EdgeEndpoint{NodeID: "node-b", Name: "B", Lat: -37.8136, Lng: 144.9631},
					DistanceKM: 94,
				},
			},
		}
	}

	routes, routesByPacket := live.BuildPublicRoutes([]live.EdgeEvent{
		edge("hash-1", 1000),
		edge("hash-2", 2000),
	})

	if len(routes) != 1 {
		t.Fatalf("routes = %d, want 1", len(routes))
	}
	if routes[0].ID != live.PublicRouteID("node-a", "node-b") || routes[0].ID != live.PublicRouteID("node-b", "node-a") {
		t.Fatalf("route ID is not stable across endpoint order")
	}
	if routes[0].PacketCount != 2 {
		t.Fatalf("packet count = %d, want 2", routes[0].PacketCount)
	}
	if routes[0].LastHeard != 2000 {
		t.Fatalf("last heard = %d, want 2000", routes[0].LastHeard)
	}
	if routes[0].FrequencyBucket < 0 || routes[0].FrequencyBucket > 4 {
		t.Fatalf("frequency bucket = %d, want 0..4", routes[0].FrequencyBucket)
	}
	if got := routesByPacket["hash-1"]; len(got) != 1 || got[0] != routes[0].ID {
		t.Fatalf("routes by packet = %#v, want route ID", got)
	}
}

func TestPublicMessageAnchorsChooseSourceThenObserverFallback(t *testing.T) {
	edge := live.EdgeEvent{
		ID:              77,
		IATA:            "SYD",
		PayloadTypeName: "PLAIN_TEXT",
		MessageSender:   "Alice",
		MessageText:     "hello from route",
		HeardAt:         1747665456000,
		MessageAnchor:   &live.MessageAnchor{Kind: "source", Endpoint: live.EdgeEndpoint{NodeID: "node-a", Name: "Sender", Lat: -33.8688, Lng: 151.2093}},
		Segments: []live.EdgeSegment{
			{
				From:       live.EdgeEndpoint{NodeID: "node-a", Name: "Sender", Lat: -33.8688, Lng: 151.2093},
				To:         live.EdgeEndpoint{NodeID: "node-b", Name: "Receiver", Lat: -37.8136, Lng: 144.9631},
				DistanceKM: 94,
			},
		},
	}
	pulse, ok := live.PublicRoutePulseFromEdge(edge)
	if !ok {
		t.Fatalf("route pulse not built")
	}
	if pulse.MessageAnchor == nil || pulse.MessageAnchor.Kind != "source" || pulse.MessageAnchor.Label != "Sender" {
		t.Fatalf("route message anchor = %#v, want sender/source", pulse.MessageAnchor)
	}
	edge.MessageAnchor = &live.MessageAnchor{Kind: "observer", Endpoint: live.EdgeEndpoint{NodeID: "observer-key", Name: "SYD observer", Lat: -37.8136, Lng: 144.9631}}
	pulse, ok = live.PublicRoutePulseFromEdge(edge)
	if !ok {
		t.Fatalf("route pulse with observer anchor not built")
	}
	if pulse.MessageAnchor == nil || pulse.MessageAnchor.Kind != "observer" || pulse.MessageAnchor.Label != "SYD observer" || pulse.MessageAnchor.NodeID != "" {
		t.Fatalf("route observer fallback anchor = %#v, want observer without node id", pulse.MessageAnchor)
	}

	observerLocation := &live.PublicObserverLocation{Label: "SYD observer", IATA: "SYD", Lat: -37.8136, Lng: 144.9631}
	activity := live.PublicActivityFromPacket(live.PacketObservation{
		ID:               88,
		IATA:             "SYD",
		PayloadTypeName:  "PLAIN_TEXT",
		MessageSender:    "Alice",
		MessageText:      "hello from observer",
		ResolutionStatus: resolve.StatusNoPath,
		HeardAt:          1747665457000,
	}, nil, observerLocation)
	if activity.MessageAnchor == nil || activity.MessageAnchor.Kind != "observer" || activity.MessageAnchor.Label != "SYD observer" {
		t.Fatalf("observer fallback message anchor = %#v, want observer", activity.MessageAnchor)
	}
	if live.PublicActivityFromPacket(live.PacketObservation{
		ID:               89,
		IATA:             "SYD",
		PayloadTypeName:  "PLAIN_TEXT",
		MessageText:      "no anchor",
		ResolutionStatus: resolve.StatusMissingCoords,
	}, nil, nil).MessageAnchor != nil {
		t.Fatalf("missing observer/source should not create message anchor")
	}
}

func TestPublicIATAAllowlistFiltersStateAndReportsAnomalies(t *testing.T) {
	lat := -37.8136
	lng := 144.9631
	filter := live.NewPublicIATAFilter([]string{"SYD"})
	if live.NewPublicIATAFilter([]string{"Y*"}).Allows("SYD") {
		t.Fatalf("wildcard IATA entries must not allow public traffic")
	}
	if filter.Allows("") {
		t.Fatalf("blank IATA must not be public when the allowlist is enabled")
	}
	state := live.State{
		ServerTime: 1747665456000,
		Nodes: []live.Node{
			{
				NodeID:       "node-syd",
				PublicKey:    "AA00000000000000000000000000000000000000000000000000000000000000",
				Name:         "Sydney",
				Role:         "repeater",
				Latitude:     &lat,
				Longitude:    &lng,
				IATAsHeardIn: []string{"SYD", "PRG"},
			},
		},
		Observers: []live.Observer{
			{PublicKey: "observer-y", IATA: "SYD", Latitude: &lat, Longitude: &lng},
			{PublicKey: "observer-prg", IATA: "PRG", Latitude: &lat, Longitude: &lng},
		},
		RecentPackets: []live.PacketObservation{
			{ID: 1, IATA: "SYD", PayloadTypeName: "ADVERT", ResolutionStatus: resolve.StatusNoPath, ObserverPublicKey: "observer-y"},
			{ID: 2, IATA: "PRG", PayloadTypeName: "ADVERT", ResolutionStatus: resolve.StatusNoPath, ObserverPublicKey: "observer-prg"},
		},
		RecentEdgeEvents: []live.EdgeEvent{
			{ID: 1, IATA: "SYD", PacketHash: "hash-y", PayloadTypeName: "ADVERT", HeardAt: 1747665456000, Segments: []live.EdgeSegment{{From: live.EdgeEndpoint{NodeID: "node-syd", Name: "Sydney", Lat: lat, Lng: lng}, To: live.EdgeEndpoint{NodeID: "node-2", Name: "Node", Lat: lat + 0.1, Lng: lng - 0.1}, DistanceKM: 12}}},
			{ID: 2, IATA: "PRG", PacketHash: "hash-prg", PayloadTypeName: "ADVERT", HeardAt: 1747665456000, Segments: []live.EdgeSegment{{From: live.EdgeEndpoint{NodeID: "node-prg", Name: "Prague", Lat: lat, Lng: lng}, To: live.EdgeEndpoint{NodeID: "node-3", Name: "Node", Lat: lat + 0.1, Lng: lng - 0.1}, DistanceKM: 12}}},
		},
	}

	filtered, excluded := filter.FilterState(state)
	publicState := live.BuildPublicLiveState(filtered, live.PublicStats{})
	cache := live.NewPublicStateCache(filter)
	cache.Replace(publicState, excluded)
	snapshot, ok := cache.Snapshot()
	if !ok {
		t.Fatalf("cache snapshot not ready")
	}
	if len(snapshot.RecentActivity) != 1 || snapshot.RecentActivity[0].IATA != "SYD" {
		t.Fatalf("public activity = %#v, want only SYD", snapshot.RecentActivity)
	}
	if len(snapshot.RecentPulses) != 1 || snapshot.RecentPulses[0].IATA != "SYD" {
		t.Fatalf("public pulses = %#v, want only SYD", snapshot.RecentPulses)
	}
	if got := snapshot.Stats.ExcludedIATAs["PRG"]; got == 0 {
		t.Fatalf("excluded PRG anomaly counter = %d, want > 0", got)
	}
	if got := snapshot.Nodes[0].IATAsHeardIn; len(got) != 1 || got[0] != "SYD" {
		t.Fatalf("filtered node IATAs = %#v, want SYD only", got)
	}
}

func TestPublicModeBlocksInternalAPIRoutes(t *testing.T) {
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	server := api.Server{
		Config:        api.Config{PublicMode: true},
		PublicHub:     live.NewHub(log, 4),
		MQTTConnected: func() bool { return false },
		MQTTTotal:     func() int64 { return 0 },
	}

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/v1/debug/stats", nil)
	server.Routes().ServeHTTP(response, request)

	if response.Code == http.StatusOK {
		t.Fatalf("debug endpoint returned 200 in public mode")
	}
}
