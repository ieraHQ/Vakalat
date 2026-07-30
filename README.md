# Vakalat — Enterprise Legal Operating System

**InkOS** is a production-grade Legal Case Management Platform designed to replace existing legal practice management software. Built for law firms, it supports local-first, offline-first, and cross-platform deployment with enterprise-grade security, performance, and scalability.

## Core Principles
- **Local-first**: Runs locally with support for multiple users over LAN/VPN.
- **Offline-first**: Full functionality without an internet connection.
- **Fast**: Optimized for performance with sub-300ms page loads.
- **Secure**: AES-256 encryption, audit logs, and zero-trust architecture.
- **Cross-platform**: Web, desktop (Tauri), and mobile (Flutter).
- **Modular**: Clean Architecture with repository pattern, dependency injection, and service layers.

## Technology Stack
### Backend
- **Go** (Fiber/Gin)
- **PostgreSQL 17** (pgvector, pg_trgm, uuid-ossp)
- **Redis** (caching, sessions, queues)

### Frontend
- **Next.js** (React, TypeScript, TailwindCSS)
- **shadcn/ui** (UI components)
- **TanStack Query** (data fetching)
- **React Hook Form + Zod** (validation)

### Desktop
- **Tauri** (macOS, Windows, Linux)

### Mobile
- **Flutter** (Android, iOS, iPad)

### AI
- **Local LLM** (LM Studio, Ollama, OpenAI-compatible endpoints)
- **OCR** (OCRmyPDF, PaddleOCR, Tesseract)

### Search
- **PostgreSQL Full Text Search** + **pgvector** (semantic search)

## Modules
- Authentication, Users, Roles, Permissions
- Clients, Matters, Hearings, Calendar, Tasks
- Documents, Evidence, Orders, Billing, Reports
- AI Assistant, OCR, Search, Notifications

## Milestones

### Phase 1 – Foundation & Architecture (✅ Complete)
- **Monorepo Structure**: Finalized directory layout for `/apps`, `/backend`, `/packages`, `/database`, and `/docs`.
- **Makefile**: Added common tasks (`make build`, `make test`, `make lint`).
- **Go Workspace**: Configured `go.work` for multi-module development and `golangci-lint` for linting.
- **Next.js**: Configured `eslint`, `prettier`, and `typescript` for strict mode.
- **Go Backend**: Added `viper` for environment variables and `app.env` for configuration.

### Phase 2 – Database
- Design complete schema (60–90 tables).
- Set up PostgreSQL with `pgvector`, `pg_trgm`, and `uuid-ossp`.
- Add migrations and seed data.

### Phase 3 – Backend Framework
- Configure logger, database, repositories, and services.
- Implement RBAC and middleware.

### Phase 4 – Authentication
- JWT + refresh tokens.
- Session management and password reset.

### Phase 5 – Document Engine
- Upload, storage, metadata, and versioning.
- OCR queue and digital signatures.

### Phase 6 – Matter Engine
- API endpoints for matters, hearings, and timeline.
- Frontend UI for matter screen.

### Phase 7 – Search
- PostgreSQL full-text search + pgvector.
- Hybrid search (keyword + semantic).

### Phase 8 – AI Workspace
- Local LLM integration (Ollama/LM Studio).
- Summarization, drafting, and research.

### Phase 9 – Deployment
- Docker and native deployment.
- Cloud and high availability.

## Documentation
- [Architecture](./docs/architecture.md)
- [Database Schema](./database/schema/schema.md)
- [API Documentation](./docs/api.md)
- [Deployment Guide](./docs/deployment.md)

## License
MIT