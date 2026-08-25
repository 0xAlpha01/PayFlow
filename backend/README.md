# PayFlow API

Go/Gin authentication API backed by MongoDB and Redis.

## Run locally

1. Copy `.env.example` to `.env` and set strong JWT secrets.
2. Start infrastructure: `docker compose up -d mongodb redis`.
3. Install modules: `go mod tidy`.
4. Start API: `go run ./cmd/server`.

The API listens on `http://localhost:8080` and exposes authentication routes under `/api/v1/auth`.
