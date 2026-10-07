# API Go v0.1

Base path `/api`, JSON UTF-8 kecuali upload/download. Cookie HttpOnly digunakan web melalui proxy; API juga menerima `Authorization: Bearer <session-token>` jika klien sudah mempunyai session aplikasi. Google access token bukan session Rantaya.

Mutasi menerima Origin yang cocok dengan APP_URL; browser cross-site ditolak. CLI yang langsung memakai Go API dapat tanpa Origin. HTTP client yang memakai proxy SvelteKit sebaiknya menyertakan Origin APP_URL. JSON field tak dikenal/objek tambahan ditolak. Upload maksimal 8 MB, body HTTP 10 MB, timeout 20 detik, rate 600 request/IP/menit.

Error berbentuk `{"error":{"code":"invalid_input","message":"…"}}`. Status umum: 400 input, 401 login, 403 izin/origin, 404 tidak ada/tidak boleh melihat, 409 konflik status/kuota/reuse, 429 rate limit, 500 server. List adalah array langsung, detail object langsung. Mutasi umumnya 200, create post/upload 201. Discovery 20 item per halaman, `page=1..1000`.

## Identitas

| Method/path               | Akses            | Payload/query dan hasil                                            |
| ------------------------- | ---------------- | ------------------------------------------------------------------ |
| GET /health               | Publik           | status/service                                                     |
| GET /auth/config          | Publik           | demo/google flags, cities                                          |
| GET /auth/me              | Publik           | user/null, follows/votes/bookmarks/unread                          |
| POST /auth/demo           | Dev only         | `{role:customer\|organizer\|admin}` → cookie                       |
| POST /auth/logout         | Session jika ada | `{}` → hapus session/cookie                                        |
| GET /auth/google/start    | Publik           | `role=customer\|organizer&next=/path` →302 Google                             |
| GET /auth/google/callback | OAuth state      | code,state →303 tujuan aman/onboarding                                    |
| PATCH /profile            | Login            | name,city,bio,avatar_url,interests[],preferences{},onboarding_done |

## Ruang, agenda dan komunitas

| Method/path                  | Akses                   | Payload/query                                                                                   |
| ---------------------------- | ----------------------- | ----------------------------------------------------------------------------------------------- |
| GET /organizers              | Publik                  | city,q,page                                                                                     |
| GET /organizers/{id}         | Publik                  | ID atau slug                                                                                    |
| PATCH /organizer             | Organizer               | name,description,about,city,category,avatar_url,cover_url                                       |
| POST /organizers/{id}/follow | Login                   | `{active:true\|false}`                                                                          |
| GET /events                  | Publik/pemilik draft    | city,category,organizer,q,period=upcoming\|past\|all,from,to,page                               |
| GET /events/{id}             | Publik/pemilik draft    | ID atau slug                                                                                    |
| POST /events                 | Organizer               | EventInput; returns full event 200                                                              |
| PUT /events/{id}             | Pemilik                 | EventInput; full replacement                                                                    |
| GET /posts                   | Publik                  | city,q,category,sort=latest\|popular,tab=following,organizer,mode=official\|mentions,event,page |
| GET /posts/{id}              | Publik                  | Nonhidden post                                                                                  |
| POST /posts                  | Login                   | title?,body,image_url?,city,link_url?,link_label?,mentions[{kind,id}]; link hanya organizer     |
| DELETE /posts/{id}           | Penulis                 | Soft hide                                                                                       |
| PATCH /posts/{id}/pin        | Organizer pemilik       | `{pinned:true\|false}`                                                                          |
| POST /posts/{id}/vote        | Login                   | `{active:true\|false}`                                                                          |
| GET /posts/{id}/comments     | Publik                  | Flat list parent_id untuk reply                                                                 |
| POST /posts/{id}/comments    | Login                   | `{body,parent_id?}`; satu tingkat reply                                                         |
| GET /mentions                | Publik                  | q; saran event/ruang                                                                            |
| GET /bookmarks               | Login                   | Array `{kind,item}` tersimpan                                                                                 |
| POST /bookmarks              | Login                   | `{kind:post\|event,id,active}`                                                                  |
| GET /reviews                 | Publik                  | event,organizer; max 100                                                                        |
| POST /reviews                | Customer                | `{event_id,body,attended:true}`; upsert setelah event selesai                                   |
| POST /reviews/{id}/reply     | Organizer pemilik event | `{body}`                                                                                        |
| POST /reviews/{id}/helpful   | Login                   | `{active:true\|false}`                                                                          |

`from/to` YYYY-MM-DD menggunakan hari WIB. City berupa teks eksplisit, bukan radius. `period=upcoming` mencakup event yang belum selesai, termasuk sedang berlangsung.

Contoh EventInput:

```json
{
  "title": "Malam Panggung",
  "description": "Pertunjukan dan obrolan bersama komunitas.",
  "category": "Teater",
  "city": "Karawang",
  "venue": "Gedung seni",
  "address": "Alamat venue",
  "maps_url": "https://www.google.com/maps?q=Karawang",
  "duration": "90 menit",
  "language": "Bahasa Indonesia",
  "age": "15+",
  "flyer_url": "",
  "trailer_url": "",
  "layout_url": "",
  "lineup": [{ "name": "Nama seniman", "role": "Sutradara", "photo": "" }],
  "published": true,
  "sessions": [
    {
      "id": "",
      "label": "Sabtu malam",
      "starts_at": "2027-01-10T19:00:00+07:00",
      "ends_at": "2027-01-10T20:30:00+07:00",
      "price": 45000,
      "capacity": 60
    }
  ]
}
```

