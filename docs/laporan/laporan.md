***PLATFORM LOAD* & *STUFFING CALCULATOR* DENGAN ALGORITMA *3D BIN
PACKING* UNTUK OPTIMASI DAN VISUALISASI PEMUATAN**

**LAPORAN PROYEK III**

Diajukan untuk memenuhi kelulusan pada mata kuliah Proyek III

Program Studi DIV Teknik Informatika

![A blue and orange logo AI-generated content may be
incorrect.](media/image1.png){width="4.222064741907261in"
height="1.5416666666666667in"}

Disusun oleh:

MUHAMAD SALADIN EKA SEPTIAN (714230037)

DINA OKTAFIANI (714230047)

**PROGRAM STUDI DIPLOMA IV TEKNIK INFORMATIKA**

**SEKOLAH VOKASI**

**UNIVERSITAS LOGISTIK BISNIS** **INTERNASIONAL**

**BANDUNG**

**2026**[\
]{.mark}

# LEMBAR PENGESAHAN

Laporan Proyek II ini telah diperiksa, disetujui dan disidangkan di\
Bandung, 24 Januari 2025

**Oleh:**

  -----------------------------------------------------------------------
          Penguji Pendamping,                   Penguji Utama,
  ----------------------------------- -----------------------------------
                                        [Roni Andarsyah, S.T., M.Kom.,
                                               SFPC]{.underline}

                 NIK :                         NIK : 115.88.193

              Pembimbing,                    Koordinator Proyek II

    [Roni Andarsyah, S.T., M.Kom.,        [Roni Habibi, S.Kom., M.T.,
           SFPC]{.underline}                  SFPC,]{.underline}

           NIK : 115.88.193                    NIK : 117.88.233

             Menyetujui,\             
    Ketua Program Studi D-IV Teknik   
              Informatika             

    [Roni Andarsyah, S.T., M.Kom.,    
           SFPC]{.underline}          

            NIK: 115.88.193           
  -----------------------------------------------------------------------

# SURAT PERNYATAAN TIDAK MELAKUKAN PLAGIARISME

Nama : Muhamad Saladin Eka Septian

NPM : 714230047

Program Studi : DIV Teknik Informatika

Judul : Platform Load & Stuffing Calculator Dengan Algoritma 3d Bin
Packing Untuk Optimasi Dan Visualisasi Pemuatan

Menyatakan bahwa :

1.  Proyek pemrograman Sistem Informasi (Proyek 3) saya ini adalah asli
    dan belum pernah diajukan untuk memenuhi kelulusan proyek 3 pada
    program studi DIV Teknik Informatika baik di Universitas Logistik &
    Bisnis Internasional maupun di perguruan tinggi lainnya.

2.  Proyek pemrograman Sistem Informasi (Proyek 3) ini adalah murni
    gagasan, rumusan, dan penelitian saya sendiri tanpa bantuan orang
    lain, kecuali arahan pembimbing.

3.  Dalam proyek pemrograman Sistem Informasi (Proyek 3) ini tidak
    terdapat karya atau pendapat yang telah ditulis ataupun dipublikasi
    orang lain, kecuali secara tertulis dengan jelas dicantumkan sebagai
    acuan dalam naskah dengan disebutkan nama pengarang dan dicantumkan
    dalam daftar pustaka.

4.  Pernyataan ini saya buat dengan sesungguhnya dan apabila dikemudian
    hari terdapat penyimpangan-penyimpangan dan ketidakbenaran dalam
    pernyataan ini, maka saya bersedia menerima sanksi akademik berupa
    pencabutan gelar yang telah diperoleh karena karya ini, sanksi
    lainnya sesuai norma yang berlaku di perguruan tinggi lain.

> Bandung, 8 Januari 2026
>
> Yang membuat pernyataan
>
> Muhamad Saladin Eka Septian
>
> NPM : 714230047

# SURAT PERNYATAAN TIDAK MELAKUKAN PLAGIARISME

Nama : Dina Oktafiani

NPM : 714230047

Program Studi : DIV Teknik Informatika

Judul : Platform Load & Stuffing Calculator Dengan Algoritma 3d Bin
Packing Untuk Optimasi Dan Visualisasi Pemuatan

Menyatakan bahwa :

1.  Proyek pemrograman Sistem Informasi (Proyek 3) saya ini adalah asli
    dan belum pernah diajukan untuk memenuhi kelulusan proyek 3 pada
    program studi DIV Teknik Informatika baik di Universitas Logistik &
    Bisnis Internasional maupun di perguruan tinggi lainnya.

2.  Proyek pemrograman Sistem Informasi (Proyek 3) ini adalah murni
    gagasan, rumusan, dan penelitian saya sendiri tanpa bantuan orang
    lain, kecuali arahan pembimbing.

3.  Dalam proyek pemrograman Sistem Informasi (Proyek 3) ini tidak
    terdapat karya atau pendapat yang telah ditulis ataupun dipublikasi
    orang lain, kecuali secara tertulis dengan jelas dicantumkan sebagai
    acuan dalam naskah dengan disebutkan nama pengarang dan dicantumkan
    dalam daftar pustaka.

4.  Pernyataan ini saya buat dengan sesungguhnya dan apabila dikemudian
    hari terdapat penyimpangan-penyimpangan dan ketidakbenaran dalam
    pernyataan ini, maka saya bersedia menerima sanksi akademik berupa
    pencabutan gelar yang telah diperoleh karena karya ini, sanksi
    lainnya sesuai norma yang berlaku di perguruan tinggi lain.

> Bandung, 3 Januari 2025
>
> Yang membuat pernyataan
>
> Dina Oktafiani
>
> NPM : 714230047

# KATA PENGANTAR

Puji syukur kehadirat Tuhan Yang Maha Esa. Atas rahmat dan hidayahnya,
penulis dapat menyelesaikan Proyek II yang berjudul "Aplikasi Gamifikasi
Untuk Kesehatan Dan Kebugaran Berbasis Web" dengan tepat waktu dan tanpa
halangan apapun. Tak lupa shalawat beserta salam tunjukkan kepada Nabi
besar Muhammad SAW.

Penulis menyadari bahwa proyek ini bisa diselesaikan berkat dukungan dan
bantuan dari macam pihak. Penulis mengucapkan terima kasih kepada
seluruh pihak yang sudah membantu pada penyelesaian proyek ini dan dalam
kesempatan ini penulis memberikan ucapan terima kasih pada yang
terhormat:

1.  Bapak Roni Andarsyah, S.T.,M.Kom.,SFPC, selaku ketua program studi
    DIV Teknik Informatika

2.  Bapak Roni Habibi, S.Kom., M.T., SFPC selaku koordinator proyek 3

3.  Bapak Roni Andarsyah, S.T.,M.Kom.,SFPC, selaku Dosen penguji utama
    sekaligus pembimbing

4.   , selaku Dosen penguji pendamping

5.  Orang tua dan teman-teman yang telah membantu dalam penyelesaian
    proyek ini.

Demikian tugas proyek ini disusun, semoga dapat berguna bagi seluruh
pihak dan penulis sendiri. Akhir istilah penulis ucapkan terima kasih.

Bandung, 8 Januari 2026

**\
**

#  DAFTAR ISI 

# Contents {#contents .TOC-Heading}

