# Pengujian

## Pemeriksaan rutin

```bash
npm run check
npm run build
npm test
```

check menjalankan Svelte/TypeScript, go vet dan unit tests. build menghasilkan frontend adapter-node dan binary Go. Unit mencakup URL policies, JSON strict/trailing content, token opaque. TestCommerceIntegration otomatis skip bila TEST_DATABASE_URL kosong; `npm test` tanpa variabel itu **bukan bukti flow database sudah diuji**.

## Integration dengan database terpisah

Buat database khusus uji dengan user yang punya DDL; jangan gunakan DB pengguna/production. Tests membuat account, organizer, event, metode dan transaksi sintetis dengan ID unik; data tetap ada setelah tes untuk inspeksi. Gunakan DB disposable.

macOS/Linux:

```bash
TEST_DATABASE_URL='postgres://ruang:ruang_local@localhost:5432/ruang_test?sslmode=disable' npm test
```

Windows PowerShell:

```powershell
$env:TEST_DATABASE_URL='postgres://ruang:ruang_local@localhost:5432/ruang_test?sslmode=disable'
npm test
Remove-Item Env:TEST_DATABASE_URL
```

Test meliputi: idempotent order, role/ownership, resume dan snapshot rekening, proof privacy, larangan cancel setelah proof, konfirmasi dana, correction/reproof, repeated approve, qty QR, gate event/organizer, reuse, delapan check-in bersamaan untuk satu QR (satu berhasil, tujuh ditolak), cancel/expiry, delapan checkout untuk empat kursi, reminder sekali, review before-end ditolak, badge check-in vs self-report, serta larangan hapus post orang lain. Tambahkan `-race` pada `go test` untuk race detector.

## Browser

Jalankan aplikasi dengan .env demo/seed pada database uji, bukan data pengguna. `cd web && npx playwright install chromium`, kembali ke root, lalu `npm run test:e2e -- --project=mobile390`. E2E_BASE_URL opsional bila host/port berbeda; API APP_URL dan origin build harus sama dengan browser URL.

Tersedia sebelas test per viewport: mobile 360×800, mobile 390×844, desktop 1440×1000. Test mencakup public SSR/status/h1/overflow/console, composer dengan mention inline/click, onboarding skip, navigasi dashboard dan penyimpanan rekening/event/lineup/foto/Maps/merch/profil. Flow pembayaran memakai browser context terpisah: history/resume, upload bukti, akses privat, notifikasi, correction/reproof, approval, dua QR, check-in sekali dan reuse ditolak. E2E_BROWSER_PATH opsional untuk Chromium yang sudah terpasang; pengguna normal tidak perlu mengisi.

Untuk viewport lain ganti `--project=mobile360` atau `--project=desktop`. Tanpa argumen semua proyek dijalankan. Burst tiga proyek dapat mencapai limiter Go 600 request/IP/menit karena seluruh SSR/proxy berasal dari loopback. Jika mendapat 429, beri jeda sampai jendela rate limit lewat atau restart API uji sebelum proyek berikutnya, dengan DB yang sama. Jangan menaikkan/mematikan limiter hanya untuk mengubah hasil tes menjadi hijau.

Regresi checkout HTTP LAN ada pada `web/tests/checkout-lan.spec.ts`: pembelian tetap menuju pembayaran ketika `crypto.randomUUID` tidak tersedia, retry setelah respons gagal tetap mengembalikan order pertama, dan kegagalan pembuat kunci melepaskan status tombol sibuk. Jalankan pada database uji dengan `E2E_BASE_URL` berupa URL HTTP IP LAN, misalnya `E2E_BASE_URL=http://192.168.1.10:5174 npm run test:e2e -- checkout-lan.spec.ts`. Untuk LAN, browser harus melaporkan `isSecureContext=false`; ini berbeda dengan pengecualian secure context pada localhost.

## Hasil build ini

Recovery 7 Oktober 2026: pemeriksaan Svelte/TypeScript 0 error, 0 warning; build adapter-node dan binary Go berhasil; go vet/unit serta `go test -race ./... -count=1 -v` dengan integration SQL berhasil. SQL memakai PostgreSQL WASM/PGlite via wire protocol dengan pool satu koneksi. Native PostgreSQL 17.11 berhasil diekstrak tetapi proses non-root tidak dapat dijalankan karena lingkungan hanya memetakan UID 0. Delapan request tetap bersamaan pada level handler; **uji row lock beberapa koneksi native perlu diulang di laptop/CI**.

