# authentication-containerized

Build, run, lint, test
- Build: go build ./...
- Run (local): go run ./...
- Run (Docker): docker compose up --build

Project layout and entrypoints
- Main entry: main.go -> server.StartAuthenticationServer()
- HTTP: Gin router (server/router.go); protected routes under /api with JWT middleware
- DB: Postgres connector in db/postgres.go; CRUD in server/handlers/users.go
- JWT: middleware/jwt.go issues HS256 access/refresh tokens via env secrets

Environment and configuration
- Required env (loaded via github.com/joho/godotenv): PSQL_HOST, PSQL_PORT, PSQL_USER, PSQL_PASSWORD, PSQL_DBNAME, ACCESS_SECRET, REFRESH_SECRET
- .env is mandatory locally; do not commit secrets. In CI, provide via environment or secret store

