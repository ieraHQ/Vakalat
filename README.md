# Vakalat — Enterprise Legal Operating System

**InkOS** is a Legal Case Management Platform for law firms, built local-first and offline-first with a Clean Architecture Go backend and a Next.js frontend. It's under active development — see [Known Limitations](#known-limitations) below for an honest read on what's implemented versus planned.

## Core Principles
- **Local-first**: Runs locally with support for multiple users over LAN/VPN.
- **Offline-first**: Full functionality without an internet connection (planned — not yet implemented on the frontend).
- **Fast**: Optimized for performance with sub-300ms page loads.
- **Secure**: Argon2id password hashing, TOTP MFA, RBAC backed by real roles/permissions, per-user rate limiting on AI endpoints.
- **Cross-platform**: Web today; desktop (Tauri) and mobile (Flutter) planned.
- **Modular**: Clean Architecture with repository pattern, dependency injection, and service layers.

## Quickstart

```bash
git clone <this repo>
cd Vakalat
cp .env.example .env   # adjust JWT_SECRET etc. before anything but local use
docker compose up -d --build
```

Then:
- Frontend: http://localhost:3000
- Backend health check: http://localhost:8080/healthz
- Seed the database with sample accounts (the `migrate` service only runs schema migrations, not seed data):
  ```bash
  docker exec -i vakalat-postgres-1 psql -U postgres -d vakalat < database/seeds/seed.sql
  ```
- Log in at http://localhost:3000/login with `admin@vakalat.com` / `ChangeMe123!` (dev-only seeded account — rotate before any non-local use).

## Technology Stack
### Backend
- **Go** (Fiber)
- **PostgreSQL 17** (pgvector, pg_trgm, uuid-ossp)
- **Redis** (caching, sessions, queues)
- **golang-migrate** for schema migrations

### Frontend
- **Next.js 16** (App Router, React 19, TypeScript, TailwindCSS)
- Server Components + Server Actions for data fetching and mutations (no client-side data-fetching library wired up yet, though TanStack Query and react-hook-form are available as dependencies)

### Desktop / Mobile
- **Tauri** (macOS, Windows, Linux) and **Flutter** (Android, iOS, iPad) — planned, not started.

### AI
- **Local LLM** (LM Studio, Ollama, or any OpenAI-compatible endpoint) for summarization, Q&A, and drafting
- Embedding generation wired into the same client, feeding hybrid search
- **OCR** (OCRmyPDF, PaddleOCR, Tesseract)

### Search
- **PostgreSQL Full-Text Search** + **pgvector** (semantic search), blended into a single ranked result set

## Modules
- Authentication, Users, Roles, Permissions
- Clients, Matters, Hearings, Orders (timeline)
- Documents (upload, versioning, OCR)
- AI Assistant, Hybrid Search
- Billing, Tasks, Calendar, Notifications — schema exists, no backend/frontend yet

## Milestones

### Phase 1 – Foundation & Architecture (✅ Complete)
Monorepo layout, Go workspace, Next.js tooling, structured config.

### Phase 2 – Database (✅ Complete)
60+ table schema across users/auth, clients, legal entities, documents/OCR, tasks, billing, AI/search, and audit logging. Every migration has a matching down-migration (`database/migrations/*.down.sql`), applied with `golang-migrate`.

### Phase 3 – Backend Framework (✅ Complete)
Repository/service/handler layers, Zap logging, pgxpool connection pooling, JWT + RBAC middleware, WebSocket hub.

### Phase 4 – Authentication & RBAC (✅ Complete)
- Argon2id password hashing with a random per-password salt.
- JWT access + refresh tokens.
- RBAC backed by real `roles`/`permissions`/`role_permissions` tables (not hardcoded checks).
- Password reset via hashed, expiring, single-use tokens.
- TOTP-based MFA (`/api/auth/mfa/enable`, `/api/auth/mfa/verify`) — not yet enforced at login time.

### Phase 5 – Document Engine (✅ Complete)
Upload, local storage, OCR background worker (English + Hindi via Tesseract/PaddleOCR/OCRmyPDF), and document versioning.

### Phase 6 – Matter Engine (✅ Complete)
Hearings and orders scoped to a matter, a merged chronological timeline endpoint, and matter/client browse + detail pages in the Next.js app.

### Phase 7 – Search (✅ Complete)
Hybrid full-text + pgvector search (`GET /api/search?q=...`), gated by a `manage_search` permission. No search UI yet — API only.

### Phase 8 – AI Workspace (✅ Complete)
Provider-agnostic LLM client (OpenAI-compatible `/chat/completions` + `/embeddings`) powering summarization, Q&A, and drafting endpoints, plus the embedding pipeline for search. Rate-limited (20 req/min/user) and gated by a `manage_ai` permission. `cmd/backfill-embeddings` indexes pre-existing rows.

### Phase 9 – Deployment (✅ Complete)
`docker-compose.yml` wires Postgres (pgvector), Redis, a one-shot `golang-migrate` migration runner, the Go API, and the Next.js web app, with a real `/healthz`-based healthcheck. No TLS/ingress layer yet — see limitations.

## Known Limitations

This is a working, internally-consistent app (login → RBAC → data → search → AI all function end to end, verified against a live stack), but it is **not** production-ready for real users yet:

- **No matters/clients create-forms in the UI** — browsing works, creating currently requires calling the API directly.
- **No outbound email/SMS provider** — password reset tokens are written to the server log, not emailed.
- **MFA isn't enforced at login** — the enable/verify endpoints work standalone but aren't yet wired into the login flow as a second step.
- **No TLS/ingress** — `docker-compose.yml` serves plain HTTP; `COOKIE_SECURE` must stay `false` until that changes, or session cookies will be silently dropped by the browser.
- **Test coverage is a start, not exhaustive** — unit tests cover auth/RBAC/validation; no integration tests against a real Postgres yet.
- **No multi-tenancy, HA, or observability** beyond structured stdout logging.

## Documentation
- [Architecture](./docs/architecture.md)
- [Database Schema](./database/schema/schema.md)
- [API Documentation](./docs/api.md)
- [Deployment Guide](./docs/deployment.md)

## License
MIT
