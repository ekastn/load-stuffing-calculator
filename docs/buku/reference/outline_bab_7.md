# Outline Bab 7: Deployment dan Produksi

## Tujuan Bab

Mempersiapkan aplikasi untuk lingkungan produksi menggunakan teknologi kontainerisasi. Pembaca akan belajar cara membungkus aplikasi ke dalam Docker, mengelola orkestrasi dengan Docker Compose, dan memahami konsep CI/CD modern.

**Referensi Outline Utama:**
> - **The Goal:** Memindahkan aplikasi dari "it works on my machine" ke containerized environment yang reliable.
> - **Technical Challenges:** Multi-stage builds, inter-container networking, environment configuration management, process isolation.
> - **Addressing Challenges:** Dockerfiles optimization, Docker Compose networking, Healthchecks, dan CI/CD automation.

---

## Struktur Bab

### Introduction
- Realitas software development: "It works on my machine" vs Production.
- Tantangan deployment manual: Dependency hell, environment drift.
- Solusi: Immutable Infrastructure dengan Docker.
- Preview hasil akhir: Aplikasi berjalan penuh hanya dengan satu perintah `docker-compose up`.

---

### 7.1 Persiapan Deployment & Konfigurasi Eksternal

- **The Twelve-Factor App**: Fokus pada prinsip **Config**.
- Audit Konfigurasi:
  - Backend: Pastikan DB DSN, Port, dan URL Service dibaca dari ENV.
  - Frontend: Memahami beda Build-time vars (`NEXT_PUBLIC_`) vs Runtime vars.
- Production Build Flags:
  - Menghapus debug logs.
  - Mengaktifkan production mode di Gin (Go) dan Next.js.

**Files:**
- `backend/internal/config/config.go` (review)
- `web/next.config.js` (optimization settings)

---

### 7.2 Pengenalan Docker & Containerization

- **Konsep Dasar**:
  - Virtual Machines vs Containers (Efisiensi OS Layer).
  - Image (Blueprint) vs Container (Runtime Instance).
- **Arsitektur Docker**: Daemon, Client, Registry.
- **Microservices Fit**: Mengapa Docker sempurna untuk aplikasi polyglot kita (Go + Python + Node.js).

---

### 7.3 Dockerizing Aplikasi (Implementation)

Langkah detail membuat `Dockerfile` untuk setiap layanan.

#### A. Backend API (Go)
- **Multi-stage Build**:
  - Stage 1: `golang:alpine` untuk kompilasi (besar).
  - Stage 2: `gcr.io/distroless/static` atau `alpine` untuk runtime (kecil & aman).
- Teknik optimasi binary size.

#### B. Packing Service (Python)
- Manajemen dependensi dengan `requirements.txt`.
- Penggunaan `python:slim` variant.
- Non-root user security best practice.

#### C. Frontend (Next.js)
- Standalone Output tracing (fitur Next.js untuk mengurangi node_modules).
- Multi-stage: Deps -> Builder -> Runner.

**Files:**
- `Dockerfile.backend`
- `Dockerfile.packing`
- `Dockerfile.frontend`
- `.dockerignore`

---

### 7.4 Orkestrasi dengan Docker Compose

Menggabungkan semua layanan agar berjalan sebagai satu kesatuan.

- Struktur `docker-compose.yml`.
- **Services Definition**: backend, frontend, packing-service, db.
- **Networking**: Bagaimana backend memanggil `http://packing-service:5000` (DNS internal).
- **Volume Management**: Persisting data PostgreSQL (`pgdata`) agar tidak hilang saat container restart.
- **Environment Injection**: Menggunakan file `.env`.

**Files:**
- `docker-compose.yml`

---

### 7.5 Deployment Best Practices

Memastikan sistem tangguh di lingkungan produksi.

- **Healthchecks**:
  - Mengapa container "Running" != "Ready".
  - Implementasi endpoint `/health` di Go dan Python.
  - Konfigurasi `healthcheck` di Docker Compose (depends_on condition).
- **Logging Strategy**:
  - The Golden Rule: Log to `stdout`/`stderr`.
  - Jangan menulis log ke file lokal di dalam container.
- **Resource Limits**:
  - Membatasi Memory dan CPU agar satu container liar tidak mematikan host.

**Files:**
- `backend/cmd/api/main.go` (add healthcheck route)
- `packing-service/app.py` (add healthcheck route)

---

### 7.6 Strategi Lanjutan (CI/CD & Reliability)

Pengantar konsep tingkat enterprise (High Level Concepts).

- **CI/CD Pipeline**:
  - Apa itu Continuous Integration & Deployment.
  - Contoh alur GitHub Actions: Code Push -> Run Test -> Build Docker Image -> Push to Registry.
- **Deployment Strategies**:
  - **Rolling Update**: Zero-downtime deployment (standar K8s/Docker Swarm).
  - **Canary Deployment**: Merilis versi baru ke 10% user untuk validasi risiko.
- **Automated Rollback**: Rencana darurat jika deployment gagal.

**Files:**
- `.github/workflows/ci.yml` (contoh konfigurasi)

---

## Estimasi

- **Panjang**: ~4000-4500 kata.
- **Code snippets**: 6-8 files utama (Dockerfiles, Compose, Healthcheck routes).

## Source Code Files

```
load-stuffing-calculator/
├── build/
│   ├── backend.Dockerfile
│   ├── frontend.Dockerfile
│   └── packing.Dockerfile
├── docker-compose.yml
├── .dockerignore
└── .github/
    └── workflows/
        └── build-check.yml
```

---

## Notes

### Koneksi dengan Bab Lain
- **Bab 3, 4, 5**: Source code aplikasi yang akan di-package berasal dari bab-bab ini.
- **Bab 6**: Konfigurasi ENV yang disiapkan di Bab 6 akan digunakan secara intensif di sini.

### Pendekatan Penulisan
- Fokus pada **"Reproducibility"**: Siapapun yang clone repo harus bisa menjalankan aplikasi dengan satu perintah.
- Tekankan **Security**: Jangan jalankan container sebagai root jika tidak perlu.
- Tekankan **Size Optimization**: Image kecil = deployment cepat.
