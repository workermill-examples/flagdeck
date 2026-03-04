# FlagDeck

**A full-stack feature flag management platform built entirely by AI agents.**

FlagDeck is a showcase application demonstrating [WorkerMill](https://workermill.com) — an autonomous AI coding platform that takes Jira/Linear/GitHub tickets and ships production code. Every line of code in this repository was written, tested, and deployed by WorkerMill's AI workers.

[Live Demo](https://flagdeck-app.workermill.com) | [WorkerMill Platform](https://workermill.com) | [Documentation](https://workermill.com/docs)

[![CI](https://github.com/workermill-examples/flagdeck/actions/workflows/ci.yml/badge.svg)](https://github.com/workermill-examples/flagdeck/actions/workflows/ci.yml)

---

## What's Inside

FlagDeck is a real, functional feature flag platform — not a toy demo. It includes:

- **Feature Flags** — Boolean, string, number, and JSON flag types with per-environment overrides
- **Targeting Rules** — Conditions including eq, neq, contains, in, regex, gt, and lt operators
- **Percentage Rollouts** — Deterministic hashing (FNV-1a) for consistent user bucketing
- **A/B Experiments** — Variant definitions with traffic allocation tracking
- **Multi-Environment** — Production, staging, and development environments with independent configs
- **User Segments** — Reusable audience groups with rule-based conditions
- **Authentication** — JWT tokens for the dashboard and API keys for flag evaluation
- **Redis Caching** — Fast flag evaluation with Redis-backed cache layer
- **Audit Logging** — Full history of flag changes and configuration updates
- **Dashboard** — Stat cards, activity feed, and project-level flag management
- **Docker Compose** — One-command local dev with MongoDB 7 and Redis 7
- **CI Pipeline** — GitHub Actions running go vet, go test -race, go build, gofmt, npm lint, and svelte-check

## How It Was Built

FlagDeck was built in **3 epics over ~4.5 hours** for a total cost of **$64**. The final repo contains **20,857 lines of code across 97 files from 67 commits**.

| Epic | PR | What Was Built | Time | Lines |
|------|----|----------------|------|-------|
| FDFBS-1 | [#1](https://github.com/workermill-examples/flagdeck/pull/1) | Backend API, seed data, CI pipeline, Docker stack | ~91 min | 9,063 |
| FDFBS-2 | [#2](https://github.com/workermill-examples/flagdeck/pull/2) | SvelteKit UI — all pages, components, and auth flow | ~122 min | 11,057 |
| FDFBS-3 | [#3](https://github.com/workermill-examples/flagdeck/pull/3) | Docker infrastructure, smoke tests, go-live deployment | ~44 min | 737 |

Each epic was planned by a WorkerMill planner agent, decomposed into parallel stories, executed by specialist AI personas (frontend developer, backend developer, QA engineer), reviewed by a tech lead agent, and consolidated into a single PR.

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Backend | Go 1.24 (Chi router) |
| Frontend | SvelteKit 2 (Svelte 5 runes) |
| Database | MongoDB 7 |
| Caching | Redis 7 |
| Auth | JWT + API keys |
| Hosting | Railway |
| CI | GitHub Actions |

## Try the Demo

Visit [flagdeck-app.workermill.com](https://flagdeck-app.workermill.com) and sign in:

| | |
|-|-|
| **Email** | demo@workermill.com |
| **Password** | demo1234 |

The demo account comes with pre-configured projects, feature flags, targeting rules, and segments.

## Run Locally

**Prerequisites:** Go 1.24+, Node.js 20+, Docker

```bash
git clone https://github.com/workermill-examples/flagdeck.git
cd flagdeck
```

Start MongoDB and Redis:

```bash
docker compose up -d
```

Seed the database and start the API:

```bash
cd api
go run ./cmd/seed/main.go
go run ./cmd/server/main.go
```

The API runs on [localhost:8080](http://localhost:8080).

In a separate terminal, start the frontend:

```bash
cd web
npm install
npm run dev
```

The frontend runs on [localhost:3000](http://localhost:3000).

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `MONGODB_URI` | `mongodb://localhost:27017/flagdeck` | MongoDB connection string |
| `REDIS_URL` | `localhost:6379` | Redis connection string |
| `JWT_SECRET` | — | Secret for JWT token signing |
| `PORT` | `8080` | API server port |
| `CORS_ORIGINS` | `http://localhost:3000` | Allowed CORS origins |

---

## About WorkerMill

[WorkerMill](https://workermill.com) is an autonomous AI coding platform. Point it at a ticket, and it:

1. **Plans** — Decomposes the task into parallel stories with file targets
2. **Executes** — Specialist AI personas (frontend dev, backend dev, QA) work in parallel
3. **Reviews** — Tech lead agent reviews each story for quality
4. **Ships** — Creates a consolidated PR with all changes

FlagDeck exists to demonstrate that WorkerMill can build and maintain a real application end-to-end. Every commit in this repo's history traces back to a WorkerMill task.

## For AI Agents

If you're an AI worker building on this codebase, see [AGENTS.md](./AGENTS.md) for development guidelines, version constraints, and deployment instructions.

## License

MIT
