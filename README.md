# Ping Shepherd

Go service that watches HTTP endpoints. Add targets, probe them on an interval, keep history in Postgres, and notify Telegram when something goes down.

https://github.com/bugship/ping-shepherd

## HTTP

- `GET /` status page
- `GET /health`
- `GET /ready`
- `POST /targets`, `GET /targets`, `GET /targets/{id}`, `DELETE /targets/{id}`
- `GET /targets/{id}/checks`

## Run

```bash
cp .env.example .env
make test
make run
```

```bash
curl -s localhost:8080/health
```

Without `DATABASE_URL`, the process uses an in-memory store. For persisted targets, start Postgres and set `DATABASE_URL`.

```bash
docker compose up --build
```

Telegram alerts need `TELEGRAM_BOT_TOKEN` and `TELEGRAM_CHAT_ID`. They fire on first down and on down/up transitions.

```bash
make vuln
```

## Layout

```
cmd/shepherd/      process entrypoint
internal/config/   env
internal/httpapi/  HTTP handlers and status page
internal/checker/  probe loop and alerts
internal/probe/    HTTP GET
internal/notify/   Telegram
internal/store/    memory + Postgres
```

## License

MIT
