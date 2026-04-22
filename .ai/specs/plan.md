# ilia-users: Implementation Plan

## Context

This is the **Users microservice (Part 2)** of the ília Financial Microservices Challenge.
The wallet service (`../ilia-wallet`) is already complete and serves as the reference
architecture. This service manages users (CRUD) and authentication, running on port **3002**.
It communicates with the wallet internally via REST using a separate JWT secret.

The implementation **mirrors ilia-wallet's Clean Architecture**: gin-gonic + GORM + PostgreSQL
+ golang-jwt/jwt/v5, following the Uber Go style guide and the unit-test rules in
`.ai/rules/developing-unit-tests.md`.

---

## Directory Structure

```
ilia-users/
├── cmd/server/main.go
├── internal/
│   ├── domain/user.go
│   ├── usecase/
│   │   ├── create_user.go + _test.go
│   │   ├── get_user.go + _test.go
│   │   ├── list_users.go + _test.go
│   │   ├── update_user.go + _test.go
│   │   ├── delete_user.go + _test.go
│   │   └── authenticate_user.go + _test.go
│   ├── adapter/
│   │   ├── http/
│   │   │   ├── handler/
│   │   │   │   ├── user.go + _test.go + mock_test.go
│   │   │   │   └── auth.go + _test.go + mock_test.go
│   │   │   └── middleware/auth.go + auth_test.go
│   │   ├── repository/postgres/user.go
│   │   └── client/wallet/client.go
│   └── infrastructure/
│       ├── bootstrap/user.go
│       ├── config/config.go
│       └── database/
│           ├── postgres.go
│           └── migrations/
│               ├── 000001_create_users.up.sql
│               └── 000001_create_users.down.sql
├── pkg/
│   ├── apperrors/errors.go
│   └── jwtutil/token.go
├── docker/Dockerfile
├── docs/swagger/ms-users.yaml
├── .ai/specs/plan.md  (this file, committed to repo)
├── docker-compose.yml
├── .env.example
├── go.mod
└── README.md
```

---

## Endpoints (from docs/swagger/ms-users.yaml)

| Method | Path        | Auth            | Description        |
|--------|-------------|-----------------|--------------------|
| GET    | /health     | None            | Health check       |
| POST   | /users      | None            | Create user        |
| GET    | /users      | Bearer JWT      | List all users     |
| GET    | /users/:id  | Bearer JWT      | Get user by ID     |
| PATCH  | /users/:id  | Bearer JWT      | Update user fields |
| DELETE | /users/:id  | Bearer JWT      | Delete user        |
| POST   | /auth       | None            | Login → JWT token  |

---

## JWT Rules

| Scenario                    | Env Var                  | Notes                                  |
|-----------------------------|--------------------------|----------------------------------------|
| External auth (clients)     | `ILIACHALLENGE`          | Used in POST /auth response + middleware |
| Internal wallet calls       | `ILIACHALLENGE_INTERNAL` | Added to outbound HTTP headers to wallet |

Token claims: `user_id` (UUID string), `email`, `exp` (24h), `iat`.

---

## Implementation Steps (Gitflow: feature branch → PR → merge to main)

---

### Step 0 — GitHub Repository & Gitflow Init

**Branch:** main (direct init commit)

```bash
gh repo create silvioubaldino/ilia-users --public --clone
cd ilia-users
git init
# initial commit: .gitignore, go.mod placeholder, README skeleton
git push -u origin main
```

---

### Step 1 — Project Scaffolding & Config

**Branch:** `feature/project-setup`

**Files:**
- `go.mod` — module `github.com/silvioubaldino/ilia-users`, Go 1.22+
- `go.sum`
- `.env.example`
- `pkg/apperrors/errors.go` — same sentinel errors as wallet (ErrNotFound, ErrUnauthorized, ErrConflict, ErrInvalidInput)
- `internal/infrastructure/config/config.go`

**Config env vars:**
```
SERVER_PORT=3002
DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME, DB_SSLMODE
ILIACHALLENGE          # external JWT secret (required)
ILIACHALLENGE_INTERNAL # internal JWT secret (required)
WALLET_BASE_URL        # e.g. http://ilia-wallet:3001
```

**Go dependencies:**
```
github.com/gin-gonic/gin
gorm.io/gorm
gorm.io/driver/postgres
github.com/golang-jwt/jwt/v5
github.com/golang-migrate/migrate/v4
github.com/google/uuid
github.com/joho/godotenv
github.com/stretchr/testify
golang.org/x/crypto        ← bcrypt
```

