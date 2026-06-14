# Dummy SMTP

A development mail-catcher. It runs a fake SMTP server, accepts mail from any app
you point at it, and lets you inspect what was sent through a web UI and a JSON
API. Mail is **never relayed** — it stays local, for development and testing.

Single Go binary with the web UI (Svelte) embedded, no external services to run.

![dummy-smtp](docs/screen.png)

## Features

- Minimal SMTP (no AUTH, no TLS) — point any client at it
- Web UI to browse messages: HTML / Text / Source / Raw / Headers / JSON
- HTML preview with light/dark background and responsive widths
- JSON API to list, fetch and delete messages
- Live updates over Server-Sent Events
- Pluggable storage: in-memory (default), SQLite or MongoDB

## Quick start

Run the published image:

```bash
docker run --rm -p 1025:1025 -p 8025:8025 michaelcardoza/dummy-smtp
```

SMTP listens on `localhost:1025`, the web UI on http://localhost:8025.

Storage defaults to in-memory and resets when the container stops. For
persistence, use SQLite with a volume:

```bash
docker run --rm -p 1025:1025 -p 8025:8025 \
  -e STORAGE=sqlite -v dummy-smtp-data:/data \
  michaelcardoza/dummy-smtp
```

## Running from source

Needs Go 1.26+. The web UI is already embedded, so this is all it takes:

```bash
go run ./cmd/dummy-smtp
```

SMTP starts on `:1025`, the web UI / API on `:8025`. To build a binary:

```bash
go build ./cmd/dummy-smtp
```

> Changing the UI under `web/` means rebuilding it with Node + pnpm:
> `cd web && pnpm install && pnpm build`.

## Docker Compose

Builds from source and runs one storage backend per service. They share ports
`1025`/`8025`, so run one at a time (Compose builds the image on first run):

```bash
docker compose up dummy-smtp-memory -d         # in-memory
docker compose up dummy-smtp-sqlite -d         # SQLite, persisted
docker compose up mongodb dummy-smtp-mongo -d  # MongoDB, persisted
```

## Configuration

Each setting is a CLI flag that defaults to its env var:

| Setting      | Env var      | Flag          | Default                                 |
| ------------ | ------------ | ------------- | --------------------------------------- |
| SMTP address | `SMTP_ADDR`  | `-smtp-addr`  | `:1025`                                 |
| HTTP address | `HTTP_ADDR`  | `-http-addr`  | `:8025`                                 |
| Storage      | `STORAGE`    | `-storage`    | `memory` (`memory`, `sqlite`, `mongo`)  |
| Mongo URI    | `MONGO_URI`  | `-mongo-uri`  | `mongodb://127.0.0.1:27017/dummysmtp`   |

```bash
# via flags
go run ./cmd/dummy-smtp -smtp-addr :2025 -http-addr :9025 -storage sqlite

# same thing via env vars
SMTP_ADDR=:2025 HTTP_ADDR=:9025 STORAGE=sqlite go run ./cmd/dummy-smtp
```

## Sending mail

Point any SMTP client at `localhost:1025` — no auth, no TLS. Quick check with
curl:

```bash
curl smtp://localhost:1025 \
  --mail-from alice@example.com \
  --mail-rcpt bob@example.com \
  -T <(printf 'Subject: Hello\n\nHi, how are you?\n')
```

More client examples (Python, Go, Node, PHP, Ruby) live in [`example/`](example/)
and in the UI's **connect** panel.

## API

| Method   | Path                       | Description           |
| -------- | -------------------------- | --------------------- |
| `GET`    | `/api/v1/messages`         | List messages         |
| `GET`    | `/api/v1/messages/{id}`    | Get one message       |
| `DELETE` | `/api/v1/messages/{id}`    | Delete one message    |
| `DELETE` | `/api/v1/messages`         | Delete all messages   |

## License

MIT
