package app

import "testing"

func TestLoadConfigDefaultsToPublicMode(t *testing.T) {
	t.Setenv("MQTT_ENABLED", "false")
	t.Setenv("PUBLIC_MODE", "")

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatal(err)
	}
	if !cfg.PublicMode {
		t.Fatalf("PublicMode = false, want true by default")
	}
}

func TestLoadConfigAllowsLocalDebugMode(t *testing.T) {
	t.Setenv("MQTT_ENABLED", "false")
	t.Setenv("PUBLIC_MODE", "false")

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.PublicMode {
		t.Fatalf("PublicMode = true, want false when PUBLIC_MODE=false")
	}
}
