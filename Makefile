DATABASE_URL ?= postgres://jobfinder:jobfinder_local@localhost:55432/jobfinder?sslmode=disable

.PHONY: help setup up down logs dry-run test browser-install browser-build browser-test browser-health verify fmt

help:
	@echo "ai-analysis Phase 1"
	@echo "  make setup          Download Go and browser-worker dependencies"
	@echo "  make up             Start PostgreSQL and pgweb"
	@echo "  make dry-run        Validate config/profile, migrate DB, lock, and save a run"
	@echo "  make verify         Run Go and TypeScript tests"
	@echo "  make down           Stop local services"
	@echo "  make browser-health Check the TypeScript worker protocol"

setup:
	go mod download
	cd browser-worker && npm install

up:
	docker compose up -d --wait

down:
	docker compose down

logs:
	docker compose logs -f

dry-run:
	DATABASE_URL="$(DATABASE_URL)" go run ./cmd/jobfinder run --dry-run

test:
	go test ./...

browser-install:
	cd browser-worker && npm install

browser-build:
	cd browser-worker && npm run build

browser-test:
	cd browser-worker && npm test

browser-health:
	cd browser-worker && npm run build && npm run health

verify: test browser-test

fmt:
	gofmt -w cmd internal
