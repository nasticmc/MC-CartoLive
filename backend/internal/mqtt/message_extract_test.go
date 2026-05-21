package mqtt

import (
	"testing"
	"time"
)

func TestNormalizeRawRSSISNRAndTimestamp(t *testing.T) {
	payload := []byte(`{"origin":"YKF Observer","raw":"0900","RSSI":"-93","SNR":"4.5","timestamp":"2025-03-16T00:07:11.191561Z"}`)
	msg, err := Normalize("meshcore/YKF/ABCDEF012345/packets", payload, time.Unix(0, 0))
	if err != nil {
		t.Fatal(err)
	}
	if msg.RawHex != "0900" {
		t.Fatalf("raw = %s", msg.RawHex)
	}
	if msg.RSSI == nil || *msg.RSSI != -93 {
		t.Fatalf("RSSI = %v", msg.RSSI)
	}
	if msg.SNR == nil || *msg.SNR != 4.5 {
		t.Fatalf("SNR = %v", msg.SNR)
	}
	if msg.ObserverName != "YKF Observer" {
		t.Fatalf("observer name = %s", msg.ObserverName)
	}
	if msg.HeardAtMs == 0 {
		t.Fatal("expected parsed timestamp")
	}
}
