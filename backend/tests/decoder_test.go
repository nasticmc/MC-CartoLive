package tests

import (
	"context"
	"testing"

	"meshcore-australia-live-map/backend/internal/meshcore"
	"meshcore-australia-live-map/backend/internal/resolve"
	"meshcore-australia-live-map/backend/internal/store"
)

func TestPacketPathLengthHashSizeParsing(t *testing.T) {
	header := byte((meshcore.PayloadPlainText << 2) | meshcore.RouteFlood)
	parsed, err := meshcore.ParsePacket([]byte{header, 0x42, 0xAA, 0xBB, 0xCC, 0xDD, 0x01})
	if err != nil {
		t.Fatal(err)
	}
	if parsed.HashSize != 2 {
		t.Fatalf("hash size = %d, want 2", parsed.HashSize)
	}
	if parsed.HopCount != 2 {
		t.Fatalf("hop count = %d, want 2", parsed.HopCount)
	}
	if got := parsed.PathChunks[0]; got != "AABB" {
		t.Fatalf("first chunk = %s, want AABB", got)
	}
	if got := parsed.PathChunks[1]; got != "CCDD" {
		t.Fatalf("second chunk = %s, want CCDD", got)
	}
}

func TestReservedFourBytePathNonTraceInvalidForMap(t *testing.T) {
	header := byte((meshcore.PayloadPlainText << 2) | meshcore.RouteFlood)
	parsed, err := meshcore.ParsePacket([]byte{header, 0xC1, 0xAA, 0xBB, 0xCC, 0xDD})
	if err != nil {
		t.Fatal(err)
	}
	if !parsed.InvalidForMap {
		t.Fatalf("expected non-trace 4-byte path to be invalid for map")
	}
}

func TestFourByteTraceAllowed(t *testing.T) {
	header := byte((meshcore.PayloadTrace << 2) | meshcore.RouteFlood)
	parsed, err := meshcore.ParsePacket([]byte{header, 0xC1, 0xAA, 0xBB, 0xCC, 0xDD})
	if err != nil {
		t.Fatal(err)
	}
	if parsed.InvalidForMap {
		t.Fatalf("trace 4-byte path should be allowed for storage/resolution gates")
	}
}

func TestGroupTextPayloadDecryptsWithChannelSecret(t *testing.T) {
	parsed, err := meshcore.ParseHexPacket("15833fa002860ccae0eed9ca78b9ab0775d477c1f6490a398bf4edc75240")
	if err != nil {
		t.Fatal(err)
	}
	text := meshcore.DecodePublicMessageText(parsed.PayloadType, parsed.Payload, "", []string{"eb50a1bcb3e4e5d7bf69a57c9dada211"})
	if text != "P" {
		t.Fatalf("decoded group text = %q, want P", text)
	}
	message := meshcore.DecodePublicMessage(parsed.PayloadType, parsed.Payload, "", []string{"eb50a1bcb3e4e5d7bf69a57c9dada211"})
	if message.Sender != "Roy B V4" {
		t.Fatalf("decoded group sender = %q, want Roy B V4", message.Sender)
	}
}

func TestPublicMessageTextPrefersDecodedJSONText(t *testing.T) {
	text := meshcore.DecodePublicMessageText(meshcore.PayloadGroupText, []byte{0x01, 0x02, 0x03}, `{"decoded":{"message":"hello mesh"}}`, nil)
	if text != "hello mesh" {
		t.Fatalf("decoded json text = %q, want hello mesh", text)
	}
}

func TestPublicMessageTextReadsNestedDecoderOutput(t *testing.T) {
	raw := `{"payload":{"decoded":{"decrypted":{"sender":"TUNA BASE CAMP 4","message":"real decoded message"}}}}`
	message := meshcore.DecodePublicMessage(meshcore.PayloadGroupText, []byte{0x01, 0x02, 0x03}, raw, nil)
	if message.Sender != "TUNA BASE CAMP 4" {
		t.Fatalf("decoded json sender = %q, want TUNA BASE CAMP 4", message.Sender)
	}
	if message.Text != "real decoded message" {
		t.Fatalf("decoded json text = %q, want real decoded message", message.Text)
	}
}

