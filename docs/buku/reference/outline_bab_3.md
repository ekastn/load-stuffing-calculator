# Outline Bab 3: Pengembangan Backend API dengan Bahasa Go

## Overview

**Tujuan Bab**: Membangun Backend API dari nol menggunakan bahasa Go, mulai dari fondasi bahasa hingga integrasi dengan Packing Service.

**Prasyarat**: Pembaca telah memahami arsitektur sistem dari Bab 2.

---

## Intro Bab 3

- Tujuan bab: membangun Backend API yang menjadi "pintu gerbang" sistem
- Tantangan teknis:
  - Mendesain database schema untuk entitas logistik
  - Membangun API yang scalable dan maintainable
  - Mengintegrasikan dengan Packing Service (Python)
- Topik yang akan dibahas (roadmap)

---

## 3.1 Fondasi Bahasa Go

### Prinsip Desain Go
- Simplicity over cleverness
- Readability over terseness
- Explicit over implicit
- Composition over inheritance
- *(Dalam narasi: referensi ke Effective Go, Go Proverbs, Go Code Review Comments)*

### Concurrency Model
- Goroutines: lightweight threads
- Channels: communication between goroutines
- "Don't communicate by sharing memory; share memory by communicating"
- Contoh: bagaimana HTTP server Go menangani ribuan koneksi

### Error Handling Philosophy
- Errors are values (bukan exceptions)
- Multiple return values: `func DoSomething() (Result, error)`
- Wrapping errors dengan `fmt.Errorf("context: %w", err)`
- Sentinel errors vs custom error types

### Packages dan Modules
- Package organization dan naming
- `go.mod` dan dependency management
- Internal packages untuk encapsulation

---

## 3.2 Idiom dan Konvensi Go

### Naming Conventions
- Short, meaningful names (i, n, err, ctx)
- MixedCaps, bukan snake_case
- Exported vs unexported (capital letter)
- Acronyms: HTTP, URL, ID (bukan Http, Url, Id)

### Interface Design
- "Accept interfaces, return structs"
- Small interfaces: io.Reader, io.Writer sebagai contoh ideal
- Interface segregation: define interfaces where they are used

### Struct Design
- Embedding untuk composition
- Constructor functions: `NewXxx()` pattern
- Options pattern untuk konfigurasi kompleks
- Zero value yang berguna

### Error Handling Patterns
- Check error immediately setelah function call
- Wrap dengan konteks yang bermakna
- Don't panic in library code

### Testing Patterns
- Table-driven tests
- Test fixtures dan test data
- Subtests dengan `t.Run()`
- Example functions untuk dokumentasi

### Code Examples: Bad vs Good
- Refactoring real code dari proyek ini
- Menunjukkan transformasi dari kode "works" ke kode "idiomatic"

---

## 3.3 Setup Proyek dan Tooling

### Inisialisasi Module
- `go mod init github.com/username/load-stuffing-calculator`
- Import path conventions
- `go mod tidy` untuk cleanup dependencies

### Struktur Direktori
- Referensi kembali ke Bab 2.5
- Penjelasan detail:
  - `cmd/` - entry points
  - `internal/` - private packages
  - `pkg/` - public packages (jika ada)

### Development Tools
- `air`: hot reload untuk development workflow
- `go test`: built-in testing framework
- `go generate`: untuk SQLC code generation
- IDE setup (VS Code dengan Go extension)

---

## 3.4 Database Schema Design

### Entity-Relationship Diagram
- **Container**: id, name, length, width, height, max_weight
- **Product**: id, label, length, width, height, weight
- **Plan**: id, container_id, status, created_at, calculated_at
- **Placement**: plan_id, product_id, pos_x, pos_y, pos_z, rotation, step_number
- Relasi antar entitas dengan diagram Mermaid

### Migrasi Database dengan Goose
- Mengapa Goose?
  - Simple dan straightforward
  - Support SQL dan Go migrations
  - Embedded migrations
- Struktur file migrasi
- Contoh: `20240101000001_create_containers_table.sql`
- Commands: `goose up`, `goose down`, `goose status`

### PostgreSQL Driver: pgx
- Mengapa pgx vs database/sql
  - Native PostgreSQL support
  - Better performance
  - JSON/JSONB support
- Connection pooling dengan `pgxpool`
- Context dan timeout untuk query

---

## 3.5 Type-safe Repository Layer dengan SQLC

### Mengapa SQLC?
- Write SQL, generate Go code
- Compile-time type safety
- No runtime reflection overhead
- SQL tetap SQL (tidak perlu belajar ORM DSL)