---

### Step 2 — Domain Layer

**Branch:** `feature/domain-user`

**File:** `internal/domain/user.go`

```go
type User struct {
    ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
    FirstName string    `gorm:"not null"`
    LastName  string    `gorm:"not null"`
    Email     string    `gorm:"uniqueIndex;not null"`
    Password  string    `gorm:"not null"` // bcrypt hash
    CreatedAt time.Time `gorm:"not null;default:now()"`
}
```

---

### Step 3 — Database & Migrations

**Branch:** `feature/database-setup`

**Files:**
- `internal/infrastructure/database/postgres.go` — GORM init + migrate runner (same pattern as wallet)
- `internal/infrastructure/database/migrations/000001_create_users.up.sql`:
  ```sql
  CREATE EXTENSION IF NOT EXISTS pgcrypto;
  CREATE TABLE IF NOT EXISTS users (
      id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
      first_name VARCHAR NOT NULL,
      last_name  VARCHAR NOT NULL,
      email      VARCHAR NOT NULL UNIQUE,
      password   VARCHAR NOT NULL,
      created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
  );
  ```
- `internal/infrastructure/database/migrations/000001_create_users.down.sql`:
  ```sql
  DROP TABLE IF EXISTS users;
  ```

---

### Step 4 — Repository Layer

**Branch:** `feature/repository-user`

**File:** `internal/adapter/repository/postgres/user.go`

Repository interface (defined in usecase package, implemented here):
```go
type UserRepository interface {
    Create(ctx, user User) (User, error)
    GetByID(ctx, id UUID) (User, error)
    GetByEmail(ctx, email string) (User, error)
    List(ctx) ([]User, error)
    Update(ctx, id UUID, updates User) (User, error)
    Delete(ctx, id UUID) error
}
```

Error mapping: GORM's `ErrRecordNotFound` → `apperrors.ErrNotFound`; unique violation → `apperrors.ErrConflict`.

---

### Step 5 — Use Cases

**Branch:** `feature/usecases`

One file + one test file per use case:

| Use Case           | Key Logic                                                              |
|--------------------|------------------------------------------------------------------------|
| `CreateUser`       | Hash password with bcrypt cost 12; call repo.Create; return sanitized user |
| `GetUser`          | repo.GetByID; propagate ErrNotFound                                   |
| `ListUsers`        | repo.List; return slice (never nil)                                    |
| `UpdateUser`       | If password field set, re-hash; repo.Update                           |
| `DeleteUser`       | repo.GetByID to verify exists; repo.Delete                            |
| `AuthenticateUser` | repo.GetByEmail; bcrypt.CompareHashAndPassword; call jwtutil.GenerateToken |

**File:** `pkg/jwtutil/token.go`
```go
func GenerateToken(userID uuid.UUID, email string, secret string) (string, error)
// Claims: user_id, email, exp (+24h), iat
// Algorithm: HMAC-HS256
```

**Test rules (from .ai/rules/developing-unit-tests.md):**
- Table-driven with `map[string]struct{input; mocks; expected}`
- Testify assertions + mock.AssertExpectations
- Package `<feature>_test`; mocks in `mock_test.go`

---

### Step 6 — HTTP Middleware & Handlers

**Branch:** `feature/http-handlers`

**Middleware:** `internal/adapter/http/middleware/auth.go`
- Identical to wallet's middleware — validates Bearer JWT signed with `ILIACHALLENGE`
- Injects `user_id` into gin context

**Handlers:**

`internal/adapter/http/handler/user.go` — `UserHandler`:
- `Create(c)` → calls CreateUser UC, returns 201 with `UsersResponse` (no password)
- `List(c)` → calls ListUsers UC, returns 200 with `[]UsersResponse`
- `Get(c)` → parses `:id`, calls GetUser UC, returns 200 or 404
- `Update(c)` → parses `:id`, binds partial body, calls UpdateUser UC
- `Delete(c)` → parses `:id`, calls DeleteUser UC, returns 204

`internal/adapter/http/handler/auth.go` — `AuthHandler`:
- `Login(c)` → binds `{email, password}`, calls AuthenticateUser UC, returns `{user, access_token}`

