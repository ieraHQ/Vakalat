# Vakalat — Project Documentation for Claude

This document outlines the architecture, coding standards, and guidelines for interacting with the **Vakalat** project. Follow these rules to maintain consistency and production-grade quality.

---

## Architecture
### Clean Architecture
The project follows **Clean Architecture** with the following layers:

1. **Domain Layer**
   - Business logic and entities (e.g., `User`, `Matter`, `Client`).
   - Repository interfaces (e.g., `UserRepository`).
   - Use cases (e.g., `CreateUser`, `GetMatter`).

2. **Application Layer**
   - Service layer (e.g., `UserService`, `MatterService`).
   - Dependency injection.
   - Transaction management.

3. **Infrastructure Layer**
   - Database implementations (e.g., `PostgresUserRepository`).
   - External services (e.g., `RedisCache`, `OCRService`).
   - API clients.

4. **Presentation Layer**
   - REST API (Fiber/Gin).
   - WebSocket handlers.
   - Frontend (Next.js, Flutter, Tauri).

### Modular Monorepo
The project is structured as a **modular monorepo** with the following directories:

```
/vakalat
  ├── /apps                # Frontend and desktop/mobile apps
  │   ├── /web            # Next.js frontend
  │   ├── /desktop        # Tauri desktop app
  │   └── /mobile         # Flutter mobile app
  │
  ├── /backend           # Backend services
  │   ├── /api           # Go backend (Fiber/Gin)
  │   └── /worker        # Background workers (OCR, AI)
  │
  ├── /packages          # Shared packages
  │   ├── /ui            # Shared UI components
  │   ├── /shared        # Shared business logic
  │   ├── /types         # Shared types (Go/TypeScript)
  │   └── /sdk           # SDK for external integrations
  │
  ├── /database         # Database migrations and schema
  │   ├── /migrations    # PostgreSQL migrations
  │   ├── /schema        # Database schema (SQL + ER diagrams)
  │   └── /seeds         # Seed data
  │
  ├── /docs             # Documentation
  ├── /scripts          # Deployment and utility scripts
  ├── /tools            # CLI tools (e.g., OCR, AI)
  └── /storage          # Local storage (clients, matters, documents)
```

---

## Coding Standards
### Backend (Go)
- Use **golangci-lint** for linting.
- Follow **idiomatic Go** (e.g., `err` handling, no globals).
- Avoid circular dependencies.
- Use **strict interfaces**.
- Format code with `gofmt`.

### Frontend (TypeScript)
- Use **strict mode** (`"strict": true` in `tsconfig.json`).
- Avoid `any` type.
- Use **ESLint** and **Prettier** for formatting.
- Follow **React best practices** (e.g., hooks, functional components).

### Mobile (Flutter)
- Use **Clean Architecture** for Flutter.
- Follow **repository pattern** for state management.
- Support **adaptive UI** (Material + Cupertino).

### Database (PostgreSQL)
- Use **migrations** for schema changes.
- Never hardcode SQL.
- Support **transactional integrity**.
- Use **pgvector** for semantic search.

---

## Module Development
Every module must include:
1. **Database Migration**: SQL file in `/database/migrations`.
2. **Backend Implementation**: Go code in `/backend/api`.
3. **API Endpoints**: REST/WebSocket endpoints.
4. **Frontend UI**: Next.js/Flutter/Tauri components.
5. **Permission Checks**: Role-based access control.
6. **Validation**: Zod (TypeScript) or Go struct tags.
7. **Audit Logging**: Track changes to entities.
8. **Tests**: Unit, integration, and API tests.
9. **Documentation**: Module-specific docs in `/docs`.

---

## Key Guidelines
### Local-First
- The app must work **offline** and sync when online.
- Use **local storage** for files (e.g., `/storage/clients`).
- Database is the **source of truth**.

### Security
- Use **AES-256 encryption** for sensitive data.
- Implement **audit logs** for all changes.
- Support **soft delete** and **version history**.
- Use **Argon2id** for password hashing.

### Performance
- Optimize queries to avoid **N+1** issues.
- Target **<300ms** page loads.
- Use **Redis** for caching and sessions.

### AI and OCR
- Support **local LLMs** (e.g., Ollama, LM Studio).
- Use **OCRmyPDF** and **PaddleOCR** for document processing.
- Never tie logic to a single model.

---

## Future Work
- **Multi-tenancy**: Support unlimited offices and users.
- **Cloud Deployment**: Optional cloud hosting.
- **Enterprise Features**: Advanced reporting, billing, and compliance.

---

## How to Interact with This Project
1. **Follow Clean Architecture**: Always separate domain, application, infrastructure, and presentation layers.
2. **Write Tests**: Every module must include unit, integration, and API tests.
3. **Document Changes**: Update `/docs` for new modules or features.
4. **Use Task Tools**: Track progress with `TaskCreate` and `TaskUpdate`.
5. **Avoid Technical Debt**: No shortcuts, no duplicate logic, no tightly coupled modules.