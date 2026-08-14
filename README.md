# Ping Shepherd

Go service that watches HTTP endpoints. Add targets, probe them on an interval, keep history in Postgres, and notify Telegram when something goes down.

## Status

Working now:

- `GET /` status page
- `GET /health`
- `GET /ready`
- `POST /targets`, `GET /targets`, `GET /targets/{id}`, `DELETE /targets/{id}`
- `GET /targets/{id}/checks`
- background probes on `CHECK_INTERVAL`

Still to build: Telegram.

## Run

```bash
cp .env.example .env
make test
make run
```

```bash
curl -s localhost:8080/health
```

Postgres is optional for tests (in-memory store). For `/ready` and persisted targets, set `DATABASE_URL` and start the Compose database.

## Layout

```
cmd/shepherd/     process entrypoint
internal/config/  env
internal/httpapi/ HTTP handlers
internal/store/   memory + Postgres
```

## License

MIT
