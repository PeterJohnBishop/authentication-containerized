# redesigned-telegram

Build, run, lint, test
- Build: go build ./...
- Run (local): go run ./...
- Run (Docker): docker compose up --build
- Lint (stdlib): go vet ./...
- Format: go fmt ./...
- Recommended linter: golangci-lint run (if installed)
- Test all: go test ./...
- Test with race: go test -race ./...
- Single package: go test ./server/handlers -v
- Single test: go test -run ^TestName$ ./server/handlers -v

Project layout and entrypoints
- Main entry: main.go -> server.StartAuthenticationServer()
- HTTP: Gin router (server/router.go); protected routes under /api with JWT middleware
- DB: Postgres connector in db/postgres.go; CRUD in server/handlers/users.go
- JWT: middleware/jwt.go issues HS256 access/refresh tokens via env secrets

Environment and configuration
- Required env (loaded via github.com/joho/godotenv): PSQL_HOST, PSQL_PORT, PSQL_USER, PSQL_PASSWORD, PSQL_DBNAME, ACCESS_SECRET, REFRESH_SECRET
- .env is mandatory locally; do not commit secrets. In CI, provide via environment or secret store

Code style guidelines (Go)
- Imports: standard -> third-party -> internal; use goimports formatting; alias only on conflict
- Formatting: enforce go fmt; no custom comments or banners; keep lines concise
- Types & structs: exported names in CamelCase; json tags are lower_snake or lowerCamel consistently (current code uses lowerCamel)
- Errors: return errors, wrap with context using fmt.Errorf("...: %w", err); do not log.Fatal inside libraries (prefer returning errors). Only main/server may exit
- Logging: prefer log.Println for server start and notable events; avoid printing secrets; do not use fmt.Println for operational logs in handlers (prefer log)
- Handlers: validate input with c.ShouldBindJSON; respond with proper HTTP status; never leak internal error details
- Context: use gin.Context-aware DB calls (QueryContext/ExecContext) as present; propagate context when possible
- JWT: HS256 with ACCESS_SECRET/REFRESH_SECRET; access ~15m, refresh ~7d; read secrets from env; never hardcode
- SQL: use parameterized queries (as done with $1, $2...); check sql.ErrNoRows distinctly
- Naming: package names are lower_snake; identifiers in lowerCamel for locals, UpperCamel for exported

Testing guidance
- Create *_test.go per package; table-driven tests preferred
- For handlers, use httptest + gin.New(); stub DB with a test database or interfaces; assert status codes and JSON
- Seed a temporary Postgres for integration via docker compose; run go test with DSN pointed to test DB

# notes

https://github.com/PeterJohnBishop/cautious-dollop/blob/main/README.md

docker build -t peterjbishop/redesigned-telegram:latest . 
docker push peterjbishop/redesigned-telegram:latest 
docker pull peterjbishop/redesigned-telegram:latest 
docker-compose down 
docker-compose build --no-cache 
docker-compose up

Endpoints

• POST /auth/register: Email/password signup
• POST /auth/login: Username/email + password to issue tokens
• POST /auth/refresh: Rotate refresh tokens (one-time use)
• POST /auth/logout: Revoke refresh token(s) <-- TBD
• POST /auth/forgot-password: Send reset email <-- TBD
• POST /auth/reset-password: Complete reset with token <-- TBD
• GET /.well-known/jwks.json: JWKs for public key discovery (if JWT signed with asymmetric keys) <-- TBD
• GET /auth/me: Return current user profile (access token required) <-- TBD
    • Optional SSO (later): /auth/oidc/google, /auth/oidc/github  <-- TBD

Security

• Password policy: min length, common password blacklist, breach check (k‑Anon HIBP optional)
• Tokens
• Access token: short-lived JWT (5–10 min), asymmetric signing (RS256/ES256), include aud/iss/sub/jti/iat/exp
• Refresh token: long-lived, opaque, one-time use rotation with family tracking and reuse detection
• Token binding: store refresh token hash (not plaintext) with device info and IP fingerprint
• Sessions/Revocation
• Maintain session store with status, last used, and revocation timestamp; invalidate on logout/compromise
• Detect refresh token reuse and revoke entire token family
• Email verification & resets
• Time-limited, single-use tokens stored hashed; rate limit sends per user/IP
• Transport & headers
• Enforce TLS, HSTS, secure cookies (HttpOnly, SameSite=Lax/Strict), CSRF for cookie flows
• Brute force protection
• IP + account-based rate limiting and lockouts with exponential backoff
• Auditing
• Append-only audit log for auth events: login success/failure, password changes, token reuse, etc.
• Observability
• Structured logging without PII; metrics for auth events; tracing for critical flows
• Key management
• Key rotation for JWT signing; JWKS publishing; KMS/HSM-backed keys (AWS KMS) if possible

Postgres Data model 

• users
• id (uuid), email (unique), email_verified_at, password_hash, password_salt, password_algo, password_peppered
(implicit), created_at, updated_at, mfa_enabled, mfa_secret (if TOTP)
• sessions (refresh token family)
• id (uuid), user_id, token_hash, family_id, status (active/revoked/compromised), created_at, expires_at,
last_used_at, user_agent, ip_hash, device_name
• email_verifications
• user_id, token_hash, expires_at, created_at, used_at
• password_resets
• user_id, token_hash, expires_at, created_at, used_at
• audit_log
• id, user_id (nullable), event_type, metadata(json), created_at

Tech stack (Go-based)

• HTTP: Gin
• Crypto: golang.org/x/crypto (argon2id), crypto/rand, jose or go-jose/v4 (JWT JWS/JWK)
• DB: Postgres (sqlc/GORM) for strong consistency or DynamoDB for serverless; Redis for rate limiting/session cache
• Email: AWS SES
• Observability: zap/zerolog, Prometheus, OpenTelemetry
• Secrets/Config: Viper or custom env loader, AWS SSM/Secrets Manager in prod
• Migrations (if Postgres): golang-migrate

• Registration
• Validate email/password → hash with Argon2id (salt + pepper) → store → issue verification email
• Login
• Verify password → create session record with device/IP → issue access JWT + new refresh (store hash) → set secure
cookies (or JSON body)
• Refresh
• Validate refresh token → check hash + session status → rotate: invalidate old, create new with same family → issue
new access + refresh
• If reuse detected: mark family compromised; delete/revoke all; require re-login
• Logout
• Revoke current session (and optionally all sessions for “logout all”)
• Password reset
• Generate one-time token, store hashed → email link → verify and allow new password; revoke existing sessions on
success

Operational features

• Healthz/Readyz endpoints
• Graceful shutdown
• Configuration via env with strong validation
• Dockerfile + docker-compose (db, redis, mailhog)
• CI: go fmt/vet/test, golangci-lint, trivy scan