# Audit pemulihan, 7 Oktober 2026

Audit dilakukan terhadap source dan data yang tersimpan dari proses sebelumnya. Proyek tidak dimulai ulang. Database lama tidak di-reset; salinannya dibuka untuk pemeriksaan. Tidak ditemukan repo Git/history sehingga perubahan tidak bisa dibandingkan dengan commit lama.

## Checkpoint yang benar-benar tersimpan

- Source Go API, SvelteKit, seluruh route publik/private/dashboard, lockfile, 15 aset lokal, dokumentasi dan panduan AGENTS.md.
- Binary Go dan build frontend dari proses lama; keduanya bukan bukti source terbaru sudah lolos tes.
- Dua direktori PostgreSQL WASM/PGlite berisi 25 tabel dengan migrasi 1 dan 2. Database integration lama: 11 akun, 5 event, 8 order, 3 tiket, 2 upload. Database browser lama: 7 akun, 4 event, 1 order, 1 tiket, tanpa upload.
- Trace browser membuktikan upload gagal dengan HTTP 403 “Cross-site POST form submissions are forbidden” sebelum mencapai Go.

## Yang saat audit masih parsial

Dashboard tersedia dan navigasinya pernah diuji; penyimpanan semua form belum dibuktikan lewat browser. Commerce API pernah diuji, tetapi browser berhenti saat upload bukti. verification.json masih menandai browser pending. Log sementara di /tmp dan ZIP source belum tersimpan.

## Perbaikan dari checkpoint itu

Origin di SvelteKit 3 ditetapkan melalui paths.origin saat build dari APP_URL. Variabel ORIGIN runtime lama tidak berlaku untuk adapter-node 6, sehingga request multipart HTTP dianggap berbeda origin. Proteksi CSRF dipertahankan. Runner preview memeriksa origin build dan menyediakan batas body 10 MB; batas file Go tetap 8 MB.

File robots.txt statis bawaan scaffold dihapus karena bertentangan dengan route robots dinamis yang mengecualikan transaksi/profil/dashboard. Migrasi terpasang tidak diubah.

Tes diperluas untuk penyimpanan rekening/event/lineup/Maps/merch/profil, resume pembayaran, koreksi/upload ulang, notifikasi, bukti privat, persetujuan, dan QR reuse. Integration Go juga mengirim delapan request check-in bersamaan untuk satu QR; tepat satu harus berhasil.

## Batas verifikasi

Hasil final aktual dicatat di TESTING.md dan verification.json. Paket PostgreSQL native berhasil diekstrak, tetapi lingkungan hanya memetakan UID 0; tidak dapat menjalankan proses non-root yang diwajibkan PostgreSQL. Tes SQL dilanjutkan dengan PGlite melalui PostgreSQL wire protocol dan satu koneksi. Ini belum membuktikan lock beberapa koneksi native. Google OAuth nyata, kamera fisik, Lighthouse dan deployment production tetap belum diverifikasi.

## Kelanjutan dari checkpoint terakhir

Source dipulihkan dari ZIP sebelumnya, lalu perubahan beranda/filter yang lebih baru di workspace disatukan. SELA v5 ditemukan dalam checkout aslinya pada commit a9c43a25a116e695df3a2f73ac978521ed092acd; prototype tidak diedit. Referensi penuh dibekukan pada docs/reference/sela-v5 agar Codex CLI dapat membacanya di laptop tanpa bergantung akses Site.

Susunan dashboard, detail event/ruang, search, merch, saved, notifications, profil, review dan aside dikembalikan ke source tersebut. Copy demo, lineup serta ukuran tote dipulihkan. Migrasi 3 menambah kelanjutan login aman, reminder deadline dan metrik klik aktual. Saat inspeksi visual ditemukan pilihan varian pertama kosong; state awal varian diperbaiki dan dicek ulang.

Pengujian akhir mencakup check/build, Go race/SQL, 32 kombinasi flow browser lulus (satu tes ukur mobile sengaja skip di desktop), dan 204 state route. Burst suite pernah mencapai 600 request/menit; tes dibagi per file dengan restart API dan database tetap dipertahankan, tanpa mematikan/menaikkan limiter. Tes tombol search/back disinkronkan ke navigasi filter yang selesai sebelum membuka search. Hasil terbaru di verification.json dan evidence/final-checks.json menggantikan status checkpoint historis di atas.