Response type `usersResponse`:
```go
type usersResponse struct {
    ID        string `json:"id"`
    FirstName string `json:"first_name"`
    LastName  string `json:"last_name"`
    Email     string `json:"email"`
}
```

---

### Step 7 — Wallet REST Client

**Branch:** `feature/wallet-integration`

**File:** `internal/adapter/client/wallet/client.go`

- Generates short-lived JWT signed with `ILIACHALLENGE_INTERNAL` for each request
- Methods: `GetBalance(ctx, userID) (int64, error)`, `GetTransactions(ctx, userID) ([]Transaction, error)`
- Used by future use cases that need wallet data enrichment (e.g., user profile summary)

---

### Step 8 — Bootstrap / Dependency Injection

**Branch:** merged with `feature/http-handlers`

**File:** `internal/infrastructure/bootstrap/user.go`

```go
func SetupUser(db *gorm.DB, cfg *config.Config, public gin.IRouter, auth gin.IRouter) {
    repo    := postgresrepo.NewUserRepository(db)
    walletClient := wallet.NewClient(cfg.WalletBaseURL, cfg.JWTInternalSecret)
    // wire each use case
    // wire handlers
    // register routes on public and auth groups
}
```

**File:** `cmd/server/main.go`

```go
func main() {
    cfg := config.Load()
    db  := database.NewPostgres(cfg.DSN(), cfg.DatabaseURL())
    r   := gin.Default()
    r.GET("/health", ...)
    public := r.Group("/")
    authGroup := r.Group("/")
    authGroup.Use(middleware.Auth(cfg.JWTSecret))
    bootstrap.SetupUser(db, cfg, public, authGroup)
    r.Run(":" + cfg.ServerPort)
}
```

---

### Step 9 — Docker & Compose

**Branch:** `feature/docker`

**`docker/Dockerfile`** — multi-stage:
```dockerfile
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o server ./cmd/server

FROM alpine:3.19
COPY --from=builder /app/server /server
EXPOSE 3002
ENTRYPOINT ["/server"]
```

**`docker-compose.yml`:**
```yaml
services:
  app:
    build: { context: ., dockerfile: docker/Dockerfile }
    ports: ["3002:3002"]
    env_file: .env
    depends_on: { postgres: { condition: service_healthy } }
  postgres:
    image: postgres:16-alpine
    environment: { POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_DB }
    healthcheck: { test: pg_isready, interval: 5s }
    ports: ["5432:5432"]
```

---

### Step 10 — Documentation & Final PR

**Branch:** `feature/docs`

- `README.md` — prerequisites, env vars table, `docker compose up`, manual API examples with curl
- Verify `docs/swagger/ms-users.yaml` matches implementation (already exists)
- Copy plan to `.ai/specs/plan.md` in repo

---

## Critical Files (reference from ilia-wallet)

| Pattern            | Source File (wallet)                                              |
|--------------------|-------------------------------------------------------------------|
| Domain model       | `../ilia-wallet/internal/domain/transaction.go`                  |
| Config             | `../ilia-wallet/internal/infrastructure/config/config.go`        |
| DB + migrations    | `../ilia-wallet/internal/infrastructure/database/postgres.go`    |
| Auth middleware    | `../ilia-wallet/internal/adapter/http/middleware/auth.go`        |
| Handler pattern    | `../ilia-wallet/internal/adapter/http/handler/transaction.go`    |
| Bootstrap DI       | `../ilia-wallet/internal/infrastructure/bootstrap/transaction.go`|
| Use case pattern   | `../ilia-wallet/internal/usecase/create_transaction.go`          |
| Apperrors          | `../ilia-wallet/pkg/apperrors/errors.go`                         |
| main.go            | `../ilia-wallet/cmd/server/main.go`                              |

---

## Verification Checklist

1. `docker compose up` — both services start, migrations run automatically
2. `POST /users` body `{first_name, last_name, email, password}` → 201 (no password in response)
3. `POST /auth` body `{email, password}` → 200 with `{user, access_token}`
4. `GET /users` with `Authorization: Bearer <token>` → 200 list
5. `GET /users/:id` → 200 user or 404
6. `PATCH /users/:id` with partial body → 200 updated user
7. `DELETE /users/:id` → 204
8. `go test ./...` — all unit tests pass
9. Wallet client: verify internal JWT header is sent on outbound requests
10. Invalid token → 401 on protected routes