[LEMBAR PENGESAHAN [2](#lembar-pengesahan)](#lembar-pengesahan)

[SURAT PERNYATAAN TIDAK MELAKUKAN PLAGIARISME
[3](#surat-pernyataan-tidak-melakukan-plagiarisme)](#surat-pernyataan-tidak-melakukan-plagiarisme)

[SURAT PERNYATAAN TIDAK MELAKUKAN PLAGIARISME
[4](#surat-pernyataan-tidak-melakukan-plagiarisme-1)](#surat-pernyataan-tidak-melakukan-plagiarisme-1)

[KATA PENGANTAR [5](#kata-pengantar)](#kata-pengantar)

[DAFTAR ISI [6](#daftar-isi)](#daftar-isi)

[DAFTAR TABEL [8](#daftar-tabel)](#daftar-tabel)

[DAFTAR GAMBAR [9](#daftar-gambar)](#daftar-gambar)

[BAB I PENDAHULUAN [10](#bab-i-pendahuluan)](#bab-i-pendahuluan)

[1.1 Latar Belakang [10](#latar-belakang)](#latar-belakang)

[1.2 Rumusan Masalah [11](#rumusan-masalah)](#rumusan-masalah)

[1.3 Tujuan Penelitian [11](#tujuan-penelitian)](#tujuan-penelitian)

[1.4 Ruang Lingkup [12](#ruang-lingkup)](#ruang-lingkup)

[BAB II LANDASAN TEORI [13](#landasan-teori)](#landasan-teori)

[2.1 Tinjauan Pustaka [13](#tinjauan-pustaka)](#tinjauan-pustaka)

[2.1.1 Manajemen Logistik dan Proses *Stuffing*
[13](#manajemen-logistik-dan-proses-stuffing)](#manajemen-logistik-dan-proses-stuffing)

[2.1.2 Container Loading Problem (CLP)
[13](#container-loading-problem-clp)](#container-loading-problem-clp)

[2.1.3 Algoritma 3D Bin Packing
[14](#algoritma-3d-bin-packing)](#algoritma-3d-bin-packing)

[2.1.4 Visualisasi 3D Berbasis Web
[14](#visualisasi-3d-berbasis-web)](#visualisasi-3d-berbasis-web)

[2.1.5 *Go* (Golang) [15](#go-golang)](#go-golang)

[*2.1.6 REST API* [15](#rest-api)](#rest-api)

[*2.1.7 PostgreSQL* [16](#postgresql)](#postgresql)

[2.1.8 *JSON Web Token* (JWT)
[16](#json-web-token-jwt)](#json-web-token-jwt)

[*2.1.9 JavaScript* [17](#javascript)](#javascript)

[*2.1.10 Tailwind CSS* [17](#tailwind-css)](#tailwind-css)

[BAB III ANALISIS & PERANCANGAN
[18](#analisis-perancangan)](#analisis-perancangan)

[3.1 Analisis [18](#analisis)](#analisis)

[3.1.1 Analisis Sistem Berjalan
[18](#analisis-sistem-berjalan)](#analisis-sistem-berjalan)

[3.1.2 Analisis Sistem Usulan
[19](#analisis-sistem-usulan)](#analisis-sistem-usulan)

[3.2 Metode Pengembangan Sistem
[21](#metode-pengembangan-sistem)](#metode-pengembangan-sistem)

[3.2.1 Alur Penelitian [21](#alur-penelitian)](#alur-penelitian)

[3.2.2 Agile Model [22](#agile-model)](#agile-model)

[3.3 Analisis Kebutuhan Sistem
[25](#analisis-kebutuhan-sistem)](#analisis-kebutuhan-sistem)

[3.3.1 Kebutuhan Fungsional
[25](#kebutuhan-fungsional)](#kebutuhan-fungsional)

[3.3.2 Kebutuhan Non-Fungsional
[27](#kebutuhan-non-fungsional)](#kebutuhan-non-fungsional)

[3.4 Perancangan Sistem [28](#perancangan-sistem)](#perancangan-sistem)

[3.4.1 Diagram Use Case [29](#diagram-use-case)](#diagram-use-case)

[3.4.2 Sequence Diagram [41](#sequence-diagram)](#sequence-diagram)

[3.4.3 Arsitektur Sistem [46](#arsitektur-sistem)](#arsitektur-sistem)

[3.4.4 Perancangan Basis Data
[48](#perancangan-basis-data)](#perancangan-basis-data)

[BAB IV IMPLEMENTASI & PENGUJIAN
[51](#implementasi-pengujian)](#implementasi-pengujian)

[4.1 Lingkungan Implementasi
[51](#lingkungan-implementasi)](#lingkungan-implementasi)

[4.2 Pembahasan Hasil Implementasi
[51](#pembahasan-hasil-implementasi)](#pembahasan-hasil-implementasi)

[4.3 Pengujian Code Coverage
[63](#pengujian-code-coverage)](#pengujian-code-coverage)

[BAB V KESIMPULAN & SARAN [64](#kesimpulan-saran)](#kesimpulan-saran)

[5.1 Kesimpulan [64](#kesimpulan)](#kesimpulan)

[5.2 Saran [65](#saran)](#saran)

[BAB VI DAFTAR PUSTAKA [66](#daftar-pustaka)](#daftar-pustaka)

# DAFTAR TABEL

[Tabel 3. 1 Definisi Aktor [29](#_Toc218781103)](#_Toc218781103)

[Tabel 3. 2 Definisi Use Case [30](#_Toc218781104)](#_Toc218781104)

[Tabel 3. 3 Skenario Use case Login
[30](#_Toc218781105)](#_Toc218781105)

[Tabel 3. 4 Skenario Use Case Kelola Data Master
[31](#_Toc218781106)](#_Toc218781106)

[Tabel 3. 5 Skenario Use Case Kelola Profil Kontainer
[32](#_Toc218781107)](#_Toc218781107)

[Tabel 3. 6 Skenario Use Case Kelola Katalog Produk
[32](#_Toc218781108)](#_Toc218781108)

[Tabel 3. 7 Skenario Use Case Kelola Workspace
[33](#_Toc218781109)](#_Toc218781109)

[Tabel 3. 8 Skenario Use Case Kelola Anggota & Undangan
[34](#_Toc218781110)](#_Toc218781110)

[Tabel 3. 9 Skenario Use Case Kelola Role & Permission
[35](#_Toc218781111)](#_Toc218781111)

[Tabel 3. 10 Skenario Use Case Membuat Rencana Pengiriman
[36](#_Toc218781112)](#_Toc218781112)

[Tabel 3. 11 Skenario Use Case Kelola Barang
[37](#_Toc218781113)](#_Toc218781113)

[Tabel 3. 12 Skenario Use Case Kalkulasi Muatan
[37](#_Toc218781114)](#_Toc218781114)

[Tabel 3. 13 Skenario Use Case Lihat Visualisasi 3D
[39](#_Toc218781115)](#_Toc218781115)

[Tabel 3. 14 Skenario Use Case Export PDF
[39](#_Toc218781116)](#_Toc218781116)

[Tabel 3. 15 Skenario Use Case Akses Trial
[40](#_Toc218781117)](#_Toc218781117)

# DAFTAR GAMBAR

[Gambar 3. 1 Alur Penelitian [21](#_Toc218780631)](#_Toc218780631)

[Gambar 3. 2 Metode Agile [22](#_Toc218782999)](#_Toc218782999)

[Gambar 3. 3 Diagram Use Case [29](#_Toc218780632)](#_Toc218780632)

[Gambar 3. 4 Sequence Diagram Authentication & Workspace Context
[42](#_Toc218783001)](#_Toc218783001)

[Gambar 3. 5 Sequence Diagram Create Plan
[43](#_Toc218783002)](#_Toc218783002)

[Gambar 3. 6 Sequence Diagram Calculate Load & Visualize 3D
[44](#_Toc218783003)](#_Toc218783003)

[Gambar 3. 7 Arsitektur Sistem [46](#_Toc218783004)](#_Toc218783004)

[Gambar 3. 8 Perancangan Basis Data
[48](#_Toc218783005)](#_Toc218783005)

[Gambar 4. 1 Tampilan Halaman Login
[52](#_Toc218847794)](#_Toc218847794)

[Gambar 4. 2 Tampilan Halaman Register
[53](#_Toc218847795)](#_Toc218847795)

[Gambar 4. 3 Tampilan Halaman Dashboard
[54](#_Toc218847796)](#_Toc218847796)

[Gambar 4. 4 Tampilan Halaman Container Profiles
[55](#_Toc218847797)](#_Toc218847797)

[Gambar 4. 5 Tampilan Halaman Producting Catalog
[55](#_Toc218847798)](#_Toc218847798)

[Gambar 4. 6 Tampilan Halaman All Shipments
[56](#_Toc218847799)](#_Toc218847799)

[Gambar 4. 7 Tampilan Halaman Visualisasi 3D
[57](#_Toc218847800)](#_Toc218847800)

[Gambar 4. 8 Tampilan Halaman Create Shipment
[58](#_Toc218847801)](#_Toc218847801)

[Gambar 4. 9 Tampilan Halaman Loading Shipments
[59](#_Toc218847802)](#_Toc218847802)

[Gambar 4. 10 Tampilan Halaman Members
[60](#_Toc218847803)](#_Toc218847803)

[Gambar 4. 11 Tampilan Halaman Invites
[61](#_Toc218847804)](#_Toc218847804)

[Gambar 4. 12 Tampilan Halaman User Management
[61](#_Toc218847805)](#_Toc218847805)

[Gambar 4. 13 Tampilan Halaman Workspaces Management
[62](#_Toc218847806)](#_Toc218847806)

[Gambar 4. 14 Tampilan Halaman Roles
[63](#_Toc218847807)](#_Toc218847807)

[Gambar 4. 15 Tampilan Halaman Permissions
[64](#_Toc218847808)](#_Toc218847808)

[Gambar 4. 16 Pengujian Code Coverage
[64](#_Toc218847809)](#_Toc218847809)

# BAB I PENDAHULUAN

## Latar Belakang

> Dalam era perdagangan global yang semakin kompetitif, efisiensi
> logistik memegang peranan vital dalam menentukan daya saing suatu
> perusahaan. Biaya pengiriman, khususnya dalam transportasi laut
> menggunakan kontainer, merupakan salah satu komponen biaya terbesar
> dalam rantai pasok. Optimalisasi kapasitas muat kontainer tidak hanya
> berdampak pada pengurangan biaya operasional, tetapi juga
> berkontribusi pada efisiensi logistik maritim secara keseluruhan\[1\].
> Namun, tantangan utama yang sering dihadapi oleh industri logistik
> adalah ketidakefisienan dalam pemanfaatan ruang kargo, di mana
> seringkali terdapat ruang kosong (void space) yang terbuang akibat
> perencanaan pemuatan yang kurang\[2\].
>
> Masalah optimasi pemuatan ini dikenal secara akademis sebagai
> Container Loading Problem (CLP). CLP dikategorikan sebagai masalah
> NP-hard, yang berarti kompleksitas perhitungan meningkat secara
> eksponensial seiring bertambahnya jumlah barang dengan dimensi yang
> berbeda-beda\[3\]. Metode perencanaan manual yang masih banyak
> diterapkan oleh perusahaan freight forwarder atau eksportir seringkali
> mengandalkan intuisi manusia, yang tidak hanya memakan waktu lama
> tetapi juga rentan terhadap human error dan jarang menghasilkan solusi
> optimal. Ketidakakuratan dalam estimasi muatan ini dapat menyebabkan
> kerugian finansial akibat kelebihan kontainer atau kerusakan barang
> karena penumpukan yang tidak stabil\[4\].
>
> Untuk mengatasi permasalahan tersebut, pendekatan algoritmik menjadi
> solusi yang tidak terelakkan. Penerapan algoritma optimasi seperti 3D
> Bin Packing telah terbukti mampu meningkatkan persentase pengisian
> kontainer (fill rate) secara signifikan. Penggunaan algoritma
> komputasi dalam simulasi packing dapat memberikan solusi penataan yang
> jauh lebih efisien dibandingkan metode manual\[5\]. Namun, sekadar
> perhitungan matematis seringkali tidak cukup bagi pengguna lapangan;
> mereka membutuhkan representasi visual untuk memahami instruksi
> pemuatan\[5\].
>
> Meskipun berbagai algoritma telah dikembangkan, masih terdapat
> kesenjangan dalam ketersediaan platform yang mengintegrasikan
> perhitungan algoritma yang presisi dengan visualisasi interaktif yang
> mudah dipahami, serta fitur manajemen kolaboratif. Berdasarkan
> analisis kebutuhan tersebut, penulis mengembangkan sebuah sistem
> berbasis web berjudul \"Platform Load & Stuffing Calculator dengan
> Algoritma 3D Bin Packing untuk Optimasi dan Visualisasi Pemuatan\".
> Sistem ini dirancang menggunakan arsitektur modern yang menggabungkan
> keandalan bahasa pemrograman Go dan Python untuk pemrosesan data,
> serta Next.js untuk antarmuka visualisasi 3D, guna memberikan solusi
> perencanaan pemuatan yang akurat, visual, dan dapat diakses secara
> real-time.

## Rumusan Masalah

> Berdasarkan latar belakang yang telah diuraikan, rumusan masalah dalam
> penelitian ini adalah sebagai berikut:

1.  Bagaimana merancang dan mengimplementasikan algoritma 3D Bin Packing
    untuk mengoptimalkan penataan barang dalam kontainer guna
    meminimalkan ruang kosong?

2.  Bagaimana membangun sistem visualisasi 3D interaktif berbasis web
    yang dapat menampilkan hasil rencana pemuatan (stuffing plan) secara
    detail kepada pengguna?

3.  Bagaimana mengembangkan platform berbasis Software as a Service
    (SaaS) yang memungkinkan pengelolaan multi-pengguna (multi-tenancy),
    manajemen produk, dan rencana pengiriman secara terpusat?

## Tujuan Penelitian

> Berdasarkan rumusan masalah diatas, maka tujuan dalam penelitian ini
> adalah:

1.  Menghasilkan layanan kalkulator pemuatan (backend service) yang
    mampu memproses data dimensi produk dan kontainer menggunakan
    algoritma 3D Bin Packing untuk menghasilkan skema pemuatan yang
    optimal.

2.  Membangun antarmuka pengguna (frontend) yang interaktif dengan fitur
    visualisasi 3D, sehingga pengguna dapat melihat posisi, rotasi, dan
    urutan pemuatan barang di dalam kontainer.

3.  Menyediakan platform manajemen logistik yang memiliki fitur
    autentikasi, manajemen ruang kerja (workspace), dan peran pengguna
    (Role-Based Access Control) untuk mendukung kolaborasi tim dalam
    perencanaan pengiriman.).

## Ruang Lingkup

> Agar pembahasan dan pengembangan sistem lebih terarah, penulis
> menetapkan batasan masalah dan ruang lingkup sebagai berikut:

1.  Sistem dibangun sebagai aplikasi web yang dapat diakses melalui
    browser, sehingga pengguna tidak perlu melakukan instalasi perangkat
    lunak khusus

2.  Fokus utama sistem adalah menghitung susunan barang dalam kontainer
    secara otomatis menggunakan algoritma 3D Bin Packing dan
    menampilkannya dalam bentuk visualisasi 3D interaktif.

3.  Pengguna menginputkan data dimensi barang (panjang, lebar, tinggi)
    dan jenis kontainer. Sistem akan menghasilkan output berupa rencana
    pemuatan (stuffing plan) dan laporan manifest.

4.  Perhitungan saat ini dibatasi pada barang berbentuk balok atau kubus
    (box) dan menggunakan standar ukuran kontainer umum (20ft, 40ft).

5.  Aplikasi dirancang untuk digunakan oleh admin logistik atau planner
    yang bertugas merencanakan pengiriman barang.

#  LANDASAN TEORI

## Tinjauan Pustaka

> Dalam pengembangan sistem ini, kajian literatur dilakukan untuk
> memahami urgensi efisiensi logistik serta pendekatan teknologi yang
> tepat dalam menangani permasalahan pemuatan (stuffing). Berikut adalah
> pemaparan teori yang relevan berdasarkan penelitian terdahulu:

### **Manajemen Logistik dan Proses *Stuffing***

> Manajemen logistik modern menuntut efisiensi tinggi dalam setiap
> rantai pasok, terutama pada tahap pemuatan barang (stuffing) ke dalam
> kontainer. *Stuffing* didefinisikan sebagai proses fisik memindahkan
> kargo dari area gudang ke dalam kontainer untuk pengiriman. Efisiensi
> dalam proses ini sangat krusial karena biaya transportasi laut
> seringkali dihitung berdasarkan volume kontainer yang digunakan, bukan
> hanya berat barang.
>
> Menurut penelitian Oktavia (2023), pemanfaatan teknologi informasi dan
> sistem manajemen logistik terintegrasi merupakan strategi kunci dalam
> meningkatkan efisiensi operasional perusahaan\[6\]. Tanpa adanya
> bantuan perangkat lunak, proses *stuffing* seringkali menghadapi
> kendala berupa *void space* (ruang kosong) yang tidak terpakai secara
> optimal akibat keterbatasan estimasi manual manusia. Sejalan dengan
> itu, Widodo (2023) dalam studinya mengenai analisis proses *stuffing*
> menekankan bahwa perencanaan yang buruk tidak hanya menyebabkan
> pemborosan ruang, tetapi juga berisiko merusak barang selama
> perjalanan akibat penataan yang tidak stabil\[7\]. Oleh karena itu,
> diperlukan transformasi dari perencanaan manual menuju sistem
> kalkulasi berbasis digital yang mampu memberikan simulasi penataan
> muatan secara presisi.

### **Container Loading Problem (CLP)**

> Masalah optimasi pemuatan barang ke dalam kontainer dikenal dalam
> literatur akademis sebagai *Container Loading Problem* (CLP). CLP
> adalah permasalahan geometris tiga dimensi di mana tujuannya adalah
> menempatkan sejumlah kotak (items) berukuran kecil ke dalam wadah
> besar (container) sedemikian rupa sehingga total volume yang dimuat
> maksimal atau jumlah kontainer yang digunakan minimal.
>
> Hessler et al. (2024) mengklasifikasikan CLP sebagai masalah NP-hard,
> yang berarti kompleksitas komputasi untuk menemukan solusi paling
> optimal akan meningkat secara eksponensial seiring dengan bertambahnya
> jumlah barang\[8\]. Dalam skenario dunia nyata, CLP tidak hanya
> memperhitungkan volume, tetapi juga batasan-batasan (constraints)
> fisik. Thi Xuan Hoa (2024) menjelaskan bahwa orientasi barang (apakah
> boleh dibalik atau tidak) dan stabilitas tumpukan menjadi parameter
> kritis yang harus dihitung oleh sistem agar hasil packing dapat
> diterapkan secara aman di lapangan\[4\].

### **Algoritma 3D Bin Packing**

> Untuk menyelesaikan kompleksitas CLP, metode matematis murni
> seringkali terlalu lambat. Oleh karena itu, pendekatan algoritma
> heuristik atau meta-heuristik seperti 3D Bin Packing banyak digunakan.
> Algoritma ini bekerja dengan menempatkan barang berdasarkan aturan
> prioritas tertentu (misalnya, barang terbesar masuk lebih dulu) dan
> mencari posisi koordinat (x, y, z) yang paling efisien di dalam ruang
> 3D.
>
> Penelitian oleh Poerwandono dan Fiddin (2023) menunjukkan bahwa
> implementasi algoritma pada kasus 3D Bin Packing mampu memberikan
> koordinat penempatan objek yang jauh lebih baik dibandingkan metode
> intuisi manual\[9\]. Algoritma ini biasanya bekerja dengan
> langkah-langkah berikut:

1.  *Sorting*: Mengurutkan daftar barang berdasarkan volume, luas alas,
    atau dimensi terpanjang.

2.  *Positioning*: Menentukan titik koordinat pojok kiri-bawah-belakang
    (Deepest Bottom-Left Fill) yang tersedia untuk menempatkan barang.

3.  *Validation*: Memeriksa apakah penempatan tersebut melanggar batasan
    dimensi kontainer atau tumpang tindih dengan barang lain.

> Dalam konteks pengembangan perangkat lunak, algoritma ini sering
> diimplementasikan menggunakan bahasa pemrograman yang efisien dalam
> perhitungan matematika seperti Python atau Go, sebagaimana diterapkan
> dalam sistem ini.

### **Visualisasi 3D Berbasis Web**

> Hasil perhitungan algoritma yang berupa sekumpulan data koordinat
> angka sulit dipahami oleh operator lapangan. Oleh karena itu,
> visualisasi data menjadi jembatan penting antara sistem komputasi dan
> pengguna manusia. Perkembangan teknologi web modern memungkinkan
> rendering grafis 3D dilakukan langsung di dalam peramban (browser)
> tanpa perlu menginstal aplikasi desktop berat.
>
> Kartiko dan Primandari (2023) mencatat bahwa penggunaan media visual
> interaktif dalam pengenalan peti kemas logistik sangat efektif dalam
> meningkatkan pemahaman pengguna\[10\]. Dalam pengembangan aplikasi
> web, teknologi seperti WebGL dan pustaka Three.js memungkinkan
> pembuatan \"Digital Twin\" dari kontainer dan barang. Dengan
> visualisasi ini, pengguna dapat memutar, memperbesar (zoom), dan
> melihat lapisan demi lapisan (layer-by-layer) rencana pemuatan barang
> sebelum eksekusi fisik dilakukan, sehingga meminimalisir kesalahan
> bongkar-muat di lapangan.

### ***Go* (Golang)**

> Golang adalah bahasa pemrograman sistem yang dikompilasi (compiled)
> dan dirancang untuk menangani beban kerja server yang tinggi. Menurut
> penelitian terbaru oleh Muharam dan Hidayat (2024), Golang menawarkan
> keunggulan signifikan dalam hal efisiensi memori dan kecepatan
> eksekusi dibandingkan bahasa berbasis interpreter, berkat dukungannya
> terhadap konkurensi melalui fitur goroutine\[11\].
>
> Dalam pengembangan sistem *Load & Stuffing Calculator* ini, pemilihan
> Golang didasarkan pada kebutuhan untuk melakukan perhitungan algoritma
> muatan yang kompleks dengan performa tinggi. Sebagaimana dijelaskan
> oleh Fernando dan Engel (2025), arsitektur Golang sangat efektif dalam
> menangani ribuan koneksi simultan dengan latensi rendah, yang
> merupakan syarat utama agar hasil simulasi stuffing dapat ditampilkan
> secara instan kepada pengguna tanpa waktu tunggu yang lama\[12\].

### ***REST API***

> Representational State Transfer (REST) API merupakan standar
> arsitektur komunikasi data yang menggunakan protokol *HTTP* untuk
> pertukaran informasi antar sistem. Simbulan dan Aryanto (2024)
> menyatakan bahwa penerapan *REST API* memungkinkan integrasi sistem
> yang fleksibel dan skalabel karena sifatnya yang stateless, di mana
> server tidak perlu menyimpan status sesi dari klien, sehingga beban
> server menjadi lebih ringan\[13\].
>
> Pada proyek ini, REST API berfungsi sebagai jembatan penghubung antara
> antarmuka pengguna (frontend) berbasis web dengan logika pemrosesan di
> server (backend). Penggunaan format data JSON yang ringan dalam
> arsitektur ini, menurut Farhandika, Sabariah, dan Adrian (2024),
> terbukti meningkatkan efisiensi bandwidth dan mempercepat waktu
> respons aplikasi, yang sangat krusial untuk menyajikan data kalkulasi
> muatan secara akurat dan cepat kepada admin operasional\[14\].

### ***PostgreSQL***

> *PostgreSQL* adalah sistem manajemen basis data relasional (RDBMS)
> *open-source* yang dikenal memiliki keandalan tinggi dan fitur yang
> mendukung integritas data kompleks. Berdasarkan analisis performa
> terbaru yang dilakukan oleh Salunke dan Ouda (2024), *PostgreSQL*
> terbukti lebih unggul dalam menangani kueri yang kompleks dan
> konkurensi (banyak akses bersamaan) dibandingkan basis data relasional
> lainnya, serta menjamin keamanan data melalui kepatuhan terhadap
> standar ACID (Atomicity, Consistency, Isolation, Durability)\[15\].
>
> Relevansi *PostgreSQL* dalam sistem *Load & Stuffing Calculator* ini
> terletak pada kemampuannya menjaga konsistensi data profil kontainer
> dan spesifikasi barang. Selain itu, fitur penyimpanan tipe data
> *JSONB* pada *PostgreSQL* memungkinkan fleksibilitas skema data, yang
> menurut studi Han dan Choi (2024) sangat berguna untuk menyimpan
> riwayat hasil simulasi stuffing dan manifest muatan yang memiliki
> struktur bervariasi, tanpa mengorbankan kecepatan performa pencarian
> data saat dibutuhkan kembali oleh admin\[16\].

### ***JSON Web Token* (JWT)** 

> *JSON Web Token* (JWT) adalah standar otentikasi berbasis token yang
> dirancang untuk mengamankan pertukaran informasi antar pihak secara
> ringkas. Dalimunthe, Putra, dan Ridha (2023) menjelaskan bahwa
> penggunaan *JWT* dengan algoritma enkripsi *HMAC-SHA256* memberikan
> mekanisme keamanan yang kuat karena setiap token ditandatangani secara
> digital, sehingga integritas data pengguna terjamin dan terhindar dari
> manipulasi pihak yang tidak berwenang\[17\].
>
> Implementasi JWT dalam proyek ini mendukung sistem keamanan stateless
> yang efisien. Dengan mekanisme ini, otorisasi hak akses antara peran
> Admin, Planner, dan Operator dapat dikelola langsung melalui payload
> token, sehingga mempermudah pengelolaan akses pada arsitektur
> terdistribusi tanpa membebani basis data untuk pengecekan sesi
> berulang kali.

### ***JavaScript***

> *JavaScript* merupakan bahasa pemrograman utama di sisi klien yang
> memungkinkan interaktivitas tinggi pada aplikasi berbasis web.
> Perkembangan teknologi web modern menempatkan JavaScript sebagai
> fondasi dalam pembangunan Single Page Application (SPA). Menurut
> Bismoputro, Huda, dan Brata (2024), penggunaan JavaScript modern
> memungkinkan aplikasi web memiliki responsivitas setara aplikasi
> desktop, di mana pembaruan konten dilakukan secara dinamis tanpa perlu
> memuat ulang seluruh halaman\[18\].
>
> Dalam konteks sistem visualisasi muatan, *JavaScript* berperan vital
> untuk merender simulasi 3D susunan barang di dalam kontainer secara
> interaktif. Teknologi ini memungkinkan Planner untuk memutar sudut
> pandang dan memverifikasi posisi barang secara visual langsung melalui
> peramban, serta menerima notifikasi status validasi sensor secara
> real-time.

### ***Tailwind CSS***

> *Tailwind CSS* adalah kerangka kerja *CSS* yang mengusung konsep
> *utility-first*, yang menyediakan kelas-kelas utilitas tingkat rendah
> untuk membangun antarmuka pengguna dengan cepat. Penelitian oleh
> Azhariyah dan Mukhlis (2024) menunjukkan bahwa pendekatan ini secara
> signifikan mempercepat proses pengembangan frontend dan menghasilkan
> kode CSS yang lebih ringkas dibandingkan metode CSS tradisional,
> sekaligus menjamin konsistensi desain di seluruh aplikasi\[19\].
>
> Penerapan Tailwind CSS dalam proyek ini bertujuan untuk memastikan
> antarmuka pengguna, terutama pada modul Operator, bersifat responsif
> (mobile-friendly) dan mudah digunakan. Dengan pendekatan ini,
> penyesuaian tampilan instruksi muat pada berbagai ukuran layar
> perangkat operator dapat dilakukan dengan efisien dan konsisten.

**\
**

#  ANALISIS & PERANCANGAN

## Analisis

> Tahap analisis merupakan tahap penelitian dengan melakukan suatu
> percobaan yang menghasilkan kesimpulan dari penguraian suatu sistem
> aplikasi, sehingga dapat diketahui mekanisme sistem, masalah-masalah
> yang terjadi. Dari proses penelitian tersebut, dapat diusulkan
> perbaikan-perbaikan yang dapat membangun dan mempertinggi sistem
> kinerja alat yang akan dibuat. Analisis sistem yang akan dibangun
> disesuaikan dengan kebutuhan, berdasarkan hasil evaluasi terhadap
> sistem yang sedang berjalan.

### **Analisis Sistem Berjalan**

> Berdasarkan hasil pengamatan terhadap proses manajemen muatan logistik
> yang berjalan saat ini, ditemukan bahwa sebagian besar perusahaan
> logistik masih mengandalkan metode manual dalam merencanakan pemuatan
> kontainer. Belum adanya sistem perencanaan digital yang dilengkapi
> dengan algoritma optimasi dan visualisasi 3D menyebabkan utilisasi
> ruang kontainer tidak maksimal dan proses perencanaan memakan waktu
> lama. Secara garis besar, prosedur sistem yang sedang berjalan dapat
> diuraikan sebagai berikut:

1.  **Perancangan Muatan**

> Perencanaan penyusunan barang ke dalam kontainer dilakukan oleh
> Planner berdasarkan intuisi atau estimasi kasar menggunakan
> perhitungan sederhana. Tidak adanya visualisasi 3D dan algoritma
> optimasi menyebabkan sering terjadinya ruang kosong di dalam kontainer
> yang seharusnya bisa dimanfaatkan. Selain itu, perhitungan manual
> sulit untuk mengantisipasi batasan teknis seperti keseimbangan berat,
> urutan tumpukan, dan stabilitas muatan.

2.  **Eksekusi Pemuatan**

> Pada tahap eksekusi, instruksi muat yang diberikan kepada operator
> lapangan seringkali hanya berupa daftar barang tanpa panduan posisi
> yang spesifik. Tidak adanya visualisasi 3D atau panduan
> langkah-demi-langkah menyebabkan operator harus mengandalkan
> pengalaman pribadi dalam menata barang, yang dapat menghasilkan
> susunan yang berbeda dari rencana awal dan berpotensi tidak optimal.

3.  **Pelaporan dan Arsip**

> Data riwayat pengiriman tersimpan dalam file-file terpisah yang tidak
> terhubung dengan data master barang. Jika terjadi perubahan
> spesifikasi barang di kemudian hari, sulit untuk melacak data historis
> pengiriman yang akurat karena tidak adanya mekanisme snapshot data
> yang sistematis. Hal ini dapat menimbulkan inkonsistensi data dalam
> laporan historis.

### **Analisis Sistem Usulan**

> Untuk mengatasi permasalahan yang telah diuraikan pada analisis sistem
> berjalan, peneliti mengembangkan platform web \"Load & Stuffing
> Calculator\" yang dilengkapi algoritma optimasi 3D Bin Packing dan
> visualisasi interaktif. Solusi yang ditawarkan dalam sistem usulan ini
> mencakup lima aspek utama:

1.  **Platform Berbasis Web dengan Arsitektur Multi-Tenant**

> Sistem dibangun sebagai aplikasi web yang dapat diakses melalui
> browser tanpa memerlukan instalasi perangkat lunak khusus. Sistem
> mendukung arsitektur multi-tenant yang memungkinkan multiple
> organisasi atau perusahaan menggunakan platform yang sama dengan data
> yang terisolasi secara logis. Setiap workspace memiliki katalog
> produk, profil kontainer, dan rencana pengiriman yang independen.
> Arsitektur ini mendukung deployment sebagai Software-as-a-Service
> (SaaS) dengan pengelolaan hak akses berbasis role (Role-Based Access
> Control) yang granular.

2.  **Algoritma Optimasi 3D Bin Packing dengan Pertimbangan Batasan
    Fisik**

> Berbeda dengan perencanaan manual yang mengandalkan intuisi, sistem
> dilengkapi dengan mesin kalkulasi yang menerapkan algoritma 3D Bin
> Packing. Algoritma yang digunakan merupakan pengembangan dari library
> py3dbp yang telah ditingkatkan dengan kemampuan memperhitungkan
> gravitasi, verifikasi stabilitas muatan, dan perhitungan rasio
> permukaan penyangga. Sistem mampu menghitung susunan barang yang
> paling optimal dengan mempertimbangkan batasan fisik seperti dimensi
> kontainer, berat maksimum, keseimbangan distribusi berat, dan
> stabilitas tumpukan. Implementasi dilakukan menggunakan layanan
> terpisah yang memungkinkan skalabilitas dan pemilihan algoritma yang
> optimal.

3.  **Visualisasi 3D Interaktif dan Panduan Pemuatan**

> Hasil perhitungan algoritma disajikan dalam bentuk visualisasi 3D
> interaktif menggunakan teknologi Three.js yang memungkinkan Planner
> untuk melihat simulasi posisi barang secara real-time. Visualisasi
> dilengkapi dengan fitur step-by-step playback yang menampilkan urutan
> pemuatan barang secara bertahap, sehingga operator lapangan dapat
> memahami sequence yang optimal. Sistem juga menyediakan fitur export
> ke format PDF yang berisi instruksi pemuatan lengkap dengan snapshot
> 3D untuk setiap langkah, koordinat posisi, dan informasi detail
> barang. Dokumen ini dapat dicetak dan digunakan sebagai panduan di
> lapangan.

4.  **Manajemen Data Master dan Katalog Produk**

> Sistem menyediakan fitur manajemen data master yang komprehensif,
> mencakup katalog produk (dimensi, berat, warna visualisasi, opsi
> rotasi) dan profil kontainer standar (dimensi internal, kapasitas
> berat maksimum). Data master ini dapat digunakan berulang kali untuk
> mempercepat proses perencanaan pengiriman. Sistem mendukung input data
> secara manual melalui form yang tervalidasi, memastikan konsistensi
> dan akurasi data yang tersimpan dalam database.

5.  **Integritas Data Historis dengan Snapshot Pattern**

> Sistem menerapkan design pattern Snapshot Data, di mana setiap detail
> barang dan spesifikasi kontainer yang masuk ke dalam rencana
> pengiriman akan diduplikasi dari data master pada saat pembuatan
> rencana. Hal ini menjamin bahwa data riwayat pengiriman tetap akurat
> dan tidak terpengaruh jika terjadi perubahan spesifikasi pada data
> master di masa mendatang. Mekanisme ini sangat penting untuk menjaga
> integritas laporan historis dan audit trail, memastikan bahwa
> dokumentasi pengiriman yang telah dilakukan tidak mengalami distorsi
> akibat perubahan data master.
>
> Dengan penerapan sistem usulan ini, diharapkan proses perencanaan
> muatan menjadi lebih akurat dan efisien, utilisasi ruang kontainer
> dapat dimaksimalkan melalui algoritma optimasi, dan dokumentasi
> pengiriman terjaga integritasnya untuk keperluan analisis dan
> pelaporan di masa mendatang.

## Metode Pengembangan Sistem

### **Alur Penelitian**

> Untuk memastikan pengembangan sistem berjalan sistematis dan terarah,
> penulis merancang alur penelitian yang terdiri dari tiga tahapan
> utama: Pengumpulan Data (Data Collection), Analisis (Analysis), dan
> Pengembangan Sistem (Agile Model). Alur penelitian ini digambarkan
> pada gambar berikut:

![A diagram of a data collection AI-generated content may be
incorrect.](media/image2.png){width="3.584613954505687in"
height="4.990964566929134in"}

[]{#_Toc218780631 .anchor}Gambar 3. 1 Alur Penelitian

> Tahapan penelitian dijelaskan sebagai berikut:

1.  Pengumpulan data (Data Collection): Tahap ini merupakan langkah awal
    penelitian yang bertujuan untuk mengidentifikasi permasalahan pada
    sistem manual yang sedang berjalan serta merumuskan kebutuhan sistem
    yang akan dibangun

2.  Analisis (System Analysis): Data yang terkumpul dianalisis untuk
    merumuskan solusi.

3.  Tahap Agile Model Ini adalah fase utama pembuatan perangkat lunak
    dan keras. Proses ini dilakukan secara berulang mulai dari System
    Planning hingga Feedback Review.

### **Agile Model**

![A diagram of a software development process AI-generated content may
be incorrect.](media/image3.png){width="3.3037959317585304in"
height="2.853312554680665in"}

[]{#_Toc218782999 .anchor}Gambar 3. 2 Metode Agile

> Metode Agile adalah pendekatan pengembangan perangkat lunak yang
> fleksibel dan iteratif, di mana produk dikembangkan dalam siklus
> pendek yang disebut sprint. Setiap sprint berfokus pada pembuatan
> fitur tertentu, memungkinkan tim untuk menyesuaikan pengembangan
> dengan perubahan kebutuhan pengguna dan stakeholder. Dengan pendekatan
> ini, pengembang dapat bekerja lebih efektif dalam menangani proyek
> yang kompleks, menghindari masalah perencanaan jangka panjang yang
> kaku, serta memastikan bahwa produk akhir benar-benar sesuai dengan
> kebutuhan operasional logistik.
>
> Dalam pengembangan sistem perencanaan muatan ini, Agile memberikan
> fleksibilitas yang memungkinkan perubahan pada algoritma atau
> arsitektur sistem dilakukan dengan cepat tanpa mengganggu keseluruhan
> proyek. Tim bekerja dalam siklus pendek yang memungkinkan evaluasi dan
> perbaikan terus-menerus, menjadikan hasil pengembangan lebih optimal
> dan relevan bagi Planner maupun Operator di lapangan. Pendekatan ini
> juga meningkatkan komunikasi dalam tim, memastikan bahwa setiap aspek
> dari kalkulasi muatan hingga visualisasi 3D dipahami dengan jelas.

1.  **Requirements**

> Tahap pertama adalah menentukan kebutuhan pengguna dan merancang
> Product Backlog, yaitu daftar fitur yang akan dikembangkan dalam
> beberapa sprint ke depan. Dengan pendekatan ini, kebutuhan operasional
> dapat disesuaikan secara berkala, memungkinkan tim untuk fokus pada
> fitur-fitur krusial terlebih dahulu.
>
> Dalam proyek ini, fitur-fitur utama yang diidentifikasi meliputi
> Sistem Autentikasi dan Manajemen Pengguna (RBAC), Manajemen Master
> Data (Kontainer & Produk), Algoritma 3D Bin Packing dengan
> Pertimbangan Batasan Fisik, Visualisasi 3D Interaktif, Arsitektur
> Multi-Tenant untuk isolasi data antar organisasi, serta fitur Export
> PDF untuk panduan pemuatan. Semua fitur ini dirancang berdasarkan
> observasi masalah di operasional logistik dan studi literatur yang
> mendukung implementasi algoritma optimasi ruang kontainer.

2.  **Design**

> Setelah kebutuhan dikumpulkan, tim melakukan perancangan sistem yang
> mencakup arsitektur perangkat lunak, skema database, serta alur
> pengguna. Semua aspek ini dirancang untuk memastikan bahwa sistem
> berjalan dengan efisien, aman, dan mudah digunakan. Tahapan ini
> melibatkan pembuatan Use Case Diagram, Activity Diagram, Entity
> Relationship Diagram, serta System Architecture yang menggambarkan
> komunikasi data antara frontend, backend API, layanan kalkulasi, dan
> database. Arsitektur dirancang dengan pemisahan layanan (backend Go,
> layanan kalkulasi Python, frontend Next.js) untuk memungkinkan
> skalabilitas dan pemeliharaan yang lebih baik.
>
> Desain basis data dilakukan dengan memperhatikan aturan bisnis khusus
> seperti Snapshot Pattern untuk memastikan integritas data historis
> pengiriman. Perancangan antarmuka juga dilakukan, meliputi dashboard
> perencanaan berbasis web untuk Admin/Planner dan tampilan visualisasi
> 3D yang interaktif untuk memahami hasil kalkulasi pemuatan.

3.  **Development**

> Pengembangan sistem dilakukan dalam beberapa sprint dengan
> masing-masing sprint berfokus pada modul tertentu. Pendekatan ini
> memungkinkan tim untuk membangun, menguji, dan memperbaiki fitur
> secara berkala. Pengembangan dimulai dari setup environment dan
> database menggunakan PostgreSQL, kemudian berlanjut ke implementasi
> backend API menggunakan Go dengan framework Gin, integrasi layanan
> kalkulasi terpisah menggunakan Python untuk algoritma 3D Bin Packing,
> pengembangan frontend dengan Next.js dan TypeScript, implementasi
> visualisasi 3D menggunakan Three.js, hingga akhirnya penyempurnaan
> fitur multi-tenancy dan sistem autentikasi berbasis JWT.
>
> Teknologi yang digunakan mencakup Go untuk REST API backend, Python
> dengan library py3dbp (fork jerry800416) untuk layanan kalkulasi,
> Next.js untuk web frontend, PostgreSQL untuk database, dan Docker
> untuk containerization dan deployment. Semua komponen dikembangkan
> dengan pendekatan modular untuk memudahkan testing dan maintenance.

4.  **Testing**

> Setelah fitur dikembangkan, sistem menjalani serangkaian pengujian
> untuk memastikan bahwa sistem bekerja dengan baik dan sesuai dengan
> kebutuhan pengguna. Pengujian dilakukan dengan pendekatan Unit Testing
> untuk logika bisnis dan algoritma, serta Integration Testing untuk
> endpoint API dan komunikasi antar layanan. Khusus pada sistem ini,
> dilakukan juga pengujian akurasi algoritma kalkulasi untuk memastikan
> hasil penataan barang optimal dan feasible secara fisik.
>
> Target pengujian adalah memastikan sistem mampu menghitung posisi
> barang dengan akurat, visualisasi 3D menampilkan hasil dengan benar,
> dan performa sistem dapat menangani beban kalkulasi untuk rencana
> muatan dengan jumlah barang yang banyak. Pengujian juga mencakup
> validasi keamanan autentikasi dan otorisasi untuk memastikan isolasi
> data antar workspace.

5.  **Deployment**

> Setelah sistem melewati tahap pengujian dan dianggap stabil, platform
> akan memasuki tahap deployment. Deployment dilakukan menggunakan
> Docker Compose yang mengorkestra beberapa container: PostgreSQL untuk
> database, backend API Go, layanan kalkulasi Python, dan frontend
> Next.js. Konfigurasi deployment mencakup setup database schema melalui
> migration tools (Goose), optimalisasi performa database dengan
> indexing yang tepat, serta konfigurasi environment variables untuk
> keamanan credential.
>
> Platform di-deploy pada server dengan akses melalui web browser,
> memudahkan pengguna untuk mengakses sistem dari berbagai perangkat
> tanpa instalasi khusus. Dokumentasi API juga dipublikasikan melalui
> Swagger UI untuk memudahkan integrasi dengan sistem lain jika
> diperlukan.

6.  **Review**

> Setelah setiap sprint selesai, dilakukan tinjauan untuk
> mendemonstrasikan fitur yang telah dikembangkan. Melalui evaluasi yang
> terus-menerus, tim dapat meningkatkan akurasi algoritma penataan
> barang, memperbaiki performa visualisasi 3D, mengoptimalkan response
> time API, serta meningkatkan pengalaman pengguna pada antarmuka.
> Dengan adanya pendekatan Agile ini, platform Load & Stuffing
> Calculator berhasil memberikan solusi yang efektif dalam
> mengoptimalkan kapasitas muatan dan menyediakan visualisasi yang
> intuitif untuk operasional logistik.

## Analisis Kebutuhan Sistem

> Analisis kebutuhan bertujuan untuk mendefinisikan spesifikasi
> fungsional dan non-fungsional yang harus dipenuhi oleh sistem Load &
> Stuffing Calculator. Identifikasi kebutuhan ini didasarkan pada hasil
> analisis masalah dan kebutuhan pengguna yang telah dirumuskan pada
> tahap perencanaan.

### **Kebutuhan Fungsional**

> Kebutuhan fungsional menggambarkan proses-proses layanan yang harus
> disediakan oleh sistem bagi pengguna. Berdasarkan hak aksesnya,
> kebutuhan fungsional dibagi menjadi empat kategori:

1.  **Kebutuhan Fungsional - Admin**

    a.  Sistem harus menyediakan fitur pengelolaan data pengguna (User
        Management), termasuk pembuatan akun dengan pengaturan role dan
        permission yang granular.

    b.  Sistem harus mendukung pembuatan dan pengelolaan workspace untuk
        isolasi data antar organisasi, memungkinkan deployment sebagai
        platform multi-tenant.

    c.  Sistem harus dapat mengelola keanggotaan workspace dan mengirim
        undangan kepada user baru untuk bergabung ke dalam workspace
        tertentu dengan role yang ditentukan.

    d.  Sistem harus dapat mengelola data master profil kontainer, yang
        mencakup nama, dimensi dalam (panjang, lebar, tinggi), dan
        kapasitas berat maksimum.

    e.  Sistem harus menyediakan fitur pengelolaan katalog produk
        (Master Product), mencakup spesifikasi fisik (dimensi dan
        berat), warna visualisasi, dan opsi rotasi untuk keperluan
        kalkulasi.

    f.  Sistem harus dapat mengelola role dan permission secara
        fleksibel untuk mengatur hak akses user pada berbagai fitur
        sistem.

    g.  Admin memiliki beberapa level (founder untuk akses
        platform-wide, workspace owner untuk pengelolaan organisasi, dan
        personal workspace untuk user individual). Implementasi role ini
        menggunakan sistem Role-Based Access Control (RBAC) dengan
        wildcard permission matching untuk fleksibilitas pengelolaan hak
        akses.

2.  **Kebutuhan Fungsional - Planner**

    a.  Sistem harus memungkinkan Planner membuat rencana pengiriman
        baru (Shipment Plan) dengan memilih jenis kontainer yang
        tersedia dari katalog.

    b.  Sistem harus mendukung input data barang secara hibrida, yaitu
        melalui pemilihan dari katalog master atau input manual untuk
        barang non-standar dengan spesifikasi khusus.

    c.  Sistem harus memiliki mesin kalkulasi yang mampu menghitung
        susunan barang paling optimal menggunakan algoritma 3D Bin
        Packing dengan pertimbangan batasan fisik seperti gravitasi,
        stabilitas, dan distribusi berat.

    d.  Sistem harus menyediakan opsi konfigurasi parameter kalkulasi
        seperti strategi packing, verifikasi stabilitas, dan rasio
        permukaan penyangga untuk menghasilkan hasil yang sesuai dengan
        kebutuhan spesifik.

    e.  Sistem harus menampilkan hasil perhitungan dalam bentuk
        visualisasi 3D interaktif yang dapat diputar, diperbesar, dan
        dilihat urutan penataannya secara step-by-step.

    f.  Sistem harus dapat mengekspor hasil kalkulasi ke format PDF yang
        berisi instruksi pemuatan lengkap dengan visualisasi 3D per
        langkah, koordinat posisi, dan informasi detail barang.

3.  **Kebutuhan Fungsional - Operator**

    a.  Sistem harus menyediakan antarmuka untuk melihat daftar rencana
        pengiriman yang telah dikalkulasi beserta statusnya.

    b.  Sistem harus menampilkan visualisasi 3D interaktif dari hasil
        kalkulasi dengan fitur step-by-step playback untuk memahami
        urutan pemuatan yang optimal.

    c.  Sistem harus menyediakan akses ke dokumen instruksi pemuatan
        dalam format PDF untuk digunakan sebagai panduan di lapangan.

    d.  Sistem harus menampilkan informasi detail setiap barang termasuk
        dimensi, berat, posisi koordinat, dan urutan pemuatan.

4.  **Kebutuhan Fungsional - Guest/Trial**

    a.  Sistem harus memungkinkan user Guest mencoba platform tanpa
        registrasi penuh dengan batasan jumlah rencana pengiriman yang
        dapat dibuat.

    b.  Sistem harus menyediakan mekanisme untuk mengkonversi akun Guest
        menjadi akun penuh dengan kemampuan mengklaim data trial yang
        telah dibuat sebelumnya.

### **Kebutuhan Non-Fungsional**

> Kebutuhan non-fungsional mendefinisikan batasan, standar kualitas, dan
> elemen teknis yang mendukung operasional sistem:

1.  Reliability: Sistem harus menerapkan mekanisme Snapshot Data, di
    mana setiap data barang dan spesifikasi kontainer yang masuk ke
    dalam rencana pengiriman akan diduplikasi dari data master pada saat
    pembuatan rencana. Hal ini menjamin bahwa data riwayat pengiriman
    tetap akurat dan tidak terpengaruh oleh perubahan pada data master
    di masa depan.

2.  Performance: Sistem harus memiliki response time yang baik untuk
    operasi-operasi kritis seperti kalkulasi bin packing dan rendering
    visualisasi 3D. API endpoints harus responsif untuk mendukung
    pengalaman pengguna yang lancar, dan proses kalkulasi harus dapat
    menangani rencana pengiriman dengan jumlah barang yang banyak dalam
    waktu yang reasonable.

3.  Usability: Antarmuka pengguna harus didesain responsif dan intuitif,
    dapat diakses dari berbagai perangkat (desktop, tablet, mobile)
    melalui web browser tanpa memerlukan instalasi aplikasi khusus.
    Dashboard harus disesuaikan dengan role pengguna untuk menampilkan
    informasi yang relevan.

4.  Scalability: Sistem harus dapat menangani multiple workspace dengan
    multiple user secara concurrent. Arsitektur layanan terpisah untuk
    proses kalkulasi memungkinkan scaling independen berdasarkan beban
    kerja.

5.  Interoperability: Sistem harus menggunakan format pertukaran data
    standar (JSON) untuk komunikasi antara frontend dan backend API.
    Dokumentasi API harus tersedia melalui Swagger UI untuk memudahkan
    integrasi dengan sistem eksternal jika diperlukan.

6.  Security: Akses ke dalam sistem harus dibatasi melalui otentikasi
    berbasis token JWT dengan mekanisme refresh token. Password harus
    disimpan dalam bentuk hash. Sistem harus menerapkan otorisasi
    berbasis role dan permission untuk memastikan user hanya dapat
    mengakses fitur sesuai hak aksesnya.

7.  Data Isolation: Sistem harus memastikan data antar workspace
    terisolasi secara logis dalam database. Setiap query harus otomatis
    di-scope berdasarkan workspace aktif user untuk mencegah kebocoran
    data antar organisasi.

8.  Maintainability: Sistem harus dibangun dengan arsitektur modular
    yang memisahkan concern (frontend, backend API, layanan kalkulasi)
    untuk memudahkan pengembangan, testing, dan maintenance. Kode harus
    menggunakan type-safe queries untuk mengurangi runtime error.

## Perancangan Sistem

> Perancangan sistem bertujuan untuk memberikan gambaran visual mengenai
> arsitektur, alur data, dan interaksi antar komponen dalam aplikasi
> Load & Stuffing Calculator. Perancangan ini menjadi acuan utama dalam
> tahap implementasi.

### **Diagram Use Case** 

> ![](media/image4.png){width="5.793170384951881in"
> height="2.929319772528434in"}

[]{#_Toc218780632 .anchor}Gambar 3. 3 Diagram Use Case

> Use case diagram pada Gambar 3.3 menggambarkan interaksi antara
> pengguna dan sistem. Aktor utama dalam sistem ini adalah Admin,
> Planner, Operator, dan Guest. Use case utama mencakup Manajemen Data
> Master, Manajemen Workspace, Perencanaan Pengiriman (Shipment Plan),
> Kalkulasi Muatan dengan Algoritma 3D Bin Packing, dan Visualisasi 3D
> Interaktif.

1.  **Definisi Use Case**

[]{#_Toc218781103 .anchor}Tabel 3. 1 Definisi Aktor

  -----------------------------------------------------------------------------
   **No.**  **Aktor**   **Keterangan**
  --------- ----------- -------------------------------------------------------
     1\.    Admin       Aktor yang memiliki hak akses penuh untuk melakukan
                        manajemen sistem, meliputi pengelolaan user, workspace,
                        role & permission, profil kontainer, dan katalog
                        produk.

     2\.    Planner     User yang bertugas membuat rencana pengiriman (Load
                        Planning), melakukan input data barang dari katalog
                        atau manual, menjalankan kalkulasi optimasi muatan, dan
                        mengekspor hasil ke PDF.

     3\.    Operator    User yang bertugas melihat rencana pengiriman yang
                        telah dikalkulasi, mengakses visualisasi 3D dengan
                        step-by-step playback, dan melihat instruksi pemuatan
                        untuk eksekusi di lapangan.

     4\.    Guest       User trial yang dapat mencoba platform dengan batasan
                        (maksimal 3 rencana pengiriman) sebelum melakukan
                        registrasi penuh.
  -----------------------------------------------------------------------------

  ----------------------------------------------------------------------------
   **No.**  **Use Case**    **Deskripsi**
  --------- --------------- --------------------------------------------------
      1     Login           Proses autentikasi pengguna (Admin, Planner,
                            Operator, Guest) untuk masuk ke dalam sistem
                            sesuai hak akses masing-masing.

      2     Manage Master   Use case induk untuk pengelolaan data referensi
            Data            sistem (container profiles, product catalog).

      3     Manage          Proses pengelolaan (tambah, ubah, hapus) data
            Container       jenis kontainer, termasuk dimensi dan batas berat.
            Profiles        

      4     Manage Product  Proses pengelolaan data barang yang sering dikirim
            Catalog         untuk mempercepat proses input saat perencanaan.

      5     Manage          Proses pembuatan dan pengelolaan workspace untuk
            Workspace       isolasi data multi-tenant antar organisasi.

      6     Manage Members  Proses pengelolaan keanggotaan workspace dan
            & Invites       pengiriman undangan kepada user baru dengan role
                            tertentu.

      7     Manage Roles &  Proses pengelolaan role dan permission secara
            Permissions     granular untuk mengatur hak akses user pada
                            berbagai fitur.

      8     Create Shipment Proses pembuatan dokumen rencana pengiriman baru
            Plan            oleh Planner dengan memilih container dari
                            katalog.

      9     Manage Items    Proses input atau pengelolaan daftar barang yang
                            akan dimasukkan ke dalam rencana pengiriman (dari
                            katalog atau manual).

     10     Calculate Load  Proses menjalankan algoritma 3D Bin Packing dengan
            (Algo)          parameter strategi, gravity simulation, dan
                            stability checking.

     11     View 3D Result  Fitur untuk menampilkan hasil perhitungan
                            algoritma dalam bentuk visualisasi 3D interaktif
                            dengan step-by-step playback.

     12     Export PDF      Fitur untuk mengekspor hasil kalkulasi dan
            Report          instruksi pemuatan ke dalam format PDF untuk
                            digunakan operator di lapangan.

     13     Trial Access    Fitur yang memungkinkan Guest mencoba platform
                            dengan batasan 3 rencana pengiriman sebelum
                            registrasi penuh.
  ----------------------------------------------------------------------------

  : []{#_Toc218781104 .anchor}Tabel 3. 2 Definisi Use Case

2.  **Skenario Use Case**

  -----------------------------------------------------------------------
  **Identifikasi**                  
  --------------------------------- -------------------------------------
  Nomor                             UC-01

  Nama                              Login

  Tujuan                            Melakukan autentikasi untuk masuk ke
                                    dalam sistem sesuai hak akses

  Deskripsi                         Pengguna masuk ke sistem dengan
                                    memasukkan Username dan Password yang
                                    terdaftar

  Aktor                             Admin, Planner, Operator

  **Skenario**                      

  Kondisi Awal                      Menampilkan halaman Login

  **Aksi Aktor**                    **Reaksi Sistem**

  1\. Aktor memasukkan email dan    2\. Sistem melakukan validasi format
  password                          email dan kelengkapan input.

  3\. Aktor menekan tombol "Login"  4\. Sistem memverifikasi kredensial
                                    dengan database (password hash
                                    comparison).

                                    5\. Jika berhasil, sistem
                                    menghasilkan access token (JWT) dan
                                    refresh token.

                                    6\. Sistem mengambil data workspace
                                    aktif dan daftar role/permission
                                    user.

                                    7\. Mengarahkan ke dashboard sesuai
                                    dengan role user (Admin Dashboard,
                                    Planner Dashboard, atau Operator
                                    Dashboard).

  Kondisi Akhir                     Menampilkan dashboard sesuai dengan
                                    hak akses aktor
  -----------------------------------------------------------------------

  : []{#_Toc218781105 .anchor}Tabel 3. 3 Skenario Use case Login

  -----------------------------------------------------------------------
  **Identifikasi**                  
  --------------------------------- -------------------------------------
  Nomor                             UC-02

  Nama                              Manage Master Data (Kelola Data
                                    Master)

  Tujuan                            Mengakses menu pengelolaan data
                                    referensi sistem

  Deskripsi                         Use case induk di mana Admin memilih
                                    jenis data master yang ingin dikelola

  Aktor                             Admin

  **Skenario**                      

  Kondisi Awal                      Menampilkan Dashboard Admin

  **Aksi Aktor**                    **Reaksi Sistem**

  1\. Aktor memilih menu "Master    2\. Menampilkan sub-menu pilihan:
  Data" pada sidebar Co             ntainer Profiles dan Product Catalog.

  3\. Aktor memilih salah satu      4\. Sistem mengarahkan ke halaman
  sub-menu                          pengelolaan sesuai sub-menu yang
                                    dipilih (lanjut ke UC-03 atau UC-04).

  Kondisi Akhir                     Menampilkan halaman daftar data (List
                                    View) sesuai pilihan
  -----------------------------------------------------------------------

  : []{#_Toc218781106 .anchor}Tabel 3. 4 Skenario Use Case Kelola Data
  Master

  -----------------------------------------------------------------------
  **Identifikasi**                  
  --------------------------------- -------------------------------------
  Nomor                             UC-03

  Nama                              Manage Container Profiles (Kelola
                                    Profil Kontainer)

  Tujuan                            Menambah, mengubah, atau menghapus
                                    jenis kontainer standar perusahaan

  Deskripsi                         Admin menginput spesifikasi fisik
                                    kontainer (Dimensi & Berat)

  Aktor                             Admin

  **Skenario**                      

  Kondisi Awal                      Menampilkan halaman daftar Container
                                    Profiles

  **Aksi Aktor**                    **Reaksi Sistem**

  1\. Aktor menekan tombol \"Tambah 2\. Menampilkan form input data
  Kontainer\"                       kontainer.

  3\. Menginput Nama (misal: \"20ft 4\. Validasi input (angka harus
  Dry\"), Dimensi Dalam (PxLxT mm), positif).
  dan Berat Maks (kg).              

  5\. Menekan tombol \"Simpan\"     6\. Menyimpan data profil kontainer
                                    baru ke database.

                                    7\. Menampilkan pesan sukses dan
                                    memperbarui daftar kontainer.

  Kondisi Akhir                     Data jenis kontainer baru tersedia
                                    untuk digunakan dalam Planning
  -----------------------------------------------------------------------

  : []{#_Toc218781107 .anchor}Tabel 3. 5 Skenario Use Case Kelola Profil
  Kontainer

  -----------------------------------------------------------------------
  **Identifikasi**                  
  --------------------------------- -------------------------------------
  Nomor                             UC-04

  Nama                              Manage Product Catalog (Kelola
                                    Katalog Produk)

  Tujuan                            Mengelola data induk barang agar bisa
                                    digunakan berulang kali oleh Planner

  Deskripsi                         Admin melakukan tambah, ubah, atau
                                    hapus data spesifikasi barang (SKU,
                                    Dimensi, Berat, Opsi Rotasi)

  Aktor                             Admin

  **Skenario**                      

  Kondisi Awal                      Menampilkan halaman daftar Product
                                    Catalog

  **Aksi Aktor**                    **Reaksi Sistem**

  1\. Aktor menekan tombol "Tambah  2\. Menampilkan form input data
  Produk"                           produk.

  3\. Menginput SKU, Nama Barang,   4\. Melakukan validasi input
  Dimensi (PxLxT), Berat, Warna     (pastikan SKU unik dan dimensi
  Visualisasi, dan Opsi Rotasi.     valid).

  5\. Menekan tombol "Simpan"       6\. Menyimpan data produk baru ke
                                    database.

                                    7\. Menampilkan notifikasi sukses dan
                                    memperbarui tampilan daftar produk.

  Kondisi Akhir                     Produk baru terdaftar di sistem dan
                                    siap dipilih saat perencanaan
  -----------------------------------------------------------------------

  : []{#_Toc218781108 .anchor}Tabel 3. 6 Skenario Use Case Kelola
  Katalog Produk

  -----------------------------------------------------------------------
  **Identifikasi**                  
  --------------------------------- -------------------------------------
  Nomor                             UC-05

  Nama                              Manage Workspace (Kelola Workspace)

  Tujuan                            Membuat dan mengelola workspace untuk
                                    isolasi data multi-tenant antar
                                    organisasi

  Deskripsi                         Admin membuat workspace baru dengan
                                    tipe tertentu
                                    (personal/team/organization) untuk
                                    memisahkan data antar tenant

  Aktor                             Admin

  **Skenario**                      

  Kondisi Awal                      Menampilkan Dashboard Admin

  **Aksi Aktor**                    **Reaksi Sistem**

  1\. Aktor memilih menu "Workspace 2\. Menampilkan daftar workspace yang
  Management"                       sudah ada beserta statusnya.

  3\. Aktor menekan tombol "Create  4\. Menampilkan form pembuatan
  Workspace"                        workspace baru.

  5\. Menginput Nama Workspace      6\. Melakukan validasi kelengkapan
  (misal: "PT Logistik Indonesia")  data dan keunikan nama workspace
  dan memilih Tipe                  dalam scope user.
  (Personal/Team/Organization).     

  7\. Menekan tombol "Create"       8\. Menyimpan data workspace baru ke
                                    database.

                                    9\. Membuat relasi workspace_member
                                    dengan user sebagai founder/owner
                                    dengan role admin penuh.

                                    10\. Menampilkan pesan sukses dan
                                    memperbarui daftar workspace.

  Kondisi Akhir                     Workspace baru terdaftar di sistem
                                    dan siap digunakan untuk mengelola
                                    data (container, product, shipment
                                    plan)
  -----------------------------------------------------------------------

  : []{#_Toc218781109 .anchor}Tabel 3. 7 Skenario Use Case Kelola
  Workspace

  -----------------------------------------------------------------------
  **Identifikasi**                  
  --------------------------------- -------------------------------------
  Nomor                             UC-06

  Nama                              Manage Members & Invites (Kelola
                                    Anggota & Undangan)

  Tujuan                            Mengelola keanggotaan workspace dan
                                    mengirim undangan kepada user baru
                                    dengan role tertentu

  Deskripsi                         Admin mengundang user baru untuk
                                    bergabung ke workspace dan mengatur
                                    role mereka (admin/planner/operator)

  Aktor                             Admin

  **Skenario**                      

  Kondisi Awal                      Menampilkan halaman detail Workspace

  **Aksi Aktor**                    **Reaksi Sistem**

  1\. Aktor memilih tab "Members"   2\. Menampilkan daftar member yang
                                    sudah tergabung beserta role-nya.

  3\. Aktor menekan tombol "Invite  4\. Menampilkan form undangan.
  Member"                           

  5\. Menginput email address user  6\. Melakukan validasi format email
  yang diundang dan memilih role    dan memeriksa apakah user sudah
  (Admin/Planner/ Operator).        terdaftar di sistem.

  7\. Menekan tombol "Send          8\. Menyimpan data undangan ke
  Invitation"                       database dengan status "pending".

                                    9\. Mengirimkan email undangan ke
                                    alamat yang dituju dengan link
                                    konfirmasi.

                                    10\. Menampilkan notifikasi sukses
                                    dan memperbarui daftar pending
                                    invitations.

  Kondisi Akhir                     Undangan terkirim dan menunggu
                                    konfirmasi dari user yang diundang
  -----------------------------------------------------------------------

  : []{#_Toc218781110 .anchor}Tabel 3. 8 Skenario Use Case Kelola
  Anggota & Undangan

  -----------------------------------------------------------------------
  **Identifikasi**                  
  --------------------------------- -------------------------------------
  Nomor                             UC-07

  Nama                              Manage Roles & Permissions (Kelola
                                    Role & Permission)

  Tujuan                            Mengkonfigurasi permission secara
                                    granular untuk custom role

  Deskripsi                         Admin membuat atau mengedit role
                                    dengan mengatur permission spesifik
                                    untuk setiap resource
                                    (read/write/delete)

  Aktor                             Admin

  **Skenario**                      

  Kondisi Awal                      Menampilkan Dashboard Admin

  **Aksi Aktor**                    **Reaksi Sistem**

  1\. Aktor memilih menu "Roles &   2\. Menampilkan daftar role yang
  Permissions"                      tersedia (default roles + custom
                                    roles).

  3\. Aktor menekan tombol "Create  4\. Menampilkan form pembuatan role
  Role" atau memilih role existing  baru.
  untuk di-edit                     

  5\. Menginput nama role (misal:   6\. Menampilkan daftar permission
  "Warehouse Manager") dan          yang dapat diatur dalam bentuk
  deskripsi.                        checklist atau toggle untuk setiap
                                    resource (Users, Workspaces,
                                    Containers, Products, Plans, Items).

  7\. Mengatur permission untuk     8\. Melakukan validasi kelengkapan
  setiap resource (centang read,    data.
  write, delete sesuai kebutuhan).  

  9\. Menekan tombol "Save"         10\. Menyimpan konfigurasi role dan
                                    permission ke database (tabel roles
                                    dan role_permissions).

                                    11\. Menampilkan notifikasi sukses
                                    dan memperbarui daftar role.

  Kondisi Akhir                     Role baru tersedia dan dapat
                                    di-assign ke user melalui UC-06
  -----------------------------------------------------------------------

  : []{#_Toc218781111 .anchor}Tabel 3. 9 Skenario Use Case Kelola Role &
  Permission

  -----------------------------------------------------------------------
  **Identifikasi**                  
  --------------------------------- -------------------------------------
  Nomor                             UC-08

  Nama                              Create Shipment Plan (Buat Rencana
                                    Pengiriman)

  Tujuan                            Membuat dokumen rencana muatan baru
                                    dan menentukan jenis kontainer

  Deskripsi                         Planner menginisialisasi rencana
                                    pengiriman dengan memilih profil
                                    kontainer yang akan digunakan

  Aktor                             Planner

  **Skenario**                      

  Kondisi Awal                      Menampilkan Dashboard Planner

  **Aksi Aktor**                    **Reaksi Sistem**

  1\. Aktor memilih menu \"Buat     2\. Menampilkan form pembuatan
  Rencana Baru\" (Create New Plan)  rencana pengiriman.

  3\. Menginput Kode Rencana (Plan  4\. Melakukan validasi kelengkapan
  Code) dan memilih Jenis Kontainer data input.
  dari dropdown.                    

  5\. Menekan tombol \"Simpan\"     6\. Mengambil data detail kontainer
  (Create)                          (dimensi & berat maks) dari Master
                                    Data.

                                    7\. Menyimpan data rencana baru ke
                                    database sekaligus menyalin data
                                    kontainer sebagai snapshot (agar
                                    perubahan data kontainer di masa
                                    depan tidak merusak rencana ini).

                                    8\. Menampilkan halaman detail
                                    rencana untuk mulai memasukkan
                                    barang.

  Kondisi Akhir                     Menampilkan halaman detail rencana
                                    pengiriman (status: Draft)
  -----------------------------------------------------------------------

  : []{#_Toc218781112 .anchor}Tabel 3. 10 Skenario Use Case Membuat
  Rencana Pengiriman

  -----------------------------------------------------------------------
  **Identifikasi**                   
  ---------------------------------- ------------------------------------
  Nomor                              UC-09

  Nama                               Manage Items (Kelola Barang)

  Tujuan                             Menambahkan daftar barang yang akan
                                     dimuat ke dalam rencana pengiriman

  Deskripsi                          Planner memasukkan detail barang
                                     (dimensi, berat, jumlah) ke dalam
                                     sistem, baik secara manual maupun
                                     dari katalog

  Aktor                              Planner

  **Skenario**                       

  Kondisi Awal                       Menampilkan halaman detail Rencana
                                     Pengiriman (Shipment Plan Detail)

  **Aksi Aktor**                     **Reaksi Sistem**

  1\. Aktor menekan tombol "Tambah   2\. Menampilkan form input barang.
  Barang" (Add Item)                 

  3\. Aktor memilih produk dari      4\. Sistem mengambil data dimensi
  Katalog Master dan memasukkan      dan berat dari Master Data produk
  jumlah (Qty), atau memilih "Custom tersebut untuk mengisi form
  Input" untuk barang yang tidak ada otomatis. Jika Custom Input, user
  di katalog.                        harus mengisi semua field manual.

  5\. Menekan tombol "Simpan Item"   6\. Sistem melakukan validasi data
                                     (angka dimensi/berat harus positif).

                                     7\. Sistem melakukan Snapshot Data
                                     (menyalin atribut dimensi dan berat
                                     saat ini ke tabel transaksi) untuk
                                     menjaga integritas data historis.

                                     8\. Menambahkan barang ke dalam
                                     daftar tabel rencana di layar.

  Kondisi Akhir                      Barang berhasil ditambahkan ke dalam
                                     daftar rencana pengiriman
  -----------------------------------------------------------------------

  : []{#_Toc218781113 .anchor}Tabel 3. 11 Skenario Use Case Kelola
  Barang

  -----------------------------------------------------------------------
  **Identifikasi**                    
  ----------------------------------- -----------------------------------
  Nomor                               UC-10

  Nama                                Calculate Load (Kalkulasi Muatan)

  Tujuan                              Mendapatkan susunan barang 3D yang
                                      paling optimal

  Deskripsi                           Sistem memproses daftar barang dan
                                      dimensi kontainer menggunakan
                                      algoritma 3D Bin Packing dengan
                                      simulasi gravitasi dan pengecekan
                                      stabilitas

  Aktor                               Planner

  **Skenario**                        

  Kondisi Awal                        Menampilkan detail rencana
                                      pengiriman dengan daftar barang
                                      yang sudah terisi

  **Aksi Aktor**                      **Reaksi Sistem**

  1\. Aktor dapat mengatur parameter  2\. Menampilkan form konfigurasi
  kalkulasi (opsional) dengan menekan parameter (Strategi Packing,
  tombol "Configure"                  Stability Ratio, Gravity
                                      Simulation).

  3\. Aktor menekan tombol "Hitung    4\. Mengambil seluruh data dimensi
  Optimasi" (Calculate/Optimize)      barang dan spesifikasi kontainer
                                      dari database.

                                      5\. Menjalankan Calculation Engine
                                      (Algoritma 3D Bin Packing) dengan
                                      parameter yang dipilih.

                                      6\. Menghitung posisi koordinat (x,
                                      y, z) dan rotasi untuk setiap
                                      barang dengan mempertimbangkan
                                      gravitasi dan stabilitas tumpukan.

                                      7\. Menyimpan hasil perhitungan ke
                                      tabel plan_results dan
                                      plan_placements.

                                      8\. Menampilkan ringkasan hasil
                                      (Total Berat, Persentase Volume
                                      Terpakai, Jumlah Barang Termuat).

  Kondisi Akhir                       Status rencana berubah menjadi
                                      \"Calculated\" dan hasil
                                      perhitungan tersimpan
  -----------------------------------------------------------------------

  : []{#_Toc218781114 .anchor}Tabel 3. 12 Skenario Use Case Kalkulasi
  Muatan

  -----------------------------------------------------------------------
  **Identifikasi**                    
  ----------------------------------- -----------------------------------
  Nomor                               UC-11

  Nama                                View 3D Result (Lihat Visualisasi
                                      3D)

  Tujuan                              Memvisualisasikan hasil penataan
                                      barang di dalam kontainer

  Deskripsi                           Menampilkan model 3D interaktif
                                      yang menunjukkan posisi setiap
                                      barang sesuai hasil algoritma
                                      dengan fitur step-by-step playback

  Aktor                               Planner, Operator

  **Skenario**                        

  Kondisi Awal                        Proses kalkulasi (UC-10) telah
                                      selesai berhasil

  **Aksi Aktor**                      **Reaksi Sistem**

  1\. Aktor memilih tab "Visualisasi  2\. Mengambil data koordinat
  3D"                                 penempatan (placements) dari
                                      database.

                                      3\. Merender model 3D kontainer dan
                                      kotak-kotak barang sesuai dimensi
                                      dan warnanya menggunakan Three.js.

  4\. Aktor melakukan interaksi       5\. Aplikasi merespons interaksi
  (memutar, memperbesar, zoom in/out) visual secara real-time.

  6\. Aktor menggunakan kontrol step  7\. Sistem memfilter tampilan
  playback (tombol play/ pause,       barang berdasarkan step_number dan
  slider step)                        menampilkan proses pemuatan secara
                                      bertahap sesuai urutan yang
                                      dihitung algoritma.

  Kondisi Akhir                       Menampilkan visualisasi susunan
                                      muatan yang interaktif dengan
                                      kontrol playback
  -----------------------------------------------------------------------

  : []{#_Toc218781115 .anchor}Tabel 3. 13 Skenario Use Case Lihat
  Visualisasi 3D

  -----------------------------------------------------------------------
  **Identifikasi**                    
  ----------------------------------- -----------------------------------
  Nomor                               UC-12

  Nama                                Export PDF Report (Ekspor Laporan
                                      PDF)

  Tujuan                              Mengekspor hasil kalkulasi dan
                                      instruksi pemuatan ke format PDF
                                      untuk digunakan operator di
                                      lapangan

  Deskripsi                           Sistem menghasilkan dokumen PDF
                                      yang berisi detail rencana
                                      pengiriman, daftar placement, dan
                                      snapshot visualisasi 3D

  Aktor                               Planner, Operator

  **Skenario**                        

  Kondisi Awal                        Menampilkan halaman detail shipment
                                      plan yang sudah dikalkulasi (UC-10)

  **Aksi Aktor**                      **Reaksi Sistem**

  1\. Aktor menekan tombol "Export to 2\. Sistem mengambil data rencana
  PDF"                                pengiriman, hasil kalkulasi, dan
                                      daftar placements dari database.

                                      3\. Sistem menghasilkan dokumen PDF
                                      yang berisi: informasi kontainer,
                                      daftar barang dengan koordinat,
                                      urutan pemuatan (step-by-step), dan
                                      gambar snapshot visualisasi 3D.

                                      4\. Sistem men-trigger download
                                      file PDF ke perangkat user.

  5\. Aktor menyimpan file PDF dan    6\. File PDF tersimpan di perangkat
  dapat membukanya dengan aplikasi    lokal dan dapat dibuka dengan PDF
  PDF reader                          reader untuk digunakan di lapangan.

  Kondisi Akhir                       Dokumen PDF instruksi pemuatan
                                      tersedia untuk dibawa operator ke
                                      lapangan
  -----------------------------------------------------------------------

  : []{#_Toc218781116 .anchor}Tabel 3. 14 Skenario Use Case Export PDF

  -----------------------------------------------------------------------
  **Identifikasi**                    
  ----------------------------------- -----------------------------------
  Nomor                               UC-13

  Nama                                Trial Access (Akses Trial)

  Tujuan                              Memungkinkan calon user mencoba
                                      platform dengan batasan tertentu
                                      sebelum melakukan registrasi penuh

  Deskripsi                           Guest user dapat mengakses platform
                                      dengan batasan maksimal 3 rencana
                                      pengiriman untuk mengevaluasi fitur
                                      sistem

  Aktor                               Guest

  **Skenario**                        

  Kondisi Awal                        Menampilkan landing page atau
                                      halaman login

  **Aksi Aktor**                      **Reaksi Sistem**

  1\. Aktor memilih opsi "Try for     2\. Sistem membuat akun trial
  Free" atau "Start Trial" d          temporary engan workspace personal
                                      otomatis.

                                      3\. Sistem memberikan akses token
                                      sementara dan mengarahkan ke
                                      dashboard.

  4\. Aktor dapat menggunakan fitur   5\. Sistem menampilkan banner
  platform (buat plan, tambah item,   notifikasi di header: "Trial Mode -
  kalkulasi, view 3D)                 X plans remaining" untuk
                                      mengingatkan batasan.

                                      6\. Setiap kali membuat plan baru,
                                      sistem mengurangi counter trial
                                      (maksimal 3 plan).

  7\. Setelah mencapai batas, aktor   8\. Sistem menampilkan prompt
  tidak dapat membuat plan baru       registrasi dan menawarkan opsi
                                      untuk claim data trial setelah
                                      registrasi penuh.

  Kondisi Akhir                       Guest berhasil mencoba platform dan
                                      dapat memutuskan untuk melakukan
                                      registrasi penuh
  -----------------------------------------------------------------------

  : []{#_Toc218781117 .anchor}Tabel 3. 15 Skenario Use Case Akses Trial

### **Sequence Diagram**

> Sequence diagram menggambarkan interaksi antar komponen sistem dalam
> urutan waktu (timeline). Diagram ini menunjukkan bagaimana objek-objek
> dalam sistem berkomunikasi satu sama lain untuk menyelesaikan suatu
> use case. Pada sub-bab ini, dijelaskan tiga alur utama yang
> merepresentasikan arsitektur teknis sistem: autentikasi & workspace
> context, pembuatan rencana dengan snapshot pattern, dan kalkulasi
> muatan dengan algoritma 3D bin packing.

1.  **Sequence Diagram Authentication & Workspace Context**

> ![](media/image5.png){width="5.251094706911636in"
> height="5.244397419072616in"}

[]{#_Toc218783001 .anchor}Gambar 3. 4 Sequence Diagram Authentication &
Workspace Context

> Diagram pada Gambar \\\_\\\_\\\_ menunjukkan proses autentikasi
> pengguna dan pengaturan konteks workspace. Alur dimulai ketika user
> menginput kredensial di frontend Next.js, yang kemudian mengirim
> request POST ke endpoint \`/api/v1/auth/login\` pada Go API. Backend
> melakukan validasi kredensial dengan membandingkan password hash
> menggunakan bcrypt, kemudian menghasilkan pasangan token JWT (access
> token & refresh token).
>
> Sistem juga secara otomatis mengambil daftar workspace yang dapat
> diakses oleh user beserta role dan permission-nya dari tabel
> \`workspace_members\` dan \`roles\`. Token refresh disimpan di
> database untuk keamanan, sementara access token dikirim ke frontend
> dan disimpan di localStorage. Frontend kemudian mengatur workspace
> context aktif dan mengarahkan user ke dashboard sesuai dengan role-nya
> (Admin/Planner/Operator). Mekanisme ini memastikan isolasi data
> multi-tenant dan kontrol akses berbasis role (RBAC) yang granular.

2.  **Sequence Diagram Create Plan**

> ![](media/image6.png){width="5.24162510936133in"
> height="4.011446850393701in"}

[]{#_Toc218783002 .anchor}Gambar 3. 5 Sequence Diagram Create Plan

> Diagram ini menggambarkan alur pembuatan rencana muat (load plan) yang
> dimulai ketika user mengisi form di frontend dengan memilih container
> dan menambahkan item-item barang, lalu mengirimkan request POST ke
> endpoint \`/api/v1/plans\`. Setelah melewati validasi JWT dan
> pengecekan permission \`plan:create\`, sistem akan melakukan
> pengecekan khusus untuk user Trial yang dibatasi maksimal 3 rencana
> muat. Service kemudian mengambil dimensi container dari master data
> (jika \`container_id\` disertakan) atau menggunakan dimensi custom
> dari request.
>
> Aspek penting dalam alur ini adalah implementasi Snapshot Pattern, di
> mana sistem menyalin seluruh dimensi container dan item (panjang,
> lebar, tinggi, berat) ke tabel \`load_plans\` dan \`load_items\` tanpa
> menyimpan foreign key reference ke tabel master. Pendekatan ini
> menjaga integritas data historis sehingga perubahan atau penghapusan
> data master tidak mempengaruhi rencana yang sudah dibuat. Setelah plan
> dan items tersimpan dengan status \`DRAFT\`, sistem dapat langsung
> menjalankan perhitungan packing jika opsi \`auto_calculate\`
> diaktifkan, yang akan berkomunikasi dengan microservice Python untuk
> algoritma 3D bin packing (detail pada Sequence Diagram 3). Seluruh
> operasi menyertakan filter \`workspace_id\` untuk isolasi
> multi-tenant.

3.  **Sequence Diagram Calculate Load & Visualize 3D**

> ![](media/image7.png){width="5.184806430446194in"
> height="6.051146106736658in"}

[]{#_Toc218783003 .anchor}Gambar 3. 6 Sequence Diagram Calculate Load &
Visualize 3D

> Diagram ini menunjukkan alur perhitungan packing dan visualisasi 3D
> yang dimulai ketika user menekan tombol \"Calculate\" pada halaman
> detail rencana muat, disertai dengan opsi strategi algoritma (seperti
> BestFitDecreasing atau Parallel) dan parameter gravity simulation.
> Frontend mengirimkan request POST ke endpoint
> \`/api/v1/plans/:id/calculate\` yang kemudian diproses oleh Plan
> Service untuk mengambil data plan dan items dari database. Service
> menyiapkan input untuk algoritma dengan mengekspansi setiap item
> berdasarkan quantity-nya, misalnya item dengan quantity 3 akan
> menghasilkan 3 instance item individual yang akan ditempatkan secara
> independen oleh algoritma.
>
> Algoritma bin packing 3D menggunakan library boxpacker3 yang
> menjalankan strategi seperti BestFitDecreasing untuk menentukan posisi
> optimal setiap item dalam container dengan mempertimbangkan constraint
> rotasi dan simulasi gravitasi untuk stabilitas. Hasil perhitungan
> berupa posisi koordinat (x, y, z), kode rotasi, dan step number untuk
> setiap item yang berhasil ditempatkan, serta informasi mengenai item
> yang tidak muat (unfits) jika container tidak mencukupi. Service
> kemudian menyimpan hasil ke database dengan menghapus result lama
> terlebih dahulu, lalu membuat record baru di tabel \`plan_results\`
> untuk statistik global (volume utilization, weight utilization,
> feasibility) dan bulk insert ke tabel \`plan_placements\` untuk detail
> posisi setiap item. Status plan diupdate menjadi \`COMPLETED\` jika
> semua item berhasil dimuat, atau \`PARTIAL\` jika ada item yang tidak
> muat.
>
> Frontend menerima response berupa CalculationResult yang berisi array
> placements, kemudian melakukan konversi data placement ke format yang
> dapat dirender oleh Three.js dengan memetakan koordinat database ke
> sistem koordinat 3D scene dan menginterpretasikan rotation_code
> menjadi orientasi mesh 3D. Three.js Viewer membuat scene yang terdiri
> dari container box sebagai wireframe, item-item sebagai colored meshes
> sesuai color_hex, serta kamera dengan orbit controls untuk navigasi.
> User dapat berinteraksi dengan step-by-step playback controls yang
> memfilter placements berdasarkan step_number, sehingga visualisasi
> menampilkan animasi urutan pemuatan item satu per satu sesuai dengan
> hasil algoritma, membantu operator memahami sequence optimal untuk
> loading fisik di lapangan.

### **Arsitektur Sistem**

![](media/image8.png){width="5.926410761154855in"
height="1.6911329833770778in"}

[]{#_Toc218783004 .anchor}Gambar 3. 7 Arsitektur Sistem

> Arsitektur sistem dirancang untuk menggambarkan topologi menyeluruh
> dari aplikasi berbasis web yang mengintegrasikan perencanaan logistik,
> algoritma bin packing 3D, dan visualisasi interaktif. Desain
> arsitektur ini mengadopsi pola layanan terdistribusi (distributed
> services) dengan REST API sebagai jembatan komunikasi antara frontend
> dan backend, serta mendukung arsitektur multi-tenant melalui mekanisme
> isolasi workspace. Sistem dibagi ke dalam tiga lapisan utama untuk
> memastikan modularitas, skalabilitas, dan pemisahan tanggung jawab
> (separation of concerns).

1.  Aktor dan Client (Frontend Layer) Sisi klien dibangun menggunakan
    framework Next.js 14 sebagai aplikasi web responsif yang melayani
    empat kategori pengguna, yaitu Admin untuk mengelola data master dan
    pengaturan workspace, Planner untuk membuat rencana muat dan
    melakukan simulasi 3D, Operator untuk membaca instruksi pemuatan
    secara bertahap, serta Guest (Trial User) yang dapat mengakses
    sistem dengan kapasitas terbatas tanpa proses registrasi penuh.
    Antarmuka pengguna menyediakan halaman manajemen workspace untuk
    kolaborasi tim berbasis workspace-based access control, serta modul
    visualisasi 3D berbasis Three.js yang memungkinkan pengguna memutar
    ulang proses pemuatan barang langkah per langkah (step-by-step
    playback) dengan rotasi kamera interaktif. Komunikasi antara
    frontend dan backend dilakukan secara sinkron menggunakan protokol
    HTTP melalui REST API yang menerapkan autentikasi berbasis JSON Web
    Token (JWT).

2.  Backend System (Server Layer) Sisi backend berfungsi sebagai pusat
    pemrosesan logika bisnis yang terdiri dari beberapa lapisan
    sub-modul. Lapisan pertama adalah Gateway (Controllers) yang
    bertindak sebagai pintu gerbang API untuk menangani permintaan
    klien, meliputi kontroler autentikasi, manajemen workspace,
    perencanaan, dan data master. Sebelum permintaan mencapai kontroler,
    request melewati Middleware Layer yang melakukan autentikasi JWT,
    validasi keanggotaan workspace (workspace context), serta
    pemeriksaan hak akses berbasis Role-Based Access Control (RBAC). Di
    bawahnya terdapat Service Layer yang memuat logika inti, yaitu:
    Snapshot Manager untuk menangani validasi input dan duplikasi data
    master (dimensi dan berat) ke tabel transaksi agar data historis
    permanen dan tidak terpengaruh oleh perubahan master data; serta
    Plan Service yang mengatur pembuatan rencana muat, penambahan item,
    dan pemicu kalkulasi. Untuk proses perhitungan bin packing, sistem
    backend menggunakan arsitektur microservice dengan memanggil layanan
    terpisah berbasis Python melalui HTTP POST ke endpoint \`/pack\`.
    Layanan packing ini menggunakan pustaka py3dbp yang telah di-fork
    dengan penambahan fitur simulasi gravitasi dan pemeriksaan
    stabilitas tumpukan (gravity simulation & stability checking). Hasil
    kalkulasi berupa koordinat posisi 3D (x, y, z), kode rotasi, dan
    nomor urut langkah pemuatan dikirim kembali ke backend Go untuk
    disimpan ke basis data menggunakan SQLC (type-safe SQL query
    generator). Backend Go dibangun menggunakan framework Gin dengan
    dukungan API documentation melalui Swagger.

3.  Database Layer Penyimpanan data dilakukan secara terpusat dan
    relasional menggunakan PostgreSQL. Tabel transaksional utama
    meliputi \`load_plans\` untuk header rencana pengiriman yang
    menyimpan snapshot dimensi kontainer, \`load_items\` untuk data
    barang dengan snapshot dimensi dan berat produk, \`plan_results\`
    untuk ringkasan hasil kalkulasi (utilization percentage, total
    weight, feasibility status), dan \`plan_placements\` untuk menyimpan
    koordinat spasial 3D (pos_x, pos_y, pos_z), kode rotasi, serta nomor
    urut langkah pemuatan. Selain itu, sistem multi-tenant didukung oleh
    tabel \`workspaces\` untuk entitas ruang kerja,
    \`workspace_members\` untuk keanggotaan pengguna dalam workspace,
    dan \`workspace_invites\` untuk mengelola undangan kolaborasi. Modul
    autentikasi menggunakan tabel \`users\` untuk data akun pengguna,
    \`roles\` dan \`permissions\` untuk RBAC, \`role_permissions\` untuk
    relasi many-to-many, serta \`refresh_tokens\` untuk manajemen sesi
    JWT yang aman. Data master disimpan dalam tabel \`containers\` dan
    \`products\` yang mencakup atribut dimensi fisik, berat maksimal,
    warna visualisasi, dan batasan tumpukan.

### **Perancangan Basis Data**

> Perancangan basis data dilakukan untuk mendefinisikan struktur
> penyimpanan data yang efisien, aman, dan mendukung integritas data
> historis. Sistem ini menggunakan basis data relasional (RDBMS) dengan
> penerapan Universally Unique Identifier (UUID) sebagai Primary Key
> pada setiap tabel. Penggunaan UUID dipilih untuk menjamin keunikan
> data di seluruh sistem terdistribusi dan meningkatkan keamanan
> referensi data.

![](media/image9.png){width="5.751065179352581in"
height="4.280055774278215in"}

[]{#_Toc218783005 .anchor}Gambar 3. 8 Perancangan Basis Data

> Skema basis data dirancang dengan memisahkan tanggung jawab antar
> modul (Separation of Concerns). Pada lapisan keamanan, modul
> autentikasi mengelola data pengguna melalui tabel \`users\` yang
> menyimpan kredensial akun, serta tabel \`roles\` dan \`permissions\`
> yang menerapkan Role-Based Access Control (RBAC) untuk membatasi hak
> akses antara Admin, Planner, Operator, dan Guest. Selain itu, tabel
> \`refresh_tokens\` digunakan untuk manajemen sesi login yang aman
> menggunakan standar JWT dengan mekanisme token rotation dan
> revocation. Untuk mendukung arsitektur multi-tenant, sistem
> menggunakan tabel \`workspaces\` sebagai entitas ruang kerja yang
> mengisolasi data antar organisasi, tabel \`workspace_members\` untuk
> mengelola keanggotaan pengguna dengan role spesifik di dalam workspace
> (owner, admin, member, viewer), serta tabel \`workspace_invites\`
> untuk menangani proses undangan kolaborasi dengan status pending,
> accepted, atau expired. Di sisi manajemen referensi, modul data master
> menyimpan spesifikasi standar melalui tabel \`containers\` untuk data
> kontainer dan tabel \`products\` untuk katalog barang, yang mencakup
> atribut dimensi fisik, berat, warna visualisasi, serta batasan
> tumpukan.
>
> Inti dari pemrosesan data terletak pada modul kalkulasi yang
> menerapkan aturan bisnis Snapshot Data. Tabel \`load_plans\` berfungsi
> sebagai header rencana pengiriman dan menyimpan snapshot dimensi
> kontainer (length_mm, width_mm, height_mm, max_weight_kg) agar tidak
> terpengaruh oleh perubahan data master di masa depan, sementara tabel
> \`load_items\` menyimpan detail barang yang akan dikirim beserta
> snapshot atribut dimensi dan berat dari produk. Tabel \`load_items\`
> didesain untuk menduplikasi atribut dimensi dan berat dari data master
> saat rencana dibuat, bukan sekadar menyimpan referensi kunci asing.
> Mekanisme ini bertujuan menjaga integritas data historis agar riwayat
> pengiriman masa lalu tidak berubah meskipun data master produk
> mengalami pembaruan. Hasil perhitungan algoritma bin packing yang
> dieksekusi oleh microservice Python disimpan dalam tabel
> \`plan_results\` untuk ringkasan utilisasi volume
> (volume_utilization_pct), total berat yang dimuat
> (total_loaded_weight_kg), dan status kelayakan pemuatan (is_feasible),
> serta tabel \`plan_placements\` untuk menyimpan koordinat spasial
> (pos_x, pos_y, pos_z), kode rotasi (rotation_code), dan nomor urut
> langkah pemuatan (step_number) guna keperluan visualisasi 3D bertahap
> di antarmuka pengguna.
>
> Seluruh operasi database menggunakan migrasi skema yang dikelola
> melalui Goose migration tool untuk memastikan konsistensi struktur
> database di berbagai environment (development, staging, production).
> Query database diimplementasikan menggunakan SQLC yang menghasilkan
> type-safe Go code dari file SQL, sehingga mengurangi risiko kesalahan
> runtime dan meningkatkan performa compile-time checking. Untuk menjaga
> isolasi data antar workspace, setiap query transaksional secara
> otomatis menyertakan filter \`workspace_id\` melalui middleware
> context, sehingga mencegah kebocoran data antar organisasi dalam
> sistem multi-tenant. Sistem juga membatasi pengguna trial (guest)
> dengan maksimal 3 load plans per workspace untuk mengelola kapasitas
> infrastruktur dan mendorong konversi ke akun berbayar.

#  IMPLEMENTASI & PENGUJIAN

## Lingkungan Implementasi

> Langkah-langkah implementasi merupakan aspek yang sangat krusial untuk
> memastikan perangkat lunak *Load & Stuffing Calculator* dapat berjalan
> sesuai dengan rancangan sistem. Implementasi ini bertujuan untuk
> menerjemahkan kebutuhan logistik dan algoritma perhitungan muatan ke
> dalam aplikasi berbasis *web* yang dapat diakses oleh admin
> operasional maupun pengguna umum.
>
> Lingkungan implementasi aplikasi ini memanfaatkan integrasi antara
> perangkat keras yang memadai untuk komputasi dan perangkat lunak
> pendukung pengembangan modern. Berikut adalah spesifikasi perangkat
> keras dan perangkat lunak yang digunakan:

1.  Perangkat keras (Hardware) yang digunakan dalam pengembangan dan
    pengujian aplikasi ini adalah sebagai berikut:

    a.  Processor : *Intel® Core™ i5*

    b.  Memory (RAM) : 8.00 GB

    c.  System type : *64-bit Operating System, x64-based processor*

    d.  Storage : SSD 512 GB

2.  Perangkat lunak (Software) yang digunakan dalam pembuatan dan
    pengoperasian aplikasi ini adalah sebagai berikut:

    a.  Sistem Operasi : Windows 11 Home

    b.  Database : *PostgreSQL*

    c.  Backend Programming : *Go* (Golang) v1.20+

    d.  Code Editor : Visual Studio Code

    e.  Web Browser : Google Chrome

## Pembahasan Hasil Implementasi

> Setelah tahap pengembangan kode program selesai, implementasi sistem
> dilakukan untuk memvalidasi apakah logika perhitungan muatan dan
> desain antarmuka yang telah dirancang mampu menjawab kebutuhan
> pengguna dalam meminimalkan void space pada kontainer. Implementasi
> ini mencakup penyajian visual dari sistem, mulai dari manajemen data
> profil kontainer hingga simulasi hasil stuffing.
>
> Antarmuka aplikasi dibangun dengan konsep responsif agar mudah
> digunakan (user-friendly). Berikut adalah pembahasan mengenai tampilan
> dan fungsi dari setiap halaman utama yang telah berhasil
> diimplementasikan:

1.  **Tampilan Halaman Login**

> Halaman Login (Masuk) merupakan gerbang utama keamanan sistem yang
> berfungsi untuk memverifikasi identitas pengguna sebelum memberikan
> akses ke dalam aplikasi Load & Stuffing Calculator. Implementasi
> halaman ini dirancang dengan antarmuka yang bersih dan fokus pada
> fungsi autentikasi.
>
> ![](media/image10.png){width="3.8268471128608925in"
> height="3.6224420384951883in"}

[]{#_Toc218847794 .anchor}Gambar 4. 1 Tampilan Halaman Login

> Berdasarkan Gambar 4.1, halaman login menampilkan formulir autentikasi
> yang mewajibkan pengguna untuk memasukkan Username dan Password yang
> telah terdaftar. Sistem dilengkapi dengan mekanisme validasi di mana
> kredensial yang dimasukkan akan dicocokkan dengan data yang tersimpan
> di database. Jika data sesuai, pengguna akan diarahkan ke halaman
> dashboard; namun jika salah, sistem akan menolak akses demi menjaga
> keamanan data logistik. Di bagian bawah formulir juga tersedia tautan
> Create one untuk memfasilitasi pengguna baru yang belum memiliki akun.

2.  **Tampilan Halaman Register**

> Halaman Registrasi (Create Account) diimplementasikan untuk
> memfasilitasi pendaftaran pengguna baru yang ingin menggunakan sistem.
> Proses pendaftaran dirancang secara bertahap (multi-step) untuk
> memastikan kelengkapan data tanpa membebani pengguna dengan formulir
> yang terlalu panjang dalam satu tampilan.
>
> ![](media/image11.png){width="3.6240791776027996in"
> height="3.705184820647419in"}

[]{#_Toc218847795 .anchor}Gambar 4. 2 Tampilan Halaman Register

> Seperti terlihat pada Gambar 4.2, tahap pertama registrasi (Step 1 of
> 2) fokus pada pengumpulan informasi akun dasar. Pengguna diminta untuk
> melengkapi kolom Username, Email, dan Password. Antarmuka ini
> dirancang untuk memvalidasi format email dan keamanan password secara
> langsung. Tombol Continue berfungsi untuk menyimpan data sementara dan
> mengarahkan pengguna ke langkah selanjutnya dari proses pendaftaran.
> Pemisahan langkah ini bertujuan untuk meningkatkan pengalaman pengguna
> (user experience) agar proses pembuatan akun terasa lebih ringan dan
> terorganisir.

3.  **Tampilan Halaman Dashboard**

> Halaman Dashboard ini dirancang sebagai pusat kontrol utama (main hub)
> yang muncul pertama kali setelah admin berhasil melakukan login.
> Halaman ini bertujuan untuk memberikan gambaran ringkas mengenai
> status operasional sistem serta menyediakan akses cepat ke berbagai
> fitur manajemen muatan.

[]{#_Toc218847796 .anchor}Gambar 4. 3 *Tampilan Halaman Dashboard*

> Seperti ditunjukkan pada Gambar 4.3, halaman dashboard terdiri dari
> beberapa elemen antarmuka yang tersusun secara terstruktur untuk
> memudahkan navigasi pengguna. Di sisi kiri terdapat Menu Sidebar yang
> berisi daftar modul utama aplikasi, mulai dari Container Profiles
> untuk pengaturan spesifikasi kontainer, Product Catalog untuk data
> barang, hingga Execution Logs untuk pemantauan riwayat proses.
>
> Pada bagian utama (main content), terdapat elemen Quick Actions yang
> dirancang untuk menampilkan jalan pintas (shortcuts) menuju
> tugas-tugas yang paling sering dilakukan oleh peran pengguna tersebut.
> Di bawahnya, terdapat elemen Recent Activity yang menyajikan log
> aktivitas sistem secara real-time, seperti notifikasi inisialisasi
> sistem atau riwayat akses pengguna. Fitur ini berfungsi penting bagi
> admin untuk memantau jejak audit dan memastikan bahwa setiap aktivitas
> yang terjadi dalam sistem Load & Stuffing terpantau dengan baik.

4.  **Tampilan Halaman Container Profiles**

> Halaman Container Profiles berfungsi sebagai modul manajemen data
> master untuk spesifikasi kontainer. Mengingat perhitungan algoritma
> stuffing sangat bergantung pada batasan dimensi ruang, halaman ini
> memungkinkan admin untuk mendefinisikan berbagai tipe kontainer yang
> dimiliki perusahaan secara presisi.
>
> ![](media/image13.png){width="5.174146981627296in"
> height="2.5127832458442696in"}

[]{#_Toc218847797 .anchor}Gambar 4. 4 Tampilan Halaman Container
Profiles

> Berdasarkan Gambar 4.4, halaman ini menampilkan daftar kontainer dalam
> bentuk tabel yang memuat informasi krusial seperti Nama Kontainer,
> Dimensi (Panjang x Lebar x Tinggi dalam milimeter), dan Berat Maksimum
> (kg). Sistem menyediakan fitur pencarian (filter names) untuk
> memudahkan admin menemukan tipe kontainer tertentu. Selain itu,
> tombol + New Container di pojok kanan atas memungkinkan pengguna untuk
> menambah varian kontainer baru, semisal kontainer ukuran custom atau
> tipe High Cube, guna memastikan fleksibilitas perhitungan kalkulator.

5.  **Tampilan Halaman Producting Catalog**

> Modul Product Catalog dirancang untuk mengelola inventaris barang atau
> kargo yang akan dimuat. Halaman ini menjadi basis data utama bagi
> spesifikasi barang, sehingga pengguna tidak perlu memasukkan dimensi
> barang berulang kali setiap membuat perencanaan pengiriman.
>
> ![](media/image14.png){width="5.539931102362205in"
> height="2.6987674978127734in"}

[]{#_Toc218847798 .anchor}Gambar 4. 5 Tampilan Halaman Producting
Catalog

> Seperti terlihat pada Gambar 4.5, setiap entitas produk
> direpresentasikan dengan atribut detail meliputi Dimensi (LxWxH),
> Berat, dan Kode Warna. Fitur kode warna (Color) diimplementasikan
> untuk memberikan visualisasi yang intuitif pada hasil simulasi 3D
> nanti, membedakan satu jenis barang dengan barang lainnya. Tabel ini
> juga dilengkapi dengan menu aksi di sisi kanan setiap baris untuk
> melakukan penyuntingan atau penghapusan data produk jika terjadi
> perubahan spesifikasi fisik barang.

6.  **Tampilan Halaman All Shipments**

> Halaman All Shipments menyajikan rekapitulasi seluruh aktivitas
> perencanaan pengiriman yang telah dibuat di dalam sistem. Halaman ini
> berfungsi sebagai log riwayat yang memungkinkan pengguna memantau
> status setiap batch pengiriman.
>
> ![](media/image15.png){width="5.537566710411198in"
> height="2.680437445319335in"}

[]{#_Toc218847799 .anchor}Gambar 4. 6 Tampilan Halaman All Shipments

> Pada Gambar 4.6, data pengiriman ditampilkan dalam format kartu (card
> view) yang informatif. Setiap kartu memuat ID Pengiriman, Tanggal
> Pembuatan, Jumlah Item, serta Status pengiriman (seperti COMPLETED
> atau PARTIAL). Desain visual ini memudahkan admin operasional untuk
> dengan cepat mengidentifikasi pengiriman mana yang sudah selesai
> diproses dan mana yang masih dalam tahap perencanaan atau tertunda.

7.  **Tampilan Halaman Visualisasi 3D**

> Halaman Detail Hasil Optimasi menyajikan laporan komprehensif mengenai
> rencana pemuatan yang telah diproses oleh algoritma. Halaman ini
> berfungsi sebagai alat validasi utama bagi pengguna untuk meninjau
> efisiensi penggunaan ruang kontainer dan melihat simulasi penataan
> barang secara tiga dimensi sebelum eksekusi fisik dilakukan.
>
> ![](media/image16.png){width="5.5218153980752405in"
> height="2.731541994750656in"}

[]{#_Toc218847800 .anchor}Gambar 4. 7 Tampilan Halaman Visualisasi 3D

> Pada Gambar 4.15, antarmuka menampilkan panel 3D Simulation di bagian
> tengah yang memvisualisasikan kontainer 40ft High Cube beserta susunan
> barang di dalamnya. Model 3D ini bersifat interaktif dan dilengkapi
> dengan kontrol playback (slider langkah 210/210) di bagian bawah, yang
> memungkinkan pengguna melihat urutan pemuatan barang secara bertahap
> (step-by-step) demi memastikan stabilitas tumpukan. Di sisi kanan,
> terdapat panel Summary yang memberikan metrik efisiensi secara
> real-time, menunjukkan bahwa rencana ini memanfaatkan 52.8% volume
> kargo dan 25.7% kapasitas berat. Selain itu, panel Cargo Breakdown
> menyajikan grafik donat dan legenda daftar barang (seperti Euro
> Pallet, Large Crate, dan Medium Box) yang diberi kode warna spesifik,
> memudahkan admin untuk mengidentifikasi distribusi jenis barang dalam
> visualisasi tumpukan tersebut.

8.  **Tampilan Halaman Create Shipment**

> Tampilan Halaman Create Shipment
>
> Halaman Create Shipment merupakan antarmuka inti dari sistem
> kalkulator ini. Pada halaman ini, pengguna melakukan input parameter
> simulasi untuk menghitung susunan muatan yang optimal.
>
> ![](media/image17.png){width="5.444400699912511in"
> height="2.7378838582677165in"}

[]{#_Toc218847801 .anchor}Gambar 4. 8 Tampilan Halaman Create Shipment

> Gambar 4.7 memperlihatkan formulir perencanaan yang terbagi menjadi
> beberapa blok logis. Pada bagian kiri, pengguna mengisi Shipment Title
> dan memilih tipe kontainer melalui blok Container Selection (bisa
> memilih preset yang sudah ada atau dimensi custom). Pada bagian kanan,
> terdapat blok Add Items di mana pengguna dapat memasukkan barang dari
> katalog (From Catalog) atau input manual. Sistem secara otomatis
> menghitung ringkasan total berat dan volume pada blok Summary di
> bagian bawah secara real-time sebelum pengguna menekan tombol Create
> Shipment & Calculate untuk memproses algoritma stuffing.

9.  **Tampilan Halaman Loading Shipments**

> Tampilan Halaman Loading Instructions
>
> Halaman Loading Instructions berfungsi sebagai antarmuka output akhir
> yang ditujukan bagi tim operasional lapangan. Halaman ini dirancang
> untuk menampilkan instruksi langkah demi langkah mengenai urutan
> pemuatan barang ke dalam kontainer.
>
> ![](media/image18.png){width="5.502244094488189in"
> height="2.6791896325459317in"}
>
> []{#_Toc218847802 .anchor}Gambar 4. 9 Tampilan Halaman Loading
> Shipments
>
> Gambar 4.8 menunjukkan tampilan awal (empty state) dari halaman
> instruksi muat ketika belum ada rencana pengiriman yang dipilih atau
> diproses. Pesan \"No planned shipments available\" mengindikasikan
> bahwa sistem siap menerima instruksi baru dari perencana (planner).
> Ketika sebuah pengiriman telah diproses, halaman ini akan menampilkan
> visualisasi urutan barang yang harus dimasukkan pertama kali hingga
> terakhir, guna memastikan stabilitas muatan dan optimalisasi ruang
> sesuai perhitungan algoritma.

10. **Tampilan Halaman Members**

> Halaman Members dirancang untuk memfasilitasi kolaborasi tim dalam
> satu ruang kerja (workspace). Fitur ini memungkinkan pemilik workspace
> untuk melihat dan mengelola siapa saja yang memiliki akses terhadap
> data pengiriman dan profil kontainer di dalam organisasi mereka.
>
> ![](media/image19.png){width="5.472761373578303in"
> height="2.6266338582677164in"}
>
> []{#_Toc218847803 .anchor}Gambar 4. 10 Tampilan Halaman Members
>
> Sebagaimana ditampilkan pada Gambar 4.9, halaman ini menyajikan daftar
> anggota aktif dalam format tabel yang memuat informasi Username,
> Email, Role (peran), dan tanggal bergabung (Added). Kolom Role
> berfungsi krusial untuk membatasi hak akses pengguna, misalnya sebagai
> Owner yang memiliki akses penuh atau Planner yang hanya fokus pada
> perencanaan. Di sisi kanan, terdapat tombol aksi Remove berwarna merah
> yang memungkinkan admin untuk mencabut akses anggota jika diperlukan,
> serta tombol Add Member di pojok kanan atas untuk menambahkan personil
> baru.

11. **Tampilan Halaman Invites**

> Halaman Invites berfungsi sebagai mekanisme keamanan untuk penambahan
> anggota baru. Alih-alih mendaftar secara bebas dan masuk ke sembarang
> workspace, pengguna baru harus melalui proses undangan yang dikirimkan
> oleh admin.
>
> ![](media/image20.png){width="5.5276126421697285in"
> height="2.690317147856518in"}
>
> []{#_Toc218847804 .anchor}Gambar 4. 11 Tampilan Halaman Invites
>
> Gambar 4.10 memperlihatkan antarmuka pengelolaan undangan yang
> statusnya masih tertunda (pending). Tabel ini mencatat detail penting
> seperti alamat Email penerima, Role yang ditawarkan, siapa
> pengundangnya (Invited By), serta masa berlaku undangan tersebut
> (Expires). Fitur ini menjamin bahwa akses ke dalam sistem hanya
> diberikan kepada pihak yang terautorisasi secara resmi.

12. **Tampilan Halaman User Management**

> Halaman User Management adalah fitur level administrator sistem (Super
> Admin) yang digunakan untuk memantau seluruh basis pengguna yang
> terdaftar dalam aplikasi Load & Stuffing Calculator, lintas workspace.
>
> ![](media/image21.png){width="5.579463035870516in"
> height="2.7804593175853016in"}
>
> []{#_Toc218847805 .anchor}Gambar 4. 12 Tampilan Halaman User
> Management
>
> Berdasarkan Gambar 4.11, halaman ini menampilkan inventarisasi akun
> pengguna secara global. Informasi yang disajikan meliputi Username,
> Email, dan Role sistem (seperti admin sistem, planner, atau founder).
> Fitur ini memungkinkan pengelola aplikasi untuk melakukan audit
> terhadap pertumbuhan pengguna dan memverifikasi peran setiap akun yang
> terdaftar di dalam database.

13. **Tampilan Halaman Workspaces Management**

> Sistem ini menganut arsitektur multi-tenancy, di mana data antar
> pengguna atau organisasi dipisahkan dalam wadah yang disebut
> Workspace. Halaman Workspace Management berfungsi untuk mengelola
> ruang-ruang kerja tersebut.
>
> ![](media/image22.png){width="5.577604986876641in"
> height="2.7634667541557305in"}
>
> []{#_Toc218847806 .anchor}Gambar 4. 13 Tampilan Halaman Workspaces
> Management
>
> Pada Gambar 4.12, terlihat daftar workspace yang aktif beserta
> atributnya, yaitu Nama, Tipe (apakah Personal untuk perorangan atau
> Organization untuk perusahaan), serta informasi Owner (pemilik).
> Pemisahan data melalui workspace ini sangat vital untuk menjamin
> privasi dan keamanan data logistik masing-masing perusahaan pengguna
> agar tidak tercampur satu sama lain.

14. **Tampilan Halaman Roles**

> Untuk menjamin keamanan sistem yang granular, aplikasi ini menerapkan
> konsep Role-Based Access Control (RBAC). Halaman Roles dan Permissions
> digunakan untuk mendefinisikan kebijakan akses tersebut secara teknis.
>
> ![](media/image23.png){width="5.5526935695538056in"
> height="2.69206583552056in"}
>
> []{#_Toc218847807 .anchor}Gambar 4. 14 Tampilan Halaman Roles
>
> Gambar 4.13 menunjukkan dua tampilan konfigurasi keamanan. Pada bagian
> Roles, sistem mendefinisikan tingkatan jabatan pengguna seperti Admin,
> Operator, Planner, hingga Trial User, lengkap dengan deskripsi
> tanggung jawabnya. Sementara itu, pada bagian Permissions, sistem
> memetakan hak akses teknis secara spesifik menggunakan format
> key-value (contoh: container:create untuk izin membuat kontainer, atau
> dashboard:read untuk izin melihat dasbor). Implementasi ini memastikan
> bahwa setiap pengguna hanya dapat melakukan aksi sesuai dengan
> kewenangannya masing-masing, meminimalkan risiko penyalahgunaan fitur.

15. **Tampilan Halaman *Permissions***

> Halaman Permissions merupakan lapisan keamanan teknis yang mendasari
> sistem Role-Based Access Control (RBAC). Jika halaman Roles
> mendefinisikan \"siapa\" (jabatan), maka halaman Permissions
> mendefinisikan \"apa\" (tindakan spesifik) yang boleh dilakukan oleh
> sistem.
>
> ![](media/image24.png){width="5.422410323709537in"
> height="2.7460487751531057in"}
>
> []{#_Toc218847808 .anchor}Gambar 4. 15 Tampilan Halaman Permissions
>
> Berdasarkan Gambar 4.14, sistem ini menerapkan kontrol akses yang
> sangat granular (terperinci). Setiap izin direpresentasikan dalam
> format key-value yang teknis namun deskriptif pada kolom Name, seperti
> container:create untuk izin membuat data kontainer, dashboard:read
> untuk izin melihat dasbor, atau invite:\* untuk akses penuh terhadap
> fitur undangan. Kolom Description memberikan penjelasan fungsi dari
> setiap kode izin tersebut agar mudah dipahami oleh administrator.
> Keberadaan halaman ini membuktikan bahwa aplikasi Load & Stuffing
> Calculator dirancang dengan standar keamanan modular, di mana hak
> akses dapat dikonfigurasi secara spesifik tanpa harus merombak kode
> program.

## Pengujian Code Coverage

[]{#_Toc218847809 .anchor}Gambar 4. 16 Pengujian Code Coverage

#  KESIMPULAN & SARAN

## Kesimpulan

> Berdasarkan hasil analisis, perancangan, implementasi, dan pengujian
> yang telah dilakukan pada sistem Load & Stuffing Calculator, dapat
> ditarik beberapa kesimpulan sebagai berikut:

1.  Keberhasilan Implementasi Algoritma Optimasi Penelitian ini berhasil
    merancang dan mengimplementasikan layanan backend yang
    mengintegrasikan bahasa pemrograman Go dan layanan mikro Python
    untuk menjalankan algoritma 3D Bin Packing. Sistem mampu memproses
    data dimensi produk dan kontainer untuk menghasilkan skema pemuatan
    yang meminimalkan void space. Algoritma yang diterapkan telah mampu
    memperhitungkan batasan fisik seperti dimensi, berat maksimum, dan
    orientasi barang, sehingga hasil kalkulasi lebih presisi
    dibandingkan metode estimasi manual.

2.  Visualisasi 3D yang Interaktif dan Informatif Pembangunan antarmuka
    pengguna berbasis web menggunakan Next.js dan Three.js berhasil
    menyajikan hasil perhitungan algoritma dalam bentuk visualisasi 3D
    yang interaktif. Fitur step-by-step playback yang dikembangkan
    memungkinkan pengguna (Planner dan Operator) untuk melihat urutan
    pemuatan barang secara bertahap, mulai dari lapisan pertama hingga
    terakhir. Hal ini menjawab kebutuhan operasional akan panduan visual
    yang jelas untuk meminimalisir kesalahan penataan di lapangan.

3.  Platform SaaS dengan Manajemen Multi-Tenant Sistem berhasil
    dikembangkan sebagai platform berbasis Software as a Service (SaaS)
    yang mendukung arsitektur multi-tenancy. Melalui penerapan isolasi
    workspace dan fitur keamanan Role-Based Access Control (RBAC),
    sistem mampu memfasilitasi kolaborasi tim dengan hak akses yang
    terkelola dengan baik (Admin, Planner, Operator). Mekanisme snapshot
    data yang diterapkan juga terbukti efektif dalam menjaga integritas
    data historis pengiriman meskipun terjadi perubahan pada data master
    di kemudian hari.

4.  Efisiensi Dokumentasi Operasional Fitur ekspor laporan ke format PDF
    yang menyertakan instruksi visual dan daftar koordinat barang
    berhasil diimplementasikan. Fitur ini mempercepat proses
    administrasi dan serah terima informasi dari perencana ke tim
    operasional gudang, menggantikan instruksi manual yang rentan
    terhadap human error.

## Saran

> Meskipun sistem Load & Stuffing Calculator telah berhasil dibangun dan
> memenuhi kebutuhan fungsional utama, terdapat beberapa aspek yang
> dapat dikembangkan lebih lanjut untuk meningkatkan kinerja dan cakupan
> sistem di masa mendatang:

1.  Pengembangan Variasi Bentuk Barang Saat ini, algoritma kalkulasi
    difokuskan pada barang berbentuk kotak atau kubus (box). Untuk
    pengembangan selanjutnya, disarankan agar algoritma ditingkatkan
    untuk mendukung bentuk barang yang lebih kompleks, seperti silinder
    (drum), bentuk L, atau barang dengan bentuk tidak beraturan, guna
    memperluas jangkauan penggunaan aplikasi di berbagai industri
    logistik.

2.  Integrasi dengan Sistem Eksternal (ERP) Sistem saat ini beroperasi
    sebagai platform mandiri. Kedepannya, disarankan untuk mengembangkan
    fitur integrasi API publik (Public API) yang memungkinkan sistem ini
    terhubung secara langsung dengan perangkat lunak Enterprise Resource
    Planning (ERP) atau Warehouse Management System (WMS) yang sudah ada
    di perusahaan, sehingga data stok barang dapat ditarik secara
    otomatis tanpa input manual.

3.  Pengembangan Aplikasi Mobile Native Meskipun antarmuka web saat ini
    sudah responsif, pengembangan aplikasi mobile native (Android/iOS)
    untuk sisi Operator lapangan akan sangat bermanfaat. Aplikasi mobile
    dapat dilengkapi dengan fitur pemindai barcode atau QR Code untuk
    memvalidasi barang secara real-time saat proses pemuatan berlangsung
    ke dalam kontainer.

4.  Penambahan Parameter Kerapuhan Barang Disarankan untuk menambahkan
    parameter \"tingkat kerapuhan\" (fragility) pada data master produk.
    Algoritma dapat dikembangkan lebih lanjut untuk memastikan barang
    yang mudah pecah otomatis ditempatkan di lapisan paling atas (Top
    Load Only) demi keamanan muatan yang lebih terjamin.

#   DAFTAR PUSTAKA

\[1\] N. Triningsih, "Peran Terminal Petikemas Dalam Meningkatkan Daya
Saing Dan Efisiensi Logistik Maritim Indonesia," *J. Ilm. Wahana
Pendidikan, Desember*, vol. 2024, no. 24, pp. 304--311.

\[2\] S. Oktavia, "CENTRAL PUBLISHER PERAN TEKNOLOGI DALAM MENINGKATKAN
EFISIENSI OPERASIONAL PERUSAHAAN LOGISTIK", \[Online\]. Available:
http://centralpublisher.co.id

\[3\] G. Mihu, "University of New Hampshire Scholars ' Repository BIN
PACKING PROBLEM - AN EFFICIENCY ANALYSIS BETWEEN COMMERCIAL CONTAINER
LOADING AND OPEN-SOURCE OPTIMIZATION ALGORITHMS," 2024.

\[4\] N. T. X. Hoa, "Optimizing Container Fill Rates for the Textile and
Garment Industry Using a 3D Bin Packing Approach," *HighTech Innov. J.*,
vol. 5, no. 2, pp. 462--478, 2024, doi: 10.28991/HIJ-2024-05-02-017.

\[5\] ( Edhy Poerwandono, F. Fiddin, and E. Poerwandono, "Implementasi
3D Bin Packing Problem Menggunakan Algoritma Tabu Search," *J. Sains dan
Teknol.*, vol. 5, no. 1, pp. 477--482, 2023.

\[6\] S. Oktavia, "Peran Teknologi Dalam Meningkatkan Efisiensi
Operasional Perusahaan Logistik," *J. Cent. Publ.*, vol. 1, no. 9, pp.
1049--1056, 2023.

\[7\] D. Widodo, "ANALISIS PROSES STUFFING CONTAINER MUATAN KARET PADA
PT. SAMUDERA INDONESIA CABANG JAMBI," Politeknik Pelayaran Sumatera
Barat, 2023.

\[8\] K. Heßler, T. Hintsch, and L. Wienkamp, "A fast optimization
approach for a complex real-life 3D Multiple Bin Size Bin Packing
Problem," *Eur. J. Oper. Res.*, vol. 327, no. 3, pp. 820--837, 2025,
doi: 10.1016/j.ejor.2025.05.016.

\[9\] ( Edhy Poerwandono, F. Fiddin, and E. Poerwandono, "Implementasi
3D Bin Packing Problem Menggunakan Algoritma Tabu Search," *J. Sains dan
Teknol.*, vol. 5, no. 1, pp. 477--482, 2023, doi:
10.55338/saintek.v5i1.1393.

\[10\] V. Kartiko and P. N. Primandari, "Media Pengenalan Peti Kemas
Logistik Menggunakan Augmented Reality Berbasis Android," *JTIM J.
Teknol. Inf. dan Multimed.*, vol. 5, no. 2, pp. 134--148, 2023, doi:
10.35746/jtim.v5i2.369.

\[11\] K. Back-end and R. Api, "PENGEMBANGAN APLIKASI BACK-END
E-COMMERCE MENGGUNAKAN REST API a . Pengembangan aplikasi front-end E-
Commerce," vol. 11, pp. 7--13, 2024.

\[12\] L. Fernando and M. M. Engel, "Comparative Performance
Benchmarking of WebSocket Libraries on Node . js and Golang," vol. 9,
no. 4, pp. 2051--2060, 2025.

\[13\] D. Manusia, R. Simbulan, and J. Aryanto, "Jurnal Indonesia :
Manajemen Informatika dan Komunikasi Implementasi REST API Web Services
pada Aplikasi Sumber Abstrak Jurnal Indonesia : Manajemen Informatika
dan Komunikasi," vol. 5, no. 1, pp. 552--560, 2024.

\[14\] R. Farhandika, M. Sabariah, and M. Adrian, "Penerapan Arsitektur
REST API pada Aplikasi Backend Manajemen Informasi Fakultas Industri
Kreatif (MI- FIK) Universitas Telkom," *Log. J. Penelit. Inform.*, vol.
2, 2024, doi: 10.25124/logic.v2i1.7530.

\[15\] S. V. Salunke and A. Ouda, "A Performance Benchmark for the
PostgreSQL and MySQL Databases," *Futur. Internet*, vol. 16, no. 10,
2024, doi: 10.3390/fi16100382.

\[16\] J. Han and Y. Choi, "Analyzing Performance Characteristics of
PostgreSQL and MariaDB on NVMeVirt," 2024, \[Online\]. Available:
http://arxiv.org/abs/2411.10005

\[17\] S. Dalimunthe, E. H. Putra, M. Arif, F. Ridha, and P. C. Riau,
"RESTFUL API SECURITY USING JSON WEB TOKEN ( JWT ) WITH HMAC-SHA512
ALGORITHM IN SESSION," vol. 8, no. 1, pp. 81--94, 2023.

\[18\] I. Bismoputro, F. Al Huda, and A. H. Brata, "PENGEMBANGAN SINGLE
PAGE APPLICATION BERBASIS REACTJS UNTUK USAHA PERCETAKAN ONLINE ( STUDI
KASUS : GLOBAL GRAFIKA )," vol. 1, no. 1, pp. 1--10, 2017.

\[19\] S. Azhariyah and M. Mukhlis, "Framework CSS : Tailwind CSS Untuk
Front-End Website Store PT . XYZ," vol. 2, pp. 30--36, 2023.
