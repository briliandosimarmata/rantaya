# Panduan agent: Rantaya

## Tujuan dan konteks

Aplikasi komunitas seni berbasis Karawang dengan rencana ekspansi Jakarta. Komunitas/interaksi adalah inti; ticketing dan merch membantu pengelola. Nama produk final **Rantaya**. Istilah “ruang” untuk komunitas, slug komunitas contoh, cookie/draft key, database/module dan nama binary internal bukan branding; pertahankan identifier agar data lama kompatibel. Jangan mengganti domain atau membangun native livestream tanpa instruksi. Demo Sites terpisah hanya untuk preview UI.

Baca README.md, docs/PRODUCT_SPEC.md, docs/PAYMENTS.md dan docs/ARCHITECTURE.md sebelum mengubah perilaku lintas fitur. Detail API ada pada docs/API.md; UI mengikuti docs/UI_UX.md dan docs/PROTOTYPE_REFERENCE.md dan docs/HOME_ALIGNMENT.md. Ukuran beranda mobile yang diukur dari SELA v5 ada pada web/tests/fixtures/mobile-home-reference.json; jangan mengganti fixture agar desain yang melenceng lolos. Jangan menafsirkan daftar fitur sebagai izin merancang ulang susunan prototype. Baca docs/VISUAL_PARITY.md dan checklist parity sebelum mengubah UI. Baseline lengkap ada di docs/reference/sela-v5; web/src/prototype.css berisi aturan asli dan penyesuaian selector Svelte. Jangan mengedit baseline/fixture untuk menutupi perbedaan.

## Stack dan batas tanggung jawab

- SvelteKit 3, Svelte 5 runes, TypeScript 6. Konfigurasi Kit ada di web/vite.config.ts, bukan svelte.config.js. Alias package imports `#lib/*`; gunakan `.js` untuk import modul TypeScript, `.svelte` untuk komponen.
- API Go menggunakan net/http ServeMux dan encoding/json. Jangan menambah HTTP framework atau ORM secara rutin. pgx adalah driver PostgreSQL, bukan framework.
- Semua izin, kuota, transaksi, badge review dan QR divalidasi Go. Frontend bukan otoritas bisnis.
- Go API dapat dipakai klien native kelak. Endpoint di `/api`. Frontend browser memakai proxy same-origin; SSR membaca API_BASE_URL.
- PostgreSQL menyimpan sumber kebenaran. Jangan mengubah transaksi menjadi localStorage. localStorage hanya draft composer.
- SQL migrasi embedded, berurutan, advisory lock. Tambahkan versi migrasi baru; jangan mengedit versi terpasang untuk perubahan skema berikutnya.

## Aturan bisnis yang wajib dipertahankan

1. Customer dan organizer adalah akun terpisah, meskipun Google sub sama: UNIQUE(google_sub,role).
2. Composer fokus isi, judul opsional, tanpa dropdown jenis/identitas. Mention inline dengan scroll internal; customer tidak dapat menambah tautan tombol pembelian.
3. Feed popular tetap dibatasi kota terpilih; following memuat post resmi ruang yang diikuti.
4. Pembelian mengunci baris sesi dan menghitung approved/pending/correction serta hold unpaid yang masih aktif. Idempotency mencegah order ganda; jangan mengandalkan stok di client.
5. Awaiting_payment kedaluwarsa 30 menit (konfigurabel), dengan notifikasi in-app sekali pada lima menit terakhir. Proof submitted tidak otomatis kedaluwarsa. Pembatalan hanya sebelum proof dan konfirmasi belum transfer.
6. Metode pembayaran disalin sebagai snapshot pesanan. Edit master tidak mengubah instruksi pesanan lama.
7. Bukti transfer privat hanya pembeli dan pengelola event; QR hanya halaman privat. Jangan memasukkannya ke cache/sitemap/feed/log.
8. Approve membutuhkan konfirmasi dana diterima, atomik dengan penerbitan tiket/notifikasi. Satu token acak per tiket; persetujuan berulang tidak menambah tiket.
9. Check-in mengunci tiket, memeriksa pengelola, event, waktu dan checked_at. Token kedua kali ditolak. Tidak ada validasi offline di MVP.
10. Badge review berasal dari actual ticket check-in, bukan checkbox hadir atau sekadar approved.
11. Pengelola dapat membalas/melapor review; tidak boleh menghapusnya sendiri.
12. Demo/seed harus mati saat APP_ENV=production; server menolak kombinasi itu.

## UX dan kualitas

Mobile terlebih dahulu: 360px, 390px, 760px, desktop1440. Global kota/profil di header; empat bottom nav. Back di halaman turunan. Pencarian, kalender, composer, pembayaran dan denah adalah halaman penuh. Foto tetap berwarna. Lineup horizontal satu baris. Jangan menyempitkan deskripsi ruang menjadi kolom samping pada mobile.

Gunakan elemen semantic, label input, tombol keyboard, focus visible, empty/error/loading states. Escape teks pengguna; jangan gunakan @html untuk konten pengguna. JSON-LD dibuat dari data tepercaya dengan `<` di-escape. Hindari menambah warna dasar ke tiga warna yang sudah ditetapkan.

## Commands

Dari root: `npm run setup`, `npm run dev`, `npm run check`, `npm run build`, `npm test`. Go perubahan memerlukan restart dev. `npm run check` tidak menjalankan integration bila TEST_DATABASE_URL kosong. Unit + go vet wajib untuk API perubahan. Untuk ticketing jalankan integration terhadap DB uji terpisah (lihat TESTING). Untuk UX jalankan Playwright dan periksa viewport/console.

Jangan menaruh rahasia di .env.example, commit, gambar bukti, log, atau dokumen. .env, var/, binary/build dan node_modules dikecualikan. Bila kredensial Google/kamera/production runtime belum diuji, laporkan dengan jelas; jangan mengklaim Lighthouse atau production readiness tanpa pengukuran.

## Checkpoint akhir

Baca docs/BUILD_HANDOFF.md dan docs/verification.json. Migrasi terakhir versi 3: OAuth next_path, payment_reminders, timestamp votes dan activity_clicks. Versi 1–3 sudah diuji; perubahan skema berikutnya harus menambah migrasi 4. Login/onboarding wajib kembali ke path lokal yang aman, termasuk query session checkout. Dashboard menghitung klik, bukan menganggapnya penjualan. Pembelian merch memilih varian pertama saat membuka produk.

## Batas versi ini

Adapter Node hanya untuk local/dev/preview. Cloudflare adapter, pemindahan upload ke object storage, concurrency native PostgreSQL di CI, review keamanan production, autentikasi native, email/push dan live streaming adalah pekerjaan lanjutan. Baca docs/ROADMAP.md sebelum memperluas scope.
