# Pemetaan SELA v5 → Rantaya

Baseline: prototype SELA commit `a9c43a25a116e695df3a2f73ac978521ed092acd`, dibekukan di `docs/reference/sela-v5`. `web/src/prototype.css` berasal dari CSS baseline; selector disesuaikan dengan DOM Svelte. Aset foto, font dan denah berasal dari prototype, tetap berwarna. Backend tidak memakai simulasi localStorage prototype.

| SELA | Route Rantaya | Susunan yang dipertahankan |
| --- | --- | --- |
| home | / | Ruang obrolan; Jelajah/Diikuti; filter/sort; composer entry; post, mention/event, foto, aksi |
| agenda | /agenda | Agenda panggung; period tabs; kategori; tanggal kanan; poster kiri dan informasi kanan |
| date | /tanggal | Preset, rentang lintas bulan, calendar, Tampilkan agenda, back membatalkan |
| communities | /ruang | Temukan ruangmu; copy, kategori, kartu dan Ikuti/Masuk ke ruang |
| community | /ruang/[slug] | Cover, identitas, deskripsi, statistik, follow; enam tab dan isi masing-masing |
| event | /event/[slug] | Detail event, crumb, summary/Maps, tiga tab; media opsional, lineup, denah, jadwal, organizer |
| event actions | detail event | Panel desktop; bar mobile di atas empat nav; Beli tiket/Tulis ulasan, simpan, share |
| venue | /event/[slug]/venue | Layout venue; zoom −/+/Reset, geser, panduan tanpa pemilihan kursi |
| product | /merch/[id] | Merchandise; crumb; gambar dan informasi dua kolom; ukuran; Pesan ke penyelenggara; event asal |
| post | /post/[id] | Percakapan; post utuh dengan konteks event; komentar dan balasan inline |
| review | /event/[slug]/ulasan | Bagikan pengalamanmu/Edit ulasan; Kirim ulasan sticky; event asal; isi; hadir tanpa badge otomatis |
| saved | /disimpan | Event lalu Postingan; jumlah; empty state; aside |
| notifications | /notifikasi | Kabar yang ingin kamu ikuti; Tandai dibaca; daftar; Atur notifikasi |
| search | /cari | Cari di Rantaya; input aktif; hasil komunitas/event/obrolan; keadaan kosong |
| profile | /profil | Profil kamu/Akun pengelola; kartu identitas; Profil/Notifikasi; form dan switches |
| manage | /kelola | Kelola komunitas; Lihat ruang; Ringkasan/Event/Postingan/Merchandise/Ulasan/Pengaturan |
| dashboard content | /kelola/* | Empat metrik, grafik tujuh hari, quick actions; tabel event/merch; kartu kelola post; review/profil |
| shared shell | seluruh route | Sidebar gelap desktop, kota/profil header, discovery aside, empat bottom nav mobile, back/task layout |

Nama produk menjadi Rantaya; istilah komunitas tetap “ruang”. Nama Teater Ruang tidak diganti. Kode rute hash prototype dipindahkan ke URL SvelteKit yang dapat dibuka langsung dan dirender SSR.

## Perilaku nyata dan tambahan yang sudah diminta

- Login Google/demo menggantikan login akun contoh prototype; akun dan otorisasi ditetapkan Go.
- Transaksi & tiket, rekening, pembayaran dan scanner merupakan alur tiket manual yang disepakati. Tiga tab dashboard tambahan serta akses transaksi membuat fitur itu dapat ditemukan. Statistik penjualan disetujui memakai DB, terpisah dari klik.
- Tanggal seed dibuat relatif agar pembelian/check-in dapat dicoba lokal; hitungan follow/vote/review berasal dari row DB. Angka fixture prototype tidak dipakai sebagai angka nyata.
- Caption prototipe/simulasi diganti keterangan demo yang jujur; aplikasi lokal tidak mengklaim login/transaksinya hanya simulasi browser.
- Halaman detail memakai satu h1 semantik; heading pengantar tetap memiliki ukuran/posisi prototype melalui h2.page-title.
- Kontak yang belum punya kanal diarahkan ke percakapan mention pengelola. Merch tetap katalog + kanal eksternal. Email digest hanya preferensi tersimpan; belum mengirim email.
- Pengelola hanya membalas ulasan event sendiri sesuai izin API; badge berasal dari check-in. Tanggal kartu review adalah tanggal event asal.

## Bukti dan batas

Header, navigasi, controls, composer dan feed beranda mobile diukur terhadap render asli pada 360/390 px dengan toleransi 0,1 px; fixture acuan tidak diganti. Flow dan layout diuji pada 360/390/1440, plus CTA tablet900. Seluruh 68 state route diperiksa pada tiga viewport: HTTP200, satu h1, tanpa overflow/exception JS. Kartu/merch/profil diperiksa lagi setelah perbaikan terakhir.

Ini membuktikan komponen dan state yang tercatat, bukan sertifikasi screenshot identik untuk setiap kombinasi konten/data/perangkat. Data asli pengguna dan panjang teks dapat mengubah tinggi konten. Perbandingan offline memakai baseline; jangan merancang ulang halaman karena hanya daftar fitur yang terlihat. Hasil terbaru ada di TESTING.md dan verification.json.
