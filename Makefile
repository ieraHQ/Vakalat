# Vakalat — Makefile

# Global Variables
GO_MOD_DIR := backend/api
WEB_DIR := apps/web
DESKTOP_DIR := apps/desktop
MOBILE_DIR := apps/mobile
DATABASE_DIR := database

# Tasks
.PHONY: build
build:
	@echo "Building all modules..."
	cd $(GO_MOD_DIR) && go build -o bin/api
	cd $(WEB_DIR) && npm run build
	cd $(DESKTOP_DIR) && npm run tauri build
	cd $(MOBILE_DIR) && flutter build apk

.PHONY: test
test:
	@echo "Running tests..."
	cd $(GO_MOD_DIR) && go test ./...
	cd $(WEB_DIR) && npm test
	cd $(MOBILE_DIR) && flutter test

.PHONY: lint
lint:
	@echo "Running linters..."
	cd $(GO_MOD_DIR) && golangci-lint run
	cd $(WEB_DIR) && npx eslint .
	cd $(MOBILE_DIR) && flutter analyze

.PHONY: migrate
migrate:
	@echo "Running database migrations..."
	cd $(DATABASE_DIR) && goose postgres "user=postgres dbname=vakalat sslmode=disable" up

.PHONY: seed
seed:
	@echo "Seeding database..."
	cd $(DATABASE_DIR) && psql -U postgres -d vakalat -f seeds/seed.sql

.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	cd $(GO_MOD_DIR) && rm -rf bin/
	cd $(WEB_DIR) && rm -rf .next/
	cd $(DESKTOP_DIR) && rm -rf src-tauri/target/
	cd $(MOBILE_DIR) && flutter clean

.PHONY: run
run:
	@echo "Running all modules..."
	cd $(GO_MOD_DIR) && go run main.go &
	cd $(WEB_DIR) && npm run dev &
	cd $(DESKTOP_DIR) && npm run tauri dev &
	cd $(MOBILE_DIR) && flutter run

.PHONY: help
help:
	@echo "Available tasks:"
	@echo "  build    - Build all modules"
	@echo "  test     - Run tests"
	@echo "  lint     - Run linters"
	@echo "  migrate  - Run database migrations"
	@echo "  seed     - Seed database"
	@echo "  clean    - Clean build artifacts"
	@echo "  run      - Run all modules"
	@echo "  help     - Show this help message"