Create memakai session.id kosong; update mengirim semua sesi yang dipertahankan. Jadwal/harga sesi dengan order dikunci. Kategori: Teater/Musik/Komedi/Lainnya. Media menerima `/assets/…`, upload milik sendiri purpose media, atau HTTPS. UI trailer memainkan YouTube/mp4/webm.

## Merch dan media

| Method/path        | Akses                     | Payload/query                                                                   |
| ------------------ | ------------------------- | ------------------------------------------------------------------------------- |
| GET /products      | Publik                    | organizer,q; max 100                                                            |
| GET /products/{id} | Publik                    | Detail katalog                                                                  |
| POST /products     | Organizer                 | name,description,price,image_url,variants[],availability,purchase_url,event_id? |
| PUT /products/{id} | Pemilik                   | Payload lengkap                                                                 |
| POST /uploads      | Login                     | multipart file,purpose=media\|proof;201 `{id,url,size}`                         |
| GET /uploads/{id}  | Media publik/proof privat | JPG/PNG/WebP; proof hanya pembeli/pengelola order                               |

Produk event_id opsional harus milik organizer. Pembelian merch diarahkan ke purchase_url. SVG upload tidak diterima; denah SVG bundled digunakan hanya sebagai aset demo.

## Transaksi dan QR

| Method/path                | Akses                 | Payload/query                                                |
| -------------------------- | --------------------- | ------------------------------------------------------------ |
| GET /payment-methods       | Publik aktif/pemilik  | organizer=ID; pemilik juga melihat disabled                  |
| POST /payment-methods      | Organizer             | kind=bank\|wallet,provider,number,holder,note,enabled        |
| PUT /payment-methods/{id}  | Pemilik               | Full payment master                                          |
| GET /orders                | Login                 | Order sendiri/customer atau event sendiri/organizer; max 100 |
| POST /orders               | Customer              | `{session_id,quantity,idempotency_key}` →order200            |
| GET /orders/{id}           | Pembeli/pemilik event | event,session,customer,snapshot,history,tickets              |
| PATCH /orders/{id}/payment | Pembeli unpaid aktif  | `{payment_method_id}`                                        |
| POST /orders/{id}/proof    | Pembeli               | `{upload_id}`; proof milik sendiri                           |
| POST /orders/{id}/cancel   | Pembeli unpaid        | `{confirm_unpaid:true}`                                      |
| POST /orders/{id}/review   | Organizer pemilik     | `{action:approve\|correction,note,received_funds}`           |
| GET /tickets               | Customer              | Admission dari approved order                                |
| POST /check-in             | Organizer pemilik     | `{event_id,code}`; QR atau token manual                      |

Satu order satu sesi, qty 1–6, key8–100 karakter. Approve membutuhkan received_funds=true; correction note 10–500 karakter. QR `ruang:ticket:<64hex>`. Lihat PAYMENTS.md untuk status/kapasitas.

## Notifikasi, dashboard dan moderasi

| Method/path               | Akses     | Payload/query                          |
| ------------------------- | --------- | -------------------------------------- |
| GET /notifications        | Login     | Max100, privat                         |
| POST /notifications/read  | Login     | `{id:""}` semua atau ID tertentu       |
| GET /dashboard            | Organizer | Counts, revenue, check-ins aktual      |
| POST /activity/click      | Publik    | `{kind:ticket\|merch,target_id}`;201   |
| POST /reports             | Login     | `{kind:post\|review,target_id,reason}` |
| GET /admin/reports        | Admin     | Daftar laporan                         |
| PATCH /admin/reports/{id} | Admin     | `{status:reviewed\|hidden}`            |

Field respons publik ada pada `web/src/lib/types.ts` dan view SQL; respons order pada `orderSelect` di commerce.go. Null payment/proof/expiry/checked_at harus ditangani client. Bukan OpenAPI machine-generated contract; API versioning path dan cursor pagination belum dibuat.

## Contoh CLI lokal

```bash
curl -c cookies.txt -H 'Content-Type: application/json' \
 -d '{"role":"customer"}' http://127.0.0.1:8080/api/auth/demo
curl -b cookies.txt http://127.0.0.1:8080/api/auth/me
curl 'http://127.0.0.1:8080/api/events?city=Karawang&period=upcoming'
```

Cookie file adalah session aktif; jangan commit/bagikan.

## Statistik aktual dan konteks review

Dashboard: followers, posts, events, interactions, ticket_clicks, merch_clicks, activity, pending_payments, tickets, checked_in, revenue. Interactions menjumlah vote yang masih tersimpan dan komentar pada post resmi yang tidak disembunyikan. Activity berisi tujuh `{date,count}` hari WIB: post, vote dan komentar pada ruang milik akun. Click menerima target event published atau produk tersedia; tidak menyimpan account/IP dan bukan event penjualan. Rate/origin/body guard tetap berlaku.

GET reviews juga menyertakan event_starts_at dan organizer_name agar tanggal/nama ruang pada kartu review selalu merujuk event, bukan tanggal review. GET bookmarks berupa array `{kind,item}`; item null untuk target yang tidak lagi publik dan tidak ditampilkan. Detail post dan bookmark menyelesaikan event mention melalui endpoint detail, termasuk event selesai di luar discovery mendatang.

Profile preferences mendukung events/replies/reminders/digest. Digest disimpan sebagai preferensi; pengiriman email belum diimplementasikan.