Seluruh 15 kombinasi test/viewport berhasil: 10 mobile lolos pada run pertama; burst desktop sempat terkena limiter dan kelima test desktop lolos setelah restart API pada DB uji baru. Tidak ada perubahan kode bisnis untuk menghindari limiter. Test restart API/SSR dengan data tetap membuktikan migrasi tidak digandakan dan status order/tiket dipertahankan. robots privat dan penolakan multipart cross-origin HTTP 403 juga diperiksa. Catatan terstruktur ada pada verification.json, ringkasan audit di RECOVERY_AUDIT.md, dan hasil tanpa token/proof di evidence/recovery-checks.json.

Google OAuth end-to-end belum dijalankan tanpa credential pengguna. Kamera fisik/izin perangkat, keyboard mobile nyata, email/push, Worker production dan native app belum diuji/diimplementasikan sesuai scope. QR backend/manual check-in diuji; camera decoder UI disediakan.

Lighthouse belum diklaim memiliki skor tertentu. Jalankan pada production build ketika server API dan database native stabil; ukur Performance/Accessibility/Best Practices/SEO pada halaman publik. Skor bergantung data/gambar/perangkat/hosting, bukan framework saja.


## Penyelarasan UI Rantaya — 7 Oktober 2026

Checklist desain di UI_UX.md dan prototype historis di PROTOTYPE_REFERENCE.md. Svelte app dan preview 0 error/0 warning; go vet/unit dan Go race SQL integration lulus; kedua build produksi berhasil. 18 tes browser (6 × 360/390/1440) lulus termasuk dashboard, composer, payment dua role, proof privat dan QR reuse. Sesudah inspeksi screenshot menemukan avatar organizer yang melebar, selector diperbaiki dan tiga tes layout diulang, semua lulus; breakpoint tablet 900px juga diperiksa. Foto tetap berwarna, CTA tiket tampak pada viewport awal, filter kalender menyimpan draft hingga Terapkan, dan deskripsi/follow mobile sesuai urutan prototype.

Screenshot local production diperiksa untuk composer/agenda/event/pengelola mobile, event desktop dan kalender mobile. Test mention juga mensimulasikan visualViewport keyboard 300px; ini bukan bukti keyboard fisik. Runtime browser dipulihkan memakai Chromium headless 141 yang tersedia; tidak mengubah dependency aplikasi. SQL kembali memakai PGlite satu koneksi, migrasi 1 dan 2. Tidak ada perubahan kode Go/skema/alur otorisasi. Bukti ringkas tanpa token/private proof: evidence/ui-alignment-checks.json.

Preview Worker memverifikasi 15 halaman SSR, role demo dan payment mutation 403; data tetap contoh/read-only. Tidak ada klaim pixel-identical untuk semua halaman, pengujian Google OAuth/kamera fisik/Lighthouse, atau native PostgreSQL beberapa koneksi.

## Hasil akhir kelanjutan checkpoint (otoritatif)

Build adapter-node/API Go dan build preview terpisah lulus; Svelte/TypeScript 0 error dan 0 warning; go vet, unit dan Go race/SQL integration lulus. Migrasi 1, 2 dan 3 terpasang pada database uji. Deadline payment diulang dua kali hanya menghasilkan satu notifikasi; metrik klik dibatasi organizer. Flow login/onboarding menjaga query sesi.

Sebelas kasus × tiga viewport menghasilkan **32 kombinasi lulus + 1 skip sengaja**: pengukuran mobile tidak dijalankan di desktop. Kasus mencakup dashboard, semantic/responsiveness, composer, isolasi akun, pembayaran/proof/correction/approval/QR, continuity, review/preferensi, geometri SELA, upvote/rollback dan kontrol home, plus layout/kalender. Burst pernah terkena 429: pengujian kemudian dibagi per file dengan restart API pada database tetap, tanpa perubahan limiter.

Crawl tambahan pada 68 keadaan route/tab × 360/390/1440 menghasilkan **204 HTTP200 dengan satu h1, tanpa overflow dan tanpa pageerror**. Detail post dan bookmark mempertahankan event mention. Setelah inspeksi terakhir, 30 pemeriksaan ulang pada profil/review/notifikasi/dashboard/merch lulus, termasuk varian produk pertama tidak kosong. Kasus continuity yang berubah dicek ulang; tidak mengulang flow lain yang inputnya tetap.

Preview menguji 30 halaman SSR, preferensi view vote/bookmark per role, category mention dan penolakan approval403. Itu tidak membuktikan pembayaran production; fixture bukan DB pengguna. Source Git demo juga menyertakan local-app agar backend/dokumentasi dapat dipulihkan.

Hasil terstruktur terbaru: verification.json, evidence/final-checks.json, evidence/home-alignment-checks.json. Evidence recovery/ui-alignment sebelumnya adalah checkpoint historis. Batas native PostgreSQL/OAuth/kamera/Lighthouse di atas tetap berlaku. Foto, layout, header dan komponen diperiksa visual; tidak ada sertifikasi semua screenshot identik untuk setiap konten/perangkat.
