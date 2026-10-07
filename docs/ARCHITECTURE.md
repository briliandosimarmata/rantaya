# Arsitektur

## Jalur lokal

Browser → SvelteKit (`localhost:5173`) → proxy `/api/*` → Go (`127.0.0.1:8080`) → PostgreSQL (`localhost:5432`). SSR halaman publik memanggil API Go langsung dengan cookie request; rendering menggunakan Svelte, bukan template Go. Backend tidak mengirim HTML produk.

Frontend dibangun menjadi JavaScript/CSS/assets dan server SSR lokal adapter-node. **Node tetap berjalan untuk dev/SSR lokal**, sesuai keputusan menunda Worker sampai production. Ini bukan static-only build; event/post baru harus bisa menghasilkan HTML SEO tanpa rebuild. Jangan mengganti deployment ke SPA statis sambil mengklaim perilaku SSR sama.

## Backend

`net/http` ServeMux method/path, request context timeout20 s, HTTP timeouts, JSON decode strict, validasi input, role dan ownership. pgxpool maksimum8 koneksi; SQL parameterized; transaksi langsung, tanpa ORM/framework. Dependensi direct hanya pgx; dependensi indirect driver tercatat go.mod/go.sum.

- app.go: config, migration, router, middleware, common helpers, expiry/reminders.
- auth.go: sessions, Google OAuth PKCE, onboarding/profile.
- community.go: organizers, events, feed, mentions, votes, comments, follows, bookmarks, reviews.
- commerce.go: orders, payment masters, proof submission, approval/tickets/check-in.
- support.go: upload/download, products, notifications, dashboard, reports/admin.
- seed.go: data demo idempotent, tidak boleh production.

Go process melayani API + scheduler per menit. Scheduler aman pada lebih dari satu instance karena update status transaksional dan PK reminder, tetapi production scheduling/observability perlu desain ulang sesuai hosting.

## Frontend

Route `+page.server.ts` mengambil data SSR. `api()` membedakan status/error server, `requireUser()` melindungi halaman; API tetap memeriksa sendiri. Route `/api/[...path]` meneruskan method/body/cookie/header dan multi Set-Cookie, tanpa menyimpan token di localStorage. Body dan foto publik melalui same-origin; bukti dan tiket no-store privat.

Context root berisi account/follows/votes/bookmarks/city/unread, toast dan invalidateAll. Query-tab tidak mereset komponen; perubahan pathname meremount form supaya berpindah resource tidak memakai draft resource lama. Draft composer saja disimpan per user di perangkat. QR generator dan camera decoder di-import ketika diperlukan.

Tidak memakai UI framework/CSS utility library. Komponen semantic sederhana dan CSS responsive. SvelteKit 3 memakai `sveltekit({...})` pada vite.config.ts; package imports `#lib/*` harus memiliki ekstensi `.js` pada import TS.

## Penyimpanan dan waktu

PostgreSQL timestamptz menyimpan instant; API ISO8601, frontend tampilkan Asia/Jakarta. Form organizer mengubah datetime-local ke `+07:00`. Harga integer rupiah (bigint); tidak memakai floating point uang. Upload file disimpan di UPLOAD_DIR, metadata DB; proof bukan media publik. Replikasi/object storage belum ada.

Entity IDs opaque random 128 bit. Session/token QR random 256 bit. Session DB menyimpan SHA256 token. QR token disimpan privat karena diperlukan menampilkan tiket kembali; endpoint publik tidak menyertakannya.

## SEO dan performa

Event/ruang/post/merch menggunakan SSR, satu h1 utama, semantic header/main/article/nav, meta description/canonical/OG. Event JSON-LD dibuat dan escape `<`. Sitemap memasukkan event published dan ruang; halaman private memakai noindex, robots disallow. Sitemap saat ini tidak memuat seluruh post/merch—perlu perluasan sebelum indeks skala besar.

Font Noto Sans lokal, WebP demo, gambar malas-load kecuali media utama, import dinamis scanner/QR. Tidak ada analytics/tracking eksternal. Font normal/bold only. Lighthouse harus diukur pada build production, bukan Vite dev; skor tidak dijamin hanya berdasarkan pemilihan framework.

## Jalur ke native dan Worker

Native menggunakan API yang sama dengan bearer session token (server sudah menerima bearer), tetapi alur Google native/token exchange/secure credential storage **belum dibuat**. API web saat ini menerbitkan cookie.

Worker kelak membutuhkan adapter-cloudflare, pembacaan runtime env mengganti process.env di helper API, domain HTTPS/cookie/CORS, upload external storage dan koneksi Go API yang dihosting terpisah. Go binary bukan Cloudflare Worker. Tidak ada config Worker/deployment yang dinyatakan siap pada v0.1.
