# Load Stuffing Calculator

## Requirements

- [Go](https://go.dev/dl/): versi 1.20 atau lebih baru
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

```bash
# Manual go run
go build -o bin/api ./cmd/api
./bin/api

# Opsi jika menggunakan make
make run

# Opsi jika menggunakan air
air
```

Server API akan berjalan di `http://localhost:8080` secara default (sesuai konfigurasi di `.env`).
**Swagger UI:** `http://localhost:8080/docs/index.html`
