# redesigned-telegram

https://github.com/PeterJohnBishop/cautious-dollop/blob/main/README.md

docker build -t peterjbishop/redesigned-telegram:latest . 
docker push peterjbishop/redesigned-telegram:latest 
docker pull peterjbishop/redesigned-telegram:latest 
docker-compose down 
docker-compose build --no-cache 
docker-compose up

Endpoints

• POST /auth/register: Email/password signup with email verification
• POST /auth/login: Username/email + password to issue tokens
• POST /auth/refresh: Rotate refresh tokens (one-time use)
• POST /auth/logout: Revoke refresh token(s)
• POST /auth/forgot-password: Send reset email
• POST /auth/reset-password: Complete reset with token
• GET /.well-known/jwks.json: JWKs for public key discovery (if JWT signed with asymmetric keys)
• GET /auth/me: Return current user profile (access token required)
    • Optional SSO (later): /auth/oidc/google, /auth/oidc/github

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