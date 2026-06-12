# AGENTS.md — SES Dashboard Monitoring

## Project Structure

- `ses-dashboard-monitoring/` — Go backend (Gin, Go 1.25, module `ses-monitoring`)
- `ses-dashboard-frontend/` — React 19 + TypeScript frontend (Vite, Tailwind CSS)
- `docker-compose.yml` — orchestrates Postgres, backend, frontend
- `Makefile` — root-level build/test/migration commands
- `install.sh` — one-shot setup script (builds images, starts services, runs migrations)

## Quick Start (Docker)

```bash
./install.sh          # builds, starts, migrates, waits for health
```

After install:
- App: http://localhost
- Swagger: http://localhost/swagger/index.html
- Default login: `admin` / `password`

## Local Development

### Backend

```bash
cd ses-dashboard-monitoring
make deps             # go mod download + tidy
make build            # build binary
make run              # build + run
make dev              # hot reload via air (requires .air.toml)
make test             # go test -v ./...
make fmt              # go fmt ./...
make lint             # golangci-lint run
make check            # fmt + lint + test
make swagger          # regenerate swagger docs
```

**Migrations** (root Makefile):
```bash
make migrate-up       # run up migrations
make migrate-down     # rollback one
make migrate-version  # current version
make migrate-create   # interactive create
```

Migration runner: `ses-dashboard-monitoring/cmd/migrate/main.go`
Migration files: `ses-dashboard-monitoring/internal/infrastructure/database/migration/`

### Frontend

```bash
cd ses-dashboard-frontend
npm install
npm run dev           # Vite dev server (port 5173), proxies /api to localhost:8080
npm run build         # tsc -b && vite build
npm run lint          # eslint .
npm start             # node server.js (production Express server)
```

**Toolchain quirks:**
- Vite is overridden to `rolldown-vite@7.2.5` (see `package.json` `overrides`)
- Uses project references: `tsconfig.json` references `tsconfig.app.json` + `tsconfig.node.json`
- Express proxy server (`server.js`) serves `dist/` and proxies `/api`, `/sns`, `/swagger` to backend

## Configuration

Single source of truth: `ses-dashboard-monitoring/config/config.yaml`

Generate `.env` for Docker Compose:
```bash
./generate-env.sh     # parses config.yaml → .env in repo root
```

`.env.example` documents all variables. Key ones:
- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_SSLMODE`
- `APP_PORT` (backend port, default 8080)
- `JWT_SECRET`
- `BACKEND_URL` (frontend proxy target, e.g. `http://backend:8080`)
- `VITE_API_URL` (frontend dev API URL, e.g. `http://localhost:8080`)

## Architecture Notes

- Backend entrypoint: `ses-dashboard-monitoring/cmd/api/main.go`
- Clean architecture: `internal/delivery/http/` → `usecase/` → `domain/` → `infrastructure/`
- Background services: `sync_service.go` (AWS SES suppression sync), `cleanup_service.go` (data retention)
- Frontend dev proxy: `vite.config.ts` proxies `/api` to `localhost:8080`
- Production: frontend Express server proxies `/api`, `/sns`, `/swagger` to `BACKEND_URL`

## Testing & Verification

- Backend: `make test` (root) or `go test -v ./...` (in backend dir)
- Frontend: no test script defined in `package.json`
- No CI workflows found in repo

## Lint / Typecheck

- Backend: `make lint` (golangci-lint), `make fmt` (gofmt)
- Frontend: `npm run lint` (eslint), `npm run build` runs `tsc -b` for type checking

## Important Gotchas

- **Migrations must be run after services start.** `install.sh` does this automatically; for manual Docker, run `make migrate-up` after `docker-compose up -d`.
- **Postgres image:** `postgres:18-alpine` (not 15 as README claims).
- **Backend Makefile loads `.env`:** `ses-dashboard-monitoring/Makefile` includes `../.env` and exports variables. If `.env` is missing, `make migrate-up` may fail.
- **CORS is wide open** in `main.go` (`AllowOriginFunc` returns `true` for all origins) — intentional for Swagger + ngrok, do not tighten without checking docs consumption.
- **No `.air.toml` in repo** but `make dev` references it; install `air` and create config if using hot reload.
- **JWT secret** defaults to hardcoded value in `config.yaml` — must be changed for production.

## Commands Summary

| Task | Command |
|------|---------|
| Full setup | `./install.sh` |
| Start services | `docker-compose up -d` |
| Run migrations | `make migrate-up` |
| Backend build | `make build` |
| Backend test | `make test` |
| Backend lint | `make lint` |
| Backend all checks | `make check` |
| Frontend dev | `cd ses-dashboard-frontend && npm run dev` |
| Frontend build | `cd ses-dashboard-frontend && npm run build` |
| Frontend lint | `cd ses-dashboard-frontend && npm run lint` |
| Generate swagger | `make swagger` |
| Generate .env | `./generate-env.sh` |
