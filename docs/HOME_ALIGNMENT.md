# Acuan beranda mobile Rantaya

Revisi 7 Oktober 2026. Acuan: [SELA v5](https://sela-komunitas-seni.sandramoored074.chatgpt.site), source commit `a9c43a25a116e695df3a2f73ac978521ed092acd`. Prototype asli tidak diubah. Nama produk tetap **Rantaya**; logo Rantaya, identitas akun, tanggal, isi post dan angka interaksi mengikuti data aplikasi, bukan disalin menjadi data palsu dari prototype.

Ukuran di bawah diambil dari render prototype lokal dengan Noto Sans yang sama, pada viewport 360 dan 390 px. Toleransi pengujian koordinat/ukuran komponen utama 0,1 px pada browser uji. Ini menguji bagian yang tercantum, bukan klaim semua halaman identik pixel.

| Komponen | Acuan mobile |
| --- | --- |
| Header | Tinggi 76 px, margin kiri/kanan 18 px; garis bawah mengikuti lebar konten. Logo 27 px, disembunyikan pada ≤370 px sesuai prototype. |
| Kota | Tombol berbingkai, tinggi 44 px, teks 13 px bold, chevron 16 px; pin disembunyikan di mobile. “Jelajahi kota”, check kota aktif, “Segera hadir” untuk Jakarta. Escape/klik luar menutup pilihan. |
| Akses lain | Cari dan notifikasi 44 px, avatar profil 40 px lime; ikon stroke 1,65. Profil berada di header. Search memiliki label “Cari komunitas atau event”. |
| Navigasi bawah | Beranda / Agenda / Ruang / Disimpan. Tinggi 75 px, ikon 21 px, teks 12 px regular. Seluruh tombol aktif berlatar lime dengan radius 10 px; bukan hanya kotak di belakang ikon. Warna bertransisi singkat dan tap memiliki feedback. |
| Heading | “Karawang · Ruang seni lokal” dengan pin, “Ruang obrolan”, “Cerita dan percakapan dari komunitasmu.” Judul 1,812 rem, copy 14 px. Padding atas 25/bawah 23 px. |
| Feed | Jelajah / Diikuti, teks 14 px bold, gap 24 px, underline 3 px. Tab “Diikuti” tetap memakai parameter backend `following`. |
| Filter | Semua / Teater / Musik, kontrol 44 px, teks 12 px. Sort kanan: Terbaru / Paling didukung, tinggi 38 px, tanpa ikon pengaturan tambahan. |
| Composer beranda | “Mau cerita apa hari ini?”, avatar lime 34 px, ikon edit; padding 14 px, gap 10 px, tinggi 74 px, radius 16 px, margin bawah 20 px. Langsung menuju composer yang sudah ada. |
| Post | Bingkai radius 16 px; header padding 17/16/0 px, gap 10 px, avatar 42 px (square radius 12 px untuk pengelola). Author link 13 px; label pengelola dan metadata 12 px, urutan kota → tanggal → jam. Tombol “Pilihan postingan” dan akses “Laporkan konten”. |
| Isi post | Judul opsional 1,188 rem, isi 16 px/line-height 1,75; padding 14/16/17 px. Body juga menjadi tautan percakapan. Mention, kartu event tertaut, dan tautan pembelian mendahului foto. Event mention di luar agenda mendatang di-resolve secara terpisah tanpa mengubah pagination discovery. |
| Foto post | Memenuhi lebar kartu, berwarna, aspect ratio 1,55 pada mobile; tidak memakai inset/radius foto tambahan. |
| Aksi post | Upvote + angka, “N komentar”, spacer, simpan, bagikan. Tinggi sentuh 44 px, teks 12 px, ikon 17 px, gap 4 px, padding footer 9 px. Dukungan/simpan aktif lime. |
| Upvote | Optimistic feedback segera; pop ikon 280 ms dan transisi warna 160 ms. State/count dikonfirmasi API, bertahan setelah reload, dapat dibatalkan; error memulihkan state/count. Tombol mencegah request ganda ketika sibuk. Reduced motion menghilangkan animasi. |
| Tombol menulis | FAB edit 54 px, radius 17 px, kanan 20 px/bawah 94 px, shadow sesuai prototype. |
| Empty state | “Belum ada obrolan di sini”, “Ikuti komunitas atau mulai cerita pertamamu.”, tombol “Jelajah komunitas”. Tidak menambahkan blok event pada beranda mobile. |

## Checklist dan bukti

Checklist terverifikasi serta hasil terbaru tersedia di `evidence/home-alignment-checks.json`. Uji reproduksi ada di `web/tests/home.spec.ts`; ukuran referensi ada di `web/tests/fixtures/mobile-home-reference.json`. Uji shared header/composer/event ada di test yang sudah tersedia. Jangan memperbarui fixture acuan agar tes menjadi hijau tanpa perubahan prototype yang disetujui pengguna.

Checkpoint awal ini mencakup beranda. Build akhir juga memindahkan halaman lain sesuai VISUAL_PARITY.md dan memeriksa seluruh state route. Ketinggian CTA tiket disesuaikan dengan menu bawah 75 px dan padding detail event dipertahankan agar konten tidak tertutup.

Demo Sites memakai komponen produk yang sama. Upvote dan bookmark post di demo merupakan preferensi cookie per perangkat dan per role untuk menunjukkan state aktif, tanpa penyimpanan publik atau transaksi nyata. Customer dan organizer tetap terpisah. Reservasi, proof, approval dan check-in di Sites tetap ditolak; otoritas transaksi sebenarnya tetap Go/PostgreSQL.
