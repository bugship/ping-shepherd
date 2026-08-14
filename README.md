# Ping Shepherd

> Repository: `ping-shepherd` · local only until we publish

A **Go** flock-watcher for HTTP endpoints. You give it URLs; it shepherds them — probes on a schedule, remembers who went down, shows a status page, and pings you on Telegram when a sheep wanders off.

Not a generic “uptime API.” Not an Uptime Kuma clone. One binary you can run and defend in an interview.

## What it will do

1. Add / list / remove HTTP targets
2. Probe them on an interval and keep last status + history in Postgres
3. Public status page
4. Telegram alert on down / recovered
5. Docker Compose, health checks, tests

## What it will not do (on purpose)

Teams, SSO, multi-region, incident theatre, a mobile app.

## Current slice

**0 — scaffold.** Config from env, `GET /health`, `GET /ready` (Postgres ping), Makefile, Compose for the database.

## Tech

| Layer | Choice |
|-------|--------|
| Language | Go |
| HTTP | `net/http` (stdlib) |
| DB | PostgreSQL via `database/sql` + `pgx` |
| Ops | Docker Compose, Makefile |

## Run (once Go is installed)

```bash
cp .env.example .env
# start Postgres (Docker) or point DATABASE_URL at a local instance
make test
make run
```

```text
curl -s localhost:8080/health
```

## Layout

```
ping-shepherd/
├── cmd/shepherd/        # one binary: API now; worker later
├── internal/
│   ├── config/          # env
│   ├── httpapi/         # HTTP handlers
│   └── store/           # Postgres
├── docker-compose.yml
├── Makefile
└── README.md
```

## License

MIT
