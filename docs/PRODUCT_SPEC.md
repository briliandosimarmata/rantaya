# Spesifikasi produk v0.1

## Masalah, tujuan dan aktor

Event seni lokal sulit ditemukan, informasi tiket terpencar, dan hubungan komunitas sering berhenti setelah pertunjukan. Rantaya menggabungkan percakapan, profil pengelola, agenda dan pemesanan tiket. Launch awal Karawang; filter kota memungkinkan Jakarta/Bandung dan kota lain melalui data API. Kota adalah pilihan eksplisit, bukan GPS/jarak meter.

Aktor: tamu membaca konten publik; customer mengikuti ruang/berinteraksi/membeli tiket; organizer memublikasikan konten dan mengelola event miliknya; admin meninjau laporan. Login customer/organizer terpisah. Tidak ada switch role dalam composer. Tidak ada multi-staff organizer dalam v0.1.

## Fitur dan perilaku

| Fitur          | Implementasi awal                                          | Kriteria yang harus dijaga                                                                  |
| -------------- | ---------------------------------------------------------- | ------------------------------------------------------------------------------------------- |
| Beranda        | SSR feed Jelajah dan Langganan                             | Jelajah dibatasi kota; popular berdasarkan jumlah upvote lalu waktu, bukan global           |
| Agenda         | Mendatang/selesai, kategori, rentang tanggal, paginasi20   | Filter tanggal mempertahankan kota/kategori; event draft hanya terlihat pemilik             |
| Lokasi         | Tombol header global + pilihan kota                        | Dapat diganti dari tab mana pun; cookie lokal setahun; tidak meminta GPS                    |
| Composer       | Halaman penuh, isi wajib1–5000, judul opsional≤120         | Fokus ke tulisan, @ inline, panel terbatas viewport dan scroll internal                     |
| Attachment     | Satu JPG/PNG/WebP≤8 MB                                     | Foto tetap berwarna; pengelola juga dapat mencantumkan tautan tiket/merch                   |
| Mention        | Event publik/organizer                                     | Chip tertaut; customer post muncul pada Mention ruang yang dituju; organizer menerima notif |
| Upvote         | Satu vote/account/post, dapat dibatalkan                   | State disimpan DB, sort popular tetap dalam kota                                            |
| Komentar       | Teks1–1500; satu tingkat balasan                           | Balasan harus dari post yang sama; posting yang disembunyikan tidak dapat dikomentari       |
| Follow         | Ikuti organizer                                            | Feed langganan hanya post resmi; notif event baru mengikuti preferensi                      |
| Bookmark/share | Post dan event                                             | Bookmark privat; share Web Share API atau salin tautan                                      |
| Ruang          | Linimasa/Mention/Agenda/Merch/Ulasan/Tentang               | Linimasa post resmi, Mention post customer; deskripsi mobile selebar konten                 |
| Event          | Deskripsi, kategori, venue, Maps, bahasa/usia/durasi, sesi | Lineup satu baris horizontal; media opsional tanpa placeholder kosong                       |
| Trailer/denah  | YouTube atau mp4/webm; denah halaman penuh                 | Pengelola mengisi URL trailer atau gambar; denah bukan pemilihan kursi                      |
| Ticketing      | Transfer/wallet manual dan QR                              | Ikuti state machine PAYMENTS; capacity/idempotency server-side                              |
| Review         | Teks10–3000, satu per customer/event, dapat diubah         | Setelah seluruh event selesai; asal event ditampilkan di ruang; badge check-in saja         |
| Respons review | Pengelola membalas, customer vote membantu                 | Pengelola dapat melaporkan, tidak menghapus review                                          |
| Merch          | Katalog gambar/harga/varian/status/link                    | Tetap dapat dilihat setelah event; pembelian di kanal pengelola                             |
| Notifikasi     | In-app persistent, polling30 detik                         | Checkpoint order wajib; organizer payment hanya setelah proof/reproof                       |
| Reminder       | Tiket approved, sesi mulai dalam 24 jam                     | Sekali/account/sesi, preferensi reminders; scheduler setiap menit                           |
| Profil         | Nama, kota, bio, foto, minat, prefs                        | Onboarding halaman penuh yang bisa skip, bukan modal wajib                                  |
| Moderasi       | Report post/review dan dashboard admin                     | Hide publik melalui admin, bukan delete pengelola                                           |

## Dashboard pengelola

Ringkasan transaksi/tiket/check-in/pendapatan approved; daftar event dan formulir lengkap; post sendiri (pin/hapus); queue pembayaran dan detail bukti; master bank/wallet; katalog merchandise; ulasan dan respons; profil ruang; scanner QR dengan pilihan event/manual fallback.

Session berisi label, starts_at, ends_at, harga dan kapasitas. Satu order1–6 tiket, bukan keranjang campuran sesi. Event berbayar membutuhkan deskripsi ruang yang lengkap dan minimal satu metode aktif sebelum terbit. Jadwal/harga sesi yang pernah dipesan tidak dapat diedit; kapasitas tidak boleh turun di bawah jumlah reservasi. Sesi dengan order tidak bisa dihapus. Unpublish menghentikan pembelian baru, tidak membatalkan order/tiket yang ada.

## Cakupan dan batas eksplisit

Tidak ada realtime chat/DM, live streaming bawaan, pembayaran gateway, refund otomatis, kursi bernomor, native app, login OTP, email/push, ticket transfer/resale, multi-staff atau sinkronisasi offline. Tautan external live dapat dibagikan pada post pengelola, namun aplikasi tidak menyelenggarakan stream. Katalog merch belum punya order internal/pengiriman.

Verifikasi organizer belum memiliki alur dokumen/verifikasi mandiri; field verified disediakan tetapi data demo tidak mengklaim komunitas sungguhan. Filter kategori menggunakan Teater/Musik/Komedi/Lainnya. Customer umum dapat memosting tanpa membuat event.

## Acceptance walkthrough

1. Tamu mengganti kota, membuka event, lokasi Maps, lineup, flyer/trailer/denah bila ada.
2. Customer masuk, skip onboarding, isi post tanpa judul, mention organizer/event, unggah gambar, lihat pada Mention ruang.
3. Follow ruang → post resmi tampil di Langganan; upvote/bookmark tetap setelah reload; komentar/reply tertaut benar.
4. Pilih sesi/qty → order tersimpan → tutup/reload → lanjutkan dari history → pilih rekening → proof → correction/reproof → approval → QR.
5. Organizer hanya melihat order miliknya; QR event lain/sudah dipakai ditolak; review setelah event menyertakan badge hanya bila check-in.
6. Merch tetap dapat dibuka pada event selesai; customer menghubungi kanal order organizer.
