# PostgreSQL dan data

PostgreSQL 17 direkomendasikan untuk local. Go memakai pgx v5, parameterized SQL, pool max 8. Tidak ada ORM. Native PostgreSQL diperlukan untuk development user; PGlite hanya dipakai di lingkungan build untuk verifikasi sementara, bukan dependency aplikasi.

## Tabel

| Kelompok            | Tabel                                                     | Hubungan/kebijakan                                                                |
| ------------------- | --------------------------------------------------------- | --------------------------------------------------------------------------------- |
| Identitas           | accounts, auth_sessions, oauth_states                     | Google sub+role unique; session hash; state sekali pakai                          |
| Ruang/agenda        | organizers, events, event_sessions                        | Organizer satu per account; slug unik; sesiFK event, harga≥0, capacity>0              |
| Komunitas           | posts, post_mentions, votes, comments, follows, bookmarks | Vote/follow/bookmark unique per account-target; mention polimorfik divalidasi API |
| Media/merch         | uploads, products                                         | Purpose media/proof; produk bisa tertaut event, katalog eksternal                 |
| Pembayaran          | payment_methods, orders, order_history                    | SnapshotJSONB; account+idempotency unique; qty 1–6                                 |
| Tiket               | tickets                                                   | Randomtoken unique; order + ordinal unique; checked_at/by                           |
| Review              | reviews, review_replies, review_helpful                   | Event+account unique; badge dari checked ticket; hideadmin                        |
| Notifikasi/moderasi | notifications, reports, event_reminders, payment_reminders                   | Notif privat; statusreport; reminder account + sesi unique                          |
| Aktivitas           | activity_clicks                                         | Klik CTA tanpa identitas/IP; organizer + target terverifikasi                        |
| Migrasi             | schema_migrations                                         | Versi applied_at; advisorylock startup                                            |

Views organizer_view/event_view/post_view/review_view/product_view menyusun payload publik JSONB. Event_view menghitung availability aktif, post_view jumlah vote/comment dan mention, review_view badge berdasar check-in aktual. Tidak ada tiket/QR/proof di view publik.

Index meliputi city/category/published event, jadwal sesi, post lokal/waktu, target mention, reservasi sesi/status/expiry, account order, account notification dan account session.

## Migrasi

Migrate() memakai pg_advisory_xact_lock dan satu transaksi. V1 schema.sql membangun tabel/views/index; v2 reminders.sql menambah reminder acara; v3 checkout_continuity.sql menambah OAuth next_path, pengingat deadline order, votes.created_at dan activity_clicks. Versi hanya dicatat setelah SQL berhasil. Jalankan API untuk migrasi otomatis. Tidak ada rollback otomatis destructive.

Untuk perubahan berikutnya tambahkan fileSQL embedded dan versi selanjutnya pada Migrate; jangan mengubah isi migration1/2/3 setelah source mulai dipakai. Versi perubahan berikutnya adalah 4. Untuk data production kelak lakukan backup dan latihan migrasi pada salinan sebelum upgrade.

## Demo

Seed idempotent dev mempunyai customer, organizer, admin, tiga ruang, event teater/musik/festival+event selesai, beberapa post/comment/follow/review/merch dan tiket contoh sudahcheck-in. Jadwal dibuat relatif terhadap **seed pertama**; tidak bergeser setiap restart. Event demo tidak permanen berada di masa depan. Buat event baru melalui dashboard untuk latihan setelah jadwal lewat.

Jangan hapus database untuk memperbarui jadwal kecuali memang ingin menghapus seluruh data lokal. `docker compose down` tidak menghapus volume; `down -v` destructive dan bukan langkah setup rutin.

## Backup lokal

Contoh gunakan client pg_dump/pg_restore sesuai instalasi lokal:

```bash
pg_dump --format=custom --file=ruang.backup "postgres://ruang:ruang_local@localhost:5432/ruang"
# Restore hanya ke database kosong/terpisah yang sengaja dibuat untuk restore
pg_restore --dbname="postgres://ruang:ruang_local@localhost:5432/ruang_restore" ruang.backup
```

Backup juga UPLOAD_DIR (default var/uploads). BackupDB saja tidak mencakup foto/proof. Jangan membagikan backup berisi PII/proof/QR.

Riwayat order/proof lama, retention, orphan upload cleanup dan anonymization belum mempunyai automation produksi. File proof sebelumnya tidak terhapus saat koreksi; setelah diganti bukan lagi dapat diakses organizer melalui order, tetap owner-authorized.
