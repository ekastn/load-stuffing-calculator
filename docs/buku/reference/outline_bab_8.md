# Bab 8: Fitur Lanjutan & Finalisasi

Bab ini adalah sentuhan akhir untuk menyempurnakan aplikasi kita. Setelah sistem dasar berjalan, kita akan melengkapinya dengan fitur-fitur esensial yang membuat aplikasi ini layak disebut *production-ready*: Autentikasi Pengguna, Visualisasi Data, dan pelaporan dokumen.

## Tujuan Pembelajaran
-   **Authentication**: Mengamankan aplikasi menggunakan mekanisme JWT (JSON Web Token).
-   **Dashboarding**: Menyajikan data statistik ringkas menggunakan agregasi SQL.
-   **Reporting**: Membuat fitur export PDF di sisi client (browser) untuk kebutuhan operasional.
-   **Wrap-up**: Merangkum arsitektur yang telah dibangun.

## Struktur Bab

### 8.1 Implementasi Autentikasi (JWT)
Kita akan membangun sistem login agar setiap pengguna memiliki ruang data privasinya sendiri (Multi-tenancy sederhana).
-   **Backend**: 
    -   Membuat tabel `users` dan migrasi database.
    -   Implementasi handler `Register` dan `Login`.
    -   Membuat Middleware Autentikasi untuk memvalidasi token JWT pada setiap request.
-   **Frontend**: 
    -   Integrasi form login dengan API.
    -   Manajemen session menggunakan Cookies/LocalStorage.
    -   Proteksi halaman (Redirect jika belum login).

### 8.2 Dashboard & Statistik
Halaman "Home" yang kosong akan kita ubah menjadi Dashboard yang informatif.
-   **Backend**: 
    -   Membuat API khusus `/stats` yang mengembalikan ringkasan data (Total Plan, Total Volume, dll).
    -   Penggunaan `COUNT()` dan `SUM()` dalam SQL query.
-   **Frontend**: 
    -   Membuat komponen Card Statistik.
    -   Menampilkan data ringkasan tersebut secara visual di halaman depan.

### 8.3 Fitur Export PDF "Surat Jalan"
Fitur praktis untuk mencetak hasil stuffing plan menjadi dokumen fisik.
-   **Pendekatan Client-Side**: Mengapa kita men-generate PDF di browser? (Kecepatan & Interaktivitas).
-   **Teknis**:
    -   Menggunakan library `jspdf`.
    -   Menyusun layout PDF: Header, Tabel Manifest Barang, dan Screenshot Visualisasi.
    -   Menambahkan tombol aksi di halaman detail plan.

### 8.4 Penutup
-   **Review Arsitektur**: Melihat kembali *big picture* dari gabungan Go (Clean Arch), Python (Algorithm), dan Next.js.
-   **Next Steps**: Saran untuk pengembangan mandiri (misal: deploy ke VPS, tambah fitur payment).