func TestPublicMessageTextReadsCoreScopeDecodedJSONString(t *testing.T) {
	raw := `{"decoded_json":"{\"type\":\"CHAN\",\"channel\":\"public\",\"sender\":\"otakup0pe\",\"text\":\"otakup0pe: woah my t-deck battery didn't die\",\"flags\":0}"}`
	message := meshcore.DecodePublicMessage(meshcore.PayloadGroupText, []byte{0x01, 0x02, 0x03}, raw, nil)
	if message.Sender != "otakup0pe" {
		t.Fatalf("decoded json string sender = %q, want otakup0pe", message.Sender)
	}
	if message.Text != "woah my t-deck battery didn't die" {
		t.Fatalf("decoded json string text = %q", message.Text)
	}
}

func TestResolverHighConfidenceSingleRepeater(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t, ctx)
	addNode(t, ctx, st, "AA00000000000000000000000000000000000000000000000000000000000000", "repeater", "YKF")
	resolver := resolve.New(st, []string{"repeater", "room_server"})
	parsed, err := meshcore.ParsePacket([]byte{byte((meshcore.PayloadPlainText << 2) | meshcore.RouteFlood), 0x01, 0xAA})
	if err != nil {
		t.Fatal(err)
	}
	result, err := resolver.Resolve(ctx, "YKF", parsed)
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != resolve.StatusHigh {
		t.Fatalf("status = %s, want high (%s)", result.Status, result.Reason)
	}
}

func TestResolverCollisionIsAmbiguous(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t, ctx)
	addNode(t, ctx, st, "AA00000000000000000000000000000000000000000000000000000000000000", "repeater", "YKF")
	addNode(t, ctx, st, "AA11000000000000000000000000000000000000000000000000000000000000", "repeater", "YKF")
	resolver := resolve.New(st, []string{"repeater", "room_server"})
	parsed, err := meshcore.ParsePacket([]byte{byte((meshcore.PayloadPlainText << 2) | meshcore.RouteFlood), 0x01, 0xAA})
	if err != nil {
		t.Fatal(err)
	}
	result, err := resolver.Resolve(ctx, "YKF", parsed)
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != resolve.StatusAmbiguous {
		t.Fatalf("status = %s, want ambiguous", result.Status)
	}
}

func TestResolverCompanionIsRoleInvalid(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t, ctx)
	addNode(t, ctx, st, "AA00000000000000000000000000000000000000000000000000000000000000", "companion", "YKF")
	resolver := resolve.New(st, []string{"repeater", "room_server"})
	parsed, err := meshcore.ParsePacket([]byte{byte((meshcore.PayloadPlainText << 2) | meshcore.RouteFlood), 0x01, 0xAA})
	if err != nil {
		t.Fatal(err)
	}
	result, err := resolver.Resolve(ctx, "YKF", parsed)
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != resolve.StatusRoleInvalid {
		t.Fatalf("status = %s, want role_invalid", result.Status)
	}
}

func TestResolverDuplicatePrefixRejected(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t, ctx)
	addNode(t, ctx, st, "AA00000000000000000000000000000000000000000000000000000000000000", "repeater", "YKF")
	resolver := resolve.New(st, []string{"repeater", "room_server"})
	parsed, err := meshcore.ParsePacket([]byte{byte((meshcore.PayloadPlainText << 2) | meshcore.RouteFlood), 0x02, 0xAA, 0xAA})
	if err != nil {
		t.Fatal(err)
	}
	result, err := resolver.Resolve(ctx, "YKF", parsed)
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != resolve.StatusDuplicatePrefix {
		t.Fatalf("status = %s, want duplicate_prefix", result.Status)
	}
}

func newTestStore(t *testing.T, ctx context.Context) *store.Store {
	t.Helper()
	st, err := store.OpenMemory(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st.Close() })
	return st
}

func addNode(t *testing.T, ctx context.Context, st *store.Store, publicKey string, role string, iata string) {
	t.Helper()
	lat := 43.4
	lng := -80.4
	nodeType := 2
	if role == "companion" {
		nodeType = 1
	}
	_, err := st.UpsertAdvertNode(ctx, iata, meshcore.Advert{
		PublicKey:      publicKey,
		NodeType:       nodeType,
		Role:           role,
		Latitude:       &lat,
		Longitude:      &lng,
		LocationSource: "test",
	}, 1747665456000)
	if err != nil {
		t.Fatal(err)
	}
}
