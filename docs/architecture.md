# Vakalat Architecture

This document outlines the architecture of Vakalat, an enterprise-grade Legal Case Management Platform. It covers the system's high-level design, key components, and interactions.

## Table of Contents
- [Overview](#overview)
- [Architecture Principles](#architecture-principles)
- [System Architecture](#system-architecture)
- [Component Breakdown](#component-breakdown)
  - [Backend](#backend)
  - [Frontend](#frontend)
  - [Database](#database)
  - [AI and OCR](#ai-and-ocr)
  - [Search](#search)
- [Data Flow](#data-flow)
- [Deployment](#deployment)
- [Diagrams](#diagrams)
  - [System Architecture Diagram](#system-architecture-diagram)
  - [Database Schema](#database-schema)
  - [Data Flow Diagram](#data-flow-diagram)

---

## Overview
Vakalat is a **local-first, offline-first** Legal Case Management Platform designed for law firms. It supports:
- **Cross-platform deployment**: Web, desktop (Tauri), and mobile (Flutter).
- **Enterprise-grade security**: AES-256 encryption, audit logs, and zero-trust architecture.
- **Modular design**: Clean Architecture with repository pattern, dependency injection, and service layers.
- **AI and OCR**: Local LLM integration and document processing.

---

## Architecture Principles
1. **Clean Architecture**: Separation of concerns with distinct layers (Domain, Application, Infrastructure, Presentation).
2. **Local-First**: Full functionality without an internet connection, with sync capabilities when online.
3. **Modular Monorepo**: Shared packages, backend, frontend, and mobile apps in a single repository.
4. **Performance**: Optimized for sub-300ms page loads.
5. **Security**: AES-256 encryption, audit logs, and role-based access control (RBAC).

---

## System Architecture
Vakalat follows a **client-server architecture** with the following high-level components:
- **Client**: Web (Next.js), Desktop (Tauri), and Mobile (Flutter) applications.
- **Server**: Go backend (Fiber/Gin) with REST and WebSocket APIs.
- **Database**: PostgreSQL with extensions for full-text search (`pg_trgm`) and semantic search (`pgvector`).
- **AI/OCR**: Local LLM and OCR services for document processing.
- **Storage**: Local file storage for documents and evidence.

---

## Component Breakdown
### Backend
- **Language**: Go (Fiber/Gin).
- **Layers**:
  - **Domain**: Business logic and entities (e.g., `User`, `Matter`, `Client`).
  - **Application**: Service layer (e.g., `UserService`, `MatterService`) and dependency injection.
  - **Infrastructure**: Database implementations (e.g., `PostgresUserRepository`) and external services.
  - **Presentation**: REST API and WebSocket handlers.
- **Middleware**: Authentication, RBAC, and request logging.

### Frontend
- **Web**: Next.js (React, TypeScript, TailwindCSS).
- **Desktop**: Tauri (macOS, Windows, Linux).
- **Mobile**: Flutter (Android, iOS, iPad).
- **UI Components**: `shadcn/ui` for reusable components.
- **Data Fetching**: TanStack Query for efficient data fetching.
- **Validation**: React Hook Form + Zod.

### Database
- **PostgreSQL 17**: Primary database with extensions:
  - `pgvector`: Semantic search.
  - `pg_trgm`: Full-text search.
  - `uuid-ossp`: UUID generation.
- **Schema**: 60+ tables for users, clients, matters, documents, billing, and more.
- **Migrations**: Managed with `golang-migrate` (paired `.up.sql`/`.down.sql` files in `database/migrations/`).

### AI and OCR
- **Local LLM**: Integration with Ollama, LM Studio, or OpenAI-compatible endpoints.
- **OCR**: OCRmyPDF, PaddleOCR, and Tesseract for document text extraction.
- **Background Workers**: Process OCR jobs and AI tasks asynchronously.

### Search
- **PostgreSQL Full-Text Search**: Keyword-based search.
- **pgvector**: Semantic search for documents and matters.
- **Hybrid Search**: Combines keyword and semantic search for improved results.

---

## Data Flow
1. **User Interaction**: User interacts with the client (web, desktop, or mobile).
2. **API Request**: Client sends a request to the backend (REST or WebSocket).
3. **Authentication**: Backend validates the request using JWT and RBAC.
4. **Business Logic**: Backend processes the request using the service layer.
5. **Database**: Backend interacts with PostgreSQL for data persistence.
6. **Response**: Backend sends a response to the client.
7. **Real-Time Updates**: WebSocket pushes updates to clients (e.g., hearing reminders).
8. **AI/OCR**: Background workers process documents and AI tasks.

---

## Deployment
Vakalat supports multiple deployment strategies:
- **Local Deployment**: Run the entire stack locally for a single user or office.
- **On-Premises**: Deploy on a private server for multiple users.
- **Cloud**: Optional cloud deployment for scalability and high availability.

---

## Diagrams
### System Architecture Diagram
![System Architecture Diagram](diagrams/system-architecture.png)
*High-level overview of Vakalat's architecture, including clients, backend, database, and AI/OCR services.*

### Database Schema
![Database Schema](diagrams/database-schema.png)
*Entity-relationship diagram of the PostgreSQL schema, including tables for users, clients, matters, documents, and billing.*

### Data Flow Diagram
![Data Flow Diagram](diagrams/data-flow.png)
*Flow of data between clients, backend, database, and AI/OCR services.*