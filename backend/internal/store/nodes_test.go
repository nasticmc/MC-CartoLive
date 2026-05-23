package store

import "testing"

func TestValidMapCoordsAcceptsAustraliaAndRejectsOldNorthernHemisphereBounds(t *testing.T) {
	t.Parallel()

	if !validMapCoords(-33.8688, 151.2093) {
		t.Fatalf("expected Sydney coordinates to be valid")
	}
	if !validMapCoords(-37.8136, 144.9631) {
		t.Fatalf("expected Melbourne coordinates to be valid")
	}
	if validMapCoords(43.65, -79.38) {
		t.Fatalf("expected old northern/western hemisphere sample to be invalid")
	}
	if validMapCoords(0, 0) {
		t.Fatalf("expected zero coordinates to be invalid")
	}
}
