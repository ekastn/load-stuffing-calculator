## 🏗️ Master Checklist Penulisan Bab (Gaya Arsitektural Shuiskov)

Aku pengen ngikutin cara penulisan yang ada di buku [Alexander Shuiskov](./reference/1012_Microservices-with-Go.pdf).

### 1. Struktur Pembuka Bab (Focus & Challenges)

Setiap bab harus dibuka dengan narasi yang efisien tanpa basa-basi pemasaran:

- **The Goal (Tujuan):** Mulai langsung dengan apa yang akan dibangun secara teknis.
    - _Pola:_ "Pada bab ini, kita akan [aktivitas teknis]. Tujuan utamanya adalah untuk [hasil teknis yang ingin dicapai]."
- **Technical Challenges (Daftar Tantangan):** Sebutkan tantangan spesifik yang akan dihadapi pengembang dalam topik tersebut menggunakan poin-poin.
    - _Contoh:_ "Terdapat beberapa tantangan teknis yang akan kita pecahkan, meliputi:
        - Bagaimana mengelola migrasi skema basis data secara otomatis.
        - Bagaimana menghasilkan kode Go yang _type-safe_ dari _query_ SQL mentah."
- **Addressing Challenges (Penjelasan Solusi):** Paragraf yang merinci bagaimana isi bab ini akan menjawab tantangan-tantangan tadi secara berurutan.
- **Jembatan Progres (Opsional):** Jika bab ini sangat bergantung pada hasil bab sebelumnya, gunakan satu atau dua paragraf terpisah untuk merangkum posisi sistem saat ini.

### 2. Roadmap Topik

- **Format:** "**Dalam bab ini, kita akan membahas topik-topik berikut:**"
- Ikuti dengan poin-poin singkat yang mencerminkan sub-bab utama.

### 3. Technical Requirements 

Bagian ini tidak hanya mencantumkan bahasa pemrograman, tetapi semua ekosistem pendukung:

- **Languages & Frameworks:** Versi spesifik (misal: Go 1.21, Python 3.11).
- **CLI Tools:** Alat bantu baris perintah (misal: `sqlc` untuk _code generation_, `goose` untuk _database migration_).
- **Databases & Infrastructure:** Perangkat lunak pendukung (misal: PostgreSQL 15, Docker).
- **IDE/Software:** Rekomendasi alat kerja (misal: VS Code dengan ekstensi tertentu).
- **Source Code:** Tautan ke GitHub yang merujuk pada direktori spesifik bab tersebut.

### 4. Detail Isi Materi (The Implementation Deep Dive)

Isi materi harus modular dan tidak bertele-tele. Aku harus mengikuti pola **"Diagram-First, Code-Second"**:

- **Diagram Alir & Arsitektur:** Setiap komponen baru wajib diawali dengan visualisasi. Jika membahas API, tampilkan diagram sekuens. Jika membahas database, tampilkan ERD atau diagram aliran data.
    - Pakai Graphviz atau mermaid untuk membuat diagram nya.
- **Logical Blocks (Modular Code):**
    - Jangan menampilkan file utuh.
    - Gunakan potongan kode yang fokus pada satu fungsi atau satu struktur data.
    - Setiap potongan kode harus memiliki judul blok (misal: `Listing 2.1: Implementasi Struct Product`).
- **Justification (The "Why"):** Shuiskov sangat kuat dalam justifikasi arsitektural. Aku harus menjelaskan:
    - "Kenapa menggunakan `sqlc`?" (Karena mengurangi _boilerplate_ dan _error-prone_ pada penulisan manual).
    - "Kenapa migrasi menggunakan `goose`?" (Karena mendukung versi skema yang terkontrol).
- **Interoperability & Contracts:** Karena ini _microservices_, fokuskan penjelasan pada kontrak data (JSON Schema). Tunjukkan bagaimana struktur data di Go akan diterjemahkan ke Python.
- **In-code Comments:** Pastikan komentar dalam kode menjelaskan logika arsitekturalnya, bukan sekadar menjelaskan cara kerja bahasa pemrograman.

### 5. Summary

- Rangkuman fungsional mengenai komponen yang baru saja selesai dibangun.
- Penegasan bahwa sistem sekarang sudah selangkah lebih dekat ke versi final.

## 🔤 Aturan Bahasa & Terminologi murni

- **Istilah Teknis (Murni Bahasa Inggris):** Jangan dipaksakan ke Bahasa Indonesia. Gunakan: _Microservices, Backend, API Gateway, 3D Bin Packing, Scaffolding, Type-safe, Boilerplate, Migration, Deployment, Concurrency, Middleware, Handlers, Service Layer, Repository Pattern, Payload, Endpoint._
- **No Buzzwords:** Hindari kata: _Canggih, luar biasa, inovatif, revolusioner._
- **Gaya Penulisan:** Objektif, berbasis data/hasil pengujian, dan mengarahkan pembaca seolah sedang membangun _blueprinting_ industri.

## 🚀 Simulasi Penerapan: Bab 3 (Database Management) (contoh)

**Bab 3: Persistensi Data dan Manajemen Skema**

Pada bab ini, kita akan mengimplementasikan lapisan persistensi data untuk _Backend API_ menggunakan PostgreSQL. Tujuan bab ini adalah untuk memastikan data produk dan kontainer tersimpan secara terstruktur dan dapat diakses dengan performa tinggi. Meskipun Go menyediakan paket `database/sql`, terdapat beberapa tantangan teknis yang akan kita hadapi:

- Bagaimana mengelola perubahan skema basis data secara konsisten tanpa kehilangan data.
- Bagaimana melakukan pemetaan (_mapping_) hasil _query_ SQL ke dalam _struct_ Go secara otomatis dan efisien.

Dalam bab ini, kita akan menjawab tantangan tersebut dengan mengintegrasikan alat migrasi skema dan _compiler_ SQL-ke-Go. Pertama, kita akan mengatur skema database menggunakan migrasi terversi. Kemudian, kita akan menggunakan alat pembuat kode untuk menghasilkan _repository layer_ yang _type-safe_.

**Dalam bab ini, kita akan membahas topik-topik berikut:**
- _Database Schema Design_
- _Versioned Migration with Goose_
- _Generating Type-safe Code with SQLC_

**Technical Requirements** Untuk menyelesaikan bab ini, Anda akan memerlukan alat dan perangkat lunak berikut:
- _Go version 1.21 or above_
- _PostgreSQL version 15_
- _SQLC CLI_ (sqlc.dev)
- _Goose CLI_ ([github.com/pressly/goose](https://github.com/pressly/goose))
- _VS Code with Go extension_ Anda dapat menemukan contoh kode dan file migrasi untuk bab ini di GitHub: [https://www.merriam-webster.com/dictionary/link](https://www.merriam-webster.com/dictionary/link)


