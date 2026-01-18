# Outline Bab 4: Packing Service dengan Python

## Tujuan Bab

Membangun microservice Python yang menerima data geometris dan mengembalikan koordinat penempatan menggunakan algoritma 3D Bin Packing, lalu mengintegrasikannya dengan backend Go.

---

## Struktur Bab

### 4.1 Setup Project Python

- Struktur direktori `cmd/packing/`
- Virtual environment dan `requirements.txt`
- Pengenalan Flask sebagai lightweight web framework
- Mengapa Python untuk service ini (ekspresi sintaksis, py3dbp library)

### 4.2 Flask Application Entry Point

- `app.py`: Application factory pattern `create_app()`
- `/health` endpoint untuk monitoring dan load balancer
- `/pack` endpoint dengan error handling
- `_json_error()` helper untuk standardized error response
- Environment variables: PACKING_HOST, PACKING_PORT, PACKING_DEBUG
- Running development server

### 4.3 Schema dan Validasi Request

- TypedDict untuk type hints tanpa runtime overhead
- Request types: ContainerIn, ItemIn, OptionsIn, RequestIn
- Response types: PlacementOut, UnfittedOut, StatsOut, PackSuccessResponse
- `parse_request()` function dengan validasi ketat
- Error handling dengan ValueError (yang ditangkap di app.py)
- Validasi: units (mm/cm/m), positive numbers, non-empty items

### 4.4 Unit Conversion

- `units.py`: Konversi mm/cm/m ke internal representation
- Mengapa integer cm untuk kalkulasi:
  - Precision: menghindari floating point errors
  - Performance: integer arithmetic lebih cepat
- `to_cm()`, `from_cm()`, `cm_int()` functions

### 4.5 Integrasi py3dbp Library

- py3dbp sebagai vendored dependency (bukan pip install)
  - Alasan: isolation, reproducibility, versi terkontrol
- `_load_py3dbp()`: dynamic import dari local directory
- Bin initialization dengan WHD (Width, Height, Depth)
- Item initialization dan expandQuantity
- Packer options yang didukung:
  - `bigger_first`: Sort by volume (default: True)
  - `fix_point`: Corner placement strategy (default: True)
  - `check_stable`: Gravity simulation (default: True)
  - `support_surface_ratio`: Minimum support area (default: 0.75)

### 4.6 Packing Logic

- `packing.py`: Orchestrating the entire flow
- Flow: parse → convert → initialize → pack → transform → respond
- Coordinate system mapping:
  - py3dbp WHD → API L,W,H
  - Position (x,y,z) transformation
- `rotation.py`: Mapping rotation codes (0-5) untuk visualisasi 3D
- Aggregating unfitted items by item_id
- Timing statistics (pack_time_ms, total_time_ms)

### 4.7 Integrasi dengan Go Backend

- Update `api.go`: MockGateway → HTTPPackingGateway
- Environment variables yang diperlukan
- Menjalankan kedua service bersamaan:
  - Terminal 1: Python Packing Service (port 5000)
  - Terminal 2: Go API Server (port 8080)
- Flow end-to-end: Frontend → Go API → Packing Service → Response

### 4.8 Testing dan Verifikasi

- Test `/health` endpoint
- Test `/pack` endpoint dengan curl:
  - Simple request dengan 1 item
  - Complex request dengan multiple items
  - Request dengan items yang tidak muat
- Test end-to-end: Go `/calculate` → Python `/pack`
- Verifikasi response format dan placement coordinates

### Summary

- Dua service berjalan: Go (API) dan Python (Packing)
- Separation of concerns: CRUD vs Computation
- Type safety dengan TypedDict
- Vendored dependency untuk reproducibility

### Further Reading

- py3dbp GitHub: https://github.com/enzoruiz/3dbinpacking
- Flask Application Factories: https://flask.palletsprojects.com/patterns/appfactories/
- TypedDict PEP 589: https://peps.python.org/pep-0589/
- 3D Bin Packing Problem: Bab 1 reference papers

---

## Estimasi

- **Panjang**: ~4000-5000 kata
- **Code snippets**: 8 files utama
- **Implementasi reference**: `docs/buku/source/bab_4/cmd/packing/`
