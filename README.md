# ilia-users

Users microservice for the ília Financial Challenge.

## Stack

- **Language:** Go 1.22+
- **HTTP Framework:** Gin Gonic
- **ORM:** GORM
- **Database:** PostgreSQL
- **Auth:** JWT (HMAC-HS256)
- **Container:** Docker + Docker Compose

## Overview

This microservice manages user accounts and authentication. It runs on port **3002** and communicates internally with the Wallet microservice (`ilia-wallet`) using a separate JWT secret.

## Setup

### Prerequisites

- Docker & Docker Compose
- Go 1.22+ (for local development)
- `gh` CLI (for gitflow)

### Environment Variables

Copy `.env.example` to `.env` and fill in the values:

```bash
cp .env.example .env
```

| Variable               | Description                          | Required |
|------------------------|--------------------------------------|----------|
| `SERVER_PORT`          | HTTP port (default: 3002)            | No       |
| `DB_HOST`              | PostgreSQL host                      | Yes      |
| `DB_PORT`              | PostgreSQL port (default: 5432)      | No       |
| `DB_USER`              | PostgreSQL user                      | Yes      |
| `DB_PASSWORD`          | PostgreSQL password                  | Yes      |
| `DB_NAME`              | PostgreSQL database name             | Yes      |
| `DB_SSLMODE`           | SSL mode (default: disable)          | No       |
| `ILIACHALLENGE`        | JWT secret for external auth         | Yes      |
| `ILIACHALLENGE_INTERNAL` | JWT secret for internal comms      | Yes      |
| `WALLET_BASE_URL`      | Wallet service base URL              | Yes      |

### Running with Docker

```bash
docker compose up --build
```

The service will be available at `http://localhost:3002`.

### Running locally

```bash
go mod download
go run ./cmd/server
```

## API

| Method | Path        | Auth       | Description        |
|--------|-------------|------------|--------------------|
| GET    | /health     | None       | Health check       |
| POST   | /users      | None       | Create user        |
| GET    | /users      | Bearer JWT | List all users     |
| GET    | /users/:id  | Bearer JWT | Get user by ID     |
| PATCH  | /users/:id  | Bearer JWT | Update user        |
| DELETE | /users/:id  | Bearer JWT | Delete user        |
| POST   | /auth       | None       | Login → JWT token  |

See `docs/swagger/ms-users.yaml` for full OpenAPI spec.

## Testing

```bash
go test ./...
```
