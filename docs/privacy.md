# Privacy Model

## Private Inputs

Keep these out of git, logs, issues, screenshots, and public artifacts:

- MQTT username and password
- MeshCore private keys
- group/channel secrets
- `.env` files with real values
- live SQLite databases, WAL files, and SHM files
- local `data/config.yaml`
- raw packet captures copied from live traffic

## Public Outputs

Public endpoints should expose only sanitized live-map data needed for display.
They must not expose:

- public keys
- observer public keys
- packet hashes
- raw packet summaries
- path hex
- resolver debug reasons
- raw packet payloads

Decoded message text is exposed only as sanitized public bubble text when the
backend can decode it from public packet data or from private channel secrets
provided locally by the operator.

## IATA Allowlist

The public map filters state and live events through `PUBLIC_IATAS`. Unsupported
or unexpected IATA traffic is counted as an anomaly and excluded from the public
map.

Keep the allowlist to supported Australia IATA region codes unless there is an
explicit product decision to publish another region.

## Route Truth

Only high-confidence RF paths become public route animations. Ambiguous,
duplicate-prefix, missing-location, missing-RF, distance-gated, invalid, and
unresolved observations do not create guessed public routes.

When an observation cannot safely draw a route but the observer has a public
location, the frontend can show observer-only live activity instead.

## Tests

Privacy-sensitive changes must keep backend public-state tests passing:

```bash
cd backend
go test ./...
```

Frontend changes that affect message bubbles, live scheduling, routes, clusters,
or labels should keep the web test suite passing:

```bash
cd web
npm test -- --run
```
