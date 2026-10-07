# Demo VPS Rantaya

Deployment demo memakai NGINX host, Certbot, dan project Docker Compose `rantaya-demo`. Web berada pada loopback `127.0.0.1:3100`; API dan PostgreSQL hanya tersedia di jaringan project Docker. Situs: https://rantaya.briliando.dev.

Ini demo dengan akun dan data fiktif yang dipakai bersama. Runtime API menggunakan `APP_ENV=demo`, `DEMO_LOGIN=true`, `SEED_DEMO=true`, dan `COOKIE_SECURE=true`. `APP_ENV=production` tetap menolak demo/seed dan selalu mengaktifkan Secure cookie. Adapter Node dipakai untuk demo ini sesuai permintaan; belum merupakan deployment production untuk transaksi nyata.

Login penonton: `/masuk/customer`, lalu **Coba akun demo penonton**. Login pengelola: `/masuk/organizer`, lalu **Coba akun demo pengelola**. Gunakan profil browser terpisah untuk kedua peran. Tidak ada password akun demo yang perlu dibagikan.

Source ada di `/opt/rantaya-demo/source`. Konfigurasi privat berada di `/opt/rantaya-demo/.env`, permission 0600, dengan password database acak yang dibuat di VPS. File tersebut tidak masuk Git, image, atau log. `.dockerignore` mengecualikan kredensial, data lokal, build, dan hasil browser test dari konteks build.

Volume `rantaya-demo_db_data` dan `rantaya-demo_uploads` mempertahankan database/upload saat container diganti. Jangan memakai `docker compose down -v` untuk restart.

Perintah pengelolaan, dijalankan pada VPS:

```sh
cd /opt/rantaya-demo/source
docker compose --env-file /opt/rantaya-demo/.env -f deploy/compose.demo.yaml ps
docker compose --env-file /opt/rantaya-demo/.env -f deploy/compose.demo.yaml logs --tail=100 web api
```

Build dilakukan satu service per giliran untuk mengurangi beban VPS:

```sh
docker compose --env-file /opt/rantaya-demo/.env -f deploy/compose.demo.yaml build api
docker compose --env-file /opt/rantaya-demo/.env -f deploy/compose.demo.yaml build web
docker compose --env-file /opt/rantaya-demo/.env -f deploy/compose.demo.yaml up -d --wait
```

APP_URL disisipkan ke frontend saat build. Perubahan domain memerlukan rebuild web dan penyesuaian vhost serta sertifikat. Google OAuth tetap memerlukan kredensial dan konfigurasi callback sendiri; demo memakai login akun contoh.

Seed awal menjadwalkan sesi pertama “Di Balik Layar” satu jam setelah database dibuat agar check-in dapat diuji. Pembelian sesi itu ditutup setelah mulai, walaupun kuota tersisa; restart tidak mengubah tanggal seed yang tersimpan. Untuk demo VPS, tambahkan sesi mendatang setelah startup:

```sh
sh deploy/refresh-demo-session.sh
```

Script memeriksa API berjalan dalam mode demo dengan login/seed aktif. Script hanya menambahkan sesi baru pada event fiktif asli jika tidak ada sesi lebih dari 24 jam ke depan: jadwal 30 hari mendatang, pukul 19.00 WIB, kapasitas 60, harga dari sesi contoh. Pemanggilan ulang tidak menggandakan sesi. Jalankan kembali saat jadwal demo mendatang mendekati habis. Jadwal lama, pesanan, bukti, dan tiket dipertahankan; script tidak membuka event yang belum diterbitkan atau menambah kuota sesi yang terjual habis. Check-in tetap mengikuti jendela waktu sesi, sehingga sesi mendatang digunakan untuk checkout/pembayaran, bukan check-in hari ini.

Verifikasi penambahan sesi dan checkout publik: [evidence/vps-demo-session-checks.json](evidence/vps-demo-session-checks.json).

NGINX Rantaya ada di `/etc/nginx/sites-available/rantaya-demo`, dengan symlink khusus di `sites-enabled`. Sebelum perubahan, backup/checksum konfigurasi lama disimpan di `/opt/rantaya-demo/backups`. Setiap perubahan vhost harus lulus `nginx -t` sebelum `systemctl reload nginx`. Tidak perlu restart NGINX atau Docker host.

Untuk menghentikan demo, jalankan `docker compose --env-file /opt/rantaya-demo/.env -f deploy/compose.demo.yaml stop` dari direktori source. Data tetap tersimpan. Untuk menghapus vhost demo, lepaskan hanya symlink `/etc/nginx/sites-enabled/rantaya-demo`, lalu uji dan reload NGINX; konfigurasi situs lain tidak perlu diubah. Backup awal berisi konfigurasi seluruh host, jadi jangan menimpa seluruh `/etc/nginx` ketika hanya memperbaiki Rantaya.

HTTPS memakai Certbot `certonly --webroot` dengan webroot `/var/www/rantaya-acme`. Sertifikat dan private key tetap dikelola Certbot di `/etc/letsencrypt/live/rantaya.briliando.dev`. HTTP diteruskan ke HTTPS, kecuali challenge ACME. Timer renewal bawaan Certbot tetap dipakai; deploy hook khusus sertifikat Rantaya menguji konfigurasi sebelum reload NGINX. Periksa renewal dengan:

```sh
certbot renew --cert-name rantaya.briliando.dev --dry-run
systemctl list-timers --all | grep certbot
```

NGINX menambahkan `X-Robots-Tag: noindex, nofollow`, menerima multipart sampai 10 MB, dan mematikan access log khusus Rantaya agar URL privat tidak dicatat. Login, CSRF, ownership, proof privat, kuota, dan check-in tetap divalidasi API Go. Rahasia dan QR/proof privat jangan disalin ke laporan atau Git.

Hasil verifikasi deployment ada pada [evidence/vps-demo-checks.json](evidence/vps-demo-checks.json): login kedua peran, cookie Secure/HttpOnly, sembilan regresi checkout, pembayaran sampai check-in, renewal Certbot, dan integritas konfigurasi NGINX lama. Google OAuth nyata, kamera fisik, Lighthouse, dan kesiapan production belum diverifikasi.
