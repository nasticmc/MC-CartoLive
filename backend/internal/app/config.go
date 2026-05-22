package app

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	ListenAddr              string
	PublicBaseURL           string
	DataDir                 string
	DBPath                  string
	LogLevel                string
	MQTTEnabled             bool
	MQTTBrokerURL           string
	MQTTTopic               string
	MQTTClientID            string
	AuthMode                string
	MQTTUsername            string
	MQTTPassword            string
	MeshcorePublicKey       string
	MeshcorePrivateKey      string
	MeshcoreChannelSecrets  []string
	MQTTTokenAudience       string
	StrictRFOnly            bool
	RequireRSSIOrSNRForEdge bool
	MaxUnverifiedEdgeKM     float64
	AllowLongTraceEdges     bool
	DefaultCenterLat        float64
	DefaultCenterLng        float64
	DefaultZoom             float64
	DefaultRegion           string
	PublicMode              bool
	RecentPacketLimit       int
	RecentEdgeEventLimit    int
	WSClientQueueSize       int
	MQTTIngestQueueSize     int
	PublicIATAs             []string
	PublicCacheRefreshSec   int
	ConfigYAML              string
	FixtureReplayPath       string
	FixtureRecordEnabled    bool
}

func LoadConfig() (Config, error) {
	_ = loadDotEnv(".env")
	cfg := Config{
		ListenAddr:              envString("LISTEN_ADDR", ":8080"),
		PublicBaseURL:           envString("PUBLIC_BASE_URL", "http://localhost:8080"),
		DataDir:                 envString("DATA_DIR", "./data"),
		DBPath:                  envString("DB_PATH", "./data/meshcore-live.db"),
		LogLevel:                envString("LOG_LEVEL", "info"),
		MQTTEnabled:             envBool("MQTT_ENABLED", true),
		MQTTBrokerURL:           envString("MQTT_BROKER_URL", "wss://mqtt1.meshcore.ca:443/mqtt"),
		MQTTTopic:               envString("MQTT_TOPIC", "meshcore/#"),
		MQTTClientID:            envString("MQTT_CLIENT_ID", "meshcore-canada-live-map-local"),
		AuthMode:                envString("MESHCORE_AUTH_MODE", "subscriber"),
		MQTTUsername:            os.Getenv("MQTT_USERNAME"),
		MQTTPassword:            os.Getenv("MQTT_PASSWORD"),
		MeshcorePublicKey:       os.Getenv("MESHCORE_PUBLIC_KEY_HEX"),
		MeshcorePrivateKey:      os.Getenv("MESHCORE_PRIVATE_KEY_HEX"),
		MeshcoreChannelSecrets:  envList("MESHCORE_CHANNEL_SECRETS"),
		MQTTTokenAudience:       envString("MQTT_TOKEN_AUDIENCE", "mqtt1.meshcore.ca"),
		StrictRFOnly:            envBool("STRICT_RF_ONLY", true),
		RequireRSSIOrSNRForEdge: envBool("REQUIRE_RSSI_OR_SNR_FOR_EDGE", true),
		MaxUnverifiedEdgeKM:     envFloat("MAX_UNVERIFIED_EDGE_KM", 150),
		AllowLongTraceEdges:     envBool("ALLOW_LONG_TRACE_EDGES", true),
		DefaultCenterLat:        envFloat("DEFAULT_CENTER_LAT", -25.2744),
		DefaultCenterLng:        envFloat("DEFAULT_CENTER_LNG", 133.7751),
		DefaultZoom:             envFloat("DEFAULT_ZOOM", 4.35),
		DefaultRegion:           envString("DEFAULT_REGION", "AUSTRALIA"),
		PublicMode:              envBool("PUBLIC_MODE", true),
		RecentPacketLimit:       envInt("RECENT_PACKET_LIMIT", 1000),
		RecentEdgeEventLimit:    envInt("RECENT_EDGE_EVENT_LIMIT", 2000),
		WSClientQueueSize:       envInt("WS_CLIENT_QUEUE_SIZE", 512),
		MQTTIngestQueueSize:     envInt("MQTT_INGEST_QUEUE_SIZE", 4096),
		PublicIATAs:             envListFallback("PUBLIC_IATAS", defaultPublicIATAs()),
		PublicCacheRefreshSec:   envInt("PUBLIC_CACHE_REFRESH_SECONDS", 10),
		ConfigYAML:              envString("CONFIG_YAML", "./data/config.yaml"),
		FixtureReplayPath:       os.Getenv("FIXTURE_REPLAY_PATH"),
		FixtureRecordEnabled:    envBool("FIXTURE_RECORD_ENABLED", false),
	}
	if cfg.AuthMode == "subscriber" && cfg.MQTTEnabled && (cfg.MQTTUsername == "" || cfg.MQTTPassword == "") {
		return cfg, fmt.Errorf("MQTT subscriber auth requires MQTT_USERNAME and MQTT_PASSWORD or MQTT_ENABLED=false")
	}
	return cfg, nil
}

func Logger(level string) *slog.Logger {
	var slogLevel slog.Level
	switch strings.ToLower(level) {
	case "debug":
		slogLevel = slog.LevelDebug
	case "warn":
		slogLevel = slog.LevelWarn
	case "error":
		slogLevel = slog.LevelError
	default:
		slogLevel = slog.LevelInfo
	}
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slogLevel}))
}

func loadDotEnv(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		value = strings.Trim(strings.TrimSpace(value), `"`)
		if _, exists := os.LookupEnv(key); !exists {
			_ = os.Setenv(key, value)
		}
	}
	return scanner.Err()
}

func envString(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}

func envList(key string) []string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return nil
	}
	fields := strings.FieldsFunc(value, func(r rune) bool {
		return r == ',' || r == ';' || r == '\n' || r == '\r' || r == '\t' || r == ' '
	})
	out := make([]string, 0, len(fields))
	for _, field := range fields {
		if item := strings.TrimSpace(field); item != "" {
			out = append(out, item)
		}
	}
	return out
}

func envListFallback(key string, fallback []string) []string {
	out := envList(key)
	if len(out) == 0 {
		return append([]string{}, fallback...)
	}
	return out
}

func defaultPublicIATAs() []string {
	return []string{
		"SYD", "MEL", "BNE", "PER", "ADL", "CBR", "HBA", "DRW", "CNS", "TSV", "OOL", "NTL", "MCY", "ABX", "LST", "MKY",
		"ASP", "AYQ", "ISA", "ROK", "PPP", "BNK", "GLT", "BDB", "BHQ", "TMW", "WGA", "SHE", "MQL", "PHE", "KGI", "BME",
		"AKL", "CHC", "WLG", "ZQN",
	}
}

func envBool(key string, fallback bool) bool {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		parsed, err := strconv.ParseBool(v)
		if err == nil {
			return parsed
		}
	}
	return fallback
}

func envInt(key string, fallback int) int {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		parsed, err := strconv.Atoi(v)
		if err == nil {
			return parsed
		}
	}
	return fallback
}

func envFloat(key string, fallback float64) float64 {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		parsed, err := strconv.ParseFloat(v, 64)
		if err == nil {
			return parsed
		}
	}
	return fallback
}
