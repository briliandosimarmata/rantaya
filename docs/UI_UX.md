# UI/UX Rantaya — acuan wajib

Nama produk **Rantaya**. “Ruang” adalah istilah fitur komunitas. Cookie, draft key, database, Go module/binary dan slug Teater Ruang tetap kompatibel.

Susunan mengikuti [prototype terakhir](PROTOTYPE_REFERENCE.md) dan feedback pengguna berikut. Jangan menganggap keberadaan fitur sebagai izin merancang ulang susunan. Svelte/Go tetap implementasi nyata; simulasi auth/transaksi prototype tidak menggantikan backend.

## Sistem visual

Noto Sans regular/bold lokal. Canvas `#F4F5F7`, ink `#172322`, accent `#D3E96B`; border/muted memakai opacity ink. Semua foto berwarna. Main/aside desktop 1.618:1, mobile satu kolom. Body 16px, kontrol utama 44px. Poster agenda memakai overlay agar judul terbaca; flyer detail memakai contain agar utuh.

## Checklist kesesuaian

| Bagian | Perilaku/susunan wajib |
| --- | --- |
| Header | Kota global dan avatar profil; kembali di header halaman turunan. Profil bukan bottom tab. |
| Navigasi | Beranda/Agenda/Ruang/Disimpan, FAB membuat post. Composer, kalender, login/onboarding, checkout, detail transaksi, review, denah dan editor event/merch tanpa bottom nav. |
| Composer | Header ringkas dengan Posting kanan atas sticky. Avatar/identitas lalu isi autofocus; judul opsional setelah isi. Toolbar gambar/@/tautan pengelola/tambah judul/draft di bawah. Mobile tanpa bingkai kartu atau judul promosi. Tanpa dropdown jenis/identitas. |
| Mention | @ di caret, panel kontekstual fixed dengan scroll internal; batas visualViewport dan header composer aktual. Dapat muncul di atas caret bila ruang bawah habis. Ikon, nama, metadata, empty state. Arrow/Enter/Escape, satu tap; email tidak memicu. Hasil fetch lama/Escape tidak boleh membuka panel lagi. |
| Agenda | Daftar horizontal: poster kiri 130px desktop/90px mobile, informasi dan aksi kanan. Mendatang/Selesai/Semua agenda. Filter tanggal di kanan seperti sort beranda; kategori bisa digulir saat sempit. Label tanggal Indonesia, tanpa panel tanggal inline. |
| Kalender | Halaman khusus, preset hanya mengubah draft. Pilih hari/rentang lintas bulan. Buka bulan tanggal filter terpilih. Tampilkan agenda menerapkan; Back membatalkan. Semua tanggal juga menunggu Terapkan. Hari ini memakai Asia/Jakarta. |
| Pengelola | Cover/identitas, deskripsi full-width, statistik, follow. Mobile wajib identitas → deskripsi → statistik → Ikuti full-width. Linimasa resmi dan Mention customer. Statistik event/ulasan memakai + bila hasil mencapai page size 20; bukan total lengkap. |
| Event | Ringkasan, penyelenggara, tanggal/venue/Maps lalu Tentang/Obrolan/Ulasan. Trailer dapat diputar → flyer landscape → tanpa area media. Denah hanya bila tersedia, halaman penuh. |
| Lineup | Satu baris horizontal scroll-snap, foto 88px desktop/80px mobile, nama/peran. Region dapat difokuskan dan digulir panah keyboard. |
| Aksi tiket | Desktop panel harga/tiket/follow sticky di aside. Mobile fixed di atas bottom nav; harga di atas satu row Beli tiket/bookmark/share. Tablet saat aside hilang fixed bawah. CTA terlihat dari viewport awal dan padding bawah mencegah konten tertutup. |
| Pembayaran | Sesi/qty → order → rekening/wallet → proof → review/correction/approve → QR. History/resume/cancel sebelum proof tetap. Proof privat/QR single-use otoritas backend. |

## Acuan beranda mobile terbaru

[HOME_ALIGNMENT.md](HOME_ALIGNMENT.md) memuat ukuran, copywriting, state tombol, dan sumber SELA v5 yang diukur langsung. Ikuti acuan ini untuk header, navigasi bawah, composer entry, kartu post dan upvote. Blok “Di sekitar kamu” tidak berada di beranda mobile; label feed adalah Jelajah/Diikuti. Seluruh bottom-nav item aktif lime, bukan hanya ikon. Foto post memenuhi lebar kartu dan mention berada sebelum foto. Jangan menyamakan data pengguna dengan teks/angka contoh prototype.

## Source dan preview Sites

Komponen produk/CSS sama. Pengecualian: fixture API, login peran demo, notice dan warning tidak transfer nyata. Notice preview sesudah konten, bukan panel besar sebelum composer. Mention fixture harus mempunyai metadata. Transaksi preview tetap read-only; upvote/bookmark post dan event hanya preferensi perangkat/peran untuk menampilkan state aktif.

## Verifikasi

Check/build Go/Svelte, browser 360/390/1440, dan breakpoint tablet. Periksa CTA tiket awal, staging/cancel tanggal, layout agenda, focus composer/mention viewport, deskripsi/order pengelola, optional media, foto, overflow/hydration. Ulang flow pembayaran kedua role, proof anonim 401/pemilik dan pengelola 200 no-store, QR reuse ditolak. Screenshot local production untuk inspeksi, bukan otomatis bukti identik pixel. Keyboard/camera fisik, OAuth dan Lighthouse hanya diklaim setelah diuji. Hasil di TESTING.md dan evidence/ui-alignment-checks.json.

## Parity seluruh halaman

Pemetaan dan pengecualian bisnis di VISUAL_PARITY.md. Header, sidebar, aside discovery, grid, tab, komunitas, event, review, merchandise, search, saved, notifications, profil dan tabel dashboard memakai CSS sumber SELA v5. Penyesuaian terakhir mempertahankan padding/radius/typography kartu sekunder dari prototype dan memilih varian produk pertama saat dibuka. Source asli disertakan untuk pembandingan offline.
