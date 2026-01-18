---
aliases:
---
# Load & Stuffing Calculator

**Implementasi Arsitektur dan Visualisasi 3D untuk Sistem Optimasi Kargo**

---

## Daftar Isi

1. [Bab 1: Digitalisasi Logistik dan Tantangan Pemuatan](#bab-1-digitalisasi-logistik-dan-tantangan-pemuatan)
2. [Bab 2: Perancangan Arsitektur Layanan Mikro](#bab-2-perancangan-arsitektur-layanan-mikro)
3. [Bab 3: Pengembangan Backend API dengan Bahasa Go](#bab-3-pengembangan-backend-api-dengan-bahasa-go)
4. [Bab 4: Membangun Mesin Kalkulasi (Packing Service) dengan Python](#bab-4-membangun-mesin-kalkulasi-packing-service-dengan-python)
5. [Bab 5: Visualisasi 3D Interaktif dengan Three.js](#bab-5-visualisasi-3d-interaktif-dengan-threejs)
6. [Bab 6: Tantangan Sinkronisasi dan Transformasi Geometris](#bab-6-tantangan-sinkronisasi-dan-transformasi-geometris)
7. [Bab 7: Evaluasi Performa dan Skalabilitas Sistem](#bab-7-evaluasi-performa-dan-skalabilitas-sistem)
8. [Bab 8: Integrasi IoT dan Masa Depan Pemuatan](#bab-8-integrasi-iot-dan-masa-depan-pemuatan)

---

<!-- BAB 1 STARTS HERE -->

## Bab 1: Digitalisasi Logistik dan Tantangan Pemuatan

Pada bab ini, kita akan memahami landasan teoritis dari permasalahan optimasi pemuatan kontainer dan mengapa sistem manual tidak lagi memadai untuk industri logistik modern. Tujuan utamanya adalah untuk memberikan fondasi pemahaman sebelum kita membangun sistem *Load & Stuffing Calculator*.

Terdapat beberapa tantangan teknis yang akan kita bahas, meliputi:

- Memahami kompleksitas matematis *NP-hard* pada *3D Bin Packing Problem* (3D-BPP).
- Mengidentifikasi inefisiensi *manual planning* yang menyebabkan *void space* dan kerugian operasional.
- Menentukan pendekatan algoritmik yang tepat untuk menghasilkan solusi dalam waktu komputasi yang wajar.

Dalam bab ini, kita akan menjawab tantangan tersebut dengan mengeksplorasi teori dasar *Container Loading Problem*, memahami mengapa pendekatan heuristik menjadi solusi praktis, dan mendefinisikan visi sistem yang akan kita bangun.

**Dalam bab ini, kita akan membahas topik-topik berikut:**

- Tantangan Logistik dalam Era Perdagangan Global
- *Container Loading Problem* (CLP) dan Kompleksitas NP-hard
- Pendekatan Algoritmik: Dari Optimal ke Heuristik
- Visi Sistem *Load & Stuffing Calculator*



### 1.1 Tantangan Logistik dalam Era Perdagangan Global

Pada tahun 1956, seorang pengusaha truk bernama Malcolm McLean melakukan eksperimen sederhana yang mengubah wajah perdagangan dunia: ia memuat 58 trailer truk ke dalam kapal tanker bekas di pelabuhan Newark, New Jersey. Eksperimen ini melahirkan kontainerisasi modern—sebuah inovasi yang mengurangi biaya pengiriman barang hingga 90% dan memungkinkan globalisasi ekonomi seperti yang kita kenal hari ini.

Namun, enam dekade kemudian, sebuah ironi masih terjadi di hampir setiap pelabuhan dunia. Meskipun proses pengangkutan kontainer telah sepenuhnya terotomasi dengan crane raksasa dan sistem pelacakan digital, proses *stuffing*—memuat barang ke dalam kontainer—masih dilakukan dengan cara yang nyaris tidak berubah sejak era McLean. Seorang operator berdiri di depan kontainer kosong, memegang clipboard atau tablet, lalu mengarahkan forklift berdasarkan intuisi dan pengalaman.

Pendekatan manual ini bukanlah masalah ketika sebuah kontainer berisi barang homogen—ribuan unit produk identik yang dapat ditumpuk dengan pola sederhana. Namun, realitas perdagangan modern jauh lebih kompleks. Sebuah eksportir furnitur, misalnya, mungkin harus memuat meja konferensi berukuran 3 meter berdampingan dengan kursi kantor, rak buku, dan aksesori kecil—semuanya dengan dimensi berbeda, batasan orientasi berbeda, dan kebutuhan penopang berbeda.

Di sinilah keterbatasan fundamental metode manual mulai terlihat. Otak manusia, sehebat apapun, memiliki batas dalam memproses kombinasi spasial tiga dimensi. Ketika jumlah barang meningkat dari puluhan menjadi ratusan, jumlah kemungkinan konfigurasi penempatan meledak secara eksponensial—dan tidak ada operator yang dapat mengevaluasi semua opsi dalam waktu yang wajar.

Dampak dari keterbatasan ini terukur dan multidimensi. Pertama, ada masalah **waktu**—perencanaan muatan dengan puluhan SKU berbeda dapat memakan waktu berjam-jam, menyebabkan delay pengiriman dan biaya overtime pekerja. Kedua, masalah **akurasi**—utilisasi volume rata-rata hanya mencapai 50-70%, sangat bergantung pada skill individu operator, yang langsung meningkatkan biaya pengiriman per unit produk.

Ketiga, ada tantangan **konsistensi**. Dua operator berbeda yang merencanakan muatan yang sama mungkin menghasilkan konfigurasi berbeda dengan utilisasi berbeda. Ketidakpastian ini menyulitkan perusahaan memprediksi kapasitas dan biaya secara akurat. Keempat, metode manual tidak **scalable**—ketika volume pengiriman meningkat, bottleneck terjadi karena jumlah planner berpengalaman terbatas.

Terakhir, ada risiko institusional yang sering diabaikan. Dokumentasi perencanaan manual cenderung tidak lengkap, menyulitkan audit dan klaim asuransi ketika terjadi kerusakan barang. Lebih kritis lagi, pengetahuan dan intuisi yang dibangun selama bertahun-tahun tersimpan di kepala individu—ketika karyawan berpengalaman resign atau pensiun, perusahaan kehilangan aset intelektual yang tidak ternilai.

Tantangan ini semakin kritis dalam konteks e-commerce dan *just-in-time manufacturing*, di mana perusahaan harus menangani ratusan hingga ribuan SKU (*Stock Keeping Unit*—kode unik untuk setiap varian produk) dengan dimensi berbeda-beda. Setiap pengiriman mungkin melibatkan kombinasi produk yang berbeda, dan perencanaan harus dilakukan dalam hitungan jam—bukan hari. Kontainer yang tersedia pun bervariasi: 20 feet, 40 feet standar, 40 feet High Cube—masing-masing dengan dimensi internal dan batasan berat maksimum yang berbeda.

Industri logistik membutuhkan solusi yang dapat memproses kompleksitas ini secara sistematis, konsisten, dan cepat. Solusi tersebut harus mampu mengevaluasi ribuan kemungkinan konfigurasi dalam hitungan detik, bukan jam. Solusi tersebut harus menghasilkan output yang dapat dipahami oleh operator lapangan—bukan hanya angka abstrak, tetapi panduan visual langkah demi langkah. Dan solusi tersebut harus robust terhadap variasi input: dari selusin barang homogen hingga ratusan SKU heterogen.

---

### 1.2 Container Loading Problem dan Kompleksitas NP-hard

Masalah yang tampak sederhana—menempatkan kotak-kotak ke dalam wadah lebih besar—ternyata telah menarik perhatian matematikawan dan ilmuwan komputer selama lebih dari setengah abad. Dalam literatur akademis, permasalahan ini dikenal sebagai *Container Loading Problem* (CLP), dan ia memiliki sejarah panjang yang terkait dengan perkembangan teori kompleksitas komputasi.

Definisi formal CLP dapat dirumuskan sebagai berikut:

> **Diberikan sebuah wadah (kontainer) dengan dimensi tertentu dan sekumpulan objek berbentuk kotak dengan dimensi masing-masing, tempatkan objek-objek tersebut ke dalam wadah sedemikian rupa sehingga total volume yang dimuat maksimal, tanpa terjadi tumpang tindih (overlap) antar objek, dan dengan memenuhi batasan fisik yang berlaku.**

Rumusan ini terdengar straightforward. Namun, di balik kesederhanaan definisinya tersembunyi kompleksitas matematis yang luar biasa.

Pada tahun 1972, Richard Karp mempublikasikan makalah seminal yang mengidentifikasi 21 masalah komputasi yang ia buktikan sebagai "NP-complete"—sebuah kelas masalah yang, hingga hari ini, tidak ada algoritma efisien yang diketahui dapat menyelesaikannya. *Bin Packing Problem*, versi satu dimensi dari CLP, termasuk dalam daftar Karp. Versi tiga dimensi—yang kita hadapi dalam pemuatan kontainer—bahkan lebih kompleks.

Apa artinya sebuah masalah bersifat **NP-hard** bagi praktisi? Mari kita telusuri implikasi praktisnya.

Pendekatan paling intuitif untuk menyelesaikan CLP adalah *brute force*: coba semua kemungkinan konfigurasi penempatan, evaluasi masing-masing, lalu pilih yang terbaik. Pendekatan ini dijamin menemukan solusi optimal. Masalahnya terletak pada jumlah kemungkinan yang harus dievaluasi.

Untuk setiap barang, kita perlu memutuskan: (1) di mana menempatkannya dalam ruang tiga dimensi, dan (2) dengan orientasi apa. Jika kita menyederhanakan dengan mengasumsikan kontainer dapat dibagi menjadi grid unit kecil, setiap barang memiliki ribuan kemungkinan posisi. Kalikan dengan 6 kemungkinan rotasi. Kalikan lagi dengan jumlah barang. Hasilnya adalah ledakan kombinatorial yang melampaui kapasitas komputasi manapun:

```
Jumlah Barang    Perkiraan Kombinasi       Waktu Brute Force
─────────────────────────────────────────────────────────────
      5          ~10^4                     < 1 detik
     10          ~10^10                    ~3 jam
     20          ~10^18                    ~30.000 tahun
     50          ~10^50                    > usia alam semesta
    100          ~10^150                   Angka yang tak bermakna
```

Perhatikan bagaimana penambahan linier pada jumlah barang menghasilkan pertumbuhan eksponensial pada kompleksitas. Inilah ciri khas masalah NP-hard: tidak ada "trik pintar" yang diketahui dapat menghindari ledakan kombinatorial ini. Jika seseorang menemukan algoritma yang dapat menyelesaikan CLP dalam waktu polynomial untuk semua kasus, ia akan membuktikan P = NP—salah satu masalah terbuka paling terkenal dalam matematika, dengan hadiah $1 juta dari Clay Mathematics Institute.

Dalam praktik industri, kompleksitas teoritis ini diterjemahkan menjadi batasan nyata. Sebuah perusahaan logistik tidak dapat menunggu bertahun-tahun untuk mendapatkan rencana pemuatan "optimal". Mereka membutuhkan jawaban dalam hitungan menit—bahkan detik—meskipun jawaban tersebut mungkin tidak sempurna.

Lebih jauh lagi, CLP dunia nyata jarang sesederhana definisi akademis. Batasan praktis (*constraints*) menambah dimensi kompleksitas baru:

**Orientasi dan arah "atas"**: Tidak semua barang boleh diputar secara bebas. Sebuah komputer harus tetap tegak. Sebuah mesin mungkin memiliki titik pengangkatan yang menentukan orientasinya. Constraint ini mengurangi ruang solusi, tetapi juga mengurangi fleksibilitas algoritma.

**Stabilitas tumpukan**: Barang yang ditumpuk harus memiliki penopang yang memadai di bawahnya. Algoritma tidak dapat seenaknya "menggantung" barang di udara. Biasanya, minimal 75% permukaan bawah barang harus didukung oleh barang lain atau lantai kontainer.

**Batasan berat**: Setiap kontainer memiliki *payload maksimum*—berat total yang boleh dimuat. Kontainer standar 40 feet, misalnya, memiliki batas sekitar 28 ton. Algoritma harus melacak akumulasi berat dan menolak penempatan yang melanggar batas ini.

**Distribusi berat**: Selain total berat, distribusi berat juga penting. Kontainer yang terlalu berat di satu sisi dapat menyebabkan masalah stabilitas saat pengangkutan. Beberapa algoritma advanced memperhitungkan *center of gravity*.

**Urutan pemuatan/pembongkaran**: Dalam skenario multi-destinasi, barang untuk pemberhentian pertama harus dapat diakses tanpa memindahkan barang lain. Ini menambahkan constraint temporal ke masalah spasial.

Dengan semua constraint ini, ruang solusi yang "valid" menjadi jauh lebih sempit—tetapi menemukannya di tengah lautan kemungkinan tetap merupakan tantangan komputasional yang besar.

---

### 1.3 Pendekatan Algoritmik: Dari Optimal ke Heuristik

Mengingat kompleksitas NP-hard, pendekatan matematis murni (*exact algorithms*) seperti *Integer Linear Programming* atau *Branch and Bound* menjadi tidak praktis untuk data dengan skala industri. Oleh karena itu, industri dan akademisi beralih ke **algoritma heuristik**.

Algoritma heuristik adalah metode yang:

- **Tidak menjamin solusi optimal**, tetapi menghasilkan solusi "cukup baik" dalam waktu komputasi yang wajar.
- **Menggunakan aturan prioritas** (*rules of thumb*) untuk membuat keputusan penempatan.
- **Dapat disesuaikan** dengan batasan spesifik domain.

Salah satu pendekatan heuristik yang terbukti efektif adalah **3D Bin Packing** dengan strategi *Bigger First*. Ide dasarnya sederhana namun powerful: **tempatkan barang terbesar terlebih dahulu**.

Mengapa pendekatan ini masuk akal? Bayangkan Anda mengemas koper untuk bepergian. Jika Anda memasukkan baju-baju kecil terlebih dahulu, kemungkinan besar laptop atau sepatu besar tidak akan muat di akhir. Sebaliknya, jika laptop dan sepatu dimasukkan dulu, baju-baju kecil dapat mengisi celah yang tersisa.

Flowchart berikut menggambarkan langkah-langkah algoritma secara sistematis:

```mermaid
flowchart TD
    A[Daftar Barang] --> B[Urutkan berdasarkan Volume]
    B --> C{Ada barang tersisa?}
    C -->|Ya| D[Ambil barang terbesar]
    D --> E[Cari posisi valid di kontainer]
    E --> F{Posisi ditemukan?}
    F -->|Ya| G[Tempatkan barang]
    G --> C
    F -->|Tidak| H[Tandai sebagai unfitted]
    H --> C
    C -->|Tidak| I[Hasil: Daftar Penempatan]
```

Mari kita telusuri alur ini:

1. **Pengurutan awal**: Semua barang diurutkan berdasarkan volume (panjang × lebar × tinggi) dari terbesar ke terkecil. Ini adalah langkah krusial yang menentukan urutan penempatan.

2. **Iterasi penempatan**: Algoritma mengambil barang terbesar yang belum ditempatkan, lalu mencari posisi valid di dalam kontainer. "Valid" berarti tidak bertumpukan dengan barang lain dan tidak melampaui batas kontainer.

3. **Pencarian posisi**: Untuk setiap barang, algoritma mencoba berbagai kombinasi koordinat (x, y, z) dan rotasi. Proses ini adalah bagian yang paling *compute-intensive*.

4. **Penanganan kegagalan**: Jika tidak ada posisi valid ditemukan (misalnya kontainer sudah terlalu penuh), barang ditandai sebagai "unfitted" dan dilanjutkan ke barang berikutnya.

5. **Output akhir**: Setelah semua barang diproses, algoritma menghasilkan daftar penempatan (koordinat + rotasi) dan daftar barang yang tidak muat.

**Mengapa strategi *Bigger First* efektif?**

1. **Prioritas yang tepat**: Barang besar memiliki lebih sedikit opsi penempatan yang valid. Dengan memprioritaskannya, kita menghindari situasi di mana barang besar "tersingkir" oleh barang kecil.

2. **Pondasi yang stabil**: Barang besar yang ditempatkan di dasar menciptakan permukaan yang rata untuk menumpuk barang di atasnya.

3. **Pemanfaatan celah**: Barang kecil lebih fleksibel dan dapat mengisi *void space* yang tidak dapat diisi oleh barang besar.

**Alternatif Strategi Heuristik**

*Bigger First* bukanlah satu-satunya strategi yang tersedia. Dalam literatur bin packing, terdapat beberapa pendekatan alternatif yang layak dipertimbangkan:

**Bottom-Left-Back (BLB)**: Strategi ini menempatkan setiap barang di posisi yang paling "bawah-kiri-belakang" yang tersedia. Koordinat Y (tinggi) diprioritaskan paling rendah, diikuti X (kiri), lalu Z (belakang). BLB cenderung menghasilkan penumpukan yang rapi secara visual, tetapi tidak selalu optimal untuk utilisasi volume.

**Best Fit Decreasing (BFD)**: Mirip dengan *Bigger First*, tetapi saat mencari posisi, algoritma memilih posisi yang menyisakan ruang kosong paling sedikit. Ini mengurangi fragmentasi ruang, tetapi membutuhkan komputasi lebih berat karena harus mengevaluasi lebih banyak opsi.

**First Fit Decreasing (FFD)**: Barang diurutkan dari besar ke kecil, lalu ditempatkan di posisi valid pertama yang ditemukan. Lebih cepat dari BFD karena tidak mencari posisi "terbaik", tetapi hasilnya cenderung kurang optimal.

Untuk implementasi kita, *Bigger First* dengan pencarian posisi standar dipilih karena menyeimbangkan **kualitas solusi** dan **kecepatan komputasi**. Dalam pengujian dengan data industri, strategi ini mencapai utilisasi volume 50-60% untuk kargo heterogen—angka yang kompetitif dengan algoritma yang lebih kompleks, tetapi dengan waktu eksekusi yang jauh lebih singkat.

**Trade-off: Kualitas vs Waktu**

Dalam dunia praktis, ada trade-off fundamental antara kualitas solusi dan waktu komputasi:

- **Algoritma exact** (seperti Integer Linear Programming) dapat menemukan solusi optimal, tetapi waktu eksekusinya meledak eksponensial seiring bertambahnya barang. Untuk 50+ barang, waktu bisa mencapai jam atau bahkan hari.

- **Algoritma heuristik** seperti *Bigger First* mengorbankan optimalitas demi kecepatan. Solusi yang dihasilkan mungkin "hanya" 95% sebaik solusi optimal, tetapi didapat dalam hitungan detik.

Untuk aplikasi industri di mana keputusan harus dibuat dengan cepat dan volume pengiriman tinggi, trade-off ini hampir selalu menguntungkan heuristik. Selisih beberapa persen utilisasi volume tidak sebanding dengan penundaan berhari-hari menunggu solusi "sempurna".

**Ilustrasi Numerik Sederhana**

Untuk memperjelas cara kerja algoritma, mari kita lihat contoh sederhana dengan 4 barang yang akan dimuat ke kontainer berukuran 1000 × 500 × 500 mm:

```
Barang Awal (belum diurutkan):
┌─────────┬────────────────────┬─────────────┐
│ Item    │ Dimensi (mm)       │ Volume      │
├─────────┼────────────────────┼─────────────┤
│ A       │ 200 × 200 × 200    │ 8.000.000   │
│ B       │ 400 × 300 × 250    │ 30.000.000  │
│ C       │ 100 × 100 × 100    │ 1.000.000   │
│ D       │ 300 × 200 × 200    │ 12.000.000  │
└─────────┴────────────────────┴─────────────┘

Setelah Pengurutan (Bigger First):
B (30jt) → D (12jt) → A (8jt) → C (1jt)

Proses Penempatan:
1. B ditempatkan di (0, 0, 0) — pojok kiri-bawah-depan
2. D ditempatkan di (400, 0, 0) — di samping B
3. A ditempatkan di (0, 250, 0) — di atas B
4. C ditempatkan di celah yang tersisa
```

Dengan urutan ini, barang besar mendapat prioritas posisi optimal, sementara barang kecil mengisi celah. Jika urutan dibalik (C dulu), barang C mungkin mengambil posisi yang seharusnya lebih baik untuk B.

---

### 1.4 Visi Sistem Load & Stuffing Calculator

Berdasarkan pemahaman tantangan di atas, kita akan membangun sistem **Load & Stuffing Calculator**—sebuah platform yang mentransformasi proses perencanaan pemuatan dari intuisi manual menjadi komputasi algoritmik dengan visualisasi interaktif.

**Visi Sistem:**

> Menyediakan solusi perencanaan pemuatan yang **akurat**, **visual**, dan **operasional**—menghubungkan kalkulasi algoritmik dengan kebutuhan praktis pekerja lapangan.

**Tujuan Operasional:**

Sistem yang kita bangun memiliki target yang terukur. Dari sisi *Fill Rate*, kita menargetkan 100%—artinya semua barang yang diminta harus dapat dimuat jika secara fisik memungkinkan. Untuk *Volume Utilization*, target kita adalah di atas 50% untuk kargo heterogen, sesuai dengan benchmark industri untuk barang dengan dimensi bervariasi. Waktu kalkulasi harus di bawah 40 detik untuk 300 barang, memastikan penggunaan interaktif tetap responsif. Akurasi geometris harus bebas *overlap* dan stabil—hasil yang dapat diterapkan secara literal di lapangan. Terakhir, output disajikan dalam dua format: visualisasi 3D interaktif untuk planner dan PDF instruksi untuk operator.

**Siapa yang Akan Menggunakan Sistem Ini?**

Sistem ini dirancang untuk melayani berbagai peran dalam rantai pasok logistik:

**Logistics Planner** adalah pengguna utama. Mereka bertanggung jawab merencanakan muatan untuk pengiriman mendatang. Skenario tipikal: planner menerima daftar order dari tim sales, lalu perlu menentukan berapa kontainer yang dibutuhkan dan bagaimana konfigurasi optimalnya. Dengan sistem ini, mereka dapat simulasi berbagai skenario dalam hitungan menit—bukan jam.

**Warehouse Supervisor** menggunakan output sistem sebagai panduan operasional. Visualisasi step-by-step membantu mereka mengarahkan operator forklift dengan presisi. PDF instruksi dapat dicetak dan dibawa ke lapangan, mengurangi miskomunikasi antara planner di kantor dan pekerja di gudang.

**Management** membutuhkan visibility untuk decision-making. Dashboard statistik—utilisasi volume, waktu kalkulasi, trend pengiriman—membantu mereka mengidentifikasi inefisiensi dan mengoptimalkan operasi secara keseluruhan.

**Skenario Penggunaan Khas:**

1. **Pengiriman rutin**: Sebuah pabrik furnitur mengirim produk ke retailer setiap minggu. Kombinasi produk relatif konsisten. Planner menggunakan template yang sudah ada, memodifikasi quantity, dan mendapat rencana baru dalam detik.

2. **Pengiriman ad-hoc**: Sebuah eksportir menerima order besar dengan produk yang belum pernah dikirim bersama. Planner memasukkan dimensi baru, bereksperimen dengan berbagai ukuran kontainer, dan menemukan konfigurasi paling efisien.

3. **Evaluasi kapasitas**: Sales team ingin memberikan quote kepada customer baru. Mereka perlu tahu apakah order tertentu muat dalam satu kontainer atau dua. Sistem memberikan jawaban instan tanpa harus menunggu tim operasional.

**Gambaran Arsitektur Sistem:**

Untuk mewujudkan visi di atas, kita memerlukan arsitektur yang modular—di mana setiap komponen memiliki tanggung jawab spesifik dan dapat dikembangkan secara independen. Sistem yang akan kita bangun terdiri dari tiga lapisan utama:

```mermaid
graph LR
    subgraph Frontend
        A[Web Browser] --> B[Next.js App]
        B --> C[Three.js Visualizer]
    end
    
    subgraph Backend
        D[Go API Server]
        E[PostgreSQL Database]
    end
    
    subgraph "Packing Service"
        F[Python Flask]
        G[py3dbp Algorithm]
    end
    
    B --> D
    D --> E
    D --> F
    F --> G
```

Mari kita telusuri alur data dalam sistem ini:

1. **User membuka browser** dan mengakses aplikasi Next.js. Di sini, mereka memilih kontainer dan menambahkan barang yang akan dimuat.

2. **Frontend mengirim request** ke Go API Server. Request ini berisi informasi kontainer dan daftar barang lengkap dengan dimensi dan quantity.

3. **Go API Server memvalidasi data** dan menyimpannya ke PostgreSQL. Ini memastikan data persisten dan dapat diakses kembali.

4. **Ketika user meminta kalkulasi**, Go API Server meneruskan request ke Packing Service (Python Flask). Di sinilah algoritma py3dbp bekerja untuk menghitung penempatan optimal.

5. **Hasil kalkulasi dikembalikan** ke Go API Server, yang menyimpannya ke database dan mengirimkannya ke frontend.

6. **Three.js Visualizer merender hasil** dalam bentuk 3D interaktif. User dapat memutar, zoom, dan melihat animasi step-by-step pemuatan.

Setiap komponen dipilih berdasarkan kekuatan spesifiknya. **Go** digunakan untuk API Server karena performa tinggi dan *concurrency* yang kuat melalui goroutines, plus kemudahan deployment sebagai single binary. **Python** menjadi pilihan untuk algoritma packing karena ekosistem pustaka matematikanya yang matang dan sintaks yang ekspresif untuk logika kompleks. **Three.js** menangani visualisasi 3D dengan *hardware-accelerated* rendering langsung di browser tanpa plugin. **PostgreSQL** menyediakan persistensi data yang reliable dengan dukungan JSON untuk data semi-structured. Terakhir, **Next.js** membangun frontend dengan server-side rendering, routing otomatis, dan integrasi TypeScript yang kuat.

**Hasil Akhir Buku Ini:**

Setelah menyelesaikan buku ini, pembaca akan memiliki pengetahuan dan *source code* untuk membangun platform yang mampu:

1. **Optimasi Maksimal**: Mencapai *Fill Rate* 100% dan *Volume Utilization* hingga 55,26% untuk kargo heterogen.
2. **Kecepatan Operasional**: Menyelesaikan perhitungan 300 barang dalam waktu kurang dari 40 detik.
3. **Akurasi Visual**: Menyajikan panduan pemuatan langkah-demi-langkah yang akurat secara geometris.
4. **Arsitektur Standar Industri**: Memahami cara membangun sistem multi-bahasa yang berkomunikasi secara efisien.

---

### Summary

Pada bab ini, kita telah memahami mengapa optimasi pemuatan kontainer merupakan masalah yang signifikan bagi industri logistik. *Container Loading Problem* adalah masalah NP-hard yang tidak dapat diselesaikan secara optimal dalam waktu yang wajar untuk skala industri. Pendekatan heuristik, khususnya algoritma *3D Bin Packing* dengan strategi *Bigger First*, menjadi solusi praktis yang akan kita implementasikan.

Kita juga telah mendefinisikan visi dan tujuan operasional sistem *Load & Stuffing Calculator* yang akan dibangun sepanjang buku ini. Pada bab berikutnya, kita akan merancang arsitektur layanan mikro yang memisahkan tanggung jawab antara manajemen data, komputasi algoritma, dan visualisasi.

---

### Further Reading

Berikut adalah referensi untuk memperdalam pemahaman konsep-konsep yang dibahas dalam bab ini:

- **Manajemen Logistik**: [Council of Supply Chain Management Professionals (CSCMP)](https://cscmp.org/)
- **Container Loading Problem**: [Bischoff, E. & Ratcliff, M. (1995) - Issues in the Development of CLP](https://doi.org/10.1016/0377-2217(95)00021-6)
- **NP-hard & Kompleksitas**: [Introduction to Algorithms (CLRS) - Chapter 34: NP-Completeness](https://mitpress.mit.edu/9780262046305/introduction-to-algorithms/)
- **Heuristic Algorithms**: [Geeksforgeeks: Heuristic Algorithms](https://www.geeksforgeeks.org/heuristic-algorithms/)
- **3D Bin Packing**: [py3dbp: 3D Bin Packing Library](https://github.com/enzoruiz/3dbinpacking)
- **Container Standard (ISO)**: [ISO 668:2020 - Series 1 Freight Containers](https://www.iso.org/standard/76912.html)

<!-- BAB 2 STARTS HERE -->

## Bab 2: Perancangan Arsitektur Layanan Mikro

Pada bab ini, kita akan merancang struktur sistem yang terdekopling (*decoupled*) untuk memisahkan manajemen data dari mesin kalkulasi. Tujuan utamanya adalah membangun arsitektur yang memungkinkan setiap komponen berkembang secara independen.

Terdapat beberapa tantangan teknis yang akan kita hadapi, meliputi:

- Bagaimana menangani *scalability* saat memproses kalkulasi kargo dalam jumlah besar.
- Bagaimana memilih *communication protocol* yang tepat antar layanan.
- Bagaimana mendefinisikan kontrak data yang jelas agar kedua layanan dapat berkembang tanpa saling merusak.

Dalam bab ini, kita akan menjawab tantangan tersebut dengan menerapkan arsitektur *Microservices*, membagi tanggung jawab antara *Backend API* (Go) dan *Packing Service* (Python), serta menggunakan REST HTTP dengan JSON sebagai protokol komunikasi.

**Dalam bab ini, kita akan membahas topik-topik berikut:**

- Mengapa Arsitektur *Microservices*
- Alur Transformasi Data
- Diagram Arsitektur Sistem
- Kontrak Komunikasi Antar Layanan
- Struktur Direktori Proyek

---

### 2.1 Mengapa Arsitektur Microservices

Sebelum Amazon, Netflix, dan perusahaan teknologi raksasa lainnya mempopulerkan istilah "microservices" pada awal 2010-an, arsitektur dominan untuk aplikasi enterprise adalah *monolith*. Sebuah aplikasi monolith adalah sistem di mana seluruh kode—dari antarmuka pengguna, logika bisnis, hingga akses database—dikemas dan di-deploy sebagai satu unit tunggal. Pendekatan ini memiliki keunggulan yang tidak bisa diabaikan: sederhana untuk dimulai, mudah di-debug karena semua kode berada di satu tempat, dan tidak memerlukan kompleksitas komunikasi jaringan internal.

Banyak aplikasi sukses dibangun dengan arsitektur monolith, dan bagi banyak kasus sederhana, monolith tetap merupakan pilihan yang valid. Namun, ketika kita menganalisis karakteristik unik dari sistem *Load & Stuffing Calculator*, beberapa pertanyaan arsitektural muncul yang membuat monolith menjadi kurang ideal.

Pertanyaan pertama adalah tentang **bahasa pemrograman**. Dunia pengembangan perangkat lunak tidak memiliki bahasa universal yang unggul untuk semua tugas. Go, misalnya, sangat baik untuk membangun API server yang menangani ribuan koneksi simultan—model *goroutines*-nya memungkinkan *concurrency* ringan yang sulit ditandingi. Namun, Go bukanlah bahasa pilihan untuk implementasi algoritma matematika kompleks. Di sisi lain, Python memiliki ekosistem pustaka ilmiah yang luar biasa kaya—NumPy, SciPy, dan khususnya untuk kasus kita, library *3D bin packing* yang sudah teruji. Dalam arsitektur monolith, kita dipaksa memilih satu bahasa untuk seluruh sistem, mengorbankan kekuatan salah satu domain.

Pertanyaan kedua menyangkut **karakteristik beban kerja**. Perhatikan perbedaan fundamental antara dua operasi berikut: "tampilkan daftar produk" versus "hitung penempatan 300 barang dalam kontainer". Operasi pertama adalah query database sederhana yang selesai dalam hitungan milidetik. Operasi kedua bisa memakan waktu 30 detik atau lebih karena melibatkan kalkulasi geometris intensif. Dalam arsitektur monolith, kedua operasi ini bersaing untuk CPU yang sama. Ketika seorang user menjalankan kalkulasi berat, user lain yang hanya ingin melihat daftar produk mungkin mengalami *slowdown*—karena CPU sibuk menghitung koordinat penempatan, bukan melayani query sederhana.

Pertanyaan ketiga berkaitan dengan **siklus pengembangan dan deployment**. Tim yang mengerjakan antarmuka pengguna biasanya memiliki ritme yang berbeda dengan tim yang mengembangkan algoritma. UI mungkin di-update harian dengan perbaikan kecil, sementara algoritma direvisi bulanan dengan perubahan signifikan. Dalam monolith, setiap perubahan—sekecil apapun—memerlukan deployment ulang seluruh sistem. Ini menciptakan coupling yang tidak perlu antara komponen yang seharusnya independen.

Arsitektur *microservices* menjawab ketiga pertanyaan di atas dengan prinsip sederhana: **pisahkan komponen berdasarkan tanggung jawab, lalu biarkan masing-masing berkomunikasi melalui antarmuka yang terdefinisi dengan jelas**. Untuk sistem kita, ini berarti dua layanan utama:

```mermaid
graph TB
    subgraph "Service 1: Backend API"
        A[Go API Server]
        B[PostgreSQL]
    end
    
    subgraph "Service 2: Packing Service"
        C[Python Flask]
        D[py3dbp Algorithm]
    end
    
    A <-->|REST/JSON| C
    A --> B
```

Diagram di atas mungkin terlihat sederhana—dan memang itulah intinya. Kompleksitas bukanlah tujuan; pemisahan tanggung jawab adalah tujuannya. **Backend API** menggunakan Go untuk menangani semua interaksi dengan dunia luar: menerima request dari frontend, memvalidasi input, mengelola autentikasi, dan menyimpan data ke PostgreSQL. Ia adalah "pintu gerbang" tunggal—tidak ada yang berbicara langsung ke Packing Service selain Backend API. **Packing Service** menggunakan Python untuk satu tugas spesifik: menerima data geometris dan mengembalikan koordinat penempatan. Ia tidak perlu tahu tentang autentikasi, tidak perlu tahu tentang database, tidak perlu tahu tentang konteks bisnis—ia hanya menghitung.

Pemilihan bahasa untuk masing-masing layanan bukanlah keputusan sewenang-wenang, melainkan berdasarkan kesesuaian dengan karakteristik tugas. Backend API harus menangani banyak request simultan dengan latensi rendah—karakteristik yang menjadi kekuatan utama Go. Model *goroutines* dalam Go memungkinkan ribuan koneksi concurrent dikelola dengan memory footprint yang sangat kecil, sesuatu yang sulit dicapai dengan bahasa interpretasi tradisional.

Di sisi lain, Packing Service memiliki karakteristik berbeda: kalkulasi intensif tetapi jarang dipanggil. Untuk tugas semacam ini, ekspresi sintaksis lebih penting daripada performa raw. Python dengan pustaka py3dbp yang sudah teruji menjadi pilihan natural—kode algoritma dapat ditulis dengan jelas dan mudah di-maintain, sementara overhead interpretasi tidak signifikan dibandingkan kompleksitas kalkulasi geometris itu sendiri.

> Keputusan ini bukan tentang mengikuti tren teknologi terbaru. Microservices bukanlah solusi universal—bagi banyak aplikasi, mereka menambah kompleksitas yang tidak diperlukan. Namun, untuk kasus spesifik kita dengan kebutuhan polyglot dan karakteristik beban kerja yang sangat berbeda, pemisahan ini memberikan manfaat nyata.

#### Memisahkan Manajemen Data dari Komputasi Berat

Mari kita dalami lebih jauh tentang mengapa pemisahan ini penting dari perspektif operasional. Dalam teori sistem terdistribusi, kita mengenal dua kategori beban kerja yang berbeda secara fundamental:

**Beban I/O-bound** adalah operasi yang menghabiskan sebagian besar waktunya menunggu data dari sumber eksternal—database, jaringan, atau filesystem. Operasi CRUD (Create, Read, Update, Delete) termasuk kategori ini. CPU hampir selalu idle menunggu disk atau network. Untuk jenis beban ini, *concurrency* adalah kunci: kita ingin menangani banyak operasi secara bersamaan, masing-masing menunggu I/O-nya sendiri.

**Beban CPU-bound** adalah operasi yang menghabiskan sebagian besar waktunya melakukan komputasi aktif. Algoritma 3D Bin Packing termasuk kategori ini. CPU bekerja keras menghitung koordinat, memeriksa collision, mencoba rotasi—hampir tidak ada waktu idle. Untuk jenis beban ini, raw processing power dan algoritma efisien adalah kunci.

Mencampur kedua jenis beban dalam satu proses menciptakan konflik. Bayangkan skenario ini: seorang user meminta kalkulasi untuk 300 barang (operasi CPU-bound yang memakan 30 detik), sementara 10 user lain mencoba mengakses daftar produk (operasi I/O-bound yang seharusnya selesai dalam 5 milidetik). Dalam arsitektur monolith dengan satu proses, 10 user tersebut harus menunggu kalkulasi selesai sebelum request mereka dilayani. Pengalaman pengguna memburuk secara dramatis.

Dengan memisahkan layanan, Backend API dapat melayani ratusan request ringan secara concurrent menggunakan goroutines Go, sementara Packing Service menjalankan kalkulasi berat di proses terpisah. Ketika user meminta kalkulasi, Backend API hanya mengirimkan HTTP request ke Packing Service lalu menunggu respons—CPU Backend API tidak terbebani oleh algoritma.

Diagram berikut mengilustrasikan bagaimana Backend API tetap responsif melayani banyak request ringan, sementara operasi kalkulasi berat dilakukan oleh Packing Service di proses yang terisolasi:

```mermaid
flowchart TB
    subgraph BACKEND["🟢 Backend API (Go) — Tetap Responsif"]
        direction TB
        R1["GET /products — 5ms"]
        R2["POST /plans — 10ms"]
        R3["GET /containers — 3ms"]
        R4["...ratusan request lain"]
        
        DB[(PostgreSQL)]
        
        R1 --> DB
        R2 --> DB
        R3 --> DB
    end
    
    subgraph PACKING["🔴 Packing Service (Python) — CPU Intensif"]
        direction TB
        ALG["py3dbp Algorithm<br/>Memproses 300 barang<br/>Durasi: 30 detik"]
    end
    
    BACKEND -->|"HTTP POST /pack<br/>(hanya saat kalkulasi)"| PACKING
    PACKING -->|"JSON Response<br/>(placements + stats)"| BACKEND
```

Perbedaan karakteristik kedua layanan ini menentukan strategi resource yang berbeda. Backend API menangani ratusan request per menit dengan waktu respons milidetik (5-50ms). Beban CPU-nya rendah karena sebagian besar waktu dihabiskan menunggu I/O. Scaling dilakukan berdasarkan jumlah user concurrent. Sebaliknya, Packing Service hanya menangani beberapa request per jam, tetapi setiap request memakan waktu detik hingga puluhan detik (5-40s) dengan beban CPU yang sangat tinggi. Scaling dilakukan berdasarkan antrian kalkulasi.

Dengan pemisahan ini, saat seorang user meminta kalkulasi 300 barang, user lain yang hanya ingin melihat daftar produk **tidak perlu menunggu**—request mereka tetap dilayani dengan cepat oleh Backend API.

**Keuntungan pemisahan ini:**

1. **Scaling independen**: Packing Service dapat di-scale horizontal (tambah instance) saat beban kalkulasi tinggi, tanpa mempengaruhi biaya Backend API.
2. **Fault isolation**: Jika Packing Service crash karena kalkulasi terlalu berat, Backend API tetap melayani request lain.
3. **Resource optimization**: Backend API dapat berjalan di instance murah (low CPU), sementara Packing Service menggunakan instance high-CPU saat dibutuhkan.

---

### 2.2 Alur Transformasi Data

Setelah memahami arsitektur tingkat tinggi, mari kita dalami bagaimana data berubah bentuk dari input pengguna hingga menjadi rencana pemuatan yang dapat divisualisasikan. Pemahaman ini akan membantu kita merancang kontrak data yang tepat antar layanan.

Bayangkan seorang planner logistik yang ingin merencanakan muatan untuk pengiriman. Ia membuka aplikasi, memilih kontainer 40ft, lalu menambahkan daftar barang yang akan dikirim. Setelah menekan tombol "Hitung", dalam hitungan detik ia melihat visualisasi 3D yang menunjukkan posisi setiap barang di dalam kontainer.

Di balik kesederhanaan interaksi tersebut, data melewati **enam tahap transformasi**:

```mermaid
flowchart TD
    subgraph "1. Input Pengguna"
        A["Form UI:<br/>Pilih Kontainer + Daftar Barang"]
    end
    
    subgraph "2. Normalisasi"
        B["Backend API:<br/>Validasi & Konversi ke mm"]
    end
    
    subgraph "3. Kalkulasi"
        C["Packing Service:<br/>Jalankan Algoritma 3D-BPP"]
    end
    
    subgraph "4. Hasil Mentah"
        D["Koordinat XYZ +<br/>Kode Rotasi + Step Number"]
    end
    
    subgraph "5. Persistensi"
        E["Backend API:<br/>Simpan ke Database"]
    end
    
    subgraph "6. Visualisasi"
        F["Three.js:<br/>Render 3D + Animasi Step"]
    end
    
    A --> B --> C --> D --> E --> F
```

Mari kita telusuri perjalanan data ini tahap demi tahap:

**Tahap 1-2: Dari Input ke Normalisasi**

User mengisi form dengan memilih kontainer dan menambahkan barang. Dimensi barang mungkin diinput dalam berbagai satuan (cm, m, atau mm). Backend API menerima request ini, memvalidasi kelengkapan data, dan mengkonversi semua dimensi ke milimeter. Standarisasi ini penting agar Packing Service tidak perlu "tahu" tentang preferensi satuan user.

**Tahap 3-4: Kalkulasi Algoritma**

Data yang sudah dinormalisasi dikirim ke Packing Service. Algoritma py3dbp menerima daftar barang, mengurutkannya berdasarkan volume, lalu menempatkan satu per satu ke dalam kontainer virtual. Outputnya adalah daftar koordinat (`pos_x`, `pos_y`, `pos_z`), kode rotasi, dan nomor urut pemuatan (`step_number`).

**Tahap 5-6: Dari Persistensi ke Visualisasi**

Backend API menyimpan hasil kalkulasi ke database agar user dapat mengaksesnya kembali tanpa menghitung ulang. Ketika frontend meminta data visualisasi, koordinat ini diambil dari database dan ditransformasi menjadi objek 3D menggunakan Three.js.

Tabel berikut merangkum input, proses, dan output di setiap tahap:

| Tahap               | Input                            | Proses               | Output              |
| ------------------- | -------------------------------- | -------------------- | ------------------- |
| **1. Input**        | Pilihan kontainer, daftar produk | User mengisi form    | JSON request        |
| **2. Normalisasi**  | Dimensi berbagai satuan          | Validasi & konversi  | Payload dalam mm    |
| **3. Kalkulasi**    | Container + Items + Options      | Algoritma py3dbp     | Daftar penempatan   |
| **4. Hasil Mentah** | Output algoritma                 | —                    | Koordinat + rotasi  |
| **5. Persistensi**  | Hasil kalkulasi                  | Simpan ke PostgreSQL | Record tersimpan    |
| **6. Visualisasi**  | Data dari database               | Render Three.js      | Scene 3D interaktif |

Setiap komponen dalam arsitektur kita bertanggung jawab atas satu jenis transformasi. Komunikasi antar komponen menggunakan format JSON yang terdefinisi dengan jelas—inilah "kontrak" yang akan kita bahas di bagian selanjutnya.

---

### 2.3 Diagram Arsitektur Sistem

Setelah memahami alur transformasi data, kita dapat melihat gambaran lengkap arsitektur sistem. Arsitektur ini dirancang dengan prinsip *separation of concerns*—setiap komponen memiliki tanggung jawab yang jelas dan terdefinisi.

Sistem kita terdiri dari tiga layanan utama yang berkomunikasi melalui jaringan. Dalam development, ketiga layanan ini dapat berjalan di mesin yang sama. Dalam production, mereka dapat ditempatkan di server berbeda sesuai kebutuhan scaling.

**Mengapa Tiga Layanan—Bukan Dua atau Empat?**

Keputusan untuk membagi sistem menjadi tiga layanan bukan angka arbitrer—ini hasil dari analisis karakteristik beban kerja dan kebutuhan tim.

Jika kita menggabungkan Backend API dan Packing Service menjadi **dua layanan** (Frontend + Backend+Packing), kita kehilangan fleksibilitas polyglot. Go dan Python harus berjalan dalam satu proses atau berkomunikasi melalui FFI (Foreign Function Interface) yang kompleks. Scaling juga menjadi masalah—kita tidak dapat menambah kapasitas kalkulasi tanpa juga menambah kapasitas API.

Jika kita memecah lebih lanjut menjadi **empat atau lebih layanan**—misalnya memisahkan Authentication Service, Reporting Service, dan sebagainya—kita menambah kompleksitas operasional yang belum diperlukan. Setiap layanan tambahan berarti network hop tambahan, titik kegagalan tambahan, dan koordinasi deployment tambahan. Untuk skala proyek ini, overhead tersebut tidak sebanding dengan manfaatnya.

Tiga layanan memberikan sweet spot: cukup terpisah untuk independensi teknis, tetapi tidak terlalu banyak hingga menyulitkan operasional.

```mermaid
flowchart LR
    subgraph Client
        Browser[Web Browser]
    end
    
    subgraph "Frontend"
        NextJS[Next.js App]
        ThreeJS[Three.js Visualizer]
    end
    
    subgraph "Backend"
        GoAPI[Go API Server<br/>:8080]
        PG[(PostgreSQL<br/>:5432)]
    end
    
    subgraph "Packing Service"
        Flask[Python Flask<br/>:5051]
        Py3DBP[py3dbp Engine]
    end
    
    Browser --> NextJS
    NextJS --> ThreeJS
    NextJS -->|REST API| GoAPI
    GoAPI --> PG
    GoAPI -->|POST /pack| Flask
    Flask --> Py3DBP
```

Mari kita bahas masing-masing layanan:

**Frontend** adalah titik interaksi dengan pengguna. Next.js menangani routing, state management, dan rendering UI. Di dalamnya, Three.js berjalan sebagai *library* yang merender visualisasi 3D langsung di browser pengguna. Layanan ini hanya berkomunikasi dengan Backend API—ia tidak "tahu" tentang keberadaan Packing Service.

**Backend** adalah pusat koordinasi sistem. Go API Server menerima semua request dari frontend, memvalidasi input, mengelola autentikasi, dan berinteraksi dengan database PostgreSQL. Ketika ada request kalkulasi, Backend API yang bertanggung jawab memanggil Packing Service dan memproses hasilnya. Pemilihan Go untuk komponen ini didasarkan pada kebutuhannya menangani banyak request simultan dengan latensi rendah.

**Packing Service** adalah *worker* terspesialisasi. Flask menyediakan HTTP endpoint sederhana, sementara py3dbp melakukan pekerjaan berat: menghitung penempatan optimal ratusan barang dalam kontainer. Layanan ini tidak memiliki akses ke database dan tidak mengetahui konteks bisnis—ia hanya menerima dimensi dan mengembalikan koordinat. Isolasi ini memungkinkan kita mengganti algoritma di masa depan tanpa mempengaruhi komponen lain.

#### Alur Interaksi Antar Komponen

Untuk memahami bagaimana komponen-komponen ini bekerja sama, mari kita lihat sequence diagram untuk operasi kalkulasi—inti dari fungsionalitas sistem ini:

```mermaid
sequenceDiagram
    participant FE as Frontend
    participant API as Go API
    participant DB as PostgreSQL
    participant Pack as Packing Service
    
    FE->>API: POST /api/plans/{id}/calculate
    API->>DB: Fetch container & items
    DB-->>API: Container + Items data
    API->>Pack: POST /pack (JSON payload)
    Pack->>Pack: Run py3dbp algorithm
    Pack-->>API: Placements + Stats
    API->>DB: Save calculation result
    API-->>FE: Return visualization data
```

Diagram ini mengilustrasikan prinsip penting dalam arsitektur kita: **Backend API adalah satu-satunya komponen yang berkomunikasi dengan Packing Service**. Frontend tidak pernah memanggil Packing Service secara langsung, dan Packing Service tidak pernah mengakses database.

Mengapa desain ini penting? Ada beberapa alasan:

1. **Keamanan**: Packing Service tidak perlu kredensial database. Jika layanan ini dikompromikan, penyerang tidak dapat mengakses data sensitif.

2. **Konsistensi data**: Backend API dapat memastikan bahwa data yang dikirim ke Packing Service sudah tervalidasi dan dalam format yang benar.

3. **Caching dan optimisasi**: Backend API dapat menyimpan hasil kalkulasi. Jika user meminta kalkulasi yang sama, hasilnya bisa diambil dari database tanpa memanggil Packing Service lagi.

4. **Monitoring terpusat**: Semua *traffic* ke Packing Service melewati Backend API, memudahkan logging dan observability.

---

### 2.4 Kontrak Komunikasi Antar Layanan

Komunikasi antara *Backend API* (Go) dan *Packing Service* (Python) memerlukan protokol yang tepat. Sebelum mendefinisikan kontrak data dan menjelaskan pilihan kita, mari kita pahami terlebih dahulu lanskap protokol komunikasi yang tersedia dalam pengembangan sistem terdistribusi modern.

#### Mengenal Protokol Komunikasi Antar Layanan

**REST (Representational State Transfer)**

REST bukanlah protokol dalam arti teknis, melainkan *architectural style* yang didefinisikan oleh Roy Fielding dalam disertasi doktoralnya pada tahun 2000. Fielding adalah salah satu arsitek utama protokol HTTP/1.1, dan REST lahir dari pemahamannya tentang bagaimana web bekerja dengan sukses.

Inti dari REST adalah memodelkan dunia sebagai *resources*—entitas yang dapat diidentifikasi dengan URL. Sebuah kontainer adalah resource (`/containers/123`), sebuah produk adalah resource (`/products/456`), dan operasi CRUD dilakukan menggunakan HTTP verbs: GET untuk membaca, POST untuk membuat, PUT/PATCH untuk memperbarui, DELETE untuk menghapus. Pendekatan ini intuitif karena mengikuti cara kita memikirkan data.

REST menjadi dominan karena kesederhanaannya. Tidak perlu tooling khusus—browser dapat menjadi client, curl dapat menguji endpoint, dan hampir setiap bahasa pemrograman memiliki HTTP library bawaan. Contoh penggunaan: hampir semua API publik (Twitter, GitHub, Stripe) menggunakan REST karena aksesibilitasnya.

**gRPC (Google Remote Procedure Call)**

gRPC lahir dari kebutuhan internal Google untuk komunikasi antar layanan dengan performa tinggi. Ketika Anda memiliki ratusan microservices yang saling berkomunikasi ribuan kali per detik, overhead JSON parsing dan HTTP/1.1 menjadi signifikan.

gRPC menggunakan Protocol Buffers (protobuf) sebagai format serialisasi—binary format yang jauh lebih compact dan cepat di-parse dibanding JSON. Kontrak didefinisikan dalam file `.proto`, lalu dikompilasi menjadi kode client/server dalam bahasa target. Ini memberikan type safety: jika struktur data berubah, kesalahan terdeteksi saat kompilasi, bukan saat runtime.

gRPC juga berjalan di atas HTTP/2, mendukung streaming bidirectional—client dan server dapat mengirim aliran data secara bersamaan. Contoh penggunaan: Netflix menggunakan gRPC untuk komunikasi internal antar layanannya yang berjumlah ratusan.

Namun, gRPC memiliki *learning curve* yang lebih tinggi. Anda perlu memahami protobuf, menginstal compiler, dan me-generate kode. Debugging juga lebih sulit karena payload binary tidak human-readable.

**GraphQL**

GraphQL dikembangkan oleh Facebook dan dirilis sebagai open source pada 2015. Masalah yang ingin dipecahkan GraphQL sangat spesifik: *over-fetching* dan *under-fetching* pada REST API.

Bayangkan Anda membangun aplikasi mobile yang menampilkan profil pengguna. REST API mungkin mengembalikan seluruh objek user (nama, email, alamat, riwayat pesanan, preferensi...) padahal Anda hanya butuh nama dan foto. Ini adalah over-fetching. Sebaliknya, jika Anda butuh data dari beberapa resource sekaligus (user + posts + comments), Anda harus melakukan multiple requests. Ini adalah under-fetching.

GraphQL memecahkan ini dengan membiarkan client menentukan persis data apa yang dibutuhkan dalam satu query. Server hanya mengembalikan data yang diminta—tidak lebih, tidak kurang. Contoh penggunaan: GitHub API v4 menggunakan GraphQL, begitu juga Shopify.

GraphQL sangat powerful untuk frontend yang kompleks, tetapi menambah kompleksitas di sisi server. Setiap query harus di-resolve, dan optimisasi performa (N+1 problem, caching) memerlukan pertimbangan khusus.

**Message Queues**

Message queues seperti RabbitMQ, Apache Kafka, atau Amazon SQS memecahkan masalah yang berbeda dari tiga protokol di atas. Jika REST, gRPC, dan GraphQL adalah komunikasi *synchronous* (client menunggu respons), message queues adalah *asynchronous* (client mengirim pesan lalu melanjutkan aktivitas).

Konsep dasarnya sederhana: producer mengirim pesan ke queue, consumer mengambil dan memproses pesan dari queue. Keduanya tidak perlu online bersamaan. Jika consumer sedang sibuk atau down, pesan menunggu di queue hingga dapat diproses.

Message queues ideal untuk operasi yang tidak memerlukan respons instan: mengirim email, memproses video, menggenerate laporan. Contoh penggunaan: Uber menggunakan Kafka untuk memproses jutaan event lokasi per detik; Instagram menggunakan RabbitMQ untuk background job processing.

#### Pemilihan Protokol: Mengapa REST?

Dengan pemahaman di atas, kita dapat membuat keputusan yang informed. Untuk komunikasi Backend API ↔ Packing Service, kita memilih **REST HTTP dengan JSON**. Berikut alasannya:

**1. Kesederhanaan Implementasi**

REST tidak memerlukan *tooling* khusus. Baik Go maupun Python memiliki dukungan HTTP bawaan yang matang. Tidak perlu mengkompilasi *proto files* (seperti gRPC) atau mempelajari query language baru (seperti GraphQL).

Perhatikan betapa ringkasnya kode yang diperlukan di kedua sisi. Di sisi Packing Service (Python), Flask menangani *routing* dan *parsing* JSON secara otomatis:

```python
# Python (Flask) - menerima request
@app.route('/pack', methods=['POST'])
def pack():
    data = request.get_json()  # Flask parse JSON otomatis
    result = process_packing(data)
    return jsonify(result)  # Konversi dict ke JSON
```

Di sisi Backend API (Go), mengirim request hanya membutuhkan satu baris menggunakan *standard library*:

```go
// Go - mengirim request menggunakan net/http bawaan
resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
```

Bandingkan dengan gRPC yang memerlukan: definisi `.proto`, kompilasi dengan `protoc`, dan *generated code* di kedua sisi. Untuk komunikasi antar dua layanan, overhead ini tidak sebanding dengan manfaatnya.

**2. Debuggability**

JSON adalah format yang *human-readable*. Ketika terjadi error, kita dapat dengan mudah melihat payload request/response menggunakan tools standar seperti `curl`, Postman, atau bahkan browser.

```bash
# Debug dengan curl
curl -X POST http://localhost:5051/pack \
  -H "Content-Type: application/json" \
  -d '{"units":"mm","container":{...},"items":[...]}'
```

Bandingkan dengan gRPC yang menggunakan binary encoding—debugging memerlukan tools khusus.

**3. Pola Request-Response Cocok dengan Use Case**

Operasi packing kita bersifat *synchronous*: user menekan tombol, menunggu hasil, lalu melihat visualisasi. Ini cocok dengan model request-response REST.

Jika kita menggunakan Message Queue (seperti RabbitMQ atau Kafka), arsitektur menjadi lebih kompleks karena kita perlu menangani callback atau polling untuk mendapatkan hasil. Kerumitan ini tidak sepadan untuk use case kita.

**4. Interoperabilitas**

REST adalah standar de facto untuk web API. Jika di masa depan kita ingin membuka Packing Service untuk sistem eksternal, tidak perlu perubahan arsitektur—endpoint yang sama dapat digunakan.

> Catatan: gRPC akan menjadi pilihan yang lebih baik jika kita memiliki ratusan layanan yang berkomunikasi dengan latensi ultra-rendah. Untuk sistem dengan dua layanan seperti ini, overhead REST (sekitar 1-5ms tambahan dibanding gRPC) tidak signifikan dibandingkan waktu kalkulasi algoritma (detik hingga puluhan detik).

#### Definisi Kontrak Data

Dengan REST sebagai protokol, kita perlu mendefinisikan "kontrak"—struktur JSON yang harus dipatuhi kedua sisi agar komunikasi berhasil. Kontrak ini penting karena Backend API (Go) dan Packing Service (Python) dikembangkan secara terpisah. Tanpa kontrak yang jelas, perubahan di satu sisi dapat merusak komunikasi.

**Endpoint yang Tersedia**

Packing Service menyediakan dua endpoint:

| Method | Endpoint | Deskripsi |
|--------|----------|-----------|
| POST | `/pack` | Menjalankan kalkulasi packing |
| GET | `/health` | Health check untuk monitoring |

Endpoint `/health` digunakan untuk memastikan layanan berjalan—berguna untuk monitoring dan *load balancer*. Endpoint utama adalah `/pack` yang menerima data kontainer dan barang, lalu mengembalikan hasil kalkulasi.

**Struktur Request**

Ketika Backend API ingin menghitung penempatan barang, ia mengirim POST request ke `/pack` dengan payload berikut:

```json
{
  "units": "mm",
  "container": {
    "length": 12032,
    "width": 2352,
    "height": 2698,
    "max_weight": 28000
  },
  "items": [
    {
      "item_id": "SKU-001",
      "label": "Kardus Besar",
      "length": 600,
      "width": 400,
      "height": 400,
      "weight": 25.5,
      "quantity": 50
    },
    {
      "item_id": "SKU-002",
      "label": "Kardus Kecil",
      "length": 300,
      "width": 200,
      "height": 200,
      "weight": 5.0,
      "quantity": 100
    }
  ],
  "options": {
    "fix_point": true,
    "check_stable": true,
    "support_surface_ratio": 0.75,
    "bigger_first": true,
    "put_type": 1
  }
}
```

Payload ini terdiri dari empat bagian utama:

1. **`units`**: Satuan panjang yang digunakan. Dalam sistem kita, Backend API selalu mengirim dalam milimeter (`"mm"`) untuk konsistensi.

2. **`container`**: Dimensi internal kontainer pengiriman. Nilai pada contoh di atas (12032 × 2352 × 2698 mm) adalah dimensi standar kontainer 40ft.

3. **`items`**: Daftar barang yang akan dimuat. Setiap item memiliki dimensi, berat, dan quantity. Packing Service akan "mengekspansi" quantity—misalnya 50 unit menjadi 50 penempatan terpisah.

4. **`options`**: Parameter algoritma. Opsi `fix_point` mengaktifkan simulasi gravitasi, `check_stable` memvalidasi stabilitas, dan `bigger_first` memprioritaskan barang besar.

Tabel berikut merangkum semua field yang tersedia:

| Field                           | Tipe   | Deskripsi                                  |
| ------------------------------- | ------ | ------------------------------------------ |
| `units`                         | string | Satuan panjang: `"mm"`, `"cm"`, atau `"m"` |
| `container.length/width/height` | float  | Dimensi internal kontainer                 |
| `container.max_weight`          | float  | Batas berat maksimum (kg)                  |
| `items[].item_id`               | string | Identifier unik item                       |
| `items[].quantity`              | int    | Jumlah unit yang akan dimuat               |
| `options.fix_point`             | bool   | Aktifkan simulasi gravitasi                |
| `options.check_stable`          | bool   | Validasi stabilitas tumpukan               |
| `options.support_surface_ratio` | float  | Persentase minimum penopang (0.75 = 75%)   |
| `options.bigger_first`          | bool   | Prioritaskan barang besar terlebih dahulu  |

**Struktur Response (Sukses)**

Jika kalkulasi berhasil, Packing Service mengembalikan daftar penempatan beserta statistik:

```json
{
  "success": true,
  "data": {
    "units": "mm",
    "placements": [
      {
        "item_id": "SKU-001",
        "label": "Kardus Besar",
        "pos_x": 0,
        "pos_y": 0,
        "pos_z": 0,
        "rotation": 0,
        "step_number": 1
      },
      {
        "item_id": "SKU-001",
        "label": "Kardus Besar",
        "pos_x": 600,
        "pos_y": 0,
        "pos_z": 0,
        "rotation": 0,
        "step_number": 2
      }
    ],
    "unfitted": [
      {
        "item_id": "SKU-002",
        "label": "Kardus Kecil",
        "count": 5
      }
    ],
    "stats": {
      "expanded_items": 150,
      "fitted_count": 145,
      "unfitted_count": 5,
      "pack_time_ms": 1250
    }
  }
}
```

Response ini memberikan informasi yang dibutuhkan untuk visualisasi dan pelaporan:

- **`placements`**: Setiap objek dalam array ini merepresentasikan satu barang yang berhasil ditempatkan. Koordinat `pos_x`, `pos_y`, `pos_z` menunjukkan posisi pojok kiri-bawah-depan barang. Field `step_number` menentukan urutan pemuatan—Frontend menggunakan informasi ini untuk animasi step-by-step.

- **`unfitted`**: Daftar barang yang tidak muat ke dalam kontainer. Ini terjadi jika volume atau berat melebihi kapasitas. Frontend menampilkan informasi ini sebagai peringatan kepada user.

- **`stats`**: Statistik kalkulasi. Field `pack_time_ms` berguna untuk monitoring performa.

| Field | Deskripsi |
|-------|-----------|
| `placements[]` | Daftar barang yang berhasil dimuat dengan posisi dan rotasi |
| `placements[].pos_x/y/z` | Koordinat pojok kiri-bawah-depan barang |
| `placements[].rotation` | Kode rotasi (0-5, lihat Bab 6) |
| `placements[].step_number` | Urutan pemuatan untuk visualisasi |
| `unfitted[]` | Daftar barang yang tidak muat |
| `stats.pack_time_ms` | Waktu komputasi algoritma dalam milidetik |

**Struktur Response (Error)**

Tidak semua request berhasil. Packing Service mengembalikan error response ketika terjadi masalah dengan input atau kondisi yang tidak dapat ditangani. Memahami format error ini penting agar Backend API dapat memberikan feedback yang bermakna kepada user.

Berikut contoh response ketika validasi gagal:

```json
{
  "success": false,
  "error": {
    "code": "INVALID_REQUEST",
    "message": "container.height must be > 0",
    "details": {}
  }
}
```

**Kapan Error Terjadi?**

| Kode Error | Penyebab | Contoh |
|------------|----------|--------|
| `INVALID_REQUEST` | Data tidak valid | Dimensi negatif, field wajib kosong |
| `CONTAINER_TOO_SMALL` | Kontainer tidak cukup | Semua barang lebih besar dari kontainer |
| `WEIGHT_EXCEEDED` | Berat melebihi batas | Total berat barang > `max_weight` |
| `INTERNAL_ERROR` | Kesalahan server | Bug algoritma, out of memory |

**Bagaimana Backend API Menangani Error?**

Ketika menerima error response, Backend API sebaiknya:

1. **Log error** untuk debugging dan monitoring.
2. **Translate pesan** jika perlu (misal: dari bahasa teknis ke bahasa yang dipahami user).
3. **Kembalikan ke frontend** dengan status HTTP yang sesuai (400 untuk input error, 500 untuk server error).

```go
// Contoh penanganan error di Go
if !packResponse.Success {
    log.Printf("Packing failed: %s", packResponse.Error.Message)
    return nil, fmt.Errorf("kalkulasi gagal: %s", packResponse.Error.Message)
}
```

Dengan kontrak error yang jelas, tim frontend dapat menampilkan pesan yang informatif kepada user, bukan hanya "terjadi kesalahan".

---

### 2.5 Struktur Direktori Proyek

Organisasi kode yang baik mempermudah navigasi, maintenance, dan kolaborasi tim. Dalam proyek multi-bahasa seperti ini, struktur direktori yang jelas menjadi semakin penting—kita perlu memisahkan kode Go, Python, dan TypeScript tanpa menciptakan kebingungan.

Struktur yang kita adopsi mengikuti konvensi standar komunitas Go, dengan adaptasi untuk mengakomodasi Packing Service (Python) dan Frontend (Next.js):

```
load-stuffing-calculator/
├── cmd/
│   └── packing/              # Python Packing Service
│       ├── app.py            # Flask application entry
│       ├── packing.py        # Core packing logic
│       ├── schema.py         # Request/response validation
│       └── 3D-bin-packing/   # py3dbp library (vendored)
│
├── internal/                 # Go internal packages
│   ├── gateway/
│   │   └── packing_gateway.go  # HTTP client untuk Packing Service
│   ├── packer/
│   │   └── packer.go         # Packer interface & types
│   ├── service/
│   │   └── packing_service.go  # Bisnis logic packing
│   └── handler/
│       └── plan_handler.go   # HTTP handler untuk plans
│
├── web/                      # Next.js Frontend
│   ├── components/
│   │   └── stuffing-viewer.tsx  # React wrapper untuk visualizer
│   └── lib/
│       └── StuffingVisualizer/  # Three.js visualization engine
```

Mari kita bahas alasan di balik setiap keputusan desain:

**Direktori `cmd/`**

Dalam konvensi Go, `cmd/` berisi entry point untuk aplikasi yang dapat dieksekusi. Setiap subdirektori di dalamnya adalah satu aplikasi. Meskipun Packing Service ditulis dalam Python, kita menempatkannya di sini untuk konsistensi—semua "runnable applications" berada di satu tempat.

*Manfaat*: Developer baru dapat langsung melihat apa saja yang bisa dijalankan dalam proyek ini hanya dengan melihat isi `cmd/`.

**Direktori `internal/`**

Ini adalah fitur khusus Go. Kode di dalam `internal/` tidak dapat di-import oleh package di luar modul ini. Ini memberikan "pagar" yang mencegah dependensi yang tidak diinginkan.

Perhatikan struktur di dalamnya:

- **`gateway/`**: Kode yang berkomunikasi dengan sistem eksternal (Packing Service). Jika kita menambah integrasi baru (misal: shipping API), gateway-nya ditempatkan di sini.

- **`packer/`**: Mendefinisikan interface `Packer` dan tipe data terkait. Ini adalah "kontrak" internal yang memungkinkan multiple implementasi.

- **`service/`**: Business logic. `packing_service.go` mengorkestrasikan panggilan ke gateway dan transformasi data.

- **`handler/`**: HTTP handlers. Menerima request, memanggil service, mengembalikan response.

*Manfaat*: Struktur berlapis ini memisahkan *transport layer* (handler) dari *business logic* (service) dari *external communication* (gateway). Jika kita ingin mengganti HTTP dengan gRPC, hanya `handler/` yang perlu diubah.

**Direktori `web/`**

Frontend hidup sepenuhnya terpisah dengan struktur standar Next.js. Kita tidak mencampur kode TypeScript dengan Go—keduanya memiliki toolchain dan konvensi yang berbeda.

- **`components/`**: React components. `stuffing-viewer.tsx` adalah wrapper yang mengintegrasikan Three.js ke dalam ekosistem React.

- **`lib/`**: Library internal. `StuffingVisualizer/` adalah engine visualisasi 3D yang dapat digunakan kembali di berbagai komponen.

*Manfaat*: Tim frontend dapat bekerja di `web/` tanpa menyentuh kode backend, dan sebaliknya.

**Mengapa Struktur Ini Penting untuk Pembaca**

Sepanjang buku ini, kita akan membangun kode secara bertahap. Dengan memahami struktur di atas:

1. **Navigasi lebih mudah**: Ketika membahas "HTTP client untuk Packing Service", pembaca langsung tahu mencarinya di `internal/gateway/`.

2. **Konteks yang jelas**: Setiap file memiliki tanggung jawab spesifik. `packing.py` berisi algoritma, bukan HTTP handling—itu tugas `app.py`.

3. **Extensibility**: Jika pembaca ingin menambah fitur baru, struktur ini memberikan panduan di mana kode harus ditempatkan.

---

### Summary

Pada bab ini, kita telah merancang arsitektur *Microservices* untuk sistem *Load & Stuffing Calculator*. Berikut ringkasan keputusan desain yang dibuat:

1. **Pemisahan layanan**: Backend API (Go) menangani orkestrasi dan persistensi data, sementara Packing Service (Python) fokus pada algoritma. Tiga layanan (bukan dua atau empat) memberikan keseimbangan antara independensi teknis dan kesederhanaan operasional.

2. **Alur transformasi data**: Data melewati enam tahap—dari input pengguna hingga visualisasi 3D—dengan setiap komponen bertanggung jawab atas satu jenis transformasi.

3. **Pemahaman protokol komunikasi**: Kita mempelajari karakteristik empat protokol utama—REST, gRPC, GraphQL, dan Message Queues—beserta sejarah dan use case masing-masing. Pemahaman ini memberikan fondasi untuk memilih protokol yang tepat.

4. **Pemilihan REST HTTP dengan JSON**: REST dipilih karena kesederhanaan implementasi, kemudahan debugging, dan kesesuaian dengan pola request-response synchronous yang dibutuhkan. gRPC atau Message Queues dapat dipertimbangkan di masa depan jika kebutuhan berubah.

5. **Kontrak data yang jelas**: Request dan response payload terdefinisi dengan baik, termasuk penanganan error yang konsisten. Ini memungkinkan Backend API dan Packing Service dikembangkan secara independen.

6. **Struktur direktori terorganisir**: Mengikuti konvensi Go dengan adaptasi untuk proyek multi-bahasa, memisahkan transport layer, business logic, dan external communication.

Sistem kini memiliki *blueprint* arsitektur yang solid. Pada bab berikutnya, kita akan mulai mengimplementasikan *Backend API* dengan bahasa Go, dimulai dari struktur proyek dan *database schema*.

---

### Further Reading

- **Microservices Pattern**: [Martin Fowler: Microservices](https://martinfowler.com/articles/microservices.html)
- **REST API Design**: [Microsoft REST API Guidelines](https://github.com/microsoft/api-guidelines)
- **REST vs gRPC**: [Google Cloud: Choosing an API Style](https://cloud.google.com/blog/products/api-management/understanding-grpc-openapi-and-rest-and-when-to-use-them)
- **Go Project Layout**: [Standard Go Project Layout](https://github.com/golang-standards/project-layout)
- **Flask Documentation**: [Flask Quickstart](https://flask.palletsprojects.com/en/3.0.x/quickstart/)
- **JSON Schema**: [JSON Schema Specification](https://json-schema.org/)



<!-- BAB 3 STARTS HERE -->

## Bab 3: Pengembangan Backend API dengan Bahasa Go

Pada bab sebelumnya, kita telah merancang arsitektur sistem dengan dua layanan terpisah: Backend API yang menangani orkestrasi dan persistensi data, serta Packing Service yang fokus pada algoritma. Kita juga telah mendefinisikan kontrak komunikasi REST JSON yang akan menghubungkan keduanya. Sekarang saatnya mengubah *blueprint* tersebut menjadi kode yang berjalan.

Pada bab ini, kita akan membangun Backend API dari nol menggunakan bahasa Go. Pemilihan Go bukanlah kebetulan—seperti yang dibahas di Bab 2, karakteristik I/O-bound dari Backend API membutuhkan bahasa yang kuat dalam menangani *concurrency*. Go, dengan model *goroutines*-nya, memungkinkan kita melayani ribuan request simultan tanpa kompleksitas *threading* tradisional.

Namun, menguasai *syntax* Go saja tidak cukup. Setiap bahasa pemrograman memiliki *idiom*—cara penulisan yang dianggap natural dan efektif oleh komunitasnya. Kode Go yang "benar secara sintaksis" tetapi tidak idiomatis akan sulit dibaca oleh developer Go lain, sulit di-*maintain*, dan seringkali kurang efisien. Oleh karena itu, bagian awal bab ini akan membahas fondasi dan idiom Go sebelum kita menulis kode aplikasi.

Terdapat beberapa tantangan teknis yang akan kita hadapi dalam bab ini:

- Mendesain *database schema* yang merepresentasikan entitas logistik (kontainer, produk, rencana muat)
- Membangun *repository layer* yang *type-safe* menggunakan SQLC
- Mengimplementasikan *service layer* yang mengorkestrasikan logika bisnis
- Menyediakan HTTP API menggunakan Gin framework
- Mengintegrasikan Backend API dengan Packing Service melalui HTTP client

**Dalam bab ini, kita akan membahas topik-topik berikut:**

- Fondasi Bahasa Go
- Idiom dan Konvensi Go
- Setup Proyek dan Tooling
- Database Schema Design
- Type-safe Repository Layer dengan SQLC
- Business Logic Layer (Service)
- HTTP Router dan Handler
- Integrasi dengan Packing Service

**Technical Requirements**

Untuk menyelesaikan bab ini, Anda memerlukan alat dan perangkat lunak berikut:

- **Go version 1.21 atau lebih baru** — [go.dev/dl](https://go.dev/dl/)
- **PostgreSQL version 15** — Database relasional untuk persistensi data
- **SQLC CLI** — Code generator untuk type-safe SQL queries — [sqlc.dev](https://sqlc.dev/)
- **Goose CLI** — Database migration tool — [github.com/pressly/goose](https://github.com/pressly/goose)
- **Air** — Hot reload untuk development — [github.com/air-verse/air](https://github.com/air-verse/air)
- **VS Code dengan Go extension** — IDE yang direkomendasikan

---

### 3.1 Fondasi Bahasa Go

Go—atau sering disebut Golang—diciptakan di Google pada tahun 2007 oleh Robert Griesemer, Rob Pike, dan Ken Thompson. Ketiga insinyur ini memiliki latar belakang yang mengesankan: Thompson adalah salah satu pencipta Unix dan bahasa C, Pike adalah co-creator UTF-8 dan Plan 9, sementara Griesemer pernah bekerja pada JavaScript V8 engine dan Java HotSpot VM. Pengalaman puluhan tahun mereka dengan sistem berskala besar di Google membentuk filosofi Go yang sangat pragmatis.

Go dirilis ke publik pada 2009 dan dengan cepat diadopsi oleh perusahaan-perusahaan yang membutuhkan *performance* tinggi dengan *developer experience* yang sederhana. Docker, Kubernetes, Terraform—infrastruktur modern yang mungkin sudah Anda gunakan—semuanya ditulis dalam Go.

#### Prinsip Desain Go

Filosofi Go dapat diringkas dalam beberapa prinsip yang saling terkait. Memahami prinsip-prinsip ini penting karena mereka mempengaruhi cara kita menulis kode sepanjang buku ini.

**Simplicity over Cleverness**

Go secara sengaja menghilangkan fitur-fitur yang dianggap menambah kompleksitas tanpa memberikan nilai proporsional: tidak ada *inheritance*, tidak ada *generic* (hingga Go 1.18), tidak ada *exceptions*, tidak ada *operator overloading*. Keputusan ini kontroversial bagi developer yang datang dari bahasa lain, tetapi ada alasan kuat di baliknya.

Rob Pike, dalam presentasinya "Simplicity is Complicated", menjelaskan bahwa Go dirancang untuk tim besar yang bekerja pada *codebase* yang bertahan bertahun-tahun. Dalam konteks ini, kode yang mudah dibaca lebih berharga daripada kode yang "elegan" tetapi membutuhkan pengetahuan khusus untuk dipahami.

**Readability over Terseness**

Go tidak berusaha meminimalkan jumlah karakter yang diketik. Sebaliknya, Go mengoptimalkan untuk *reading*—karena kode dibaca jauh lebih sering daripada ditulis. Ini terlihat dari nama *package* yang eksplisit (`context`, `encoding/json`), *error handling* yang verbose, dan konvensi penamaan yang mendorong kejelasan.

**Explicit over Implicit**

Go menghindari "magic". Tidak ada *dependency injection framework* yang misterius, tidak ada *annotation* yang mengubah perilaku secara silent, tidak ada *implicit conversion* antar tipe. Jika sesuatu terjadi, Anda bisa melihatnya langsung di kode.

**Composition over Inheritance**

Go tidak memiliki *class* atau *inheritance* dalam arti tradisional. Sebagai gantinya, Go menggunakan *struct embedding* dan *interfaces* untuk mencapai *code reuse*. Pendekatan ini mendorong komposisi objek-objek kecil daripada hierarki kelas yang dalam.

Untuk memperdalam pemahaman tentang prinsip-prinsip ini, ada tiga dokumen fundamental yang wajib dibaca:

- [Effective Go](https://go.dev/doc/effective_go): Panduan resmi tentang cara menulis kode Go yang idiomatis
- [Go Proverbs](https://go-proverbs.github.io/): Kumpulan prinsip filosofis Go dari Rob Pike
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments): Panduan *code review* yang menjadi standar komunitas

#### Concurrency Model

Salah satu kekuatan utama Go adalah model *concurrency*-nya yang sederhana namun powerful. Dalam konteks Backend API kita, ini berarti kemampuan menangani ribuan request HTTP secara bersamaan tanpa kompleksitas *thread management*.

**Goroutines**

*Goroutine* adalah *lightweight thread* yang dikelola oleh Go *runtime*, bukan sistem operasi. Ini membuatnya jauh lebih ringan—sebuah *goroutine* hanya membutuhkan sekitar 2KB *stack* (yang dapat tumbuh dinamis), dibandingkan dengan *thread* OS yang membutuhkan sekitar 1MB.

Membuat *goroutine* sangat sederhana—cukup tambahkan kata kunci `go` sebelum *function call*:

```go
func main() {
    go handleRequest(request1)  // Berjalan di goroutine terpisah
    go handleRequest(request2)  // Berjalan di goroutine lain
    // main() terus berjalan tanpa menunggu
}
```

HTTP server bawaan Go secara otomatis menjalankan setiap *request handler* dalam *goroutine* terpisah. Ini berarti ketika satu request menunggu respons dari database, request lain tetap dapat dilayani.

**Channels**

*Goroutines* berkomunikasi melalui *channels*—sebuah mekanisme untuk mengirim dan menerima nilai antar *goroutines* secara *thread-safe*. Filosofi Go tentang *concurrency* terangkum dalam proverb:

> "Don't communicate by sharing memory; share memory by communicating."

Artinya: daripada menggunakan *shared variable* dengan *mutex locks* (seperti di banyak bahasa lain), Go mendorong pengiriman data melalui *channels*. Ini mengurangi risiko *race condition* dan membuat kode lebih mudah dipahami.

```go
// Membuat channel
results := make(chan int)

// Mengirim nilai ke channel
go func() {
    result := expensiveCalculation()
    results <- result  // Kirim ke channel
}()

// Menerima nilai dari channel
value := <-results  // Blokir hingga nilai tersedia
```

#### Error Handling Philosophy

Jika Anda datang dari bahasa dengan *exceptions* (Java, Python, JavaScript), *error handling* Go mungkin terasa verbose. Go tidak memiliki *try-catch*. Sebagai gantinya, *functions* yang dapat gagal mengembalikan *error* sebagai nilai kembalian terakhir:

```go
file, err := os.Open("config.json")
if err != nil {
    return fmt.Errorf("gagal membuka config: %w", err)
}
defer file.Close()
```

Pola ini—check *error* segera setelah *function call*—akan Anda lihat ratusan kali dalam kode Go. Meskipun verbose, pendekatan ini memiliki keuntungan:

1. **Explicit**: Tidak ada *hidden control flow*. Anda tahu persis di mana *error* bisa terjadi.
2. **Composable**: *Errors* adalah nilai biasa yang dapat di-*wrap*, di-*log*, atau di-*transform*.
3. **Performance**: Tidak ada *overhead* dari *exception mechanism*.

**Error Wrapping**

Sejak Go 1.13, kita dapat *wrap errors* dengan konteks tambahan menggunakan `%w`:

```go
if err != nil {
    return fmt.Errorf("failed to connect to packing service: %w", err)
}
```

Ini memungkinkan kita membangun *chain of errors* yang menjelaskan konteks kegagalan dari titik terendah hingga tertinggi—sangat membantu untuk *debugging*.

#### Packages dan Modules

Go mengorganisir kode dalam *packages*—direktori yang berisi satu atau lebih file `.go`. Setiap file dalam direktori yang sama harus memiliki *package declaration* yang sama.

```go
// File: internal/service/packing_service.go
package service

type PackingService struct {
    // ...
}
```

*Package* menentukan *visibility*: nama yang diawali huruf kapital (`PackingService`) dapat diakses dari *package* lain (*exported*), sedangkan huruf kecil (`packingService`) hanya dapat diakses dari dalam *package* yang sama (*unexported*).

**Go Modules**

Sejak Go 1.11, *dependency management* ditangani oleh Go Modules. Setiap proyek memiliki file `go.mod` yang mendefinisikan *module path* dan *dependencies*:

```
module load-stuffing-calculator

go 1.21

require (
    github.com/gin-gonic/gin v1.9.1
    github.com/jackc/pgx/v5 v5.5.0
)
```

*Module path* juga berfungsi sebagai *import path* untuk *packages* dalam proyek:

```go
import "load-stuffing-calculator/internal/service"
```

**Internal Packages**

Go memiliki fitur khusus untuk *encapsulation*: direktori bernama `internal`. *Packages* di dalam `internal` hanya dapat di-*import* oleh kode dalam *parent directory* yang sama. Ini mencegah *external packages* bergantung pada implementasi internal kita.

```
load-stuffing-calculator/
├── internal/           # Tidak dapat di-import dari luar module
│   ├── service/
│   └── repository/
└── cmd/                # Entry points
    └── api/
```

Dengan pemahaman fondasi ini, kita siap mempelajari *idiom* dan konvensi yang membuat kode Go benar-benar idiomatis.

---

### 3.2 Idiom dan Konvensi Go

Mengetahui *syntax* Go dan memahami cara menulis Go yang idiomatis adalah dua hal berbeda. Kode idiomatis mengikuti pola dan konvensi yang telah menjadi standar komunitas—kode yang ketika dibaca developer Go lain, terasa "natural" dan mudah dipahami.

Bagian ini akan membahas idiom-idiom kunci yang akan kita gunakan sepanjang proyek. Tujuannya bukan hanya agar kode kita "benar", tetapi agar kode kita *maintainable* oleh siapa pun yang familiar dengan ekosistem Go.

#### Naming Conventions

Go memiliki konvensi penamaan yang sangat spesifik, dan mengikutinya sangat penting untuk *readability*.

**Nama Pendek dan Bermakna**

Berbeda dengan bahasa lain yang mendorong nama variabel deskriptif panjang, Go menyukai nama pendek—terutama untuk variabel dengan *scope* kecil:

```go
// Idiomatis Go
for i, v := range items {
    // i = index, v = value
}

// Bukan gaya Go
for itemIndex, itemValue := range items {
    // Terlalu verbose untuk scope kecil
}
```

Konvensi umum:
- `i`, `j`, `k` untuk loop counters
- `n` untuk count atau length
- `v` untuk value dalam range
- `k` untuk key dalam map
- `err` untuk error
- `ctx` untuk context
- `req`, `res` untuk request/response

Namun, untuk variabel dengan *scope* lebih luas atau *exported*, gunakan nama deskriptif: `containerRepository`, `packingService`.

**MixedCaps, Bukan snake_case**

Go menggunakan `MixedCaps` atau `mixedCaps`—tidak pernah *snake_case*:

```go
// Benar
type PackingService struct {}
func calculateVolume() {}

// Salah
type packing_service struct {}
func calculate_volume() {}
```

**Akronim Tetap Kapital**

Akronim seperti HTTP, URL, ID harus tetap dalam bentuk kapital penuh atau lower-case penuh:

```go
// Benar
var userID string
var httpClient *http.Client
type URLParser struct {}

// Salah
var userId string
var httpClient *http.Client
type UrlParser struct {}
```

#### Interface Design

Go memiliki pendekatan unik terhadap *interfaces* yang berbeda dari bahasa OOP tradisional.

**"Accept Interfaces, Return Structs"**

Prinsip ini mendorong fleksibilitas input dan kejelasan output:

```go
// Baik: menerima interface, mengembalikan struct konkret
func NewPackingService(repo PackRepository, gateway PackingGateway) *PackingService {
    return &PackingService{
        repo:    repo,
        gateway: gateway,
    }
}

// Kurang baik: mengembalikan interface
func NewPackingService(repo PackRepository) PackingServiceInterface {
    // Menyembunyikan tipe konkret tanpa alasan jelas
}
```

Mengapa? Karena mengembalikan *struct* konkret memberikan kejelasan tentang apa yang akan digunakan *caller*, sementara menerima *interface* memberikan fleksibilitas untuk *testing* dan *decoupling*.

**Small Interfaces**

Go mendorong *interfaces* kecil—idealnya satu atau dua *methods*. Ini adalah pola yang sangat berbeda dari bahasa seperti Java:

```go
// Interface kecil - mudah di-implement dan di-mock
type PackRepository interface {
    SavePlacements(ctx context.Context, planID string, placements []Placement) error
}

// Interface kecil lainnya
type PackingGateway interface {
    Calculate(ctx context.Context, req PackRequest) (*PackResponse, error)
}
```

*Standard library* Go penuh dengan contoh: `io.Reader` hanya memiliki satu method (`Read`), `io.Writer` hanya `Write`. Kesederhanaan ini memungkinkan komposisi yang powerful.

**Interface di Sisi Consumer**

Dalam Go, *interface* biasanya didefinisikan di *package* yang menggunakannya (*consumer*), bukan di *package* yang mengimplementasikannya (*provider*):

```go
// Di package service (consumer)
type PackRepository interface {
    SavePlacements(ctx context.Context, planID string, placements []Placement) error
}

type PackingService struct {
    repo PackRepository  // Menerima interface
}

// Di package repository (provider)
type PostgresRepository struct {
    db *pgxpool.Pool
}

// PostgresRepository secara implisit implement PackRepository
func (r *PostgresRepository) SavePlacements(ctx context.Context, planID string, placements []Placement) error {
    // implementasi...
}
```

Ini memungkinkan *loose coupling*—repository tidak perlu tahu tentang interface yang didefinisikan service.

#### Struct Design

**Constructor Functions: `NewXxx()`**

Go tidak memiliki *constructors* seperti bahasa OOP. Sebagai gantinya, kita menggunakan *factory functions* dengan prefix `New`:

```go
// Constructor function
func NewPackingService(repo PackRepository, gateway PackingGateway) *PackingService {
    return &PackingService{
        repo:    repo,
        gateway: gateway,
    }
}

// Penggunaan
service := NewPackingService(repo, gateway)
```

**Zero Value yang Berguna**

Desain *struct* Anda agar *zero value* (nilai default ketika tidak diinisialisasi) tetap berguna atau setidaknya tidak menyebabkan *panic*:

```go
// sync.Mutex bisa langsung digunakan tanpa inisialisasi
var mu sync.Mutex
mu.Lock()  // Bekerja!

// bytes.Buffer juga
var buf bytes.Buffer
buf.WriteString("hello")  // Bekerja!
```

#### Error Handling Patterns

Kita sudah membahas *error handling* di section sebelumnya, tetapi berikut beberapa pola idiomatis tambahan.

**Handle Error Segera**

Jangan menunda penanganan *error*:

```go
// Idiomatis: handle langsung
file, err := os.Open(filename)
if err != nil {
    return fmt.Errorf("open config: %w", err)
}
defer file.Close()

// Tidak idiomatis: menunda handling
file, err := os.Open(filename)
// ... kode lain ...
if err != nil {
    // Terlalu jauh dari sumber error
}
```

**Error Wrapping dengan Konteks**

Tambahkan konteks yang membantu *debugging*:

```go
// Baik: konteks jelas
if err != nil {
    return fmt.Errorf("save placements for plan %s: %w", planID, err)
}

// Kurang baik: hanya pass-through
if err != nil {
    return err  // Kehilangan konteks
}
```

**Sentinel Errors untuk Kasus Spesifik**

Untuk *error* yang perlu di-check secara spesifik oleh *caller*, gunakan *sentinel errors*:

```go
// Definisikan sentinel error
var ErrNotFound = errors.New("resource not found")

// Penggunaan di function
func (r *Repository) GetByID(ctx context.Context, id string) (*Container, error) {
    // ...
    if rows.Next() == false {
        return nil, ErrNotFound
    }
}

// Caller dapat check secara spesifik
container, err := repo.GetByID(ctx, id)
if errors.Is(err, ErrNotFound) {
    // Handle 404
}
```

#### Testing Patterns

Go memiliki *testing framework* bawaan yang sederhana namun powerful.

**Table-Driven Tests**

Pola ini sangat idiomatis di Go—mendefinisikan *test cases* dalam *slice* lalu *iterate*:

```go
func TestCalculateVolume(t *testing.T) {
    tests := []struct {
        name     string
        length   float64
        width    float64
        height   float64
        expected float64
    }{
        {"cube", 10, 10, 10, 1000},
        {"rectangle", 20, 10, 5, 1000},
        {"zero dimension", 0, 10, 10, 0},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := CalculateVolume(tt.length, tt.width, tt.height)
            if result != tt.expected {
                t.Errorf("got %v, want %v", result, tt.expected)
            }
        })
    }
}
```

Keuntungan:
- Mudah menambah *test cases* baru
- Setiap *case* memiliki nama deskriptif
- *Subtests* dengan `t.Run()` dapat di-*run* secara individual

**Mocking dengan Interfaces**

Karena kita sudah mendesain dengan *small interfaces*, *mocking* menjadi sederhana:

```go
// Mock implementation
type mockPackingGateway struct {
    calculateFunc func(ctx context.Context, req PackRequest) (*PackResponse, error)
}

func (m *mockPackingGateway) Calculate(ctx context.Context, req PackRequest) (*PackResponse, error) {
    return m.calculateFunc(ctx, req)
}

// Dalam test
func TestPackingService_Calculate(t *testing.T) {
    mock := &mockPackingGateway{
        calculateFunc: func(ctx context.Context, req PackRequest) (*PackResponse, error) {
            return &PackResponse{Success: true}, nil
        },
    }
    
    service := NewPackingService(repo, mock)
    // ... test logic
}
```

---

### 3.3 Setup Proyek dan Tooling

Sebelum menulis kode aplikasi, kita perlu menyiapkan fondasi proyek: inisialisasi *module*, struktur direktori, dan *tools* yang akan membantu *development workflow*.

#### Inisialisasi Module

Setiap proyek Go modern dimulai dengan inisialisasi *module*. Buka terminal di direktori proyek dan jalankan:

```bash
mkdir load-stuffing-calculator
cd load-stuffing-calculator
go mod init load-stuffing-calculator
```

Perintah ini membuat file `go.mod` yang menjadi "identitas" proyek:

```
module load-stuffing-calculator

go 1.21
```

Ketika kita menambahkan *dependencies* (seperti Gin atau pgx), Go akan secara otomatis mengupdate file ini. Untuk menambahkan *dependency* dan membersihkan yang tidak terpakai:

```bash
go mod tidy
```

#### Struktur Direktori

Struktur direktori yang baik memudahkan navigasi dan pemahaman arsitektur. Kita akan mengikuti konvensi Go dengan adaptasi untuk proyek *multi-service*:

```
load-stuffing-calculator/
├── cmd/
│   ├── api/                  # Entry point Backend API
│   ├── db/                   # Database tools (migrations, queries)
│   │   ├── migrations/       # Goose migration files
│   │   └── queries/          # SQLC query files
│   └── packing/              # Packing Service (Python)
│
├── internal/
│   ├── config/               # Konfigurasi aplikasi
│   ├── store/                # Generated SQLC code
│   ├── service/              # Business logic
│   ├── gateway/              # External service communication
│   └── handler/              # HTTP handlers
│
├── web/                      # Frontend (Next.js) - Bab 5
│
├── go.mod
├── go.sum
├── sqlc.yaml                 # SQLC configuration
└── .air.toml                 # Air hot-reload config
```

Mari kita bahas setiap direktori:

**`cmd/`** berisi *entry points* dan *tooling*. Direktori `cmd/api/` adalah *entry point* Backend API, sementara `cmd/db/` menyimpan file-file terkait database—migrasi Goose dan *query files* SQLC.

**`internal/`** adalah direktori khusus Go—kode di dalamnya tidak dapat di-*import* dari luar *module*. Ini menciptakan "batas" yang jelas antara *public API* dan implementasi internal.

Di dalam `internal/`:
- **`config/`**: Membaca dan memvalidasi konfigurasi dari *environment variables*
- **`store/`**: Kode Go yang di-*generate* oleh SQLC dari *query files*
- **`service/`**: *Business logic*—orkestrasi antar komponen
- **`gateway/`**: Komunikasi dengan *external services* (Packing Service)
- **`handler/`**: *HTTP handlers*—menerima *request*, memanggil *service*, mengembalikan *response*

#### Development Tools

**Air: Hot Reload**

Ketika mengembangkan API, sangat tidak nyaman jika harus *restart* server setiap kali ada perubahan kode. Air mengatasi ini dengan *watching* file dan *auto-rebuild*:

```bash
# Install Air
go install github.com/air-verse/air@latest

# Inisialisasi config
air init
```

Ini membuat file `.air.toml`. Edit beberapa konfigurasi penting:

```toml
[build]
  cmd = "go build -o ./tmp/main ./cmd/api"
  bin = "./tmp/main"
  delay = 1000

[log]
  time = true
```

Sekarang cukup jalankan `air` di terminal, dan setiap perubahan file `.go` akan otomatis *trigger rebuild*.

**Go Test**

Go memiliki *testing framework* bawaan. Untuk menjalankan semua *tests*:

```bash
go test ./...
```

Flag berguna:
- `-v`: *Verbose output*
- `-cover`: Menampilkan *test coverage*
- `-run TestName`: Menjalankan *test* spesifik

#### Membuat Entry Point Awal

Mari buat *entry point* minimal untuk memastikan setup berjalan. Pertama, buat struktur direktori:

```bash
mkdir -p cmd/api
```

Kemudian buat file `main.go` di dalam direktori `cmd/api/`:

```go
package main

import (
    "log"
    "net/http"

    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()

    r.GET("/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "status": "healthy",
        })
    })

    log.Println("Server starting on :8080")
    if err := r.Run(":8080"); err != nil {
        log.Fatal(err)
    }
}
```

Install Gin dan jalankan:

```bash
go get github.com/gin-gonic/gin
air
```

Akses `http://localhost:8080/health` untuk memverifikasi server berjalan.

#### Environment Variables dan .env File

Sebelum melanjutkan, kita perlu membahas bagaimana aplikasi akan membaca konfigurasi. Di production, konfigurasi seperti *database URL* dan *port* disediakan melalui *environment variables*—ini adalah best practice yang direkomendasikan oleh [The Twelve-Factor App](https://12factor.net/config). Namun, mengetik `export DATABASE_URL=...` setiap kali membuka terminal sangat tidak nyaman saat development.

Solusinya adalah **godotenv**—library Go yang membaca file `.env` dan mengisi environment variables secara otomatis. Install terlebih dahulu:

```bash
go get github.com/joho/godotenv
```

Buat file `.env` di root proyek (jangan commit ke Git!):

```bash
# .env - Konfigurasi development lokal
DATABASE_URL=postgres://user:password@localhost:5432/loadstuff?sslmode=disable
SERVER_PORT=8080
PACKING_SERVICE_URL=http://localhost:5000
```

Tambahkan `.env` ke `.gitignore`:

```bash
echo ".env" >> .gitignore
```

Sekarang update `main.go` untuk menggunakan godotenv:

```go
package main

import (
    "log"
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
)

func main() {
    // Load .env file jika ada (opsional untuk development)
    // Di production, environment variables di-set langsung oleh orchestrator
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, using environment variables")
    }

    r := gin.Default()

    r.GET("/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "status": "healthy",
        })
    })

    log.Println("Server starting on :8080")
    if err := r.Run(":8080"); err != nil {
        log.Fatal(err)
    }
}
```

Perhatikan bahwa `godotenv.Load()` tidak meng-*error* jika file tidak ditemukan—kita hanya mencetak pesan log. Ini penting karena di production, file `.env` tidak akan ada (dan tidak seharusnya ada).

Dengan *setup* ini, kita siap melanjutkan ke desain *database schema*.

---

### 3.4 Database Schema Design

Sebelum menulis *query* atau *business logic*, kita perlu merancang struktur database yang akan menyimpan data sistem. Desain yang baik di tahap ini akan memudahkan pengembangan selanjutnya dan mencegah refaktor yang mahal di kemudian hari. Seperti kata pepatah dalam dunia *software engineering*: kesalahan desain yang ditemukan di tahap awal biaya perbaikannya jauh lebih murah daripada yang ditemukan di production.

#### Entity-Relationship Diagram

Langkah pertama dalam mendesain database adalah mengidentifikasi *entitas*—objek-objek utama yang perlu kita simpan. Untuk sistem *Load & Stuffing Calculator*, kita perlu menjawab pertanyaan: "Data apa saja yang harus bertahan (*persist*) di antara sesi pengguna?"

Setelah analisis *use cases* dari Bab 1, kita mengidentifikasi lima entitas utama:

```mermaid
erDiagram
    CONTAINER ||--o{ PLAN : "digunakan dalam"
    PLAN ||--o{ PLAN_ITEM : "memiliki"
    PRODUCT ||--o{ PLAN_ITEM : "direferensikan oleh"
    PLAN ||--o{ PLACEMENT : "menghasilkan"

    CONTAINER {
        uuid id PK
        string name
        float length_mm
        float width_mm
        float height_mm
        float max_weight_kg
        timestamp created_at
    }

    PRODUCT {
        uuid id PK
        string label
        string sku
        float length_mm
        float width_mm
        float height_mm
        float weight_kg
        timestamp created_at
    }

    PLAN {
        uuid id PK
        uuid container_id FK
        string status
        timestamp created_at
        timestamp calculated_at
    }

    PLAN_ITEM {
        uuid id PK
        uuid plan_id FK
        uuid product_id FK
        int quantity
    }

    PLACEMENT {
        uuid id PK
        uuid plan_id FK
        uuid product_id FK
        float pos_x
        float pos_y
        float pos_z
        int rotation
        int step_number
    }
```

Mari kita bahas setiap entitas dan keputusan desain di baliknya:

**Container** menyimpan dimensi kontainer yang tersedia untuk dipilih pengguna. Perhatikan bahwa semua dimensi menggunakan satuan milimeter (`mm`) dan berat dalam kilogram (`kg`). Keputusan ini bukan kebetulan—menggunakan satu satuan konsisten di seluruh sistem menghilangkan kebutuhan konversi berulang dan mengurangi risiko bug perhitungan. Milimeter dipilih karena memberikan presisi yang cukup tanpa *floating point issues* yang sering muncul pada meter.

**Product** adalah barang yang akan dimuat. Kita memisahkan `sku` (*Stock Keeping Unit*) dari `label` karena keduanya memiliki fungsi berbeda: SKU adalah kode unik untuk sistem (sering berupa kombinasi huruf-angka seperti "CHAIR-RED-001"), sementara label adalah nama yang lebih deskriptif untuk pengguna ("Kursi Merah Model A"). Constraint `UNIQUE` pada SKU memastikan tidak ada duplikasi.

**Plan** merepresentasikan satu sesi perencanaan pemuatan. Mengapa perlu entitas terpisah, bukan langsung menyimpan hasil? Karena sebuah *plan* memiliki *lifecycle*: dimulai sebagai `draft` saat pengguna memilih kontainer dan barang, berubah menjadi `calculating` saat algoritma berjalan, dan akhirnya `completed` atau `failed`. Field `calculated_at` memungkinkan kita melacak kapan kalkulasi terakhir dilakukan—berguna untuk *caching* dan audit.

**PlanItem** adalah *junction table* (atau *associative table*) yang menghubungkan Plan dengan Product. Dalam database relasional, ini adalah cara standar untuk merepresentasikan relasi many-to-many dengan atribut tambahan. Satu *plan* bisa memiliki banyak produk, dan satu produk bisa muncul di banyak *plan*. Field `quantity` mencatat berapa unit dari produk tersebut yang akan dimuat.

**Placement** adalah hasil kalkulasi algoritma—posisi setiap barang di dalam kontainer. Ini adalah "output" dari Packing Service yang kita simpan untuk visualisasi. Field `rotation` menggunakan kode integer (0-5) yang merepresentasikan 6 kemungkinan orientasi box dalam ruang 3D. Field `step_number` menentukan urutan pemuatan—esensial untuk fitur visualisasi step-by-step yang membantu operator gudang.

**Mengapa UUID, Bukan Auto-Increment?**

Anda mungkin bertanya mengapa kita menggunakan UUID sebagai *primary key* daripada *auto-increment integer* yang lebih familiar. Ada beberapa alasan:

1. **Distributed-friendly**: UUID dapat di-*generate* di mana saja tanpa koordinasi dengan database. Ini penting jika di masa depan kita ingin *scale* ke multiple database instances.

2. **Security by obscurity**: Integer sequential (`/plans/1`, `/plans/2`) memudahkan *enumeration attack*—penyerang bisa menebak ID lain. UUID membuat ini praktis mustahil.

3. **Merge-friendly**: Jika kita perlu menggabungkan data dari dua environment (misalnya development dan staging), UUID tidak akan konflik.

Trade-off-nya adalah UUID memakan lebih banyak *storage* (16 bytes vs 4 bytes untuk integer) dan sedikit lebih lambat untuk indexing. Untuk skala aplikasi kita, trade-off ini sangat layak.

#### Migrasi Database dengan Goose

Dalam pengembangan aplikasi, skema database tidak pernah statis. Kita akan menambah tabel baru, mengubah kolom, menambah index—perubahan ini harus dilakukan dengan cara yang terkontrol dan *reproducible*. Inilah mengapa kita membutuhkan *database migrations*.

**Mengapa Tidak Langsung Ubah Database?**

Bayangkan skenario ini: Anda mengubah skema di database development Anda secara langsung menggunakan SQL client. Sebulan kemudian, Anda lupa perubahan apa saja yang sudah dibuat. Saat deploy ke staging, database berbeda dari development. Bug muncul, dan Anda tidak tahu apakah itu karena skema atau kode—*chaos*.

*Versioned migrations* menyelesaikan masalah ini dengan mendokumentasikan setiap perubahan skema dalam file yang tersimpan di version control (Git). Setiap environment (development, staging, production) dapat di-*migrate* ke versi yang sama dengan menjalankan file-file tersebut secara berurutan.

**Mengapa Goose?**

Ada banyak tool migrasi: golang-migrate, sql-migrate, Atlas, dan lainnya. Kita memilih Goose karena:

1. **Simple dan focused**: Tidak ada fitur berlebihan—hanya up, down, dan status
2. **SQL-first**: File migrasi adalah SQL murni, bukan DSL khusus
3. **Embedded migrations**: Dapat di-*embed* ke dalam binary Go untuk deployment yang lebih clean
4. **Komunitas Go**: Dibuat dan di-*maintain* oleh komunitas Go (Pressly)

**Instalasi Goose**

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

**Sequential vs Timestamp Versioning**

Goose mendukung dua format penomoran file:
- **Timestamp**: `20240115143022_create_containers.sql` (waktu pembuatan)
- **Sequential**: `00001_create_containers.sql` (urutan integer)

Kita memilih *sequential* dengan flag `-s` karena lebih mudah dibaca dan dipahami urutannya. Timestamp berguna untuk tim besar di mana dua developer mungkin membuat migrasi bersamaan, tetapi untuk proyek ini, sequential sudah cukup.

**Membuat File Migrasi**

File migrasi ditempatkan di `cmd/db/migrations/`. Gunakan Goose CLI untuk membuat file:

```bash
# Buat direktori migrations
mkdir -p cmd/db/migrations

# Buat migrasi dengan sequential version
goose -dir cmd/db/migrations -s create create_containers sql
goose -dir cmd/db/migrations -s create create_products sql
goose -dir cmd/db/migrations -s create create_plans sql
```

Ini menghasilkan file dengan template kosong. Buka file `00001_create_containers.sql` dan isi dengan:

```sql
-- +goose Up
CREATE TABLE containers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    length_mm DOUBLE PRECISION NOT NULL,
    width_mm DOUBLE PRECISION NOT NULL,
    height_mm DOUBLE PRECISION NOT NULL,
    max_weight_kg DOUBLE PRECISION NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS containers;
```

Perhatikan komentar `-- +goose Up` dan `-- +goose Down`—ini menandai SQL mana yang dijalankan saat *migrate* dan *rollback*.

Isi file `00002_create_products.sql`:

```sql
-- +goose Up
CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    label VARCHAR(255) NOT NULL,
    sku VARCHAR(100) UNIQUE NOT NULL,
    length_mm DOUBLE PRECISION NOT NULL,
    width_mm DOUBLE PRECISION NOT NULL,
    height_mm DOUBLE PRECISION NOT NULL,
    weight_kg DOUBLE PRECISION NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS products;
```

Isi file `00003_create_plans.sql`:

```sql
-- +goose Up
CREATE TABLE plans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    container_id UUID NOT NULL REFERENCES containers(id),
    status VARCHAR(50) NOT NULL DEFAULT 'draft',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    calculated_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE plan_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    plan_id UUID NOT NULL REFERENCES plans(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id),
    quantity INTEGER NOT NULL DEFAULT 1
);

CREATE TABLE placements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    plan_id UUID NOT NULL REFERENCES plans(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id),
    pos_x DOUBLE PRECISION NOT NULL,
    pos_y DOUBLE PRECISION NOT NULL,
    pos_z DOUBLE PRECISION NOT NULL,
    rotation INTEGER NOT NULL DEFAULT 0,
    step_number INTEGER NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS placements;
DROP TABLE IF EXISTS plan_items;
DROP TABLE IF EXISTS plans;
```

**Menjalankan Migrasi**

Pastikan PostgreSQL sudah berjalan, lalu:

```bash
# Set database URL
export DATABASE_URL="postgres://user:password@localhost:5432/loadstuff?sslmode=disable"

# Jalankan migrasi
goose -dir cmd/db/migrations postgres "$DATABASE_URL" up

# Cek status
goose -dir cmd/db/migrations postgres "$DATABASE_URL" status
```

#### Koneksi Database dengan pgx

Go *standard library* menyediakan package `database/sql` sebagai abstraksi universal untuk berbagai database. Anda bisa menggunakan interface yang sama untuk MySQL, PostgreSQL, SQLite, dan lainnya. Namun, abstraksi ini memiliki trade-off: fitur-fitur spesifik database tidak tersedia, dan ada *overhead* dari layer tambahan.

Untuk proyek ini, kita menggunakan **pgx**—driver PostgreSQL native yang ditulis khusus untuk Go. Ini bukan wrapper di atas `database/sql`, melainkan implementasi langsung protokol PostgreSQL.

**Mengapa Memilih pgx daripada database/sql?**

1. **Dukungan tipe data native**: PostgreSQL memiliki tipe data kaya yang tidak ada di SQL standar—UUID, JSONB, arrays, ranges. Dengan `database/sql`, Anda perlu melakukan konversi manual. Dengan pgx, tipe-tipe ini langsung mapped ke tipe Go.

2. **Performance lebih baik**: Tanpa layer abstraksi tambahan, pgx dapat berkomunikasi langsung dengan PostgreSQL menggunakan protokol wire-nya. Benchmark menunjukkan perbedaan signifikan untuk high-throughput applications.

3. **Copy protocol**: PostgreSQL memiliki protokol COPY yang sangat cepat untuk bulk insert. pgx mendukung ini secara native, sementara `database/sql` tidak.

4. **Connection pooling terintegrasi**: Package `pgxpool` menyediakan connection pool yang efisien tanpa perlu library tambahan.

**Connection Pooling: Mengapa Penting?**

Membuka koneksi database adalah operasi mahal—melibatkan TCP handshake, authentication, dan setup session. Jika setiap HTTP request membuka koneksi baru lalu menutupnya, sebagian besar waktu akan dihabiskan untuk *overhead* ini.

*Connection pool* menyelesaikan masalah ini dengan mempertahankan sekumpulan koneksi yang sudah terbuka. Ketika handler butuh koneksi, ia "meminjam" dari pool. Setelah selesai, koneksi dikembalikan ke pool—bukan ditutup. Request berikutnya dapat menggunakan koneksi yang sama.

```mermaid
flowchart LR
    subgraph Pool
        C1[Conn 1]
        C2[Conn 2]
        C3[Conn 3]
    end
    
    R1[Request 1] --> |borrow| C1
    R2[Request 2] --> |borrow| C2
    R1 --> |return| Pool
    R3[Request 3] --> |reuse| C1
```

Dengan pgxpool, pool management sepenuhnya otomatis. Kita cukup configure ukuran pool dan timeout, lalu library mengurus sisanya.

**Setup Koneksi Database**

Sebelum menulis kode koneksi, kita perlu memikirkan bagaimana konfigurasi akan dikelola. Hardcoding `DATABASE_URL` dalam kode adalah ide buruk—nilai ini berbeda di setiap environment (development, staging, production), dan kita tidak ingin mengubah kode setiap kali deploy.

Solusi standar adalah **environment variables**. Konfigurasi dibaca dari environment saat aplikasi start, dengan fallback ke nilai default untuk development. Ini mengikuti [The Twelve-Factor App](https://12factor.net/config) methodology yang menjadi best practice untuk aplikasi modern.

Buat struktur direktori:

```bash
mkdir -p internal/config
```

Buat file `config.go` di dalam direktori `internal/config/`:

```go
package config

import (
    "os"
)

type Config struct {
    DatabaseURL string
    ServerPort  string
}

func Load() *Config {
    return &Config{
        DatabaseURL: getEnv("DATABASE_URL", "postgres://localhost:5432/loadstuff"),
        ServerPort:  getEnv("SERVER_PORT", "8080"),
    }
}

func getEnv(key, fallback string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return fallback
}
```

Perhatikan pattern `getEnv` dengan fallback—ini memungkinkan development berjalan tanpa setup environment variables, sementara production wajib set nilainya.

**Dependency Injection melalui main()**

Koneksi database dibuat di `main()` lalu diteruskan (*inject*) ke komponen yang membutuhkan. Pola ini memungkinkan testability dan membuat dependencies eksplisit—kita akan membahas ini lebih detail di Section 3.6.

Update file `main.go` di dalam direktori `cmd/api/`:

```go
import (
    "context"
    "log"

    "github.com/jackc/pgx/v5/pgxpool"
    "load-stuffing-calculator/internal/config"
)

func main() {
    cfg := config.Load()

    // Buat connection pool
    pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
    if err != nil {
        log.Fatalf("Unable to connect to database: %v", err)
    }
    defer pool.Close()

    // Verifikasi koneksi dengan ping
    if err := pool.Ping(context.Background()); err != nil {
        log.Fatalf("Unable to ping database: %v", err)
    }

    log.Println("Connected to database")

    // Pool ini akan diteruskan ke repository dan service
    // ... setup router dan handlers
}
```

`defer pool.Close()` memastikan koneksi ditutup dengan benar saat aplikasi shutdown. `Ping()` adalah health check cepat untuk memverifikasi database reachable sebelum menerima traffic.

Install dependency:

```bash
go get github.com/jackc/pgx/v5
```

Dengan *database schema* yang sudah terdefinisi dan koneksi yang berjalan, kita siap membangun *repository layer* menggunakan SQLC.

---

### 3.5 Type-safe Repository Layer dengan SQLC

Setelah database schema siap, langkah berikutnya adalah membangun *repository layer*—komponen yang bertanggung jawab untuk semua interaksi dengan database. Di sinilah kita menulis query SQL dan memetakan hasilnya ke struct Go.

Ada beberapa pendekatan untuk membangun *data access layer* di Go:

1. **Raw SQL dengan database/sql**: Menulis query manual, scanning hasil ke struct field-by-field
2. **ORM (GORM, Ent)**: Abstraksi tinggi yang menyembunyikan SQL di balik method calls
3. **Query Builder (squirrel, goqu)**: Membangun query secara programatik dengan Go code
4. **Code Generator (SQLC)**: Menulis SQL murni, generate Go code dari SQL tersebut

Masing-masing memiliki trade-off. ORM mempercepat development tetapi menyembunyikan query yang sebenarnya—saat performance issue muncul, Anda harus menebak query apa yang di-generate. Raw SQL memberikan kontrol penuh tetapi rentan typo dan perubahan schema yang tidak terdeteksi sampai runtime.

#### Mengapa SQLC?

SQLC mengambil pendekatan unik: Anda menulis SQL murni—query yang sebenarnya akan dieksekusi—lalu SQLC men-*generate* kode Go yang type-safe dari query tersebut.

```mermaid
flowchart LR
    A["SQL Queries<br/>(*.sql)"] -->|SQLC| B["Go Code<br/>(generated)"]
```

Mengapa pendekatan ini powerful?

1. **Compile-time type safety**: Jika query return tiga kolom, Go function yang di-generate juga return struct dengan tiga field. Tambah kolom baru? SQLC akan mengupdate struct. Typo nama kolom? Error saat generate, bukan saat runtime.

2. **SQL tetap SQL**: Tidak perlu belajar DSL baru. Jika Anda tahu PostgreSQL, Anda sudah bisa menggunakan SQLC. Complex joins, window functions, CTEs—semuanya bisa ditulis seperti biasa.

3. **No runtime reflection**: ORM biasanya menggunakan reflection untuk mapping struct ke query. SQLC men-generate kode konkret tanpa reflection, memberikan performance yang lebih baik.

4. **Query optimization**: Karena SQL tertulis eksplisit, DBA atau developer berpengalaman bisa langsung mereview dan mengoptimalkan query.

#### Setup SQLC

Install SQLC CLI:

```bash
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

Buat file konfigurasi `sqlc.yaml` di root project:

```yaml
version: "2"
sql:
  - engine: "postgresql"
    queries: "cmd/db/queries/"
    schema: "cmd/db/migrations/"
    gen:
      go:
        package: "store"
        out: "internal/store"
        sql_package: "pgx/v5"
        emit_json_tags: true
        emit_empty_slices: true
        emit_interface: true
        overrides:
          - db_type: "uuid"
            go_type:
              import: "github.com/google/uuid"
              type: "UUID"
          - db_type: "timestamptz"
            go_type:
              import: "time"
              type: "Time"
            nullable: true
```

Mari kita bahas setiap konfigurasi:

- **engine**: Database yang digunakan (postgresql, mysql, sqlite)
- **queries**: Direktori tempat file query SQL berada
- **schema**: Direktori tempat file migrasi berada—SQLC membaca ini untuk memahami struktur tabel
- **gen.go.package**: Nama package untuk kode yang di-generate
- **gen.go.out**: Direktori output
- **sql_package**: Driver database yang digunakan (pgx/v5 untuk pgx)
- **emit_json_tags**: Generate JSON tags pada struct—berguna saat struct ini juga menjadi response API
- **emit_empty_slices**: Return empty slice daripada nil untuk query yang return array kosong
- **emit_interface**: Generate interface `Querier` yang berisi semua method query—berguna untuk dependency injection
- **overrides**: Map tipe database ke tipe Go. Secara default, pgx menggunakan `pgtype.UUID` yang lebih kompleks. Dengan override, kita gunakan `github.com/google/uuid` yang lebih ergonomis


#### Menulis Query Files

Query ditulis dalam file `.sql` dengan annotation khusus yang memberitahu SQLC nama function dan tipe return.

Buat file `containers.sql` di dalam direktori `cmd/db/queries/`:

```sql
-- name: GetContainer :one
SELECT * FROM containers WHERE id = $1;

-- name: ListContainers :many
SELECT * FROM containers ORDER BY created_at DESC;

-- name: CreateContainer :one
INSERT INTO containers (name, length_mm, width_mm, height_mm, max_weight_kg)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: UpdateContainer :one
UPDATE containers
SET name = $2, length_mm = $3, width_mm = $4, height_mm = $5, max_weight_kg = $6
WHERE id = $1
RETURNING *;

-- name: DeleteContainer :exec
DELETE FROM containers WHERE id = $1;
```

Perhatikan annotation di setiap query:

- `-- name: GetContainer :one` — Function bernama `GetContainer`, return satu row
- `:many` — Return slice of rows
- `:exec` — Tidak return data, hanya execute (untuk DELETE, UPDATE tanpa RETURNING)
- `:execrows` — Return jumlah rows yang affected

`$1`, `$2`, dst adalah parameter placeholders. SQLC akan men-generate function arguments dengan tipe yang sesuai.

Buat juga file `products.sql`:

```sql
-- name: GetProduct :one
SELECT * FROM products WHERE id = $1;

-- name: GetProductBySKU :one
SELECT * FROM products WHERE sku = $1;

-- name: ListProducts :many
SELECT * FROM products ORDER BY label ASC;

-- name: CreateProduct :one
INSERT INTO products (label, sku, length_mm, width_mm, height_mm, weight_kg)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: DeleteProduct :exec
DELETE FROM products WHERE id = $1;
```

Dan file `plans.sql` untuk query yang lebih complex:

```sql
-- name: GetPlan :one
SELECT * FROM plans WHERE id = $1;

-- name: CreatePlan :one
INSERT INTO plans (container_id, status)
VALUES ($1, 'draft')
RETURNING *;

-- name: UpdatePlanStatus :one
UPDATE plans SET status = $2, calculated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: AddPlanItem :one
INSERT INTO plan_items (plan_id, product_id, quantity)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetPlanItems :many
SELECT pi.*, p.label, p.length_mm, p.width_mm, p.height_mm, p.weight_kg
FROM plan_items pi
JOIN products p ON pi.product_id = p.id
WHERE pi.plan_id = $1;

-- name: SavePlacement :exec
INSERT INTO placements (plan_id, product_id, pos_x, pos_y, pos_z, rotation, step_number)
VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: GetPlacements :many
SELECT pl.*, p.label
FROM placements pl
JOIN products p ON pl.product_id = p.id
WHERE pl.plan_id = $1
ORDER BY pl.step_number ASC;

-- name: DeletePlanPlacements :exec
DELETE FROM placements WHERE plan_id = $1;
```

Perhatikan query `GetPlanItems` dan `GetPlacements`—keduanya menggunakan JOIN untuk mengambil data dari tabel terkait dalam satu query. SQLC cukup pintar untuk men-generate struct yang mencakup kolom dari kedua tabel.

#### Generate Kode

Setelah query ditulis, generate kode Go:

```bash
sqlc generate
```

SQLC akan membaca schema dari migrations, memahami struktur tabel, lalu men-generate kode ke `internal/store/`. Hasilnya berupa beberapa file:

```
internal/store/
├── db.go           # DBTX interface dan Queries struct
├── models.go       # Struct untuk setiap tabel
├── containers.sql.go
├── products.sql.go
└── plans.sql.go
```

**models.go** berisi struct yang merepresentasikan row dari setiap tabel:

```go
// Code generated by sqlc. DO NOT EDIT.

type Container struct {
    ID          uuid.UUID
    Name        string
    LengthMm    float64
    WidthMm     float64
    HeightMm    float64
    MaxWeightKg float64
    CreatedAt   time.Time
}

type Product struct {
    ID        uuid.UUID
    Label     string
    Sku       string
    LengthMm  float64
    WidthMm   float64
    HeightMm  float64
    WeightKg  float64
    CreatedAt time.Time
}
```

**containers.sql.go** berisi function untuk setiap query:

```go
func (q *Queries) GetContainer(ctx context.Context, id uuid.UUID) (Container, error) {
    // ... implementation
}

func (q *Queries) ListContainers(ctx context.Context) ([]Container, error) {
    // ... implementation
}

func (q *Queries) CreateContainer(ctx context.Context, arg CreateContainerParams) (Container, error) {
    // ... implementation
}
```

Perhatikan `CreateContainerParams`—SQLC secara otomatis men-generate struct untuk query dengan multiple parameters:

```go
type CreateContainerParams struct {
    Name        string
    LengthMm    float64
    WidthMm     float64
    HeightMm    float64
    MaxWeightKg float64
}
```

#### Menggunakan Generated Code

Untuk menggunakan query, kita membuat instance `Queries` dengan connection pool:

```go
package main

import (
    "context"
    "log"

    "github.com/jackc/pgx/v5/pgxpool"
    "load-stuffing-calculator/internal/store"
)

func main() {
    // ... setup pool seperti sebelumnya

    queries := store.New(pool)

    // Create container
    container, err := queries.CreateContainer(ctx, store.CreateContainerParams{
        Name:        "40ft Standard",
        LengthMm:    12192,
        WidthMm:     2438,
        HeightMm:    2591,
        MaxWeightKg: 26780,
    })
    if err != nil {
        log.Fatal(err)
    }

    // Get container by ID
    c, err := queries.GetContainer(ctx, container.ID)
    
    // List all containers
    containers, err := queries.ListContainers(ctx)
}
```

Setiap function menerima `context.Context` sebagai parameter pertama—ini memungkinkan cancellation dan timeout propagation, pattern standar di Go untuk operasi I/O.

Dengan repository layer yang type-safe dan queries yang di-generate, kita siap membangun service layer yang mengorkestrasikan business logic.

---

### 3.6 Business Logic Layer (Service)

Dengan repository layer yang sudah siap, kita sekarang membangun **service layer**—lapisan yang berisi logika bisnis aplikasi. Ini adalah "otak" dari sistem, tempat keputusan dibuat dan berbagai komponen diorkestrasikan.

#### Mengapa Perlu Layer Terpisah?

Pertanyaan umum dari developer yang baru mengenal arsitektur berlapis: "Mengapa tidak langsung panggil repository dari HTTP handler? Bukankah itu menambah kompleksitas tanpa manfaat?"

Argumen ini masuk akal untuk aplikasi sederhana. Tetapi seiring aplikasi berkembang, beberapa masalah muncul:

1. **Duplikasi logika**: Jika dua endpoint perlu melakukan validasi yang sama, kita harus copy-paste kode ke kedua handler.

2. **Testing menjadi sulit**: HTTP handler terikat dengan request/response HTTP. Untuk menguji logika bisnis, kita harus membuat HTTP request palsu—rumit dan lambat.

3. **Coupling yang tinggi**: Handler yang langsung memanggil repository membuat keduanya terikat erat. Mengganti PostgreSQL dengan MySQL memerlukan perubahan di setiap handler.

Service layer menyelesaikan masalah ini dengan menjadi **abstraction layer**:

```mermaid
flowchart TB
    subgraph "HTTP Layer"
        H1[Container Handler]
        H2[Product Handler]
        H3[Plan Handler]
    end
    
    subgraph "Service Layer"
        S1[ContainerService]
        S2[ProductService]
        S3[PlanService]
    end
    
    subgraph "Data Layer"
        R[Repository / Store]
        G[Packing Gateway]
    end
    
    H1 --> S1
    H2 --> S2
    H3 --> S3
    S1 --> R
    S2 --> R
    S3 --> R
    S3 --> G
```

- **Handler** bertanggung jawab untuk: parsing request, validasi format, memanggil service, formatting response
- **Service** bertanggung jawab untuk: validasi bisnis, orkestrasi, transformasi data
- **Repository/Gateway** bertanggung jawab untuk: akses data, komunikasi external

Dengan pemisahan ini, kita bisa mengganti repository tanpa mengubah handler, menguji service tanpa HTTP, dan me-reuse logika di berbagai endpoint.

#### Dependency Injection

Seperti yang disebutkan di Section 3.4, koneksi database diteruskan dari `main()` ke komponen yang membutuhkan. Pola ini disebut **Dependency Injection** (DI)—komponen menerima dependencies-nya melalui constructor, bukan membuat sendiri.

Mengapa ini penting? Mari bandingkan dua pendekatan:

```go
// Buruk: Service membuat dependency sendiri
type PlanService struct {
    store *store.Queries  // Hardcoded concrete type
}

func NewPlanService() *PlanService {
    pool := createDatabaseConnection()  // Sulit di-mock
    return &PlanService{store: store.New(pool)}
}
```

```go
// Baik: Dependency di-inject
type PlanService struct {
    store   Querier         // Interface
    gateway PackingGateway  // Interface
}

func NewPlanService(store Querier, gateway PackingGateway) *PlanService {
    return &PlanService{
        store:   store,
        gateway: gateway,
    }
}
```

Pendekatan kedua memberikan tiga keuntungan utama:

**1. Testability** — Saat testing, kita bisa memberikan mock implementation yang mengembalikan data terprediksi. Tidak perlu database sungguhan, tidak perlu Packing Service berjalan. Test menjadi cepat dan deterministik.

**2. Explicit Dependencies** — Dengan melihat constructor `NewPlanService(store, gateway)`, langsung jelas apa yang dibutuhkan service untuk berfungsi. Tidak ada "magic" atau state tersembunyi.

**3. Lifecycle Control** — `main()` mengontrol kapan dependencies dibuat dan dihancurkan. Database pool dibuat sebelum service, dan ditutup setelah service selesai. Ini mencegah resource leak dan race conditions.

Dengan menerima interface (bukan concrete type), service tidak peduli apakah ia menerima real database atau mock—selama memenuhi kontrak interface. Di Go, ini bekerja karena interface dipenuhi secara implisit: jika sebuah type memiliki semua method yang didefinisikan interface, ia otomatis mengimplementasi interface tersebut.

#### Menggunakan Generated Interface dari SQLC

Kabar baiknya: kita tidak perlu mendefinisikan interface database sendiri. SQLC sudah men-generate interface `Querier` di package `store` yang berisi semua method query kita. Kita bisa langsung menggunakannya:

```go
// Di internal/store/querier.go (generated by SQLC)
type Querier interface {
    GetContainer(ctx context.Context, id uuid.UUID) (Container, error)
    ListContainers(ctx context.Context) ([]Container, error)
    CreateContainer(ctx context.Context, arg CreateContainerParams) (Container, error)
    // ... dan method lainnya
}
```

Interface ini sudah tersedia setelah menjalankan `sqlc generate`. Service layer cukup import dan menggunakannya.

Untuk gateway ke Packing Service, kita akan mendefinisikan interface sendiri. Berikut preview dari interface yang akan kita implementasikan di Section 3.8:

```go
// PackingGateway mendefinisikan kontrak untuk komunikasi dengan Packing Service
type PackingGateway interface {
    Pack(ctx context.Context, req PackRequest) (*PackResponse, error)
}
```

Interface ini memungkinkan service layer memanggil Packing Service tanpa tahu detail implementasinya. Untuk sekarang, service akan menerima interface ini sebagai dependency yang bisa `nil` selama development.

#### ContainerService Implementation

Sekarang kita implementasi service pertama. Buat struktur direktori untuk service layer:

```bash
mkdir -p internal/service
```

Service layer menggunakan UUID untuk ID. Install dependency uuid:

```bash
go get github.com/google/uuid
```

Buat file `container_service.go` di dalam direktori `internal/service/`:

```go
package service

import (
    "context"
    "fmt"

    "github.com/google/uuid"
    "load-stuffing-calculator/internal/store"
)

// ContainerService menangani operasi bisnis terkait Container.
// Menerima store.Querier sebagai dependency untuk akses database.
type ContainerService struct {
    store store.Querier
}

// NewContainerService membuat instance baru ContainerService.
// Pattern ini disebut "constructor function" di Go.
func NewContainerService(store store.Querier) *ContainerService {
    return &ContainerService{store: store}
}

// GetByID mengambil satu container berdasarkan ID.
// Mengembalikan pointer ke Container atau error jika tidak ditemukan.
func (s *ContainerService) GetByID(ctx context.Context, id uuid.UUID) (*store.Container, error) {
    container, err := s.store.GetContainer(ctx, id)
    if err != nil {
        // Wrap error dengan konteks untuk debugging yang lebih mudah
        return nil, fmt.Errorf("get container %s: %w", id, err)
    }
    return &container, nil
}

// List mengambil semua container yang tersedia.
func (s *ContainerService) List(ctx context.Context) ([]store.Container, error) {
    containers, err := s.store.ListContainers(ctx)
    if err != nil {
        return nil, fmt.Errorf("list containers: %w", err)
    }
    return containers, nil
}

// Create membuat container baru dengan validasi bisnis.
// Validasi memastikan dimensi dan berat bernilai positif.
func (s *ContainerService) Create(ctx context.Context, name string, length, width, height, maxWeight float64) (*store.Container, error) {
    // Validasi bisnis: dimensi harus positif
    // Ini adalah aturan domain, bukan constraint database
    if length <= 0 || width <= 0 || height <= 0 {
        return nil, fmt.Errorf("dimensions must be positive")
    }
    if maxWeight <= 0 {
        return nil, fmt.Errorf("max weight must be positive")
    }

    // Panggil repository untuk menyimpan ke database
    container, err := s.store.CreateContainer(ctx, store.CreateContainerParams{
        Name:        name,
        LengthMm:    length,
        WidthMm:     width,
        HeightMm:    height,
        MaxWeightKg: maxWeight,
    })
    if err != nil {
        return nil, fmt.Errorf("create container: %w", err)
    }
    return &container, nil
}
```

Mari kita bahas beberapa pola penting dalam kode di atas:

**Receiver Methods**

`func (s *ContainerService) GetByID(...)` menggunakan *pointer receiver* (`*ContainerService`). Ini berarti method dapat memodifikasi state dari struct jika diperlukan, dan menghindari copy struct setiap kali method dipanggil.

**Context sebagai Parameter Pertama**

Setiap method public menerima `context.Context` sebagai parameter pertama. Ini adalah konvensi Go untuk functions yang melakukan I/O. Context memungkinkan:
- **Cancellation**: Request dibatalkan? Query bisa berhenti segera
- **Timeout**: Batas waktu untuk operasi
- **Values**: Meneruskan request ID untuk tracing

**Error Wrapping**

Perhatikan pattern `fmt.Errorf("create container: %w", err)`. Verb `%w` membungkus error original sehingga caller dapat menggunakan `errors.Is()` atau `errors.As()` untuk memeriksa error spesifik. String "create container:" menambah konteks tentang *di mana* error terjadi.

#### Melengkapi CRUD dengan Update dan Delete

Untuk menyelesaikan operasi CRUD, kita perlu menambahkan method `Update` dan `Delete`. Pattern yang sama diterapkan—validasi bisnis di service, operasi database di repository:

```go
// Update memperbarui container yang sudah ada.
func (s *ContainerService) Update(ctx context.Context, id uuid.UUID, name string, length, width, height, maxWeight float64) (*store.Container, error) {
    // Validasi bisnis: dimensi harus positif
    if length <= 0 || width <= 0 || height <= 0 {
        return nil, fmt.Errorf("dimensions must be positive")
    }
    if maxWeight <= 0 {
        return nil, fmt.Errorf("max weight must be positive")
    }

    container, err := s.store.UpdateContainer(ctx, store.UpdateContainerParams{
        ID:          id,
        Name:        name,
        LengthMm:    length,
        WidthMm:     width,
        HeightMm:    height,
        MaxWeightKg: maxWeight,
    })
    if err != nil {
        return nil, fmt.Errorf("update container %s: %w", id, err)
    }
    return &container, nil
}

// Delete menghapus container berdasarkan ID.
func (s *ContainerService) Delete(ctx context.Context, id uuid.UUID) error {
    err := s.store.DeleteContainer(ctx, id)
    if err != nil {
        return fmt.Errorf("delete container %s: %w", id, err)
    }
    return nil
}
```

Method `Update` memiliki validasi yang sama dengan `Create`—ini memastikan konsistensi aturan bisnis terlepas dari bagaimana data masuk ke sistem. Method `Delete` lebih sederhana karena tidak ada nilai kembalian yang perlu di-wrap.

#### ProductService: Menerapkan Pattern yang Sama

Dengan pemahaman yang solid tentang bagaimana `ContainerService` dibangun, kita bisa menerapkan pattern yang sama untuk `ProductService`. Kedua entitas ini memiliki karakteristik serupa—keduanya adalah *master data* dengan operasi CRUD standar.

Buat file `product_service.go` di dalam direktori `internal/service/`:

```go
package service

import (
    "context"
    "fmt"

    "github.com/google/uuid"
    "load-stuffing-calculator/internal/store"
)

// ProductService menangani operasi bisnis terkait Product.
type ProductService struct {
    store store.Querier
}

func NewProductService(store store.Querier) *ProductService {
    return &ProductService{store: store}
}

func (s *ProductService) GetByID(ctx context.Context, id uuid.UUID) (*store.Product, error) {
    product, err := s.store.GetProduct(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("get product %s: %w", id, err)
    }
    return &product, nil
}

func (s *ProductService) List(ctx context.Context) ([]store.Product, error) {
    products, err := s.store.ListProducts(ctx)
    if err != nil {
        return nil, fmt.Errorf("list products: %w", err)
    }
    return products, nil
}

func (s *ProductService) Create(ctx context.Context, label, sku string, length, width, height, weight float64) (*store.Product, error) {
    if length <= 0 || width <= 0 || height <= 0 {
        return nil, fmt.Errorf("dimensions must be positive")
    }
    if weight <= 0 {
        return nil, fmt.Errorf("weight must be positive")
    }

    product, err := s.store.CreateProduct(ctx, store.CreateProductParams{
        Label:    label,
        Sku:      sku,
        LengthMm: length,
        WidthMm:  width,
        HeightMm: height,
        WeightKg: weight,
    })
    if err != nil {
        return nil, fmt.Errorf("create product: %w", err)
    }
    return &product, nil
}

func (s *ProductService) Update(ctx context.Context, id uuid.UUID, label, sku string, length, width, height, weight float64) (*store.Product, error) {
    if length <= 0 || width <= 0 || height <= 0 {
        return nil, fmt.Errorf("dimensions must be positive")
    }
    if weight <= 0 {
        return nil, fmt.Errorf("weight must be positive")
    }

    product, err := s.store.UpdateProduct(ctx, store.UpdateProductParams{
        ID:       id,
        Label:    label,
        Sku:      sku,
        LengthMm: length,
        WidthMm:  width,
        HeightMm: height,
        WeightKg: weight,
    })
    if err != nil {
        return nil, fmt.Errorf("update product %s: %w", id, err)
    }
    return &product, nil
}

func (s *ProductService) Delete(ctx context.Context, id uuid.UUID) error {
    err := s.store.DeleteProduct(ctx, id)
    if err != nil {
        return fmt.Errorf("delete product %s: %w", id, err)
    }
    return nil
}
```

Perhatikan bagaimana struktur `ProductService` hampir identik dengan `ContainerService`. Ini bukan kebetulan—konsistensi dalam pattern membuat codebase lebih mudah dipahami. Developer baru yang sudah memahami satu service bisa langsung memahami yang lain.

Perbedaan utama ada pada field yang divalidasi dan parameter yang diterima. Product memiliki `Label` dan `SKU` (bukan hanya `Name`), dan menggunakan `WeightKg` untuk berat produk individual (bukan `MaxWeightKg` untuk kapasitas container).

**Thin Layer?**

Untuk CRUD sederhana, service mungkin terlihat seperti pass-through ke repository. Nilainya akan terlihat saat:

- Ada validasi bisnis yang tidak bisa dilakukan di database (e.g., "container harus exist sebelum plan dibuat")
- Perlu mengorkestrasikan multiple repository calls dalam satu operasi
- Perlu berkomunikasi dengan external services
- Logika yang sama dibutuhkan di multiple handlers

Contoh berikutnya akan menunjukkan nilai sebenarnya dari service layer.

#### Implementasi PlanService

Berbeda dengan `ContainerService` yang sederhana, `PlanService` adalah contoh nyata mengapa service layer diperlukan. Service ini mengorkestrasikan multiple operations: mengambil data dari berbagai tabel, memanggil external service, dan menyimpan hasil. Logika ini terlalu kompleks untuk diletakkan di handler.

Buat file `plan_service.go` di dalam direktori `internal/service/`:

```go
package service

import (
    "context"
    "fmt"

    "github.com/google/uuid"
    "load-stuffing-calculator/internal/gateway"
    "load-stuffing-calculator/internal/store"
)

// PlanService mengelola lifecycle dari Plan:
// - Pembuatan plan baru
// - Penambahan items ke plan
// - Orkestrasi kalkulasi dengan Packing Service
type PlanService struct {
    store   store.Querier      // Dependency: akses database via generated interface
    gateway gateway.PackingGateway  // Dependency: komunikasi dengan Packing Service
}

// NewPlanService membuat instance PlanService dengan dependencies yang di-inject.
// Ini memungkinkan substitusi dengan mock saat testing.
func NewPlanService(store store.Querier, gw gateway.PackingGateway) *PlanService {
    return &PlanService{
        store:   store,
        gateway: gw,
    }
}

func (s *PlanService) Create(ctx context.Context, containerID uuid.UUID) (*store.Plan, error) {
    // Validasi: container harus exist
    _, err := s.store.GetContainer(ctx, containerID)
    if err != nil {
        return nil, fmt.Errorf("container not found: %w", err)
    }

    plan, err := s.store.CreatePlan(ctx, containerID)
    if err != nil {
        return nil, fmt.Errorf("create plan: %w", err)
    }
    return &plan, nil
}

func (s *PlanService) AddItem(ctx context.Context, planID, productID uuid.UUID, quantity int) error {
    // Validasi
    if quantity <= 0 {
        return fmt.Errorf("quantity must be positive")
    }

    _, err := s.store.GetProduct(ctx, productID)
    if err != nil {
        return fmt.Errorf("product not found: %w", err)
    }

    _, err = s.store.AddPlanItem(ctx, store.AddPlanItemParams{
        PlanID:    planID,
        ProductID: productID,
        Quantity:  int32(quantity),
    })
    return err
}
```

Kode di atas menunjukkan dua method sederhana dari `PlanService`. Method `Create` pertama-tama memverifikasi bahwa container yang dipilih benar-benar ada di database—ini adalah contoh validasi yang tidak bisa dilakukan di level repository. Method `AddItem` memvalidasi quantity dan keberadaan product sebelum menyimpan ke database.

#### Melengkapi PlanService dengan CRUD

Seperti halnya Container dan Product, Plan juga membutuhkan operasi CRUD lengkap. Berikut methods tambahan untuk `PlanService`:

```go
// GetByID mengambil satu plan berdasarkan ID.
func (s *PlanService) GetByID(ctx context.Context, id uuid.UUID) (*store.Plan, error) {
    plan, err := s.store.GetPlan(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("get plan %s: %w", id, err)
    }
    return &plan, nil
}

// List mengambil semua plan dengan informasi container.
func (s *PlanService) List(ctx context.Context) ([]store.ListPlansRow, error) {
    plans, err := s.store.ListPlans(ctx)
    if err != nil {
        return nil, fmt.Errorf("list plans: %w", err)
    }
    return plans, nil
}

// Update memperbarui plan (hanya container yang bisa diubah selama status draft).
func (s *PlanService) Update(ctx context.Context, id uuid.UUID, containerID uuid.UUID) (*store.Plan, error) {
    // Validasi: container harus exist
    _, err := s.store.GetContainer(ctx, containerID)
    if err != nil {
        return nil, fmt.Errorf("container not found: %w", err)
    }

    plan, err := s.store.UpdatePlan(ctx, store.UpdatePlanParams{
        ID:          id,
        ContainerID: containerID,
    })
    if err != nil {
        return nil, fmt.Errorf("update plan %s: %w", id, err)
    }
    return &plan, nil
}

// Delete menghapus plan beserta semua items dan placements (cascade).
func (s *PlanService) Delete(ctx context.Context, id uuid.UUID) error {
    err := s.store.DeletePlan(ctx, id)
    if err != nil {
        return fmt.Errorf("delete plan %s: %w", id, err)
    }
    return nil
}
```

Perhatikan bahwa method `List` mengembalikan `store.ListPlansRow`—tipe ini di-generate oleh SQLC berdasarkan query JOIN yang kita tulis. Query tersebut mengambil nama container bersama dengan data plan, sehingga frontend bisa menampilkan informasi yang lebih informatif.

#### Manajemen Plan Items

Plan items juga membutuhkan operasi update dan delete. Berikut methods untuk mengelola items dalam plan:

```go
// GetItems mengambil semua items dalam plan beserta detail produk.
func (s *PlanService) GetItems(ctx context.Context, planID uuid.UUID) ([]store.GetPlanItemsRow, error) {
    items, err := s.store.GetPlanItems(ctx, planID)
    if err != nil {
        return nil, fmt.Errorf("get plan items: %w", err)
    }
    return items, nil
}

// UpdateItem memperbarui quantity dari item yang sudah ada.
func (s *PlanService) UpdateItem(ctx context.Context, itemID uuid.UUID, quantity int) (*store.PlanItem, error) {
    if quantity <= 0 {
        return nil, fmt.Errorf("quantity must be positive")
    }

    item, err := s.store.UpdatePlanItem(ctx, store.UpdatePlanItemParams{
        ID:       itemID,
        Quantity: int32(quantity),
    })
    if err != nil {
        return nil, fmt.Errorf("update plan item: %w", err)
    }
    return &item, nil
}

// DeleteItem menghapus item dari plan.
func (s *PlanService) DeleteItem(ctx context.Context, itemID uuid.UUID) error {
    err := s.store.DeletePlanItem(ctx, itemID)
    if err != nil {
        return fmt.Errorf("delete plan item: %w", err)
    }
    return nil
}

// GetPlacements mengambil semua placements hasil kalkulasi.
func (s *PlanService) GetPlacements(ctx context.Context, planID uuid.UUID) ([]store.GetPlacementsRow, error) {
    placements, err := s.store.GetPlacements(ctx, planID)
    if err != nil {
        return nil, fmt.Errorf("get placements: %w", err)
    }
    return placements, nil
}
```

Method `GetPlacements` akan digunakan oleh handler untuk menyertakan hasil kalkulasi saat mengambil detail plan. Ini memungkinkan frontend menampilkan visualisasi 3D tanpa request tambahan.

#### Method Calculate

Method `Calculate` adalah method yang paling kompleks dalam sistem kita. Di sinilah nilai sebenarnya dari service layer terlihat—kita mengorkestrasikan tujuh langkah berbeda menjadi satu operasi atomis dari perspektif pengguna.

```go
func (s *PlanService) Calculate(ctx context.Context, planID uuid.UUID) (*CalculateResult, error) {
    // 1. Ambil plan dan container
    plan, err := s.store.GetPlan(ctx, planID)
    if err != nil {
        return nil, fmt.Errorf("get plan: %w", err)
    }

    container, err := s.store.GetContainer(ctx, plan.ContainerID)
    if err != nil {
        return nil, fmt.Errorf("get container: %w", err)
    }

    // 2. Ambil items dengan detail produk
    items, err := s.store.GetPlanItems(ctx, planID)
    if err != nil {
        return nil, fmt.Errorf("get plan items: %w", err)
    }

    if len(items) == 0 {
        return nil, fmt.Errorf("plan has no items")
    }

    // 3. Update status ke "calculating"
    _, err = s.store.UpdatePlanStatus(ctx, store.UpdatePlanStatusParams{
        ID:     planID,
        Status: "calculating",
    })
    if err != nil {
        return nil, fmt.Errorf("update status: %w", err)
    }

    // 4. Build request untuk Packing Service
    req := s.buildPackRequest(container, items)

    // 5. Panggil Packing Service
    resp, err := s.gateway.Pack(ctx, req)
    if err != nil {
        // Rollback status ke draft jika gagal
        s.store.UpdatePlanStatus(ctx, store.UpdatePlanStatusParams{
            ID:     planID,
            Status: "failed",
        })
        return nil, fmt.Errorf("packing calculation failed: %w", err)
    }

    // 6. Simpan placements
    err = s.savePlacements(ctx, planID, resp.Placements)
    if err != nil {
        return nil, fmt.Errorf("save placements: %w", err)
    }

    // 7. Update status ke "completed"
    _, err = s.store.UpdatePlanStatus(ctx, store.UpdatePlanStatusParams{
        ID:     planID,
        Status: "completed",
    })

    return &CalculateResult{
        Placements:    resp.Placements,
        UnfittedItems: resp.UnfittedItems,
        Statistics:    resp.Statistics,
    }, nil
}
```

> **Catatan:** Untuk menjaga fokus pada alur utama, beberapa helper telah disederhanakan. Tipe `CalculateResult` adalah struct sederhana yang berisi hasil kalkulasi (placements, unfitted items, statistics). Method `savePlacements` menyimpan setiap placement ke database. Implementasi lengkap dapat dilihat di repository proyek.

Method ini memulai dengan mengumpulkan semua data yang diperlukan dari database: plan itu sendiri, container yang terkait, dan semua items dalam plan. Setiap query di-handle dengan error checking—jika satu gagal, method berhenti dan mengembalikan error yang informatif.

Sebelum memulai kalkulasi yang mungkin memakan waktu, kita update status plan menjadi "calculating". Ini memberikan feedback ke pengguna bahwa proses sedang berjalan. Jika ada request duplikat, sistem bisa mendeteksi bahwa kalkulasi sudah berjalan.

Packing Service memiliki format request sendiri yang berbeda dari model database kita. Method helper `buildPackRequest` mengkonversi data dari format internal ke format yang diharapkan oleh Packing Service—ini adalah contoh mengapa service layer berguna, karena transformasi data ini tidak milik repository maupun handler.

Saat memanggil Packing Service, perhatikan error handling: jika kalkulasi gagal, kita update status ke "failed" sebelum return error. Ini memastikan state tetap konsisten meski terjadi kegagalan. Setelah kalkulasi berhasil, hasil placements disimpan ke database dan status di-update ke "completed".

Method `buildPackRequest` melakukan transformasi data:

```go
func (s *PlanService) buildPackRequest(container store.Container, items []store.GetPlanItemsRow) PackRequest {
    packItems := make([]PackItem, 0, len(items))
    for _, item := range items {
        packItems = append(packItems, PackItem{
            ID:       item.ProductID.String(),
            Label:    item.Label,
            Length:   item.LengthMm,
            Width:    item.WidthMm,
            Height:   item.HeightMm,
            Weight:   item.WeightKg,
            Quantity: int(item.Quantity),
        })
    }

    return PackRequest{
        Container: PackContainer{
            Length:    container.LengthMm,
            Width:     container.WidthMm,
            Height:    container.HeightMm,
            MaxWeight: container.MaxWeightKg,
        },
        Items: packItems,
    }
}
```

Method ini mengiterasi setiap item dalam plan dan mengkonversi dari struct database (`store.GetPlanItemsRow`) ke struct yang diharapkan Packing Service (`PackItem`). Container juga dikonversi ke format `PackContainer`. Dengan memisahkan logika transformasi ke method sendiri, kode `Calculate` tetap bersih dan fokus pada alur utama.

Dengan service layer yang solid, kita siap membangun HTTP handlers yang akan mengekspos logika ini ke dunia luar.

---

### 3.7 HTTP Router dan Handler

Setelah service layer siap, kita perlu cara untuk menerima request dari dunia luar dan mengembalikan response. Inilah tugas HTTP handler—menjadi "pintu masuk" ke aplikasi kita.

#### Memilih Web Framework

Go standard library sudah menyediakan package `net/http` yang cukup powerful untuk membangun HTTP server. Namun untuk aplikasi dengan banyak endpoint, kita membutuhkan fitur tambahan seperti routing dengan parameter, middleware, dan JSON binding. Di sinilah web framework membantu.

Kita memilih **Gin** karena beberapa alasan. Pertama, Gin memiliki API yang sederhana dan mudah dipelajari—tidak ada "magic" yang menyembunyikan cara kerja HTTP. Kedua, Gin adalah salah satu framework tercepat di Go berkat penggunaan httprouter yang efisien. Ketiga, ekosistemnya matang dengan banyak middleware siap pakai.

#### Standarisasi Response dengan Response Package

Sebelum menulis handlers, mari buat package untuk standarisasi format response API. Dengan ini, semua endpoint akan mengembalikan struktur JSON yang konsisten—memudahkan frontend untuk parsing response.

Buat struktur direktori:

```bash
mkdir -p internal/response
mkdir -p internal/dto
```

Buat file `responder.go` di dalam direktori `internal/response/`:

```go
package response

import (
    "github.com/gin-gonic/gin"
)

// ErrorDetail mendefinisikan struktur untuk error individual
type ErrorDetail struct {
    Code    string      `json:"code,omitempty"`
    Message string      `json:"message"`
    Source  interface{} `json:"source,omitempty"`
}

// APIResponse adalah struktur standar untuk semua API response
type APIResponse struct {
    Success bool          `json:"success"`
    Data    interface{}   `json:"data,omitempty"`
    Errors  []ErrorDetail `json:"errors,omitempty"`
}

// Success mengirim response sukses dengan data
func Success(c *gin.Context, status int, data interface{}) {
    c.JSON(status, APIResponse{
        Success: true,
        Data:    data,
    })
}

// Error mengirim response error dengan satu pesan
func Error(c *gin.Context, status int, message string) {
    c.JSON(status, APIResponse{
        Success: false,
        Errors: []ErrorDetail{
            {Message: message},
        },
    })
}
```

Dengan helper functions ini, semua response API akan memiliki struktur yang sama. Response sukses selalu berisi `"success": true` dan data di field `data`. Response error selalu berisi `"success": false` dan detail error di field `errors`. Konsistensi ini sangat membantu frontend dalam handling response.

#### Data Transfer Objects (DTO)

DTO adalah struct yang mendefinisikan format data untuk request dan response. Dengan memisahkan DTO dari model database, kita bisa mengontrol data apa yang diterima dari client dan apa yang dikirim ke client—tanpa mengekspos struktur internal database.

Buat file `container_dto.go` di `internal/dto/`:

```go
package dto

// CreateContainerRequest mendefinisikan payload untuk pembuatan container
type CreateContainerRequest struct {
    Name        string  `json:"name" binding:"required,min=2,max=100"`
    LengthMm    float64 `json:"length_mm" binding:"required,gt=0"`
    WidthMm     float64 `json:"width_mm" binding:"required,gt=0"`
    HeightMm    float64 `json:"height_mm" binding:"required,gt=0"`
    MaxWeightKg float64 `json:"max_weight_kg" binding:"required,gt=0"`
}

// UpdateContainerRequest mendefinisikan payload untuk update container
type UpdateContainerRequest struct {
    Name        string  `json:"name" binding:"required,min=2,max=100"`
    LengthMm    float64 `json:"length_mm" binding:"required,gt=0"`
    WidthMm     float64 `json:"width_mm" binding:"required,gt=0"`
    HeightMm    float64 `json:"height_mm" binding:"required,gt=0"`
    MaxWeightKg float64 `json:"max_weight_kg" binding:"required,gt=0"`
}

// ContainerResponse mendefinisikan format JSON yang dikembalikan ke client
type ContainerResponse struct {
    ID          string  `json:"id"`
    Name        string  `json:"name"`
    LengthMm    float64 `json:"length_mm"`
    WidthMm     float64 `json:"width_mm"`
    HeightMm    float64 `json:"height_mm"`
    MaxWeightKg float64 `json:"max_weight_kg"`
}
```

Tag `binding` digunakan oleh Gin validator. `required` berarti field harus ada, `gt=0` berarti nilai harus lebih besar dari 0. Ini memberikan validasi input dasar secara otomatis.

Selanjutnya, buat `product_dto.go` di folder yang sama untuk definisi produk:

```go
package dto

type CreateProductRequest struct {
    Label    string  `json:"label" binding:"required,min=2,max=100"`
    SKU      string  `json:"sku" binding:"required,min=2,max=50"`
    LengthMm float64 `json:"length_mm" binding:"required,gt=0"`
    WidthMm  float64 `json:"width_mm" binding:"required,gt=0"`
    HeightMm float64 `json:"height_mm" binding:"required,gt=0"`
    WeightKg float64 `json:"weight_kg" binding:"required,gt=0"`
}

type UpdateProductRequest struct {
    Label    string  `json:"label" binding:"required,min=2,max=100"`
    SKU      string  `json:"sku" binding:"required,min=2,max=50"`
    LengthMm float64 `json:"length_mm" binding:"required,gt=0"`
    WidthMm  float64 `json:"width_mm" binding:"required,gt=0"`
    HeightMm float64 `json:"height_mm" binding:"required,gt=0"`
    WeightKg float64 `json:"weight_kg" binding:"required,gt=0"`
}

type ProductResponse struct {
    ID       string  `json:"id"`
    Label    string  `json:"label"`
    SKU      string  `json:"sku"`
    LengthMm float64 `json:"length_mm"`
    WidthMm  float64 `json:"width_mm"`
    HeightMm float64 `json:"height_mm"`
    WeightKg float64 `json:"weight_kg"`
}
```

Terakhir, buat `plan_dto.go` untuk menangani struktur data yang lebih kompleks pada perencanaan muatan:

```go
package dto

type CreatePlanRequest struct {
    ContainerID string `json:"container_id" binding:"required,uuid"`
}

type UpdatePlanRequest struct {
    ContainerID string `json:"container_id" binding:"required,uuid"`
}

type AddPlanItemRequest struct {
    ProductID string `json:"product_id" binding:"required,uuid"`
    Quantity  int    `json:"quantity" binding:"required,gt=0"`
}

type UpdatePlanItemRequest struct {
    Quantity int `json:"quantity" binding:"required,gt=0"`
}

type PlanResponse struct {
    ID            string `json:"id"`
    ContainerID   string `json:"container_id"`
    ContainerName string `json:"container_name,omitempty"`
    Status        string `json:"status"`
}

type PlanItemResponse struct {
    ID        string  `json:"id"`
    ProductID string  `json:"product_id"`
    Label     string  `json:"label"`
    Quantity  int     `json:"quantity"`
    LengthMm  float64 `json:"length_mm"`
    WidthMm   float64 `json:"width_mm"`
    HeightMm  float64 `json:"height_mm"`
    WeightKg  float64 `json:"weight_kg"`
}

type PlacementResponse struct {
    ID         string  `json:"id"`
    ProductID  string  `json:"product_id"`
    Label      string  `json:"label"`
    PosX       float64 `json:"pos_x"`
    PosY       float64 `json:"pos_y"`
    PosZ       float64 `json:"pos_z"`
    Rotation   int     `json:"rotation"`
    StepNumber int     `json:"step_number"`
}

type PlanDetailResponse struct {
    ID            string              `json:"id"`
    ContainerID   string              `json:"container_id"`
    ContainerName string              `json:"container_name"`
    Status        string              `json:"status"`
    Items         []PlanItemResponse  `json:"items"`
    Placements    []PlacementResponse `json:"placements,omitempty"`
}
```

Dengan DTO yang lengkap, kita siap mengimplementasikan HTTP handler yang memanfaatkan model-model ini.

```

Perhatikan bahwa `ContainerResponse` menggunakan `string` untuk ID, bukan `uuid.UUID`. Ini karena JSON tidak memiliki tipe UUID—konversi dilakukan di handler atau service layer. DTO juga memungkinkan kita untuk tidak mengekspos field internal seperti `created_at` jika tidak diperlukan oleh frontend.

Buat juga file `plan_dto.go` di dalam direktori `internal/dto/`:

```go
package dto

type CreatePlanRequest struct {
    ContainerID string `json:"container_id" binding:"required,uuid"`
}

type AddPlanItemRequest struct {
    ProductID string `json:"product_id" binding:"required,uuid"`
    Quantity  int    `json:"quantity" binding:"required,gt=0"`
}

type PlanResponse struct {
    ID          string `json:"id"`
    ContainerID string `json:"container_id"`
    Status      string `json:"status"`
}
```

#### Implementasi Handler

Sekarang kita bisa mengimplementasikan handler dengan menggunakan response package dan DTO. Buat struktur direktori untuk handler:

```bash
mkdir -p internal/handler
```

Buat file `container_handler.go` di dalam direktori `internal/handler/`:

```go
package handler

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "load-stuffing-calculator/internal/dto"
    "load-stuffing-calculator/internal/response"
    "load-stuffing-calculator/internal/service"
)

type ContainerHandler struct {
    svc *service.ContainerService
}

func NewContainerHandler(svc *service.ContainerService) *ContainerHandler {
    return &ContainerHandler{svc: svc}
}

func (h *ContainerHandler) Create(c *gin.Context) {
    var req dto.CreateContainerRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.Error(c, http.StatusBadRequest, "Invalid request: "+err.Error())
        return
    }

    container, err := h.svc.Create(
        c.Request.Context(),
        req.Name,
        req.LengthMm,
        req.WidthMm,
        req.HeightMm,
        req.MaxWeightKg,
    )
    if err != nil {
        response.Error(c, http.StatusInternalServerError, err.Error())
        return
    }

    // Convert ke DTO response
    resp := dto.ContainerResponse{
        ID:          container.ID.String(),
        Name:        container.Name,
        LengthMm:    container.LengthMm,
        WidthMm:     container.WidthMm,
        HeightMm:    container.HeightMm,
        MaxWeightKg: container.MaxWeightKg,
    }

    response.Success(c, http.StatusCreated, resp)
}

func (h *ContainerHandler) GetByID(c *gin.Context) {
    idStr := c.Param("id")
    if idStr == "" {
        response.Error(c, http.StatusBadRequest, "Container ID is required")
        return
    }

    id, err := uuid.Parse(idStr)
    if err != nil {
        response.Error(c, http.StatusBadRequest, "Invalid container ID format")
        return
    }

    container, err := h.svc.GetByID(c.Request.Context(), id)
    if err != nil {
        response.Error(c, http.StatusNotFound, "Container not found")
        return
    }

    resp := dto.ContainerResponse{
        ID:          container.ID.String(),
        Name:        container.Name,
        LengthMm:    container.LengthMm,
        WidthMm:     container.WidthMm,
        HeightMm:    container.HeightMm,
        MaxWeightKg: container.MaxWeightKg,
    }

    response.Success(c, http.StatusOK, resp)
}

func (h *ContainerHandler) List(c *gin.Context) {
    containers, err := h.svc.List(c.Request.Context())
    if err != nil {
        response.Error(c, http.StatusInternalServerError, "Failed to list containers")
        return
    }

    // Convert slice ke DTO responses
    resp := make([]dto.ContainerResponse, 0, len(containers))
    for _, container := range containers {
        resp = append(resp, dto.ContainerResponse{
            ID:          container.ID.String(),
            Name:        container.Name,
            LengthMm:    container.LengthMm,
            WidthMm:     container.WidthMm,
            HeightMm:    container.HeightMm,
            MaxWeightKg: container.MaxWeightKg,
        })
    }

    response.Success(c, http.StatusOK, resp)
}
```

Handler menggunakan `response.Success` dan `response.Error` untuk memastikan format response yang konsisten. Setiap handler juga mengkonversi dari model internal ke DTO response—ini memberikan kontrol penuh atas data yang dikirim ke client.

Pattern yang sama diterapkan untuk `PlanHandler`. Buat file `plan_handler.go` di dalam direktori `internal/handler/`:

```go
package handler

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "load-stuffing-calculator/internal/dto"
    "load-stuffing-calculator/internal/response"
    "load-stuffing-calculator/internal/service"
)

type PlanHandler struct {
    svc *service.PlanService
}

func NewPlanHandler(svc *service.PlanService) *PlanHandler {
    return &PlanHandler{svc: svc}
}

func (h *PlanHandler) Create(c *gin.Context) {
    var req dto.CreatePlanRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.Error(c, http.StatusBadRequest, err.Error())
        return
    }

    containerID, _ := uuid.Parse(req.ContainerID)
    
    plan, err := h.svc.Create(c.Request.Context(), containerID)
    if err != nil {
        response.Error(c, http.StatusBadRequest, err.Error())
        return
    }

    resp := dto.PlanResponse{
        ID:          plan.ID.String(),
        ContainerID: plan.ContainerID.String(),
        Status:      plan.Status,
    }

    response.Success(c, http.StatusCreated, resp)
}

func (h *PlanHandler) Calculate(c *gin.Context) {
    idStr := c.Param("id")
    id, err := uuid.Parse(idStr)
    if err != nil {
        response.Error(c, http.StatusBadRequest, "Invalid plan ID")
        return
    }

    result, err := h.svc.Calculate(c.Request.Context(), id)
    if err != nil {
        response.Error(c, http.StatusBadRequest, err.Error())
        return
    }

    response.Success(c, http.StatusOK, result)
}
```

#### API Package dan Router Setup

Untuk mengorganisir initialization dan routing, kita buat package `api` yang menampung struktur aplikasi. Package ini berfungsi sebagai "composition root"—satu tempat di mana semua komponen aplikasi dibuat dan dihubungkan.

Buat direktori:

```bash
mkdir -p internal/api
```

Sebelum menulis kode, install dependency CORS yang diperlukan:

```bash
go get github.com/gin-contrib/cors
```

Buat file `api.go` di dalam direktori `internal/api/`:

```go
package api

import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
    "github.com/jackc/pgx/v5/pgxpool"
    "load-stuffing-calculator/internal/config"
    "load-stuffing-calculator/internal/handler"
    "load-stuffing-calculator/internal/service"
    "load-stuffing-calculator/internal/store"
)

// App menampung semua dependencies aplikasi.
// Dengan menyimpan semua komponen di satu struct, kita bisa dengan mudah
// melihat apa saja yang dibutuhkan aplikasi untuk berjalan.
type App struct {
    config           *config.Config
    router           *gin.Engine
    db               *pgxpool.Pool
    containerHandler *handler.ContainerHandler
    planHandler      *handler.PlanHandler
}

// NewApp menginisialisasi semua layers dan dependencies.
// Urutan inisialisasi mengikuti dependency graph:
// database → store → services → handlers → router
func NewApp(cfg *config.Config, db *pgxpool.Pool) *App {
    // Store layer: wrapper untuk database queries
    querier := store.New(db)

    // Service layer: business logic
    containerSvc := service.NewContainerService(querier)
    planSvc := service.NewPlanService(querier, nil) // gateway ditambahkan di section berikutnya

    // Handler layer: HTTP request/response handling
    containerHandler := handler.NewContainerHandler(containerSvc)
    planHandler := handler.NewPlanHandler(planSvc)

    app := &App{
        config:           cfg,
        db:               db,
        containerHandler: containerHandler,
        planHandler:      planHandler,
    }

    app.setupRouter()

    return app
}
```

Struct `App` menyimpan semua dependencies yang dibutuhkan aplikasi. Pendekatan ini memiliki beberapa keuntungan. Pertama, kita bisa melihat dengan jelas apa saja komponen yang dibutuhkan hanya dengan membaca struct definition. Kedua, semua initialization logic terpusat di `NewApp()`—tidak tersebar di berbagai tempat. Ketiga, jika perlu menambah dependency baru, cukup tambah field di struct dan initialize di `NewApp()`.

Method `setupRouter()` mengkonfigurasi Gin router dengan middleware yang diperlukan:

```go
func (a *App) setupRouter() {
    // gin.Default() membuat router dengan Logger dan Recovery middleware
    router := gin.Default()

    // CORS middleware memungkinkan frontend dari domain berbeda
    // untuk mengakses API kita
    corsConfig := cors.DefaultConfig()
    corsConfig.AllowAllOrigins = true
    corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
    corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
    router.Use(cors.New(corsConfig))

    // Setup routes di file terpisah untuk organisasi yang lebih baik
    a.setupRoutes(router)
    a.router = router
}
```

CORS (Cross-Origin Resource Sharing) adalah mekanisme keamanan browser yang memblokir request dari domain berbeda. Karena frontend kita mungkin berjalan di domain berbeda dari backend (misalnya `localhost:3000` vs `localhost:8080`), kita perlu mengkonfigurasi CORS agar request diizinkan.

Method `Run()` memulai HTTP server dengan graceful shutdown:

```go
func (a *App) Run() error {
    // Buat http.Server dengan konfigurasi kustom
    srv := &http.Server{
        Addr:    ":" + a.config.ServerPort,
        Handler: a.router,
    }

    // Jalankan server di goroutine terpisah agar tidak blocking
    go func() {
        log.Printf("Server starting on :%s", a.config.ServerPort)
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("listen: %s\n", err)
        }
    }()

    // Graceful shutdown: tunggu signal interrupt (Ctrl+C) atau terminate
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit  // Block sampai menerima signal
    log.Println("Shutting down server...")

    // Beri waktu 10 detik untuk request yang sedang berjalan
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    if err := srv.Shutdown(ctx); err != nil {
        log.Fatal("Server forced to shutdown:", err)
    }

    log.Println("Server exiting")
    return nil
}
```

Graceful shutdown adalah fitur penting untuk production. Saat server menerima signal shutdown (misalnya dari Kubernetes atau saat user menekan Ctrl+C), server tidak langsung berhenti. Sebaliknya, server berhenti menerima request baru, menunggu request yang sedang berjalan selesai (dengan batas waktu 10 detik), lalu baru benar-benar berhenti. Ini mencegah client menerima error karena koneksi tiba-tiba terputus.

Buat file `routes.go` di dalam direktori `internal/api/`:

```go
package api

import (
    "time"

    "github.com/gin-gonic/gin"
)

func (a *App) setupRoutes(r *gin.Engine) {
    // Semua routes diawali dengan /api/v1 untuk versioning
    v1 := r.Group("/api/v1")
    {
        v1.GET("/health", a.HealthCheck)

        // Container routes: CRUD operasi untuk container
        containers := v1.Group("/containers")
        {
            containers.GET("", a.containerHandler.List)
            containers.GET("/:id", a.containerHandler.GetByID)
            containers.POST("", a.containerHandler.Create)
            containers.PUT("/:id", a.containerHandler.Update)
            containers.DELETE("/:id", a.containerHandler.Delete)
        }

        // Product routes: CRUD operasi untuk product
        products := v1.Group("/products")
        {
            products.GET("", a.productHandler.List)
            products.GET("/:id", a.productHandler.GetByID)
            products.POST("", a.productHandler.Create)
            products.PUT("/:id", a.productHandler.Update)
            products.DELETE("/:id", a.productHandler.Delete)
        }

        // Plan routes: CRUD dan kalkulasi
        plans := v1.Group("/plans")
        {
            plans.GET("", a.planHandler.List)
            plans.GET("/:id", a.planHandler.GetByID)
            plans.POST("", a.planHandler.Create)
            plans.PUT("/:id", a.planHandler.Update)
            plans.DELETE("/:id", a.planHandler.Delete)

            // Plan items management
            plans.POST("/:id/items", a.planHandler.AddItem)
            plans.PUT("/:id/items/:itemId", a.planHandler.UpdateItem)
            plans.DELETE("/:id/items/:itemId", a.planHandler.DeleteItem)

            // Calculation
            plans.POST("/:id/calculate", a.planHandler.Calculate)
        }
    }
}

func (a *App) HealthCheck(c *gin.Context) {
    c.JSON(200, gin.H{
        "status":  "ok",
        "time":    time.Now().Format(time.RFC3339),
        "version": "1.0.0",
    })
}
```

Routes dipisahkan ke file sendiri untuk menjaga `api.go` tetap fokus pada initialization. Gin mendukung route grouping dengan `r.Group()`—semua routes dalam group otomatis mendapat prefix yang sama. Dengan struktur ini, kita memiliki endpoint API yang lengkap:

**Master Data:**
- `GET /api/v1/containers` — List semua containers
- `POST /api/v1/containers` — Create container baru
- `PUT /api/v1/containers/:id` — Update container
- `DELETE /api/v1/containers/:id` — Hapus container
- `GET /api/v1/products` — List semua products
- `POST /api/v1/products` — Create product baru
- `PUT /api/v1/products/:id` — Update product
- `DELETE /api/v1/products/:id` — Hapus product

**Planning:**
- `POST /api/v1/plans` — Create plan baru
- `GET /api/v1/plans/:id` — Get plan detail dengan items & placements
- `PUT /api/v1/plans/:id` — Update plan header
- `DELETE /api/v1/plans/:id` — Delete plan (cascade items)
- `POST /api/v1/plans/:id/items` — Add item ke plan
- `PUT /api/v1/plans/:id/items/:itemId` — Update quantity item
- `DELETE /api/v1/plans/:id/items/:itemId` — Hapus item dari plan
- `POST /api/v1/plans/:id/calculate` — Jalankan kalkulasi packing

Health check endpoint adalah endpoint sederhana yang mengembalikan status server. Endpoint ini berguna untuk monitoring dan load balancer—mereka bisa mengecek apakah server masih hidup tanpa membebani database.

#### Entry Point

Terakhir, update `main.go` untuk menggunakan `App` struct. Buat atau update file `main.go` di dalam direktori `cmd/api/`:

```go
package main

import (
    "context"
    "log"

    "github.com/jackc/pgx/v5/pgxpool"
    "load-stuffing-calculator/internal/api"
    "load-stuffing-calculator/internal/config"
)

func main() {
    cfg := config.Load()

    pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
    if err != nil {
        log.Fatalf("Unable to connect to database: %v", err)
    }
    defer pool.Close()

    app := api.NewApp(cfg, pool)
    if err := app.Run(); err != nil {
        log.Fatal(err)
    }
}
```

Entry point menjadi sangat sederhana—cukup load config, buat database connection, dan jalankan app. Semua kompleksitas initialization tersembunyi di dalam `api.NewApp()`.

Dengan HTTP layer yang lengkap dan terstruktur, API kita sudah bisa menerima request. Langkah terakhir adalah mengintegrasikan Packing Service untuk kalkulasi yang sebenarnya.

---

### 3.8 Integrasi dengan Packing Service

Sampai titik ini, kita sudah membangun API yang lengkap—dari database hingga HTTP handlers. Namun method `Calculate` di `PlanService` masih menerima `nil` sebagai gateway. Saatnya mengimplementasikan komponen yang menghubungkan backend Go dengan Packing Service Python.

#### Gateway Pattern

Gateway adalah pattern yang mengenkapsulasi komunikasi dengan external service. Dengan gateway, service layer tidak perlu tahu detail bagaimana komunikasi dilakukan—apakah via HTTP, gRPC, atau message queue. Service hanya memanggil method interface, dan gateway menangani sisanya.

Pattern ini memberikan beberapa keuntungan. Pertama, jika protokol komunikasi berubah (misalnya dari REST ke gRPC), hanya gateway yang perlu dimodifikasi. Kedua, gateway mudah di-mock untuk testing—kita bisa menguji service tanpa perlu Packing Service berjalan. Ketiga, semua logika terkait external service (retry, timeout, error handling) terpusat di satu tempat.

Buat direktori untuk gateway:

```bash
mkdir -p internal/gateway
```

#### Interface Definition

Pertama, definisikan interface yang akan diimplementasi oleh gateway. Buat file `gateway.go` di dalam direktori `internal/gateway/`:

```go
package gateway

import (
    "context"
)

// PackingGateway mendefinisikan kontrak untuk komunikasi dengan Packing Service.
// Interface ini memungkinkan kita mengganti implementasi (HTTP, mock, dll)
// tanpa mengubah service layer.
type PackingGateway interface {
    Pack(ctx context.Context, req PackRequest) (*PackResponse, error)
}
```

Interface hanya memiliki satu method: `Pack`. Ini menerima request dan mengembalikan response atau error. Kesederhanaan ini disengaja—gateway hanya bertanggung jawab untuk komunikasi, bukan business logic.

#### Request dan Response Types

Selanjutnya, definisikan struct untuk request dan response. Struct ini merepresentasikan format data yang diharapkan oleh Packing Service.

```go
// PackRequest adalah request yang dikirim ke Packing Service
type PackRequest struct {
    Units     string         `json:"units"`      // Satuan ukuran: "mm" atau "cm"
    Container PackContainer  `json:"container"`  // Dimensi container
    Items     []PackItem     `json:"items"`      // Items yang akan di-pack
    Options   PackOptions    `json:"options"`    // Opsi algoritma
}

// PackContainer mendeskripsikan container untuk packing
type PackContainer struct {
    Length    float64 `json:"length"`
    Width     float64 `json:"width"`
    Height    float64 `json:"height"`
    MaxWeight float64 `json:"max_weight"`
}

// PackItem mendeskripsikan satu item yang akan di-pack
type PackItem struct {
    ItemID   string  `json:"item_id"`   // ID unik untuk tracking
    Label    string  `json:"label"`     // Label untuk display
    Length   float64 `json:"length"`
    Width    float64 `json:"width"`
    Height   float64 `json:"height"`
    Weight   float64 `json:"weight"`
    Quantity int     `json:"quantity"`  // Jumlah item dengan dimensi sama
}

// PackOptions berisi opsi untuk algoritma packing
type PackOptions struct {
    FixPoint            bool    `json:"fix_point"`
    CheckStable         bool    `json:"check_stable"`
    SupportSurfaceRatio float64 `json:"support_surface_ratio"`
    BiggerFirst         bool    `json:"bigger_first"`
}
```

Response dari Packing Service berisi hasil kalkulasi:

```go
// PackResponse adalah response dari Packing Service
type PackResponse struct {
    Success bool       `json:"success"`
    Data    *PackData  `json:"data,omitempty"`
    Error   *PackError `json:"error,omitempty"`
}

// PackData berisi hasil kalkulasi jika sukses
type PackData struct {
    Units      string          `json:"units"`
    Placements []PackPlacement `json:"placements"`  // Posisi setiap item
    Unfitted   []PackUnfitted  `json:"unfitted"`    // Items yang tidak muat
    Stats      PackStats       `json:"stats"`       // Statistik packing
}

// PackPlacement mendeskripsikan posisi satu item dalam container
type PackPlacement struct {
    ItemID     string  `json:"item_id"`
    Label      string  `json:"label"`
    PosX       float64 `json:"pos_x"`       // Posisi X dalam container
    PosY       float64 `json:"pos_y"`       // Posisi Y dalam container
    PosZ       float64 `json:"pos_z"`       // Posisi Z (tinggi) dalam container
    Rotation   int     `json:"rotation"`    // Rotasi yang diterapkan
    StepNumber int     `json:"step_number"` // Urutan penempatan
}

// PackUnfitted mendeskripsikan item yang tidak muat dalam container
type PackUnfitted struct {
    ItemID string `json:"item_id"`
    Label  string `json:"label"`
    Count  int    `json:"count"`
}

// PackStats berisi statistik hasil packing
type PackStats struct {
    FittedCount   int `json:"fitted_count"`   // Jumlah item yang muat
    UnfittedCount int `json:"unfitted_count"` // Jumlah item yang tidak muat
    PackTimeMs    int `json:"pack_time_ms"`   // Waktu kalkulasi
}

// PackError berisi detail error jika gagal
type PackError struct {
    Code    string `json:"code"`
    Message string `json:"message"`
}
```

Perhatikan bahwa response memiliki field `Success` untuk mengindikasikan apakah kalkulasi berhasil. Ini berbeda dengan HTTP status code—Packing Service mungkin mengembalikan 200 OK tapi dengan `Success: false` jika ada error pada level aplikasi.

#### HTTP Gateway Implementation

Sekarang implementasi gateway yang berkomunikasi via HTTP. Buat file `packing_gateway.go` di dalam direktori `internal/gateway/`:

```go
package gateway

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "strings"
    "time"
)

// HTTPPackingGateway adalah implementasi PackingGateway yang berkomunikasi via HTTP
type HTTPPackingGateway struct {
    baseURL string       // URL base Packing Service
    client  *http.Client // HTTP client dengan timeout
}

// NewHTTPPackingGateway membuat instance gateway dengan konfigurasi yang diberikan
func NewHTTPPackingGateway(baseURL string, timeout time.Duration) *HTTPPackingGateway {
    // Hapus trailing slash untuk konsistensi
    baseURL = strings.TrimRight(baseURL, "/")
    if baseURL == "" {
        baseURL = "http://localhost:5000" // Default jika tidak dikonfigurasi
    }

    return &HTTPPackingGateway{
        baseURL: baseURL,
        client: &http.Client{
            Timeout: timeout, // Timeout untuk seluruh request
        },
    }
}
```

Constructor menerima base URL dan timeout. Timeout penting karena kalkulasi packing bisa memakan waktu—kita tidak ingin request menggantung selamanya jika Packing Service tidak merespons.

Method `Pack` melakukan komunikasi dengan Packing Service:

```go
// Pack mengirim request ke Packing Service dan mengembalikan hasil
func (g *HTTPPackingGateway) Pack(ctx context.Context, req PackRequest) (*PackResponse, error) {
    // 1. Serialize request ke JSON
    payload, err := json.Marshal(req)
    if err != nil {
        return nil, fmt.Errorf("marshal pack request: %w", err)
    }

    // 2. Buat HTTP request dengan context
    httpReq, err := http.NewRequestWithContext(
        ctx,
        http.MethodPost,
        g.baseURL+"/pack",
        bytes.NewReader(payload),
    )
    if err != nil {
        return nil, fmt.Errorf("build /pack request: %w", err)
    }
    httpReq.Header.Set("Content-Type", "application/json")

    // 3. Kirim request
    resp, err := g.client.Do(httpReq)
    if err != nil {
        return nil, fmt.Errorf("call packing service: %w", err)
    }
    defer resp.Body.Close()

    // 4. Baca response body
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("read packing response: %w", err)
    }

    // 5. Parse JSON response
    var decoded PackResponse
    if err := json.Unmarshal(body, &decoded); err != nil {
        return nil, fmt.Errorf("decode packing response (status %d): %w", resp.StatusCode, err)
    }

    // 6. Handle HTTP errors
    if resp.StatusCode < 200 || resp.StatusCode >= 300 {
        msg := fmt.Sprintf("packing service returned %d", resp.StatusCode)
        if decoded.Error != nil && decoded.Error.Message != "" {
            msg = decoded.Error.Message
        }
        return nil, fmt.Errorf("packing service error: %s", msg)
    }

    // 7. Handle application-level errors
    if !decoded.Success {
        msg := "packing service returned success=false"
        if decoded.Error != nil && decoded.Error.Message != "" {
            msg = decoded.Error.Message
        }
        return nil, fmt.Errorf("packing service error: %s", msg)
    }

    return &decoded, nil
}
```

Method ini mengikuti alur yang jelas. Pertama, request di-serialize ke JSON. Kedua, HTTP request dibuat dengan context—ini memungkinkan cancellation jika client memutus koneksi. Ketiga, request dikirim dan response dibaca. Keempat, response di-parse dan error di-handle baik di level HTTP maupun aplikasi.

#### Mock Gateway untuk Development

Untuk development dan demo, kita tidak selalu memiliki Packing Service Python berjalan. Untuk itu, kita buat mock gateway yang mensimulasikan response. Buat file `mock_gateway.go` di dalam direktori `internal/gateway/`:

```go
package gateway

import (
    "context"
)

// MockPackingGateway adalah implementasi mock PackingGateway untuk testing dan demo.
type MockPackingGateway struct{}

func NewMockPackingGateway() *MockPackingGateway {
    return &MockPackingGateway{}
}

// Pack mensimulasikan kalkulasi packing dengan response sederhana
func (g *MockPackingGateway) Pack(ctx context.Context, req PackRequest) (*PackResponse, error) {
    placements := make([]PackPlacement, 0)
    stepNumber := 1
    currentX := 0.0

    // Simulasi sederhana: letakkan items berurutan sepanjang sumbu X
    for _, item := range req.Items {
        for i := 0; i < item.Quantity; i++ {
            if currentX+item.Length <= req.Container.Length {
                placements = append(placements, PackPlacement{
                    ItemID:     item.ItemID,
                    Label:      item.Label,
                    PosX:       currentX,
                    PosY:       0,
                    PosZ:       0,
                    Rotation:   0,
                    StepNumber: stepNumber,
                })
                currentX += item.Length + 10
                stepNumber++
            }
        }
    }

    return &PackResponse{
        Success: true,
        Data: &PackData{
            Units:      req.Units,
            Placements: placements,
            Unfitted:   []PackUnfitted{},
            Stats: PackStats{
                FittedCount:   len(placements),
                UnfittedCount: 0,
                PackTimeMs:    15,
            },
        },
    }, nil
}
```

Mock gateway ini sangat berguna untuk:
- Development tanpa perlu menjalankan Packing Service
- Demo aplikasi ke stakeholder
- Unit testing tanpa external dependencies

#### Wiring Gateway ke Service

Sekarang kita perlu menghubungkan gateway ke `PlanService`. Update file `api.go` di dalam direktori `internal/api/` untuk membuat gateway:

```go
import (
    // ... imports lainnya
    "load-stuffing-calculator/internal/gateway"
)

func NewApp(cfg *config.Config, db *pgxpool.Pool) *App {
    querier := store.New(db)

    // Gunakan MockPackingGateway untuk demo tanpa Packing Service
    packingGW := gateway.NewMockPackingGateway()
    
    // Untuk production dengan Packing Service, gunakan:
    // packingGW := gateway.NewHTTPPackingGateway(
    //     "http://localhost:5000",
    //     60*time.Second,
    // )

    // Initialize services
    containerSvc := service.NewContainerService(querier)
    planSvc := service.NewPlanService(querier, packingGW)

    // ... sisanya sama
}
```

Dengan mock gateway, endpoint `/calculate` bisa ditest tanpa Packing Service berjalan.

#### Update PlanService untuk Menggunakan Gateway

Terakhir, update method `Calculate` di file `plan_service.go` di dalam direktori `internal/service/` untuk menggunakan gateway yang sebenarnya:

```go
func (s *PlanService) Calculate(ctx context.Context, planID uuid.UUID) (*CalculateResult, error) {
    // ... kode pengambilan data sama seperti sebelumnya ...

    // Build request untuk Packing Service
    req := gateway.PackRequest{
        Units: "mm",
        Container: gateway.PackContainer{
            Length:    container.LengthMm,
            Width:     container.WidthMm,
            Height:    container.HeightMm,
            MaxWeight: container.MaxWeightKg,
        },
        Items:   s.buildPackItems(items),
        Options: gateway.PackOptions{
            CheckStable: true,
            BiggerFirst: true,
        },
    }

    // Panggil Packing Service via gateway
    resp, err := s.gateway.Pack(ctx, req)
    if err != nil {
        s.store.UpdatePlanStatus(ctx, store.UpdatePlanStatusParams{
            ID:     planID,
            Status: "failed",
        })
        return nil, fmt.Errorf("packing calculation failed: %w", err)
    }

    // Simpan hasil dan update status
    // ... kode sisanya sama ...
}
```

Dengan gateway yang sudah terhubung, flow kalkulasi lengkap menjadi:

1. Frontend memanggil `POST /api/v1/plans/:id/calculate`
2. Handler meneruskan ke `PlanService.Calculate()`
3. Service mengambil data dari database via SQLC-generated queries
4. Service membangun request dan memanggil `gateway.Pack()`
5. Gateway mengirim HTTP request ke Packing Service Python
6. Packing Service melakukan kalkulasi 3D bin packing
7. Response dikembalikan melalui gateway ke service
8. Service menyimpan hasil ke database
9. Response dikirim ke frontend

Dengan ini, backend Go kita sudah lengkap—dari menerima request HTTP hingga berkomunikasi dengan Packing Service untuk kalkulasi yang sebenarnya.

#### Menjalankan dan Menguji Aplikasi

Sebelum melanjutkan ke bab berikutnya, mari pastikan semua komponen bekerja dengan benar. Pertama, pastikan PostgreSQL berjalan dan database sudah dibuat:

```bash
# Buat database (jika belum ada)
createdb loadstuff

# Jalankan migrasi
goose -dir cmd/db/migrations postgres "postgres://localhost:5432/loadstuff?sslmode=disable" up
```

Set environment variables yang diperlukan:

```bash
export DATABASE_URL="postgres://localhost:5432/loadstuff?sslmode=disable"
export SERVER_PORT="8080"
```

Sebelum menjalankan aplikasi, pastikan semua dependencies terinstall dengan benar:

```bash
go mod tidy
```

Jalankan aplikasi:

```bash
go run cmd/api/main.go
```

Jika berhasil, Anda akan melihat output seperti:

```
[GIN-debug] Listening and serving HTTP on :8080
Server starting on :8080
```

Sekarang kita bisa menguji API dengan curl. Buka terminal baru dan coba endpoint health check terlebih dahulu untuk memastikan server berjalan:

```bash
curl http://localhost:8080/api/v1/health
```

Response:

```json
{
  "status": "ok",
  "time": "2024-01-15T10:30:00Z",
  "version": "1.0.0"
}
```

Health check mengembalikan status "ok" beserta waktu server dan versi aplikasi. Jika endpoint ini merespons, berarti server sudah siap menerima request.

Sekarang mari buat container baru. Container merepresentasikan kontainer pengiriman seperti kontainer 40 kaki standar. Request menggunakan method POST dengan body JSON yang berisi dimensi dalam milimeter dan berat maksimum dalam kilogram:

```bash
curl -X POST http://localhost:8080/api/v1/containers \
  -H "Content-Type: application/json" \
  -d '{
    "name": "40ft Standard",
    "length_mm": 12192,
    "width_mm": 2438,
    "height_mm": 2591,
    "max_weight_kg": 26780
  }'
```

Response:

```json
{
  "success": true,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "40ft Standard",
    "length_mm": 12192,
    "width_mm": 2438,
    "height_mm": 2591,
    "max_weight_kg": 26780
  }
}
```

Perhatikan struktur response mengikuti format `APIResponse` yang kita buat di response package. Field `success: true` mengindikasikan operasi berhasil, dan `data` berisi container yang baru dibuat lengkap dengan UUID yang di-generate oleh database. Simpan ID ini karena akan digunakan untuk membuat plan.

Untuk memverifikasi container tersimpan di database, gunakan endpoint list:

```bash
curl http://localhost:8080/api/v1/containers
```

Response:

```json
{
  "success": true,
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "40ft Standard",
      "length_mm": 12192,
      "width_mm": 2438,
      "height_mm": 2591,
      "max_weight_kg": 26780
    }
  ]
}
```

Response berupa array yang berisi semua containers. Saat ini hanya ada satu container yang baru kita buat.

Terakhir, buat plan baru yang mereferensi container tersebut. Plan adalah "rencana packing" yang akan berisi items dan hasil kalkulasi:

```bash
curl -X POST http://localhost:8080/api/v1/plans \
  -H "Content-Type: application/json" \
  -d '{
    "container_id": "550e8400-e29b-41d4-a716-446655440000"
  }'
```

Response:

```json
{
  "success": true,
  "data": {
    "id": "660f9500-f39c-52e5-b827-557766550111",
    "container_id": "550e8400-e29b-41d4-a716-446655440000",
    "status": "draft"
  }
}
```

Plan baru dibuat dengan status "draft"—artinya plan masih dalam tahap persiapan dan belum dikalkulasi. Status akan berubah ke "calculating" saat proses kalkulasi berjalan, dan "completed" setelah selesai.

Dengan mock gateway yang sudah kita buat, endpoint `/calculate` bisa diuji tanpa Packing Service berjalan. Mock gateway akan mengembalikan hasil simulasi sederhana—berguna untuk memverifikasi bahwa flow aplikasi bekerja dengan benar sebelum mengintegrasikan dengan Packing Service Python yang sebenarnya di Bab 4.

Dengan API yang berjalan dan teruji, kita telah menyelesaikan backend Go untuk Load Stuffing Calculator. Di bab berikutnya, kita akan membangun Packing Service dengan Python yang akan melakukan kalkulasi 3D bin packing.

### Summary

Dalam bab ini, kita telah membangun backend API lengkap menggunakan Go dengan arsitektur berlapis yang clean dan maintainable. Fitur utama yang telah kita implementasikan:

1.  **Product Management**: CRUD lengkap untuk master data produk.
2.  **Container Management**: CRUD lengkap untuk jenis-jenis container.
3.  **Plan Management**: Pembuatan rencana muatan, manajemen item dalam plan, dan integrasi kalkulasi.

**Arsitektur 3-Tier**
- **Handler Layer**: Menerima HTTP request, validasi input, format response
- **Service Layer**: Business logic, orkestrasi, validasi domain
- **Repository Layer**: Akses database via SQLC-generated code

**Tool dan Library**
- **Gin**: Web framework yang performant untuk HTTP routing
- **SQLC**: Type-safe SQL queries tanpa ORM overhead
- **Goose**: Database migrations yang simple dan SQL-first
- **pgx**: PostgreSQL driver native dengan connection pooling
- **godotenv**: Manajemen konfigurasi lingkungan lokal

**Pattern yang Diterapkan**
- **Dependency Injection**: Dependencies di-inject melalui constructor, memudahkan testing
- **Gateway Pattern**: Abstraksi komunikasi dengan external service
- **Standardized Response**: Format JSON konsisten untuk semua endpoint

**Struktur Direktori Final**
```
load-stuffing-calculator/
├── cmd/
│   ├── api/main.go
│   └── db/
│       ├── migrations/
│       └── queries/
├── internal/
│   ├── api/          # Route definitions & App struct
│   ├── config/       # Configuration loader
│   ├── dto/          # Request/Response structs
│   ├── gateway/      # External service clients
│   ├── handler/      # HTTP Handlers (Container, Product, Plan)
│   ├── response/     # API Response standardizer
│   ├── service/      # Business Logic (Container, Product, Plan)
│   └── store/        # Generated database code
├── go.mod
└── sqlc.yaml
```

### Further Reading

Untuk memperdalam pemahaman tentang topik yang dibahas di bab ini:

**Go Web Development**
- [Effective Go](https://go.dev/doc/effective_go) — Panduan resmi untuk menulis kode Go yang idiomatic
- [Standard Go Project Layout](https://github.com/golang-standards/project-layout) — Konvensi struktur direktori Go (meski ini bukan standar resmi, banyak diadopsi)

**Database dan SQLC**
- [SQLC Documentation](https://docs.sqlc.dev/) — Panduan lengkap SQLC termasuk fitur advanced
- [pgx Documentation](https://github.com/jackc/pgx) — Driver PostgreSQL yang kita gunakan

**Arsitektur dan Pattern**
- [Clean Architecture in Go](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) — Artikel original Uncle Bob tentang Clean Architecture
- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/) — Pattern ports and adapters yang menginspirasi gateway pattern

**Testing (untuk eksplorasi mandiri)**
- [Testing in Go](https://go.dev/doc/tutorial/add-a-test) — Tutorial resmi testing di Go
- [testify](https://github.com/stretchr/testify) — Library assertion dan mocking populer

---

<!-- BAB 4 STARTS HERE -->
