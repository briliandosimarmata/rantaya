# Pengembangan lokal

## Prasyarat

Node.js 24+ dengan npm, Go 1.25+ dalam PATH, PostgreSQL 17 atau Docker Compose. Internet diperlukan ketika mengunduh dependencies serta memakai Google/Maps/trailer eksternal. Setelah setup, akun demo dan aset lokal tidak memerlukan payment gateway/OAuth.

## Database melalui Docker

compose.yaml menjalankan PostgreSQL 17-alpine saja: port loopback5432, user `ruang`, password `ruang_local`, database `ruang`, dan volume persistent. Jalankan `docker compose up -d db`; periksa `docker compose ps` hingga healthy. `docker compose stop` tidak menghapus data. Jangan menjalankan `down -v` bila ingin mempertahankan database.

## Tanpa Docker

Pasang PostgreSQL native sesuai OS, jalankan service dan gunakan psql admin:

```sql
CREATE USER ruang WITH PASSWORD 'ruang_local';
CREATE DATABASE ruang OWNER ruang;
```

Isi DATABASE_URL pada .env: `postgres://ruang:ruang_local@localhost:5432/ruang?sslmode=disable`. User DB perlu izin DDL untuk migrasi lokal. API membuat tabel/views sendiri; jangan import schema.sql manual sehingga migration history terlewat.

## Start dan environment

Copy .env.example menjadi .env, jalankan `npm run setup`, start DB, lalu `npm run dev`. Buka **http://localhost:5173** secara konsisten. macOS/Linux menggunakan cp, Windows PowerShell menggunakan Copy-Item. Scripts memakai Node child_process dan npm.cmd pada Windows.

Runner mengompilasi Go lalu menjalankan binary dan Vite. Ubah Go → Ctrl+C → ulangi dev. Ubah Svelte menggunakan HMR. Binary Go tidak membutuhkan runtime Go ketika dijalankan sendiri, tetapi proses development/build memerlukan compiler.

Env shell mempunyai prioritas di atas .env. Parser mendukung KEY=VALUE, quoted value dan comment satu baris penuh; tidak mengekspansi `$VAR`. Runner membuat UPLOAD_DIR absolut dari root repo.

| Env                     | Default / tujuan                                               |
| ----------------------- | -------------------------------------------------------------- |
| APP_ENV                 | development; production memakai cookie Secure dan menolak demo |
| API_ADDR                | 127.0.0.1:8080                                                 |
| DATABASE_URL            | PostgreSQL lokal DSN                                           |
| APP_URL                 | http://localhost:5173; origin utama yang diizinkan             |
| APP_ALLOWED_ORIGINS     | Kosong; origin tambahan lengkap, dipisahkan koma, tanpa wildcard |
| HOST                    | 127.0.0.1; gunakan 0.0.0.0 untuk frontend di jaringan lokal    |
| API_BASE_URL            | http://127.0.0.1:8080; SSR/proxy                               |
| UPLOAD_DIR              | var/uploads; foto/proof disk lokal                             |
| DEMO_LOGIN / SEED_DEMO  | true dalam example, khusus development                         |
| HOLD_MINUTES            | 30, valid1–120                                                 |
| GOOGLE_CLIENT_ID/SECRET | Kosong hingga dikonfigurasi                                    |
| GOOGLE_REDIRECT_URI     | http://localhost:5173/api/auth/google/callback                 |
| TEST_DATABASE_URL       | Opsional; DB uji terpisah                                      |

## Akses dari perangkat lain di Wi-Fi/LAN

Cari IPv4 laptop pada pengaturan jaringan (macOS: `ipconfig getifaddr en0`; Windows: `ipconfig`). Ubah `.env`, misalnya:

```dotenv
HOST=0.0.0.0
APP_URL=http://192.168.1.10:5173
APP_ALLOWED_ORIGINS=http://localhost:5173,http://127.0.0.1:5173
API_ADDR=127.0.0.1:8080
API_BASE_URL=http://127.0.0.1:8080
```

Restart `npm run dev`, lalu buka **http://localhost:5173 di laptop** atau **alamat APP_URL di perangkat lain**. Port Vite mengikuti APP_URL. API dan database tetap di loopback; perangkat lain memakai proxy `/api` frontend. Go mengizinkan APP_URL dan origin tambahan yang ditulis persis dalam APP_ALLOWED_ORIGINS; frontend menggunakan daftar yang sama untuk form/upload. Proteksi origin/Fetch-Metadata tetap aktif. Skema, hostname, dan port harus cocok; wildcard, path, query, fragment, dan kredensial pada origin ditolak. Tidak perlu mengubah Google redirect untuk akun demo.

