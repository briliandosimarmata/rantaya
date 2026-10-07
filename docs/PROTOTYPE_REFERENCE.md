# SELA — rancangan UI/UX peluncuran awal

Nama SELA merupakan nama kerja, belum keputusan merek. Prototipe memakai komunitas, acara, profil, statistik, dan produk fiktif. Tanggal acuan contoh: 7 Oktober 2026.

## Arah produk

SELA adalah ruang komunitas seni yang tetap berjalan di antara event. Pengguna mengikuti penyelenggara, membaca proses kreatif, memulai diskusi, menemukan pertunjukan, dan memberi ulasan. Agenda dan merchandise berada di dalam hubungan komunitas tersebut.

Peluncuran berfokus pada teater dan musik di Karawang. Pilihan Jakarta memperlihatkan keadaan “segera hadir”, tanpa mengisi kota tersebut dengan acara rekaan yang seolah benar-benar tersedia.

## Referensi dan keputusan desain

| Referensi | Pola yang diadaptasi |
| --- | --- |
| [Circle](https://circle.so/) | Ruang komunitas permanen, feed, postingan dan komentar, serta aktivitas yang berlanjut di luar event. |
| [Luma Discover](https://luma.com/discover) | Penemuan event menurut kota dan kategori, serta mengikuti kalender atau komunitas. |
| [GOERS](https://www.goersapp.com/) | Hubungan antara daftar acara, profil penyelenggara, follow, dan pengalaman penonton. |

Palet, nama, komponen, dan susunan SELA merupakan keputusan desain untuk produk ini. Noto Sans dipilih sendiri untuk keterbacaan bahasa Indonesia, bukan klaim tentang font yang dipakai platform referensi.

## Sistem visual

“Golden design principle” diinterpretasikan sebagai aturan warna 60–30–10, dengan rasio emas sekitar 1,618 untuk kolom utama dan kolom pendukung pada desktop. Keduanya dipakai sebagai pedoman komposisi, bukan kewajiban menghitung setiap piksel.

| Warna dasar | Nilai | Peran |
| --- | --- | --- |
| Putih keabu-abuan | `#F4F5F7` | Sekitar 60% bobot visual: kanvas, kartu, ruang kosong, dan latar formulir. |
| Arang gelap | `#172322` | Sekitar 30%: navigasi, teks, garis, panel informasi, dan tombol utama. |
| Lime | `#D3E96B` | Sekitar 10%: pilihan aktif, penekanan, dan aksi penting. |

Warna teks sekunder serta garis menggunakan transparansi arang, sehingga tidak menambah warna dasar. Foto mempertahankan warna aslinya tanpa filter monokrom. Aturan tiga warna berlaku pada komponen UI, bukan foto. Thumbnail pengisi menggunakan foto orang/band fiktif sebagai contoh. Informasi status tetap disampaikan lewat tulisan atau ikon.

- **Font:** Noto Sans, regular dan bold, dimuat sebagai aset lokal; Arial/sans-serif sebagai fallback.
- **Teks utama:** 1 rem; ukuran font menggunakan rem agar mengikuti pengaturan teks browser.
- **Hierarki:** judul halaman, judul kartu, isi, lalu metadata. Metadata penting seperti harga dan jadwal tidak disembunyikan di tooltip.
- **Grid desktop:** navigasi samping, kolom utama, dan kolom pendukung. Kolom utama terhadap pendukung menggunakan `1.618fr : 1fr` selama ruang cukup.
- **Tablet:** kolom pendukung disembunyikan; informasi yang sama tetap bisa dicapai melalui Agenda atau detail event.
- **HP:** satu kolom, navigasi bawah, dan tombol tiket langsung pada detail event.
- **Komponen:** sudut kartu 16 px, tombol sekitar 10 px, garis tipis, jarak yang konsisten, dan gerakan minimal.
- **Aksesibilitas dasar:** label formulir, fokus keyboard yang jelas, elemen tombol/tautan asli, kontrol sentuh utama minimal 44 px, dan pengurangan animasi saat diminta perangkat.

## Halaman dan alur

| Halaman | Isi dan interaksi |
| --- | --- |
| Beranda | Jelajah/Diikuti, filter teater/musik, urutan terbaru/paling didukung, composer, upvote, komentar, simpan, dan bagikan. |
| Ruang komunitas | Daftar komunitas lokal dan tombol follow. |
| Detail komunitas | Linimasa resmi pengelola, Mention dari pengguna yang menyebut pengelola/eventnya, Agenda, Merchandise, Ulasan, dan Tentang. Pengelola bisa menyematkan post resmi miliknya. |
| Percakapan | Isi postingan, komentar, dan balasan bertingkat. |
| Agenda | Acara mendatang/selesai, kategori, dan kalender pada halaman khusus untuk memilih satu tanggal atau rentang tanggal. Pilihan kota tersedia secara global di header. |
| Detail event | Informasi acara, akses lokasi Google Maps, pengisi dan perannya, media opsional, viewer venue, penyelenggara, jadwal/sesi, obrolan terkait, ulasan, simpan, bagikan, dan akses tiket. |
| Merchandise | Katalog tetap pada profil komunitas, detail produk, varian, harga, ketersediaan, dan alur menuju kanal pemesanan. |
| Disimpan | Event dan postingan yang dipilih pengguna. |
| Notifikasi | Kabar komunitas dan event, tautan ke konten, serta tandai dibaca. |
| Profil | Identitas pengguna, edit profil, serta pengaturan kabar event, balasan, pengingat, dan rangkuman email. |
| Panel pengelola | Ringkasan, kelola event, postingan, merchandise, balas ulasan, dan edit identitas komunitas. |
| Moderasi platform | Laporan konten, status peninjauan, dan contoh pemeriksaan profil penyelenggara. |
| Pencarian | Halaman khusus dengan input aktif dan hasil langsung untuk komunitas, event, dan obrolan. |

### Mengikuti komunitas

Pengguna membuka Ruang komunitas, melihat profil, lalu memilih Ikuti. Komunitas tersebut masuk ke daftar yang diikuti dan postingannya muncul di tab Diikuti. Pengguna dapat berhenti mengikuti dari tombol yang sama.

### Memulai percakapan

Composer berupa halaman penuh yang langsung memfokuskan isi tulisan. Judul opsional ditambahkan lewat tombol, gambar lewat pemilih berkas, dan mention pengelola/event lewat tombol @ atau mengetik @. Tidak ada dropdown identitas, ruang, atau jenis postingan. Pengelola memakai composer serupa dengan tambahan tautan pembelian tiket/merchandise dan teks tombol opsional. Identitas otomatis mengikuti akun yang login. Draft tersimpan terpisah per akun, termasuk saat pengguna kembali ke halaman sebelumnya. Pengguna lain bisa memberi dukungan, komentar, serta balasan inline.

### Menemukan dan mengakses event

Pengguna memfilter agenda, membuka detail, memilih sesi bila festival memiliki beberapa pertunjukan, lalu melanjutkan ke kanal tiket. Harga, tanggal, dan tempat terlihat sebelum tindakan itu. Dalam prototipe, kelanjutan ini disimulasikan tanpa membuka transaksi nyata atau menerbitkan tiket.

### Menulis review

Review dibuka setelah acara selesai. Pengguna menyatakan hadir dan menulis pengalaman. Satu review melekat pada satu event, lalu ditampilkan juga pada profil penyelenggara. Edit dan balasan tampil konsisten di kedua tempat. Nama event dan tanggal selalu terlihat. Pernyataan kehadiran sendiri tidak menghasilkan badge terverifikasi; badge pada data contoh merupakan ilustrasi untuk event yang kelak memiliki bukti check-in.

### Membeli merchandise setelah event

Pengguna membuka katalog komunitas atau postingan produk, memilih varian, lalu menuju kanal penyelenggara. Katalog tidak hilang ketika event selesai. Prototipe hanya mensimulasikan kelanjutan pemesanan tersebut.

## Keadaan dan batas prototipe

- Keadaan kosong tersedia untuk komunitas yang belum diikuti, filter tanpa hasil, katalog tanpa produk, review sebelum acara, serta kota yang belum diluncurkan.
- Umpan balik diberikan setelah follow, simpan, kirim komentar/review, perubahan profil, dan laporan konten.
- Hapus postingan meminta konfirmasi. Ulasan publik memakai alur laporan dan balasan; panel pengelola tidak menyediakan penghapusan review sepihak.
- Login pengguna dan pengelola berada di halaman terpisah dengan identitas, email, draft, serta preferensi terpisah. Login memakai akun contoh tanpa kata sandi. Belum ada autentikasi produksi, database bersama, verifikasi identitas sungguhan, pengiriman email, pembayaran, atau transaksi.
- Perubahan contoh disimpan di browser/perangkat yang sama. Ini bukan penyimpanan lintas pengguna atau lintas perangkat.
- Statistik adalah data simulasi. Klik tiket tidak dihitung sebagai penjualan.
- Foto dan merchandise merupakan aset ilustratif hasil generasi, bukan dokumentasi event Karawang yang nyata.

## Validasi

Pemeriksaan sintaks JavaScript lulus. Revisi terbaru diuji lewat 64 pemeriksaan logika interaksi, termasuk isolasi akun/preferensi/draft, posting tanpa judul, mention event yang muncul pada pengelola, hak kelola, review lintas halaman tanpa duplikasi, balasan inline, media opsional, penghapusan media, filter tanggal, pencarian, simulasi transaksi, escaping, dan kembali ke posisi scroll sebelumnya. Sebanyak 40 keadaan halaman/tab pengguna dan pengelola dapat dirender. Aset lokal, SVG layout, tiga warna dasar, dan struktur CSS juga diperiksa.

Pemeriksaan visual langsung di browser belum dilakukan karena lingkungan ini tidak menyediakan jalur pratinjau browser yang didukung untuk proyek statis. Dukungan desktop dan HP disiapkan melalui breakpoint responsif; tinjauan visual di perangkat nyata tetap diperlukan sebelum implementasi produksi.


## Revisi mobile berdasarkan feedback

- Navigasi bawah hanya Beranda, Agenda, Ruang, dan Disimpan. Profil diakses dari avatar header; halaman turunan memiliki tombol kembali.
- Lokasi berada di header yang tetap terlihat dan dapat diganti dari setiap halaman. Pilihan kota berupa menu kecil non-modal. Menutupnya tidak menghapus tulisan.
- Pencarian memakai halaman khusus. Filter tanggal memakai halaman kalender khusus dengan pilihan cepat dan rentang tanggal, sehingga agenda tidak dipenuhi formulir dua tanggal. Composer, review, editor event/produk, dan kelanjutan tiket/pemesanan memakai halaman penuh. Balasan komentar/review diperluas inline. Dialog tersisa untuk menu singkat, laporan, kontak contoh, dan konfirmasi penghapusan.
- Pengguna dan pengelola login melalui akun terpisah. Akses panel kelola dan post resmi hanya untuk pengelola. Pengguna tidak dapat menambahkan tombol pembelian atau menyematkan post resmi. Login pada prototipe bukan sistem keamanan produksi.
- Linimasa hanya memuat post resmi penyelenggara. Mention hanya memuat post pengguna yang menyebut penyelenggara atau event miliknya. Menyebut sebuah event menampilkan post itu juga di Obrolan event terkait. Review tetap fitur terpisah dan melekat pada event sumber.
- Pengisi acara disimpan sebagai nama dan peran bebas, sehingga dapat memuat aktor, sutradara, musisi, band, atau komedian tanpa membuat profil artis baru.
- Media event: trailer YouTube atau video MP4/WebM yang dapat diputar, lalu flyer sebagai fallback, lalu tidak ada section jika keduanya kosong. Layout hanya memiliki tombol bila pengelola menyertakan gambar. Viewer memakai halaman lebar dengan zoom/geser, tanpa memilih kursi.
- Di Balik Layar menggunakan contoh lineup, flyer landscape ilustratif, dan layout venue. Sore Mendengar memperlihatkan event tanpa media/layout. Semua nama pengisi dan denah adalah fiktif.
- Editor event memiliki informasi utama, daftar pengisi yang dapat ditambah/hapus, dan bagian media opsional yang dapat diperluas. Penghapusan trailer/flyer/layout menghilangkan bagian terkait pada detail event. Festival yang telah memiliki beberapa sesi mempertahankan sesi tersebut saat metadata diubah.


## Penyempurnaan mobile kedua

- Mention memakai autocomplete saat mengetik @. Daftar berada dekat caret tulisan dan memiliki scroll internal yang dibatasi area layar yang terlihat, termasuk saat keyboard membuka. Daftar tidak memerlukan input pencarian kedua. Fokus tetap pada tulisan; keyboard panah/Enter/Escape didukung. Jika hanya satu hasil, panel mendekat ke caret, bukan menyisakan ruang kosong setinggi daftar panjang.
- Tombol @ pada toolbar memasukkan @ di posisi kursor saat ini. Pemilihan mempertahankan teks sebelum/sesudah caret dan data mention. Email tidak membuka autocomplete. Hasil kosong diberi pesan singkat.
- Filter tanggal berada di sisi kanan baris kategori, seperti urutan pada Beranda, dengan ikon kalender kecil yang sejajar teks. Memilihnya membuka halaman kalender. Pilihan cepat tersedia untuk semua tanggal, hari ini, pekan ini, dan bulan ini. Dua tap memilih rentang, termasuk melewati pergantian bulan. Tombol kembali membatalkan pilihan yang belum diterapkan.
- Pengisi acara tampil dalam satu deretan horizontal dengan foto bulat, nama, dan peran. Deretan dapat digeser, memiliki scroll snap ringan, dan dapat difokuskan melalui keyboard. Pengelola bisa mengunggah/menghapus foto tiap pengisi pada editor; nama/peran tetap tersimpan. Data lama mendapatkan foto contoh sesuai nama, sedangkan pengisi baru tanpa foto memakai inisial.
- Aksi detail event memakai satu baris: Beli tiket/Tulis ulasan sebagai aksi utama, lalu tombol ikon simpan dan bagikan dengan label aksesibel. Pada HP, bar ini tetap terlihat di atas navigasi bawah dengan harga singkat. Ruang di akhir halaman ditambahkan agar isi tidak tertutup; tombol posting mengambang tidak menimpa bar tiket.
- Enam thumbnail berwarna tersedia sebagai ilustrasi profil fiktif. Foto bukan dokumentasi seniman atau band Karawang yang nyata.
- Validasi tambahan meliputi autocomplete, pilihan dengan Enter/Escape, penempatan saat viewport mengecil/bergeser, satu hasil dekat caret, kalender antarbulan, pembatalan pilihan, preset pekan Senin–Minggu, dan unggah/hapus foto pengisi. Uji posisi memakai ukuran viewport dalam harness logika; pengujian sentuhan, keyboard HP, dan tampilan browser perangkat nyata belum tersedia.

Rujukan implementasi: [MDN VisualViewport](https://developer.mozilla.org/en-US/docs/Web/API/VisualViewport) untuk posisi mengikuti area yang terlihat; [W3C APG Combobox](https://www.w3.org/WAI/ARIA/apg/patterns/combobox/) untuk navigasi saran dengan fokus yang tetap pada input. Composer tetap berupa textarea multiline, tidak mengklaim menjadi komponen combobox lengkap.


## Perbaikan detail ruang dan warna foto

Header detail ruang dipisahkan menjadi identitas, deskripsi, statistik, dan tombol follow. Di layar sampai 760 px, avatar dan nama memakai dua kolom; deskripsi serta statistik membentang selebar area konten; tombol Ikuti berada pada baris sendiri. Tidak ada paragraf panjang di kolom yang harus berbagi lebar dengan avatar dan tombol. Ukuran deskripsi 1 rem, line-height 1,75, dan kata panjang dapat membungkus tanpa keluar dari layar.

Semua filter grayscale dihapus dari foto postingan, sampul komunitas, poster event, produk, dan thumbnail pengisi. Foto unggahan tetap memakai warna asli. Aset contoh latihan, merchandise, potret, dan band diperbarui ke versi berwarna dengan subjek/komposisi yang dipertahankan. Overlay gelap hanya dipakai pada sampul/poster agar tulisan tetap terbaca; palet UI tetap tiga warna.

## Revisi akses lokasi event

Detail event menyediakan tombol Google Maps di bawah nama venue, sehingga dapat diakses dari semua tab event. Pengelola dapat memasukkan tautan lokasi opsional pada editor event. Tautan hanya menerima alamat HTTPS Google Maps, termasuk tautan pendek maps.app.goo.gl. Tanpa tautan dari pengelola, tombol menjadi pencarian nama venue dan kota di Google Maps; venue fiktif contoh tidak diberi pin seolah alamatnya telah terverifikasi. Mengosongkan kolom mengembalikan perilaku pencarian tersebut. Akses menggunakan tautan langsung tanpa peta embed atau API key.

Pemeriksaan sintaks dan sembilan pemeriksaan tambahan lulus: tautan resmi/pendek, penolakan alamat tidak sesuai, pencarian venue, field opsional pengelola, penyimpanan tanpa kehilangan lineup/jadwal, akses lintas tab event, penolakan perubahan invalid, penghapusan tautan, dan pembatasan hak pengguna. Pemeriksaan ini menggunakan logika serta HTML hasil render, tanpa pemeriksaan visual di browser.

Perubahan ini diperiksa lewat sintaks JavaScript, struktur CSS, rendering ketiga profil ruang, dan pemeriksaan aset berwarna. Pengujian visual langsung di browser HP belum tersedia.


## Identitas dan penggunaan sekarang

Dokumen di atas adalah rekaman prototype terakhir, dengan nama historis SELA. Nama produk sekarang **Rantaya**. Pertahankan istilah ruang untuk komunitas; jangan ganti nama organizer contoh Teater Ruang. Referensi ini menentukan susunan dan interaksi, bukan implementasi login/transaksi simulasi. Backend Go/PostgreSQL, transfer privat, approval dan QR single-use tetap sumber kebenaran.
