# PayFlow API

Go/Gin API for the PayFlow frontend. It uses MongoDB for application data and
Redis for refresh tokens and short-lived verification codes.

## Local development

1. Copy `.env.example` to `.env` and set strong, unique JWT secrets.
2. Start MongoDB and Redis: `docker compose up -d mongodb redis`.
3. Run the API: `go run ./cmd/server`.

The API listens on `http://localhost:8080`. All application endpoints use the
`/api/v1` prefix; for example, signup is `POST /api/v1/auth/signup`.

## Health checks

- `GET /` returns `{ "status": "ok" }`.
- `GET /healthz` returns `{ "status": "ok" }`.

Use `/healthz` as the Render health-check path if you configure one.

## Render deployment

Deploy the `backend` directory as the Render service root, or select the
included `Dockerfile`. Render supplies `PORT`; do not hard-code a production
port.

Configure these environment variables in Render:

| Variable | Purpose |
| --- | --- |
| `MONGODB_URI` | Production MongoDB connection string. |
| `MONGODB_DATABASE` | Database name, such as `payflow`. |
| `REDIS_ADDR` | Redis host and port. |
| `REDIS_PASSWORD` | Redis password, if required. |
| `JWT_ACCESS_SECRET` | Long random secret for access tokens. |
| `JWT_REFRESH_SECRET` | Different long random secret for refresh tokens. |
| `CLIENT_URL` | Exact deployed frontend origin, without a trailing slash. |
| `SMTP_*` | SMTP configuration for email verification and password reset. |
| `SMS_*` | Optional phone-verification provider configuration. |
| `PAYSTACK_*` | Paystack keys and base URL. |

`CLIENT_URL` supports a comma-separated list. For example:

```env
CLIENT_URL=https://your-frontend.vercel.app,http://localhost:5173
```

Never commit real credentials or copy them into frontend environment files.

## Frontend deployment

Set this frontend build variable before deploying to Vercel:

```env
VITE_API_URL=https://payflow-37vp.onrender.com/api/v1
```

The `/api/v1` portion is required. Vite embeds `VITE_*` values during its build,
so redeploy the frontend after changing the value. The frontend includes
`vercel.json` to route direct visits and reloads of client-side pages such as
`/dashboard` and `/settings/profile` back to the React application.

## Verification flow

New users are created and signed in immediately; registration does not send a
verification code. Phone verification is optional and can be completed later
from **Settings > Security**. SMTP is still used for password-reset emails.

For email delivery, configure `SMTP_HOST`, `SMTP_PORT`, `SMTP_USERNAME`,
`SMTP_PASSWORD`, and `SMTP_FROM`. When SMTP is intentionally unset locally,
verification codes are written to backend logs for development only.
