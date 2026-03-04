# AGENTS.md — FlagDeck

Instructions for AI agents working on this repository.

## Project Structure

```
flagdeck/
├── api/                    # Go backend (Chi router)
│   ├── cmd/
│   │   ├── server/         # API server entrypoint
│   │   └── seed/           # Idempotent seed script
│   ├── internal/
│   │   ├── config/         # Environment variable loading
│   │   ├── database/       # MongoDB + Redis connections
│   │   ├── handlers/       # HTTP handlers (flags, auth, segments, etc.)
│   │   ├── middleware/     # Auth (JWT + API key), rate limiting, request ID
│   │   ├── models/         # MongoDB document models
│   │   ├── routes/         # Route registration
│   │   └── services/       # Evaluation engine, caching, audit logging
│   ├── Dockerfile          # Multi-stage: golang:1.24-alpine → alpine
│   ├── go.mod              # Go 1.24
│   └── go.sum
├── web/                    # SvelteKit 2 frontend (Svelte 5 runes)
│   ├── src/
│   │   ├── lib/            # API client, auth store, components, types
│   │   └── routes/         # SvelteKit pages (login, flags, segments, etc.)
│   ├── e2e/                # Playwright E2E tests
│   ├── Dockerfile          # SvelteKit adapter-static → nginx
│   └── playwright.config.ts
├── scripts/                # Production smoke test
├── docker-compose.yml      # Local MongoDB 7 + Redis 7
└── .github/workflows/ci.yml
```

## Tech Stack

- **Backend:** Go 1.24, Chi router, MongoDB driver v2, go-redis v9
- **Frontend:** SvelteKit 2, Svelte 5 (runes only), TailwindCSS, adapter-static
- **Database:** MongoDB 7 (document store), Redis 7 (caching)
- **Deployment:** Railway (auto-deploy on merge to main)
- **CI:** GitHub Actions (go vet, go test -race, go build, gofmt, npm lint, svelte-check, Playwright E2E)

## Quality Gates

### API (Go)
```bash
cd api && go vet ./...
cd api && go test ./... -v -count=1 -race
cd api && go build -o /dev/null ./cmd/server
cd api && gofmt -l . | grep . && exit 1 || true
```

### Web (SvelteKit)
```bash
cd web && npm run lint
cd web && npm run build
```

### E2E Tests
```bash
cd web && npx playwright test
```

## Rules

- **Go 1.24** — Do NOT change the Go version. It is pinned intentionally.
- **Svelte 5 runes only** — Use `$state()`, `$derived()`, `$effect()`, `$props()`. No legacy `export let`, `$:`, or `on:event` syntax.
- **snake_case fields** — All MongoDB fields and API responses use snake_case.
- **Seed is idempotent** — Uses upsert operations. Safe to re-run.
- **adapter-static** — All data fetching must use `$effect` (client-side). No `+page.server.ts` on dynamic routes.
- **Auth endpoints** — Login/register at `/auth/*`, protected API at `/api/v1/*`.

## Environment Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `MONGODB_URI` | MongoDB connection string | `mongodb://localhost:27017/flagdeck` |
| `REDIS_URL` | Redis connection string | `redis://localhost:6379` |
| `JWT_SECRET` | JWT signing secret | `your-secret-key` |
| `PORT` | API server port | `8080` |
| `CORS_ORIGINS` | Allowed CORS origins | `http://localhost:3000` |
| `ENVIRONMENT` | Runtime environment | `development` |

## Demo Credentials

- **Email:** demo@workermill.com
- **Password:** demo1234
