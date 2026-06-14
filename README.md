# Dummy SMTP

A development mail-catcher: a fake SMTP server that accepts mail from any app
pointed at it, stores it (in memory, SQLite or MongoDB), and exposes a web UI and JSON API
to inspect captured messages. **Nothing is ever relayed outbound** — it's for
local development and testing only.

Built in **Go**, shipped as a single binary with the web UI (Svelte) embedded.

![dummy-smtp](docs/screen.png)

## Features

- Catches mail over a minimal SMTP subset (no AUTH, no TLS) — point any app at it
- Web UI to browse captured messages (HTML / Text / Source / Raw / Headers / JSON)
- Responsive HTML preview with light/dark background
- JSON API to list, fetch and delete messages
- Live updates over Server-Sent Events
- Pluggable storage: in-memory (default), SQLite or MongoDB

## Getting started

```bash
go run ./cmd/dummy-smtp
```

This starts:

- **SMTP** on `:1025` — point your app's mail config here
- **Web UI / API** on `:8025` — open http://localhost:8025

Build a standalone binary:

```bash
go build ./cmd/dummy-smtp
```

## Docker

Run the image (in-memory, ephemeral — the usual mail-catcher use):

```bash
docker build -t dummy-smtp .
docker run --rm -p 1025:1025 -p 8025:8025 dummy-smtp
# UI: http://localhost:8025 · SMTP: localhost:1025
```

Or with Docker Compose — two services share the image: one in-memory and one
backed by SQLite (persisted on the `sqlite-mail-data` volume):

```bash
docker compose build
docker compose up dummy-smtp-memory -d           # in-memory (ephemeral)
docker compose up dummy-smtp-sqlite -d           # SQLite, persisted across restarts
docker compose up monogodb dummy-smtp-mongo -d   # MongoDB, persisted across restarts
```

> Both services bind the same host ports (`1025`/`8025`), so run **one at a
> time** — or change the ports in `docker-compose.yml`.

## Configuration

Each setting can be passed as an env var or a CLI flag (the env var is the flag's
default):

| Setting      | Env var      | Flag          | Default                                 |
| ------------ | ------------ | ------------- | --------------------------------------- |
| SMTP address | `SMTP_ADDR`  | `-smtp-addr`  | `:1025`                                 |
| HTTP address | `HTTP_ADDR`  | `-http-addr`  | `:8025`                                 |
| Storage      | `STORAGE`    | `-storage`    | `memory` (`memory`, `sqlite`, `mongo`)  |
| Mongo URI    | `MONGO_URI`  | `-mongo-uri`  | `mongodb://127.0.0.1:27017/dummysmtp`   |

```bash
# custom ports + sqlite, via flags
go run ./cmd/dummy-smtp -smtp-addr :2025 -http-addr :9025 -storage sqlite

# the same via env vars
SMTP_ADDR=:2025 HTTP_ADDR=:9025 STORAGE=sqlite go run ./cmd/dummy-smtp
```

`mongo` needs a running MongoDB reachable at `MONGO_URI`:

```bash
docker run -d -p 27017:27017 mongo:7     # MongoDB on localhost:27017
STORAGE=mongo go run ./cmd/dummy-smtp    # uses the default MONGO_URI (localhost)
```

## Sending mail

Point any SMTP client at `localhost:1025` (no auth, no TLS). Quick test with
curl:

```bash
curl smtp://localhost:1025 \
  --mail-from alice@example.com \
  --mail-rcpt bob@example.com \
  -T <(printf 'Subject: Hello\n\nHi, how are you?\n')
```

More examples (Python, Go, Node, PHP, Ruby) live in [`example/`](example/) and in
the **connect** panel of the web UI.

## API

| Method   | Path                       | Description           |
| -------- | -------------------------- | --------------------- |
| `GET`    | `/api/v1/messages`         | List messages         |
| `GET`    | `/api/v1/messages/{id}`    | Get one message       |
| `DELETE` | `/api/v1/messages/{id}`    | Delete one message    |
| `DELETE` | `/api/v1/messages`         | Delete all messages   |

## License

MIT