Session cookie terpisah per hostname: login di localhost tidak otomatis membuat sesi pada IP LAN. APP_ALLOWED_ORIGINS kosong mempertahankan perilaku satu origin; localhost tidak otomatis diizinkan jika APP_URL berbeda. Google OAuth tetap memerlukan hostname callback yang sama dengan hostname saat memulai login.

Perangkat harus berada di jaringan yang sama tanpa client isolation. Bila koneksi ditolak, izinkan Node pada firewall macOS/Windows untuk jaringan lokal. Tidak perlu membuka port API/database atau port forwarding router. Bila IP laptop berubah, sesuaikan APP_URL dan restart; untuk preview, build ulang terlebih dahulu. Untuk kembali ke mode laptop saja, pakai `HOST=127.0.0.1` dan `APP_URL=http://localhost:5173`.

Akun demo penonton/pengelola dapat dicoba tanpa password. Akun demo tiap peran dipakai bersama oleh semua pengunjung; gunakan browser/profile berbeda untuk mencoba dua peran. Session memakai cookie HttpOnly/SameSite Lax, dan server menolak demo/seed dalam production. HTTP LAN belum mengenkripsi koneksi: gunakan jaringan pribadi tepercaya dan data fiktif. Kamera scanner pada HTTP IP LAN dapat tidak tersedia karena browser membutuhkan secure context; gunakan input kode manual.

Checkout tiket memakai kunci UUID acak dari `crypto.getRandomValues()`, yang tersedia pada HTTP LAN. Proses ini tidak memerlukan `crypto.randomUUID()` atau secure context. Jika pembuatan kunci atau permintaan gagal, pesan error ditampilkan dan tombol kembali bisa ditekan; retry memakai kunci yang sama agar tidak menambah pesanan ganda.

## Build dan preview

`npm run build` menghasilkan web/build dan backend/bin/ruang-api. Start DB lalu jalankan **`npm run preview` dari root** untuk menjalankan Go dan frontend hasil build bersama, menggunakan .env yang sama. Tidak perlu memasang Cloudflare.

SvelteKit 3 / adapter-node 6 memakai `paths.origin` yang diisi dari **APP_URL saat build** dan daftar CSRF trustedOrigins dari **APP_ALLOWED_ORIGINS saat build**. Variabel runtime ORIGIN dari adapter lama tidak dipakai. Bila mengganti APP_URL atau APP_ALLOWED_ORIGINS, build ulang, lalu restart preview; sesuaikan Google redirect bila memakai OAuth. Runner memeriksa web/build/origin.txt dan allowed-origins.txt agar build dengan konfigurasi origin lama tidak digunakan tanpa sengaja. Proteksi CSRF tetap aktif.

Runner preview memakai HOST dari environment/.env (default 127.0.0.1), PORT dari APP_URL, dan BODY_SIZE_LIMIT=10M supaya multipart foto sampai 8 MB bisa diteruskan ke Go. Vite preview dari folder web hanya menjalankan frontend dan tetap memerlukan API/DB dengan origin yang cocok. Binary Go langsung tidak membaca .env; gunakan runner atau export env sendiri. API_BASE_URL dibaca server frontend saat runtime.

## Troubleshooting

- 503/backend unavailable: periksa service DB, DSN dan port; restart Go. Startup API gagal akan menghentikan runner.
- Origin 403: origin browser (skema/hostname/port) harus sama persis dengan APP_URL atau salah satu APP_ALLOWED_ORIGINS. Tambahkan origin localhost bila APP_URL memakai IP LAN, lalu restart dev. Untuk preview, build ulang setelah mengubah salah satu variabel origin; jangan mematikan CSRF untuk mengatasi upload.
- Google redirect mismatch: http, hostname, port dan path console harus persis sama dengan .env.
- Command tidak ditemukan: pasang runtime dan refresh PATH terminal, lalu setup ulang.
- Port in use: tutup server lain atau ubah API_ADDR+API_BASE_URL bersama. Port Vite mengikuti APP_URL; sesuaikan Google URI juga bila memakai OAuth.
- Migrasi gagal: user DB perlu DDL; lihat log API. Jangan hapus database untuk mengatasi error tanpa backup.
- QR ditolak: pastikan approved, event benar, window gate dan belum digunakan.
- Kamera tidak tersedia: izin browser, localhost/HTTPS secure context; gunakan kode manual.
- Seed lewat: jadwal relatif terhadap seed pertama; buat event baru, jangan reset seluruh data.
- Import TS gagal: alias #lib membutuhkan ekstensi .js; konfigurasi Kit3 berada pada Vite.

Tidak memerlukan Cloudflare CLI, Redis, message broker, HTMX atau framework Go.
