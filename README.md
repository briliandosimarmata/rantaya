# Rantaya — komunitas seni dan tiket lokal

Source aplikasi web yang bisa dijalankan di laptop: **SvelteKit + TypeScript**, **Go API tanpa HTTP framework**, dan **PostgreSQL**. Nama produk **Rantaya** sudah ditetapkan. “Ruang” tetap dipakai sebagai istilah fitur komunitas dan pada beberapa identifier internal untuk kompatibilitas data.

Fokus awal Karawang: percakapan komunitas tetap berlangsung sebelum/sesudah pertunjukan, penemuan event, review, merchandise, dan pembelian tiket melalui transfer langsung kepada pengelola. Tidak membutuhkan Cloudflare atau payment gateway untuk pengembangan lokal.

## Mulai dalam lima langkah

Prasyarat: **Node.js 24+**, **Go 1.25+**, dan **PostgreSQL 17** (langsung terpasang atau melalui Docker Compose). Git opsional. Pastikan `node`, `npm`, dan `go` ada di PATH.

```bash
# Setelah ekstrak ZIP, buka folder rantaya-app
cp .env.example .env
npm run setup
docker compose up -d db
npm run dev
```

Windows PowerShell: gunakan `Copy-Item .env.example .env` untuk langkah pertama. Tanpa Docker, buat database/user mengikuti [panduan lokal](docs/LOCAL_DEVELOPMENT.md), lalu ubah `DATABASE_URL`.

Buka **http://localhost:5173**. API berada pada `127.0.0.1:8080`; SvelteKit meneruskan `/api/*` sehingga browser memakai satu origin. Startup API menjalankan migrasi otomatis dan data contoh bila `SEED_DEMO=true`.

