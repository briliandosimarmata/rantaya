# Referensi SELA v5 yang dibekukan

Direktori sela-v5 adalah salinan output prototype asli dari commit a9c43a25a116e695df3a2f73ac978521ed092acd. Ini acuan desain yang diminta pengguna, bukan backend aplikasi Rantaya. Prototype asli tidak diubah.

Untuk membandingkan di laptop: dari folder ini jalankan `python -m http.server 5174`, lalu buka http://localhost:5174/sela-v5/. Jalankan Rantaya di http://localhost:5173 pada tab lain. Nama SELA, tanggal tetap dan angka simulasi pada referensi adalah data historis.

Jangan menjalankan login, pembayaran, penyimpanan browser atau moderasi simulasi referensi sebagai logika bisnis Rantaya. Otoritas akun, kuota, proof dan QR tetap API Go/PostgreSQL. Jangan mengedit baseline ini atau fixture ukurannya agar tes parity lolos tanpa persetujuan perubahan prototype.
