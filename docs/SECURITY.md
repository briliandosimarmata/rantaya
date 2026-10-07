# Kontrol keamanan dan batas v0.1

Ini implementasi untuk pengembangan lokal. Belum merupakan audit keamanan atau deployment production.

## Sudah diterapkan

- Session opaque 256 bit, hash SHA256 di DB, cookie HttpOnly/SameSite Lax dan Secure saat production.
- Google OAuth authorization code, PKCE S256, state cookie + DB satu kali, verified email melalui Google userinfo.
- Akun/izin terpisah; ownership dicek Go pada event, produk, rekening, order, proof, review reply dan check-in.
- SQL parameterized, JSON strict, validasi panjang/enum/URL, body limit, upload whitelist MIME dan ukuran, path file random/basename, bukan nama file pengguna.
- Proof private no-store, QR/tiket hanya endpoint privat; SSR akun no-store, private routes noindex.
- Kuota/order/approval/check-in dalam transaksi dan row lock; unique idempotency/token/admission.
- Origin/Fetch-Metadata protection untuk mutasi browser; CORS hanya APP_URL dan origin tambahan eksplisit dalam APP_ALLOWED_ORIGINS (default kosong). Origin harus HTTP(S) lengkap tanpa wildcard/path/kredensial; Sec-Fetch-Site cross-site tetap ditolak. Frontend memakai daftar yang sama untuk form/upload. Rate limit dan HTTP/context timeouts.
- Konten teks dirender escaped. SVG pengguna ditolak. External link noopener/noreferrer.
- Data demo/akses demo ditolak pada APP_ENV=production; rahasia tidak ada di repo.

## Pekerjaan sebelum menerima pengguna/uang nyata

1. Uji native PostgreSQL dengan beberapa koneksi/proses, deployment HTTPS, storage durable, backup DB+upload dan restore drill.
2. Terapkan limiter per-account/per-route di proxy edge. Saat ini semua request proxy Go terlihat dari IP server yang sama; batas global600/min tidak cocok untuk trafik ramai atau beberapa instance.
3. OAuth consent/configuration nyata, account recovery/provisioning admin, session revocation dan strategi Google native jika ada native app.
4. Tambahkan CSP sesuai host trailer/media, secure proxy/header config, audit dependency, scanning/normalisasi gambar, validasi dimensi WebP, quotas upload dan cleanup orphan. JPEG/PNG divalidasi dimensi; WebP baru signature/size, belum full decode.
5. Retention proof/payment/PII, penghapusan akun, moderasi/staff, audit trail admin, log tanpa token/proof/PII, monitoring alert dan batas file storage.
6. Kebijakan refund/cancel event, sengketa manual, organizer verification dan kontak bantuan. Expired-but-paid saat ini ditangani manual oleh pengelola.

Tidak ada request dari server untuk mengambil URL gambar arbitrary: browser memuat URL media. Berbagi QR berarti memberikan akses admission, jadi URL/token jangan masuk analytics/log publik. Tidak ada offline QR verification; offline fallback tanpa otoritas server akan membuka risiko double entry.