Untuk membuka aplikasi dari HP atau laptop lain di Wi-Fi yang sama, atur `HOST=0.0.0.0`, `APP_URL=http://<IP-laptop>:5173`, dan `APP_ALLOWED_ORIGINS=http://localhost:5173,http://127.0.0.1:5173` di `.env`, lalu restart `npm run dev`. Laptop dapat memakai localhost, sedangkan perangkat lain memakai IP laptop. Login dan upload menerima hanya origin yang dikonfigurasi; API dan database tetap memakai loopback. Detail IP, firewall, dan batas akun demo ada di [panduan akses LAN](docs/LOCAL_DEVELOPMENT.md#akses-dari-perangkat-lain-di-wi-filan).

Masuk melalui **Coba akun demo penonton** atau **Coba akun demo pengelola**. Login moderator contoh ada di `/masuk/admin`. Akun-akun ini terpisah. Semua rekening, transaksi, komunitas, dan gambar contoh bersifat fiktif: **jangan transfer uang ke rekening demo**. Login Google nyata memerlukan kredensial milikmu; lihat [AUTH](docs/AUTH.md).

`npm run dev` membangun binary Go terlebih dahulu lalu menjalankan API dan Vite. Setelah mengubah Go, hentikan dengan Ctrl+C dan jalankan kembali. Perubahan frontend menggunakan HMR. Upload dan data transaksi tetap tersimpan setelah restart.

## Yang sudah tersedia

- Feed Jelajah/Diikuti, filter kota/kategori, upvote, bookmark, share, komentar dan balasan.
- Posting dengan judul opsional, gambar berwarna, mention inline event/ruang; tautan pembelian untuk pengelola.
- Profil ruang dengan Linimasa, Mention, Agenda, Merchandise, Ulasan, dan Tentang.
- Detail event dengan jadwal, Google Maps, deretan pengisi acara, trailer/flyer opsional, dan denah opsional.
- Transfer manual: buat pesanan, pilih bank/wallet, lanjutkan dari history, bukti privat, koreksi, persetujuan, tiket QR dan check-in online.
- Review tertaut ke event/ruang; badge hanya untuk tiket yang benar-benar check-in.
- Dashboard pengelola: buat/edit/publikasi event, kelola merch/rekening, profil, pin post, konfirmasi pembayaran, tanggapan ulasan, scanner kamera dan kode manual.
- Google OAuth + PKCE, onboarding bisa dilewati, profil/preferensi, notifikasi dalam aplikasi, reminder acara 24 jam dan batas pembayaran 5 menit, laporan dan moderasi.
- Halaman publik SSR, metadata/OG/canonical, JSON-LD Event, sitemap; halaman akun/transaksi noindex.

Merch menggunakan katalog + tautan pemesanan eksternal. Native livestream, aplikasi native, push/email, payment gateway, refund otomatis, dan pemilihan kursi belum termasuk versi awal. Rinciannya di [cakupan produk](docs/PRODUCT_SPEC.md) dan [roadmap](docs/ROADMAP.md).

Acuan desain asli tersedia di `docs/reference/sela-v5`, lengkap dengan CSS, HTML, JS dan aset. `web/src/prototype.css` memindahkan aturan visual SELA v5 ke komponen Svelte. Lihat [pemetaan halaman](docs/VISUAL_PARITY.md) sebelum mengubah desain.

## Pemeriksaan dan build

```bash
npm run check                 # Svelte/TypeScript + go vet + unit tests
npm run build                 # web/build + backend/bin/ruang-api
npm run preview               # Go + frontend hasil build; DB tetap harus berjalan
npm test                      # unit tests; integration skips without TEST_DATABASE_URL
```

Integration test membutuhkan **database uji terpisah**; [TESTING](docs/TESTING.md) menjelaskan perintah dan hasil aktual. Browser tests: `cd web`, `npx playwright install chromium`, kembali ke root, jalankan aplikasi pada DB uji, lalu `npm run test:e2e -- --project=mobile390`. Viewport lain dan pengaturan burst rate limit dijelaskan pada panduan tes.

## Peta repo dan dokumentasi

| Lokasi                      | Isi                                                         |
| --------------------------- | ----------------------------------------------------------- |
| `backend/cmd/api`           | Entry point, graceful shutdown, scheduler                   |
| `backend/internal/app`      | HTTP stdlib, SQL, auth, komunitas, transaksi, upload, tests |
| `web/src/routes`            | Halaman SSR dan proxy `/api/*`                              |
| `web/src/lib`               | Komponen, tipe, format, akses API                           |
| `web/static/assets`         | Foto contoh berwarna, font lokal, denah                     |
| `scripts`                   | Runner lintas platform untuk setup/dev/check/build          |
| `AGENTS.md`                 | Konteks dan aturan untuk Codex CLI/agent                    |
| `docs/PRODUCT_SPEC.md`      | Fitur, aktor, cakupan, acceptance criteria                  |
| `docs/ARCHITECTURE.md`      | Batas frontend/API, dependensi, arsitektur lokal            |
| `docs/API.md`               | Endpoint, payload, izin dan error                           |
| `docs/PAYMENTS.md`          | State machine pembayaran, kuota dan QR                      |
| `docs/AUTH.md`              | Google OAuth, akun terpisah, konfigurasi                    |
| `docs/DATABASE.md`          | Tabel, migrasi, data dan backup                             |
| `docs/UI_UX.md`             | Sistem visual dan keputusan mobile                          |
| `docs/LOCAL_DEVELOPMENT.md` | Setup Windows/macOS/Linux dan troubleshooting               |
| `docs/TESTING.md`           | Pengujian dan keterbatasan verifikasi                       |
| `docs/SECURITY.md`          | Kontrol saat ini dan pekerjaan sebelum production           |
| `docs/ROADMAP.md`           | Fitur lanjutan dan batas MVP                                |
| `docs/ASSETS.md`            | Asal aset dan lisensi dependensi                            |
| `docs/VISUAL_PARITY.md`    | Pemetaan semua halaman ke source SELA v5 dan batas data      |
| `docs/DEMO.md`             | Demo Sites, contoh data, dan batas transaksi                  |
| `docs/BUILD_HANDOFF.md`    | Checkpoint akhir dan petunjuk melanjutkan di Codex CLI        |
| `docs/reference`          | Prototype SELA v5 lengkap yang dapat dibuka offline           |
| `docs/RECOVERY_AUDIT.md`    | Checkpoint tersimpan, perbaikan recovery dan batas verifikasi |

Untuk melanjutkan di Codex CLI, buka folder ini dan minta agent membaca `AGENTS.md`, `README.md`, serta dokumen fitur yang akan diubah. Tidak ada kredensial Google, database pribadi, `node_modules`, ataupun binary di ZIP.

Demo VPS dengan Docker, NGINX, dan HTTPS Certbot dijelaskan pada [panduan deployment demo](docs/DEPLOYMENT_DEMO.md).