### Setup SQLC
- `sqlc.yaml` configuration
- Query files organization
- Generated code structure

### Menulis Queries
- Annotations: `-- name: GetContainer :one`
- CRUD patterns
- Batch operations
- Parameterized queries

### Generated Code
- Models dari schema
- Querier interface
- Type-safe function signatures
- Contoh output dan cara penggunaan

### Testing Repository Layer
- Database test fixtures
- Transaction rollback untuk isolation
- Test dengan database real vs mock

---

## 3.6 Business Logic Layer (Service)

### Arsitektur Berlapis
- Handler → Service → Repository
- Dependency injection melalui constructor
- Interface untuk loose coupling
- Mengapa Service layer penting?

### Service Implementations
- **ContainerService**: CRUD operations
- **ProductService**: CRUD operations
- **PlanService**: Plan management + orchestration

### PackingService: Orchestrator Utama
- Mengumpulkan data dari repositories
- Preparing request untuk Packing Service (Python)
- Calling PackingGateway
- Transforming dan menyimpan results

### Unit Testing Services
- Mocking repository dengan interface
- Mocking gateway untuk external calls
- Test coverage untuk happy path dan error cases
- Test fixtures

---

## 3.7 HTTP Router dan Handler

### Gin Router
- Mengapa Gin?
  - Performance tinggi
  - Middleware ecosystem matang
  - Dokumentasi lengkap
  - Komunitas besar
- Route definition: `router.GET()`, `router.POST()`, dll
- Route groups untuk versioning: `/api/v1/`
- URL parameters: `/containers/:id`
- Query strings: `?limit=10&offset=0`

### Middleware Stack
- Logging: `gin.Logger()`
- Recovery: `gin.Recovery()`
- CORS configuration
- Request ID (opsional)

### Handler Implementation
- Inject service ke handler
- `c.ShouldBindJSON()` untuk JSON decoding
- Validation dengan binding tags: `binding:"required"`
- Memanggil service layer
- Response formatting

### Response Formatting
- Standard response structure:
  ```json
  {
    "success": true,
    "data": {...}
  }
  ```
- Error responses:
  ```json
  {
    "success": false,
    "error": {"code": "...", "message": "..."}
  }
  ```
- HTTP status codes yang tepat

---

## 3.8 Integrasi dengan Packing Service

### HTTP Client Setup
- `http.Client` dengan timeout configuration
- Base URL dari environment variable

### PackingGateway Implementation
- Interface definition
- Request/Response handling
  - Marshaling `PackRequest` ke JSON
  - Unmarshaling `PackResponse` dari JSON
- Error handling: HTTP errors vs application errors

### Resilience Patterns
- Timeout: jangan biarkan request menggantung
- Retry dengan exponential backoff (opsional)
- Context cancellation untuk graceful shutdown

### Testing Gateway
- Mock HTTP server untuk testing
- Error scenario testing

---

## 3.9 Summary dan Further Reading

### Summary
- Recap keputusan desain yang dibuat
- Arsitektur final Backend API
- Bagaimana setiap layer bekerja sama
- Koneksi ke Bab 4 (Packing Service dengan Python)

### Further Reading
- Go Documentation: https://go.dev/doc/
- Effective Go: https://go.dev/doc/effective_go
- Gin Documentation: https://gin-gonic.com/docs/
- SQLC Documentation: https://docs.sqlc.dev/
- Goose Documentation: https://github.com/pressly/goose
- pgx Documentation: https://github.com/jackc/pgx

---

## Estimasi

- **Panjang**: ~20-25 halaman dengan code snippets
- **Code snippets**: 15-20 snippets
- **Diagrams**: 2-3 (ERD, layer architecture, sequence)

---

## Notes

### Pendekatan Kode untuk Buku

Buku ini akan membangun aplikasi **dari scratch** yang merupakan versi sederhana dari main app. Fitur yang **dihilangkan** atau **ditunda**:
- Multi-tenancy
- Admin panel
- **Authentication** (akan dibahas di bab akhir sebagai enhancement)
- Fitur kompleks lainnya yang tidak relevan untuk pembelajaran

**Strategi penulisan kode:**
- Logika dan implementasi dapat di-copy dari main app
- Sesuaikan dan sederhanakan sesuai kebutuhan buku
- Fokus pada core functionality: Container, Product, Plan, Packing

### Guidelines Penulisan

- Setiap konsep dijelaskan dengan "mengapa" sebelum "bagaimana"
- Referensi ke Bab 2 untuk konteks arsitektur
- Persiapan untuk Bab 4 (Python side)
- Section order: Database → Repository (SQLC) → Service → Handler

