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
* **Technical Challenges:** *Rendering* objek 3D secara *real-time* di *browser* dan implementasi urutan pemuatan kronologis.
* **Addressing Challenges:** Penggunaan *WebGL* via Three.js, representasi *BoxGeometry*, dan mekanisme *Step Playback* berbasis `step_number`.
* **Technical Requirements:** Node.js, Three.js library.

#### **Bab 6: Tantangan Sinkronisasi dan Transformasi Geometris**

* **The Goal:** Menyelaraskan sistem koordinat yang berbeda antara mesin kalkulasi dan *rendering* visual.
* **Technical Challenges:** Ketidaksamaan konvensi sumbu (Y-up vs Z-up) antar *library* dan pencegahan anomali barang melayang (*floating items*).
* **Addressing Challenges:** Implementasi *double transformation* koordinat  dari *Packing Service* ke API, lalu ke format Three.js.
* **Technical Requirements:** Pemahaman dasar aljabar linier dan sistem koordinat Kartesius.

#### **Bab 7: Evaluasi Performa dan Skalabilitas Sistem**

* **The Goal:** Memvalidasi reliabilitas sistem melalui berbagai skenario beban kerja.
* **Technical Challenges:** Analisis kompleksitas waktu  dan efektivitas penggunaan ruang kontainer.
* **Addressing Challenges:** Benchmarking skenario S1-S5 (50 hingga 300 unit), analisis korelasi volume vs berat, dan pembuktian efisiensi strategi *Bigger First*.
* **Technical Requirements:** *Testing tools*, Microsoft Excel atau Python Matplotlib untuk grafik.

#### **Bab 8: Integrasi IoT dan Masa Depan Pemuatan**

* **The Goal:** Mengeksplorasi konektivitas sistem dengan perangkat keras untuk validasi *real-time*.
* **Technical Challenges:** Sinkronisasi data sensor berat dengan *Load Plan* digital.
* **Addressing Challenges:** Arsitektur integrasi sensor IoT dan rencana pengembangan optimasi dinamis.
* **Technical Requirements:** Dasar-dasar protokol MQTT atau HTTP *request* dari perangkat IoT.

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
