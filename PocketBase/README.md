## Development

The backend is a custom [PocketBase](https://pocketbase.io) binary (Go) in this directory. It provides:

- `POST /api/igdb/search-games` and `POST /api/igdb/game-info` — IGDB proxies
- `GET/POST /api/sync/game-data` and `GET /api/sync/game-data/share/{hash}` — cross-device sync and sharing
- `POST /api/sync/game-data/pair/create` and `POST /api/sync/game-data/pair/claim` — short-lived, single-use pairing codes used to hand a sync key to another device without putting the key in a URL

### Running locally

```pwsh
copy .env.example .env        # fill in your IGDB credentials
.\scripts\loadenv.ps1         # loads .env into the session
go build -o pb.exe .
.\pb.exe serve                # http://127.0.0.1:8090 (admin UI at /_/)
```

Or with Docker: `docker compose up -d --build` (http://127.0.0.1:8090). The dev compose binds
127.0.0.1 only and CORS is limited to local origins plus the Pages origin.

### Configuration

| Variable | Purpose |
| --- | --- |
| `IGDB_CLIENT_ID` / `IGDB_CLIENT_SECRET` | IGDB API credentials (from https://api-docs.igdb.com/#account-creation) |
| `PB_HTTP` | Override the listen address (default `0.0.0.0:8090`) |
| `PB_HOST` | Domain to serve with PocketBase's built-in Let's Encrypt autocert (production). Leave unset for local dev, which uses `PB_HTTP`. |
| `PB_ORIGINS` | Comma-separated CORS `--origins` allowlist |

`VITE_PB_API` in `../SimpleGameTracker/.env.local` points the frontend at the API
(e.g. `VITE_PB_API="http://127.0.0.1:8090"`); run the frontend with `npm run dev`.
Production builds read `../SimpleGameTracker/.env.production` instead — see the
frontend README for the full env layout.

### Deploying to a server

The app ships as a single static binary, so any Linux host with Docker works. The Dockerfile
cross-compiles from any build host (the Go build stage honors `TARGETOS`/`TARGETARCH`, so
`docker buildx build --platform linux/arm64` on an amd64 machine produces a real arm64 image —
no QEMU needed).

Build and export the image plus data:

```pwsh
docker buildx build --platform linux/arm64 -t simplegametracker-pocketbase:latest --load .
# NOTE: overwrites the local amd64 :latest tag — local dev recovers later with `docker compose up -d --build`

docker compose stop                                   # graceful stop flushes the SQLite WAL
tar -czf pb-data.tar.gz pb_data                       # carries the superuser + all data rows
docker save -o pb-image.tar simplegametracker-pocketbase:latest
# (don't pipe `docker save` through `>` in PowerShell — the UTF-16 redirect corrupts the tarball)

scp pb-image.tar pb-data.tar.gz docker-compose.prod.yml .env user@server:~/app/
```

On the server, from the directory holding all four files:

```sh
docker load -i pb-image.tar
docker image inspect simplegametracker-pocketbase:latest --format '{{.Os}}/{{.Architecture}}'  # expect linux/arm64
tar -xzf pb-data.tar.gz       # creates ./pb_data
sudo chown -R 100:101 pb_data # CRITICAL: the container runs as user pb (100:101); a tarball
                              # extracted by your login user is uid 1000 and SQLite writes would fail
chmod 600 .env
docker compose -f docker-compose.prod.yml up -d
docker logs -f simplegametracker-pb           # watch for autocert issuance + listeners on :80/:443
```

`docker-compose.prod.yml` publishes ports 80 and 443 so PocketBase's built-in autocert can issue
and renew a Let's Encrypt certificate for `PB_HOST` (cached in `pb_data`, surviving restarts) and
redirect HTTP to HTTPS. PocketBase retries issuance in the background, so a late firewall/DNS fix
self-heals. Docker-published ports traverse the FORWARD/DOCKER chain, which Docker manages itself —
on hosts with a restrictive INPUT-chain firewall (e.g. some cloud images), published ports still
work; add explicit firewall rules only if curl times out without a TLS error.

**Checklist before first start:**

1. DNS: an A record for `PB_HOST` pointing at the server — required for Let's Encrypt validation.
2. Cloud firewall / security list: allow inbound TCP 80 and 443.

Rotate the superuser password after first start (an imported `pb_data` carries the old one):

```sh
docker exec simplegametracker-pb /pb superuser update your-admin@example.com "NewStrongPassword"
```

### Serving the PWA from PocketBase (optional)

Mount the built frontend `dist/` at `/pb_public` and set `PB_PUBLIC_DIR=/pb_public` to serve the
PWA from PocketBase itself instead of a static host.
