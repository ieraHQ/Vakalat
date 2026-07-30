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

### Phase 2 – Database (✅ Complete)
- **PostgreSQL Extensions**: Enabled `pgvector`, `pg_trgm`, and `uuid-ossp`.
- **Schema Design**: Designed 60+ tables for users, clients, matters, documents, billing, and more.
  - **Users & Authentication**: `users`, `roles`, `permissions`, `sessions`.
  - **Clients & Contacts**: `clients`, `contacts`, `organizations`.
  - **Legal Entities**: `courts`, `judges`, `advocates`, `matters`, `hearings`, `orders`.
  - **Documents & OCR**: `documents`, `document_versions`, `embeddings`, `ocr_jobs`.
  - **Tasks & Calendar**: `tasks`, `calendar_events`, `reminders`.
  - **Billing & Expenses**: `invoices`, `invoice_items`, `expenses`, `payments`.
  - **AI & Search**: `ai_sessions`, `ai_summaries`, `search_index`.
  - **Audit & Compliance**: `audit_logs`, `backup_logs`.
- **Indexes**: Added indexes for performance-critical columns (e.g., `matter_id`, `client_id`).
- **Migrations**: Added `goose` migration system and SQL migrations for all tables.
- **Seed Data**: Added sample data for testing (roles, users, clients).

### Phase 3 – Backend Framework (✅ Complete)
- **Logger**: Configured structured logging with Zap.
- **Database**: Set up PostgreSQL connection pooling with `pgxpool`.
- **Repositories**: Implemented `UserRepository`, `ClientRepository`, `MatterRepository`, and `DocumentRepository`.
- **Services**: Implemented `UserService`, `ClientService`, `MatterService`, and `DocumentService`.
- **Middleware**: Added authentication, RBAC, and request logging middleware.
- **WebSocket**: Set up real-time updates for hearings, tasks, and notifications.
- **API Endpoints**: Integrated services with Fiber routes for users, clients, matters, and documents.

### Phase 4 – Authentication & RBAC (✅ Complete)
- **JWT Authentication**: Secure API access with JWT tokens.
- **Refresh Tokens**: Support long-lived sessions.
- **Password Hashing**: Use Argon2id for secure password storage.
- **Role-Based Permissions**: Fine-grained access control for users and roles.
- **Session Management**: Password reset and MFA readiness.

### Phase 5 – Document Engine (✅ In Progress)
- **File Upload**: Let lawyers upload PDFs, Word docs, and other files.
- **Local Storage**: Store files in `/storage/documents`.
- **Next Steps**:
  - Add OCR to read text inside documents (English + Hindi).
  - Support document versioning (like saving different drafts).
  - Add digital signatures to verify documents.

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