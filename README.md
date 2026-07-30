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

## Getting Started
1. Clone the repository:
   ```bash
   git clone https://github.com/ieraHQ/Vakalat.git
   cd Vakalat
   ```
2. Set up the backend:
   ```bash
   cd backend/api
   go mod init github.com/ieraHQ/Vakalat/backend/api
   go mod tidy
   ```
3. Set up the frontend:
   ```bash
   cd ../../apps/web
   npm install
   ```
4. Set up the database:
   ```bash
   cd ../../database/migrations
   # Run migrations (TBD)
   ```

## Documentation
- [Architecture](./docs/architecture.md)
- [Database Schema](./database/schema/schema.md)
- [API Documentation](./docs/api.md)
- [Deployment Guide](./docs/deployment.md)

## License
MIT