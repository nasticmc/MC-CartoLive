package tests

import (
	"bufio"
	"encoding/json"
	"os"
	"testing"
	"time"

	"meshcore-canada-live-map/backend/internal/meshcore"
	imqtt "meshcore-canada-live-map/backend/internal/mqtt"
)

func TestSyntheticFixtureNormalizesAndParsesPackets(t *testing.T) {
	f, err := os.Open("../../examples/fixtures/synthetic-live.ndjson")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	packetCount := 0
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var line struct {
			Topic   string          `json:"topic"`
			Payload json.RawMessage `json:"payload"`
		}
		if err := json.Unmarshal(scanner.Bytes(), &line); err != nil {
			t.Fatalf("fixture line is invalid JSON: %v", err)
		}
		msg, err := imqtt.Normalize(line.Topic, line.Payload, time.Unix(0, 0))
		if err != nil {
			t.Fatalf("fixture line failed MQTT normalization: %v", err)
		}
		if msg.TopicInfo.Subtopic == "packets" {
			packetCount++
			if _, err := meshcore.ParseHexPacket(msg.RawHex); err != nil {
				t.Fatalf("fixture packet raw hex failed MeshCore parse: %v", err)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}
	if packetCount == 0 {
		t.Fatalf("fixture should include packet lines")
	}
}
