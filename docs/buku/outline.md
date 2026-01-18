### **1. Title (Judul)**

**"Load & Stuffing Calculator: Implementasi Arsitektur dan Visualisasi 3D untuk Sistem Optimasi Kargo"**

---

### **2. Outline Terperinci (Daftar Isi)**

#### **Bab 1: Digitalisasi Logistik dan Tantangan Pemuatan**

* **The Goal:** Memahami landasan teoritis optimasi pemuatan dan batasan sistem manual.
* **Technical Challenges:** Kompleksitas *NP-hard* pada *3D Bin Packing Problem* (3D-BPP) dan inefisiensi *manual planning*.
* **Addressing Challenges:** Pengenalan pendekatan heuristik sebagai solusi praktis untuk komputasi cepat.
* **Technical Requirements:** *Browser* (Chrome/Firefox).

#### **Bab 2: Perancangan Arsitektur Layanan Mikro**

* **The Goal:** Merancang struktur sistem yang terdekopling untuk memisahkan manajemen data dan mesin hitung.
* **Technical Challenges:** *Scalability* saat menangani kalkulasi kargo besar dan pemilihan *communication protocol* antar layanan.
* **Addressing Challenges:** Penerapan arsitektur *Microservices* dengan pembagian tugas antara *Backend API* (Go) dan *Packing Service* (Python) melalui REST HTTP.
* **Technical Requirements:** *Diagramming tools* (Draw.io atau Lucidchart).

#### **Bab 3: Pengembangan Backend API dengan Bahasa Go**

* **The Goal:** Membangun *gateway* utama sistem yang mengelola persistensi data dan keamanan.
* **Technical Challenges:** Pengelolaan migrasi *database* yang konsisten dan pembuatan *data access layer* yang *type-safe*.
* **Addressing Challenges:** Implementasi *repository pattern*, penggunaan `sqlc` untuk eliminasi *boilerplate* SQL, dan `goose` untuk *versioned migration*.
* **Technical Requirements:** Go 1.21, PostgreSQL 15, `sqlc`, `goose`, VS Code.

#### **Bab 4: Membangun Mesin Kalkulasi (Packing Service) dengan Python**

* **The Goal:** Mengimplementasikan algoritma optimasi pemuatan dengan batasan fisik dunia nyata.
* **Technical Challenges:** Penanganan kargo heterogen dan validasi stabilitas tumpukan barang.
* **Addressing Challenges:** Enkapsulasi algoritma heuristik dengan strategi *Bigger First*, fitur *Fix Point* untuk gravitasi, dan *Check Stable* (75% penopang).
* **Technical Requirements:** Python 3.11, `pip`, *3D Bin Packing library*.

#### **Bab 5: Visualisasi 3D Interaktif dengan Three.js**

* **The Goal:** Mentransformasi hasil kalkulasi menjadi panduan operasional visual yang interaktif.
* **Technical Challenges:** *Rendering* objek 3D secara *real-time* di *browser*, transformasi koordinat antara API dan Three.js, dan implementasi urutan pemuatan kronologis.
* **Addressing Challenges:** 
  - Penggunaan *WebGL* via Three.js dengan *Manager Pattern* untuk separation of concerns
  - Representasi *BoxGeometry* untuk container (wireframe) dan items (solid)
  - Double transformation koordinat (API → Three.js coordinate system)
  - Mekanisme *Step Playback* berbasis `step_number` dengan *Observer Pattern*
  - React integration dengan bidirectional state sync
* **Technical Requirements:** Node.js, Next.js, shadcn/ui, Three.js library, TypeScript.

#### **Bab 6: Integrasi Full Stack**

* **The Goal:** Menghubungkan frontend Next.js dengan backend Go API untuk membuat aplikasi yang berfungsi penuh.
* **Technical Challenges:** Konfigurasi CORS, pengelolaan environment variables, error handling, dan state management untuk data fetching.
* **Addressing Challenges:** 
  - Konfigurasi API client dengan environment variables
  - Implementasi CRUD operations (containers, items, plans) dari frontend
  - Integrasi dengan Packing Service melalui backend API
  - Error handling dan loading states
  - (Optional) React Query / SWR untuk data fetching dan caching
