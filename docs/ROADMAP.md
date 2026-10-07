# Roadmap dan keputusan yang ditunda

## Prioritas setelah local walkthrough

- Pilot bersama komunitas Karawang: apakah organizer rutin berbagi proses, apakah customer mengikuti/mention, apakah event ditemukan dan pembelian selesai.
- Pengujian native PostgreSQL di CI, integrasi OAuth nyata, perangkat kamera/mobile keyboard, aksesibilitas dan Lighthouse production build.
- Kontak bantuan organizer yang jelas untuk kasus sudah transfer tetapi expired; operasi koreksi/refund/cancel event, retensi proof dan backup.
- Perbaikan typed API contract, cursor pagination untuk order/review/catalog, native login dan admin provisioning.

## Production hosting

SvelteKit adapter-cloudflare + runtime env; API Go tetap di hosting terpisah; PostgreSQL managed/native; upload object storage. HTTPS/domain/cookie/proxy/limiter/observability/backup harus ditetapkan. Tidak ada provisioning cloud dalam source ini; adapter-node dipakai local agar instalasi sederhana.

## Pengembangan produk berikutnya

Native mobile, push/email, reminder pilihan waktu, Google Calendar, kursi bernomor, multi-staff gate, organizer verification dan laporan performa. Payment gateway baru sesudah biaya/volume masuk akal, dengan webhooks/idempotency/signature dan reconciliation; jangan menggantikan manual status hanya dari response browser.

Native livestream memerlukan provider/WebRTC/RTMP, moderasi, hak konten, biaya bandwidth/recording dan operasi terpisah. MVP dapat membagikan tautan live eksternal melalui post pengelola. Tidak ada streaming/transcoding bawaan yang diklaim sudah selesai.

Merch checkout internal membutuhkan stok per varian, ongkir/alamat/order/refund. Saat ini katalog + pemesanan ke kanal pengelola tetap bekerja setelah event; jangan menganggap ticket_orders dapat digunakan untuk shipping tanpa model baru.

Rekomendasi geografis/jarak GPS, pemeringkatan reputasi organizer, forum/DM dan monetisasi promoted post merupakan eksperimen berikutnya. Upvote v0.1 berlaku pada post; ranking organizer saat ini follower count, bukan algoritma reputasi.
