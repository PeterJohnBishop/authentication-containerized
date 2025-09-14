# authentication-containerized

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

