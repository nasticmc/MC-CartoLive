package mqtt

import "testing"

func TestParseTopic(t *testing.T) {
	info, err := ParseTopic("meshcore/YKF/ABCDEF012345/packets")
	if err != nil {
		t.Fatal(err)
	}
	if info.IATA != "YKF" || info.PublisherPK != "ABCDEF012345" || info.Subtopic != "packets" {
		t.Fatalf("unexpected topic info: %+v", info)
	}
}

func TestParseTopicRejectsMalformed(t *testing.T) {
	if _, err := ParseTopic("meshcore/YKF/ABCDEF012345"); err == nil {
		t.Fatal("expected malformed topic error")
	}
}

func TestParseTopicAllowsInternalForExplicitDrop(t *testing.T) {
	info, err := ParseTopic("meshcore/YKF/ABCDEF012345/internal")
	if err != nil {
		t.Fatal(err)
	}
	if info.Subtopic != "internal" {
		t.Fatalf("subtopic = %s, want internal", info.Subtopic)
	}
}
