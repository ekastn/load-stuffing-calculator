# Outline Bab 6: Integrasi Full Stack

## Tujuan Bab

Menghubungkan frontend Next.js dengan backend Go API untuk membuat aplikasi yang berfungsi penuh. Pembaca akan belajar mengintegrasikan semua komponen yang telah dibangun di bab sebelumnya.

**Referensi Outline Utama:**
> - **The Goal:** Menghubungkan frontend Next.js dengan backend Go API untuk membuat aplikasi yang berfungsi penuh.
> - **Technical Challenges:** Konfigurasi CORS, pengelolaan environment variables, error handling, dan state management untuk data fetching.
> - **Addressing Challenges:** Konfigurasi API client, implementasi CRUD operations, integrasi dengan Packing Service, error handling dan loading states.

---

## Struktur Bab

### Introduction

- Recap: apa yang sudah kita bangun
  - Bab 3: Backend API (Go) dengan CRUD untuk containers, items, plans
  - Bab 4: Packing Service (Python) untuk kalkulasi packing
  - Bab 5: Frontend Visualization (Next.js) dengan demo packing lokal
- Masalah saat ini: Demo page menggunakan client-side packing, bukan API
- Goal: Menghubungkan semua pieces menjadi aplikasi yang berfungsi
- Preview hasil akhir: Full working application

**Technical Requirements:**
- Backend API running di localhost:8080
- Packing Service running di localhost:5000
- PostgreSQL database dengan data

---

### 6.1 Konfigurasi CORS di Backend

- Mengapa CORS diperlukan: browser security untuk cross-origin requests
- Frontend di localhost:3000, Backend di localhost:8080
- Implementasi CORS middleware di Go:
  - Allowed origins
  - Allowed methods (GET, POST, PUT, DELETE, OPTIONS)
  - Allowed headers (Content-Type, Authorization)
- Preflight requests (OPTIONS)
- Environment-based configuration (dev vs production)

**Files:**
- `backend/internal/api/cors.go` (baru)
- `backend/internal/api/api.go` (update untuk add cors middleware)

---

### 6.2 Environment Variables di Frontend

- Mengapa environment variables: tidak hardcode URLs
- Next.js environment variables:
  - `NEXT_PUBLIC_API_URL` - untuk client-side
  - `.env.local` untuk development
  - `.env.production` untuk production
- Mengakses env vars di code: `process.env.NEXT_PUBLIC_API_URL`
- Type-safe env dengan validation

**Files:**
- `web/.env.local`
- `web/.env.example`
- `web/lib/config.ts` (baru)

---

### 6.3 API Client Setup

- Membuat reusable API client dengan fetch
- Base configuration:
  - Base URL dari environment
  - Common headers (Content-Type)
  - Error handling wrapper
- Type-safe API responses
- Request/Response interfaces yang match dengan backend DTOs

**Files:**
- `web/lib/api/client.ts` (baru)
- `web/lib/api/types.ts` (baru)

---

### 6.4 Container Operations

- Mengimplementasikan CRUD untuk containers
- API endpoints:
  - `GET /api/containers` - list all
  - `GET /api/containers/:id` - get by ID
  - `POST /api/containers` - create
  - `PUT /api/containers/:id` - update
  - `DELETE /api/containers/:id` - delete
- React hooks atau functions untuk data fetching
- Loading dan error states
- Form untuk create/edit container

**Files:**
- `web/lib/api/containers.ts` (baru)
- `web/components/container-form.tsx` (baru)
- `web/components/container-list.tsx` (baru)

---

### 6.5 Item (Product) Operations

- CRUD untuk items/products
- API endpoints:
  - `GET /api/products` - list all
  - `GET /api/products/:id` - get by ID
  - `POST /api/products` - create
  - `PUT /api/products/:id` - update
  - `DELETE /api/products/:id` - delete
- Items table dengan actions
- Form untuk add/edit item
- Color picker untuk item visualization

**Files:**
- `web/lib/api/products.ts` (baru)
- `web/components/product-form.tsx` (baru)
- `web/components/product-list.tsx` (baru)

---

### 6.6 Stuffing Plan Creation

- Workflow untuk membuat plan:
  1. Select container
  2. Add items dengan quantities
  3. Submit untuk kalkulasi
- API endpoint: `POST /api/plans`
- Request body structure matching backend expectations
- Response handling: plan ID untuk subsequent operations

**Files:**
- `web/lib/api/plans.ts` (baru)
- `web/components/plan-form.tsx` (baru)

---

### 6.7 Packing Calculation Integration

- Trigger packing calculation dari frontend
- API endpoint: `POST /api/plans/:id/calculate`
- Long-running operation handling:
  - Loading indicator selama kalkulasi
  - Polling atau webhook untuk hasil (optional)
- Error handling untuk packing failures:
  - Items tidak fit
  - Container terlalu kecil
- Menampilkan hasil: fitted, unfitted, utilization

**Files:**
- `web/lib/api/plans.ts` (update)
- `web/components/packing-result.tsx` (baru)

---

### 6.8 Visualization dengan Real Data

- Menghubungkan StuffingViewer dengan API data
- Endpoint: `GET /api/plans/:id` dengan placements
- Transform API response ke `StuffingPlanData` format
- Auto-load visualization setelah packing complete
- Refresh visualisasi saat data berubah

