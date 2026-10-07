# Demo Rantaya

Demo: https://ruang-ui-preview.sandramoored074.chatgpt.site . Akses mengikuti Site pribadi yang sudah ada. Nama/tampilan produk Rantaya; slug lama dipertahankan.

Preview memakai komponen Svelte produk, CSS dan aset yang sama, dengan fixture fiktif SELA v5. Masuk demo hanya memilih tampilan penonton/pengelola/moderator. Navigasi, pencarian, kategori, kalender, mention, katalog dan contoh status transaksi dapat dijelajahi. Upvote dan bookmark adalah preferensi cookie perangkat/peran pada konten contoh; tidak mengubah data publik.

Reservasi, proof, approval, check-in, post/review/profil/rekening nyata tidak disimpan di demo. Mutasi pembayaran ditolak. Contoh order menunjukkan unpaid, menunggu, perbaikan, approved serta QR contoh, tanpa menerima uang/tiket nyata. Perhatikan notice “Preview UI · data contoh”.

Aplikasi lengkap ada pada source `rantaya-app`: Node/SvelteKit + API Go + PostgreSQL. Jalankan langkah README untuk mencoba transfer manual dengan gambar bukti sintetis dan QR sungguhan di database lokal. Jangan transfer ke rekening fiktif. Tidak ada credential/account pengguna nyata dalam fixture.

Backend lokal tidak dipindahkan ke Workers. Preview memakai adapter terpisah; pilihan hosting production, upload object storage dan database kelak ada di ROADMAP.md. Tidak ada OAuth Google live yang tersambung ke preview.
