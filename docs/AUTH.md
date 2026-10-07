# Autentikasi

## Akun terpisah

Masuk penonton: /masuk/customer. Masuk pengelola: /masuk/organizer. Masing-masing membuat row accounts berbeda, dengan UNIQUE(google_sub,role). Email Google yang sama boleh punya kedua akun, profil/aktivitas terpisah. Satu browser memegang satu session aktif; logout/login pada halaman role lainnya untuk berganti.

Tidak ada password aplikasi, OTP email atau dropdown posting sebagai. Pengelola otomatis mempunyai satu ruang. Profile nama/kota/bio/avatar/minat/preferences dapat diubah. Pada akun baru diarahkan /onboarding dengan **Lewati, lengkapi nanti**; bukan popup wajib. Admin tidak bisa mendaftarkan diri lewat OAuth; provisioning admin production harus dilakukan secara eksplisit oleh operator di DB setelah audit akses.

## Menyiapkan Google

1. Buat Google Cloud project dan konfigurasi Google Auth Platform/OAuth consent (branding, audience, data access) untuk aplikasi sendiri.
2. Buat OAuth client tipe **Web application**.
3. Tambahkan redirect URI persis: `http://localhost:5173/api/auth/google/callback`.
4. Isi .env `GOOGLE_CLIENT_ID`, `GOOGLE_CLIENT_SECRET`, `GOOGLE_REDIRECT_URI`.
5. Dalam mode Testing, tambahkan emailmu sebagai test user sesuai konfigurasi audience Google.
6. Restart `npm run dev`; tombol **Lanjutkan dengan Google** aktif. Browser harus memakai hostname localhost yang cocok dengan .env, bukan berganti ke127.0.0.1.

Scope hanya openid/email/profile. Dokumentasi resmi: https://developers.google.com/identity/protocols/oauth2/web-server . Pengaturan/nama menu dapat berubah; ikuti console dan URI di atas.

## Protokol implementasi

GET start memvalidasi role customer/organizer, membuat state random, menyimpan SHA256 state + PKCE verifier di DB10 menit, menaruh state cookie HttpOnly SameSiteLax, lalu redirect ke Google dengan challengeS256.

Callback memeriksa state query/cookie, mengonsumsi state satu kali, menukar authorization code melalui token endpoint HTTPS dan mendapatkan userinfo melalui bearer access token. Memerlukan verified email dan stable Google sub. Tidak mempercayai email/token yang dikirim client. Access token Google tidak disimpan; session aplikasi baru diterbitkan.

Cookie ruang_session HttpOnly, SameSiteLax,14 hari, Secure saat APP_ENV=production. DB hanya hash session. Logout menghapus row token dan cookie. Header Authorization Bearer diterima API untuk klien terkontrol, tetapi endpoint native login belum dibuat. Kegagalan OAuth kembali ke halaman login; state invalid tidak menghasilkan session.

## Demo lokal

.env.example mengaktifkan DEMO_LOGIN=true dan SEED_DEMO=true agar tidak terhambat kredensial. POST /auth/demo menerima customer/organizer/admin untuk akun fiktif tetap. Demo bukan metode login production. Server menolak startup bila APP_ENV=production dan salah satu flag itu aktif.

Google login nyata belum diuji end-to-end tanpa client ID/secret milik pengguna. Unit/API/role/session demo telah diuji. Jangan menganggap pengaturan consent/verification Google sudah dilakukan oleh source ini.

## Kembali ke tujuan setelah login

Halaman yang membutuhkan akun membawa `next` berupa path lokal dan query (misalnya `/event/festival-panggung-kecil/tiket?session=s-e3b`). Login demo dan onboarding menjaga tujuan tersebut. Google start menyimpan tujuan yang sudah divalidasi bersama state di DB, sehingga callback tidak menerima tujuan bebas dari browser. Migrasi 3 menambah oauth_states.next_path.

Path absolute, protocol-relative, backslash dan control character ditolak. Tanpa next valid: customer ke `/`, organizer ke `/kelola`, admin demo ke `/admin`. Onboarding, baik Simpan maupun Lewati, melanjutkan ke tujuan valid yang sama. Unit URL policy serta browser demo/onboarding/session checkout menguji perilaku ini; OAuth nyata tetap membutuhkan kredensial pengguna.