**Files:**
- `web/lib/api/plans.ts` (update)
- `web/lib/transforms.ts` (baru) - API response to StuffingPlanData
- `web/app/plans/[id]/page.tsx` (baru) - plan detail page

---

### 6.9 Loading dan Error States

- Consistent UX untuk async operations
- Loading states:
  - Skeleton loaders untuk lists
  - Spinner untuk forms
  - Progress indicator untuk packing
- Error states:
  - Toast notifications untuk errors
  - Retry buttons
  - Fallback UI
- shadcn/ui components: Toast, Skeleton

**Files:**
- `web/components/loading-states.tsx` (baru)
- `web/components/error-boundary.tsx` (baru)

---

### 6.10 Application Layout dan Navigation

- Multi-page application structure
- Navigation:
  - Dashboard / Home
  - Containers management
  - Products/Items management
  - Plans list dan detail
- Shared layout dengan sidebar atau navbar
- shadcn/ui: NavigationMenu, Breadcrumb

**Files:**
- `web/app/layout.tsx` (update)
- `web/components/navigation.tsx` (baru)
- `web/app/containers/page.tsx` (baru)
- `web/app/products/page.tsx` (baru)
- `web/app/plans/page.tsx` (baru)

---

### Summary

- CORS configuration memungkinkan frontend-backend communication
- Environment variables menjaga konfigurasi flexible
- API client abstraction memudahkan data fetching
- CRUD operations untuk containers, products, dan plans
- Integration dengan Packing Service untuk real calculations
- StuffingViewer menampilkan hasil dari real data
- Proper loading dan error states untuk good UX

### Further Reading

- Next.js Data Fetching: https://nextjs.org/docs/app/building-your-application/data-fetching
- React Query: https://tanstack.com/query/latest
- SWR: https://swr.vercel.app/
- CORS Explained: https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS

---

## Estimasi

- **Panjang**: ~4000-5000 kata
- **Code snippets**: 12-15 files baru/update
- **Implementasi reference**: Update `docs/buku/source/bab_5/web/` atau buat `docs/buku/source/bab_6/`

---

## Source Code Files

```
docs/buku/source/bab_6/
├── backend-updates/
│   └── internal/api/cors.go           # CORS middleware
└── web/
    ├── .env.local                     # Environment variables
    ├── .env.example                   # Example env file
    ├── lib/
    │   ├── config.ts                  # Environment config
    │   ├── transforms.ts              # API to visualization transforms
    │   └── api/
    │       ├── client.ts              # Base API client
    │       ├── types.ts               # API types
    │       ├── containers.ts          # Container CRUD
    │       ├── products.ts            # Product CRUD
    │       └── plans.ts               # Plan operations
    ├── components/
    │   ├── navigation.tsx             # App navigation
    │   ├── container-form.tsx         # Container create/edit
    │   ├── container-list.tsx         # Container table
    │   ├── product-form.tsx           # Product create/edit
    │   ├── product-list.tsx           # Product table
    │   ├── plan-form.tsx              # Plan creation wizard
    │   ├── packing-result.tsx         # Packing stats display
    │   ├── loading-states.tsx         # Skeleton/spinner
    │   └── error-boundary.tsx         # Error handling
    └── app/
        ├── layout.tsx                 # Updated with navigation
        ├── page.tsx                   # Dashboard
        ├── containers/
        │   └── page.tsx               # Containers management
        ├── products/
        │   └── page.tsx               # Products management
        └── plans/
            ├── page.tsx               # Plans list
            └── [id]/
                └── page.tsx           # Plan detail with visualization
```

---

## Notes

### Koneksi dengan Bab Lain

- **Bab 3**: Backend API endpoints yang diakses
- **Bab 4**: Packing Service yang dipanggil via backend
- **Bab 5**: StuffingViewer component untuk visualization
- **Bab 7**: Deployment semua services bersama

### Pendekatan Penulisan

- Start dengan CORS karena itu blocker pertama
- Build API client secara incremental
- Show real API calls dengan curl untuk verification
- Focus pada practical integration, bukan React patterns
- Gunakan existing components dari Bab 5 sebanyak mungkin

### Keputusan Arsitektur

- **Fetch vs Axios**: Gunakan native fetch untuk simplicity
- **React Query vs SWR vs useEffect**: Mulai dengan useEffect, mention alternatives
- **Form handling**: React Hook Form atau native forms
- **State management**: Local state untuk simplicity, mention Zustand/Jotai for scale

### Fitur yang Diimplementasikan

- CORS middleware di backend
- Environment-based configuration
- Type-safe API client
- CRUD for all entities
- Packing calculation trigger
- Real data visualization
- Multi-page navigation
- Loading dan error states

### Fitur yang Tidak Dibahas

- Authentication dan authorization
- Caching strategies (React Query)
- Optimistic updates
- Real-time updates (WebSocket)
- File upload untuk batch import
- Export to PDF/Excel

### Guidelines Penulisan

- Diagram-First: Sequence diagram untuk API flow
- Modular Code Snippets: Satu file per section
- Architectural Justification: Mengapa fetch, mengapa structure ini
- In-code Comments: Explain error handling decisions
