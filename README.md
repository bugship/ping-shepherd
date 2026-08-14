# HTTP Uptime Monitor

> Repository: `canary-coop` · local only until we publish

Small **Go** service that watches HTTP endpoints, stores check history, and exposes a status page. Built to be résumé-honest: one process you can run and explain, not a clone of [Uptime Kuma](https://github.com/louislam/uptime-kuma).

## What it will do

1. Add / list / remove HTTP targets
2. Probe them on an interval and keep last status + history in Postgres
3. Public status page
4. One alert channel (Telegram) on down / recovered
5. Docker Compose, health checks, tests

## What it will not do (on purpose)

Teams, SSO, multi-region, fancy incident timelines, mobile apps.

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
canary-coop/
├── cmd/coop/            # one binary: API now; worker later
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
