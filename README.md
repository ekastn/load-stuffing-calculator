# Load Stuffing Calculator

## Requirements

- [Go](https://go.dev/dl/): versi 1.24 atau lebih baru
- [PostgreSQL](https://www.postgresql.org/download/): versi 14 atau lebih baru
- [Goose](https://github.com/pressly/goose): Migrasi database
- [swag](https://github.com/swaggo/swag): Swagger Generator (opsional)
- [SQLC](https://github.com/sqlc-dev/sqlc): Auto-gen SQL (opsional)
- [air](https://github.com/air-verse/air): Live Reload (opsional)

## Setup

```bash
# Clone Repository
git clone https://github.com/ekastn/load-stuffing-calculator.git
cd load-stuffing-calculator

# Install Dependencies
go mod tidy

# Konfigurasi Environment Variables
cp .env.example .env # edit sesuai kebutuhan

# Menjalankan Migrasi Database
goose -dir cmd/db/migrations postgres "$(grep ^DATABASE_URL .env | cut -d '=' -f2-)" up

# Generate SQLC Code (jika ada perubahan query)
sqlc generate

# Generate Dokumentasi Swagger (jika ada perubahan endpoint)
swag init -g cmd/api/main.go -o internal/docs
```

## Run

### Option A: Docker Compose (recommended)

This starts Postgres, runs migrations, starts the packing microservice (py3dbp), then starts the Go API.

```bash
# If you cloned fresh, fetch the vendored py3dbp submodule
git submodule update --init --recursive

docker compose up --build
```

API:
- `http://localhost:8080`
- Swagger UI: `http://localhost:8080/docs/index.html`

Web UI:
- `http://localhost:3000`

Packing service (internal, for debugging only):
- Health: `http://localhost:5051/health`

### Option B: Local (manual)

You must run the packing microservice separately and set `PACKING_SERVICE_URL`.

```bash
# terminal 1: packing service
python3 -m venv cmd/packing/.venv
. cmd/packing/.venv/bin/activate
pip install -r cmd/packing/requirements.txt
python cmd/packing/app.py

# terminal 2: API
go build -o bin/api ./cmd/api
./bin/api
```

Default URL if unset: `PACKING_SERVICE_URL=http://localhost:5051`
