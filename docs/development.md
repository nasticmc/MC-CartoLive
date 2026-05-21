# Development

## Local Docker

```bash
cp .env.example .env
docker compose up --build
```

Open `http://localhost:39476`.

The public example starts in fixture mode. To use live MQTT, edit `.env`, set
`MQTT_ENABLED=true`, clear `FIXTURE_REPLAY_PATH`, and add private MQTT
credentials.

Do not commit `.env`, `data/config.yaml`, live databases, or WAL/SHM files.

## Credential-Free Fixture Run

Use this when you do not have MQTT credentials or when testing UI behavior in a
repeatable way.

The committed `.env.example` already uses:

```text
MQTT_ENABLED=false
FIXTURE_REPLAY_PATH=/app/examples/fixtures/synthetic-live.ndjson
```

Then run:

```bash
docker compose up --build
```

The fixture at `examples/fixtures/synthetic-live.ndjson` contains fake public
keys, fake node names, and synthetic decoded message text.

## Backend

```bash
cd backend
go test ./...
go run ./cmd/app
```

Useful local debug APIs are available only when `PUBLIC_MODE=false`:

```bash
curl http://localhost:39476/api/v1/live/state
curl "http://localhost:39476/api/v1/debug/resolution?status=ambiguous&limit=50"
curl "http://localhost:39476/api/v1/debug/collisions?hashSize=1"
```

## Frontend

```bash
cd web
npm ci
npm test -- --run
npm run build
```

Vite dev server:

```bash
cd web
npm run dev
```

The frontend expects the Go backend for live API/WebSocket data when running
outside Docker.

## Release Checks

Run before publishing or opening a pull request:

```bash
cd backend
go test ./...
```

```bash
cd web
npm ci
npm test -- --run
npm run build
```

```bash
docker compose build
```

Check privacy before committing:

```bash
git status --short --ignored
```

Private files should appear only under ignored output.
