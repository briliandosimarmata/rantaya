# Melanjutkan Rantaya di Codex CLI

Build ini melanjutkan checkpoint yang tersimpan, bukan memulai ulang. Source asli SELA v5 disertakan dan source aplikasi tetap SvelteKit + Go stdlib + PostgreSQL. Semua dokumentasi fitur/API/otorisasi ada pada repo, ditambah bukti verifikasi tanpa session token, QR pribadi, atau proof.

## Mulai di laptop

1. Ekstrak ZIP lalu buka folder rantaya-app.
2. Baca AGENTS.md dan README.md. Salin .env.example menjadi .env.
3. Pastikan Node24+, Go1.25+, Docker/PostgreSQL17 tersedia, lalu `npm run setup`.
4. `docker compose up -d db`, lalu `npm run dev`. Buka http://localhost:5173.
5. Pilih akun demo penonton atau pengelola. Gunakan dua browser/profile untuk mencoba kedua peran bersamaan.

Untuk memulai percakapan dengan Codex CLI: “Baca AGENTS.md, BUILD_HANDOFF.md dan spesifikasi. Pertahankan desain SELA dan lanjutkan fitur yang saya minta.” Git dapat diinisialisasi setelah ekstrak; dependency, .env, binary/build dan data pengguna tidak disertakan.

## Checkpoint

- UI seluruh route dipindahkan dari baseline dan diuji responsif. Referensi offline: docs/reference/README.md.
- Commerce Go berjalan dengan transaksi SQL, proof privat, history/resume, correction, approval, kuota, QR dan scanner manual.
- Login/onboarding kembali ke tujuan aman dan query sesi. Google login membutuhkan client ID/secret sendiri.
- Migrasi terpasang sampai versi3; tambahkan versi4 untuk perubahan berikutnya.
- Read TESTING.md/verification.json untuk hasil aktual, terutama batas PostgreSQL native beberapa koneksi dan kamera/OAuth fisik.
- Demo Sites adalah preview data contoh. Seluruh penulisan sebenarnya terjadi pada Go/PostgreSQL lokal.

## Checklist uji pengguna

Coba buat/edit event dan merch, isi rekening contoh, login penonton, pilih sesi/dua tiket, pilih metode, unggah bukti sintetis, minta koreksi sebagai pengelola, unggah ulang, setujui setelah konfirmasi dana contoh, kemudian catat masuk dan coba kode yang sama lagi. Tes otomatis melakukan alur ini; pengujian native PostgreSQL di laptop/CI melengkapi verifikasi concurrency.

Production Workers, payment gateway, aplikasi native, staff scanner multipengguna, email/push, cleanup/retention dan audit production adalah pekerjaan berikutnya sesuai ROADMAP, bukan syarat menjalankan versi lokal.