* **Technical Requirements:** Pemahaman REST API, fetch API, React hooks.

#### **Bab 7: Deployment dan Produksi**

* **The Goal:** Menyiapkan sistem untuk deployment ke lingkungan produksi.
* **Technical Challenges:** Orkestrasi multi-service, konfigurasi environment, dan build optimization.
* **Addressing Challenges:** 
  - Docker Compose untuk menjalankan semua services (Go API, Python Packing, PostgreSQL, Next.js)
  - Production build configurations untuk setiap service
  - Environment management (development vs production)
  - Health checks dan logging
* **Technical Requirements:** Docker, Docker Compose, dasar-dasar networking.

#### **Bab 8: Evaluasi dan Pengembangan Lanjutan**

* **The Goal:** Memvalidasi reliabilitas sistem dan mengeksplorasi pengembangan masa depan.
* **Technical Challenges:** Analisis performa sistem dan identifikasi area pengembangan.
* **Addressing Challenges:** 
  - Benchmarking dengan berbagai skenario beban kerja
  - Analisis volume utilization dan fit rate
  - Roadmap pengembangan: multiple containers, weight distribution, IoT integration
* **Technical Requirements:** Testing tools, analisis data dasar.

---

### **3. Scope (Batasan Bahasan)**

* **In-Scope:**
* Arsitektur *Microservices* murni (Go & Python).
* Algoritma *3D Bin Packing* dengan pendekatan heuristik.
* Manajemen *database* relasional dengan teknik *code generation*.
* Visualisasi 3D interaktif berbasis web (Three.js).
* Analisis data hasil pengujian performa sistem.


* **Out-of-Scope:**
* Manajemen infrastruktur *Cloud* (AWS/Azure) secara mendalam.
* *Container Orchestration* (Kubernetes).
* Keamanan jaringan tingkat lanjut (*Network Hardening*).
* Fitur-fitur komersial seperti sistem pembayaran atau manajemen vendor.



---

### **4. Hasil Akhir (Expected Outcome)**

Pembaca akan memiliki pengetahuan dan *source code* untuk membangun sebuah platform yang mampu:

1. **Optimasi Maksimal:** Mencapai *Fill Rate* 100% dan *Volume Utilization* hingga 55,26% untuk kargo yang sangat heterogen.
2. **Kecepatan Operasional:** Menyelesaikan perhitungan 300 barang dalam waktu kurang dari 40 detik.
3. **Akurasi Visual:** Menyajikan panduan pemuatan langkah-demi-langkah yang akurat secara geometris (bebas *overlap* dan stabil).
4. **Arsitektur Standar Industri:** Memahami cara membangun sistem multi-bahasa yang saling berkomunikasi secara efisien.

---

### **5. Detail Isi Materi (Deep Dive Per Bab)**

Sesuai gaya Shuiskov, setiap bab akan diisi dengan:

* **Diagram-First:** Membuka pembahasan komponen baru dengan *Architecture Diagram*, *Sequence Diagram*, atau *ER-Diagram*.
* **Modular Code Snippets:** Menampilkan potongan kode fungsional yang spesifik (misal: hanya fungsi transformasi koordinat) dengan judul `Listing`.
* **Architectural Justification:** Menjelaskan alasan pemilihan teknologi tanpa menggunakan *buzzwords*. Contoh: *"Penggunaan Go pada API Gateway dipilih karena performa concurrency yang superior, sementara Python digunakan pada Packing Service untuk memanfaatkan ekosistem pustaka matematika yang matang."*
* **In-code Comments:** Setiap baris kode yang krusial diberi komentar yang menjelaskan peran baris tersebut dalam arsitektur sistem.
* **Summary:** Menutup dengan status fungsional sistem (misal: *"Backend API kini siap menerima payload kargo"*).
