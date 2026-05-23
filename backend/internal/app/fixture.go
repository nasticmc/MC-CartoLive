package app

import (
	"bufio"
	"context"
	"encoding/json"
	"os"
	"time"

	imqtt "meshcore-australia-live-map/backend/internal/mqtt"
)

type fixtureLine struct {
	Topic      string          `json:"topic"`
	Payload    json.RawMessage `json:"payload"`
	ReceivedAt int64           `json:"receivedAt"`
}

func (a *Application) replayFixture(ctx context.Context, path string) {
	f, err := os.Open(path)
	if err != nil {
		a.Log.Warn("fixture replay open failed", "path", path, "error", err)
		return
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return
		default:
		}
		var line fixtureLine
		if err := json.Unmarshal(scanner.Bytes(), &line); err != nil {
			a.Log.Debug("fixture line invalid", "error", err)
			continue
		}
		received := time.Now()
		if line.ReceivedAt > 0 {
			received = time.UnixMilli(line.ReceivedAt)
		}
		msg, err := imqtt.Normalize(line.Topic, line.Payload, received)
		if err != nil {
			a.Log.Debug("fixture normalize failed", "error", err)
			continue
		}
		a.HandleMQTT(ctx, msg)
		time.Sleep(150 * time.Millisecond)
	}
	if err := scanner.Err(); err != nil {
		a.Log.Warn("fixture replay read failed", "path", path, "error", err)
	}
}
