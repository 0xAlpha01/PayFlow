# PayFlow API

Go/Gin authentication API backed by MongoDB and Redis.

## Run locally

1. Copy `.env.example` to `.env` and set strong JWT secrets.
2. Start infrastructure: `docker compose up -d mongodb redis`.
3. Install modules: `go mod tidy`.
4. Start API: `go run ./cmd/server`.

The API listens on `http://localhost:8080` and exposes authentication routes under `/api/v1/auth`.

## Render and frontend deployment

Deploy the `backend` directory as the Render service root (or use its Dockerfile).
Set the frontend host's `VITE_API_URL` to `https://payflow-37vp.onrender.com/api/v1`
and set Render's `CLIENT_URL` to the exact deployed frontend origin (without a
trailing slash). If you need local development too, separate allowed origins in
`CLIENT_URL` with commas. Redeploy the frontend after changing `VITE_API_URL`:
Vite embeds `VITE_*` values at build time.

## Email verification

Set `SMTP_HOST`, `SMTP_PORT`, `SMTP_USERNAME`, `SMTP_PASSWORD`, and `SMTP_FROM` in `.env` to deliver verification and password-reset emails through an SMTP provider. When SMTP is intentionally unset, verification codes are logged by the backend for local development. Restart the API after changing `.env`, then use **Resend code** for accounts that were created before email delivery was configured.
