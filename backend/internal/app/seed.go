package app

import (
	"context"
	"encoding/json"
	"time"
)

// Seed is idempotent and development-only. Every person, venue and account is fictional.
func (a *App) Seed(ctx context.Context) error {
	tx, e := a.DB.Begin(ctx)
	if e != nil {
		return e
	}
	defer tx.Rollback(ctx)
	if _, e = tx.Exec(ctx, `SELECT pg_advisory_xact_lock(7465822)`); e != nil {
		return e
	}
	var seeded bool
	e = tx.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM accounts WHERE id='customer-demo')`).Scan(&seeded)
	if e != nil {
		return e
	}
	if seeded {
		return tx.Commit(ctx)
	}
	accounts := [][5]string{{"customer-demo", "customer", "Nadia Ananda", "nadia@example.test", "Penonton teater, pencatat cerita, warga Karawang."}, {"customer-rani", "customer", "Rani Wulandari", "rani@example.test", ""}, {"customer-bagas", "customer", "Bagas Setiawan", "bagas@example.test", ""}, {"customer-raka", "customer", "Raka Pratama", "raka@example.test", ""}, {"customer-alya", "customer", "Alya Putri", "alya@example.test", ""}, {"customer-dimas", "customer", "Dimas", "dimas@example.test", ""}, {"organizer-demo", "organizer", "Teater Ruang", "ruang@example.test", ""}, {"organizer-nada", "organizer", "Nada Kolektif", "nada@example.test", ""}, {"organizer-sudut", "organizer", "Panggung Sudut", "sudut@example.test", ""}, {"admin-demo", "admin", "Moderator Lokal", "admin@example.test", ""}}
	for _, u := range accounts {
		if _, e = tx.Exec(ctx, `INSERT INTO accounts(id,role,name,email,bio,onboarding_done) VALUES($1,$2,$3,$4,$5,$6)`, u[0], u[1], u[2], u[3], u[4], u[1] != "customer"); e != nil {
			return e
		}
	}
	orgs := [][7]string{{"ruang", "organizer-demo", "teater-ruang", "Teater Ruang", "Teater", "Ruang bertemu, berproses, dan bercerita lewat teater. Terbuka buat penonton, pemain, dan siapa pun yang penasaran.", "Kami adalah kelompok teater yang tumbuh dari pertemuan dan cerita sehari-hari di Karawang. Di sini, kami berbagi proses latihan, ngobrol soal karya, dan mendengarkan pengalaman penonton."}, {"nada", "organizer-nada", "nada-kolektif", "Nada Kolektif", "Musik", "Musisi, pendengar, dan cerita di balik lagu. Dari ruang latihan ke panggung kecil di kota sendiri.", "Ruang bersama untuk musisi dan pendengar musik independen. Kami mengadakan pertunjukan kecil, sesi dengar, dan percakapan tentang proses menulis lagu."}, {"sudut", "organizer-sudut", "panggung-sudut", "Panggung Sudut", "Teater", "Pertunjukan intim, percakapan panjang. Yuk kenalan dengan seni pertunjukan dari dekat.", "Kami ingin mempertemukan penonton dan seniman lewat pertunjukan berskala kecil, diskusi, serta lokakarya."}}
	for _, o := range orgs {
		if _, e = tx.Exec(ctx, `INSERT INTO organizers(id,account_id,slug,name,category,description,about,cover_url,created_at) VALUES($1,$2,$3,$4,$5,$6,$7,'/assets/rehearsal.webp',$8)`, o[0], o[1], o[2], o[3], o[4], o[5], o[6], time.Date(map[string]int{"ruang": 2021, "nada": 2022, "sudut": 2023}[o[0]], 1, 1, 0, 0, 0, 0, time.UTC)); e != nil {
			return e
		}
		if _, e = tx.Exec(ctx, `INSERT INTO payment_methods(id,organizer_id,kind,provider,number,holder,note) VALUES($1,$2,'bank','Bank Contoh','0000000000','DEMO - JANGAN TRANSFER','Data fiktif untuk mencoba alur lokal.')`, "bank-"+o[0], o[0]); e != nil {
			return e
		}
	}
	if _, e = tx.Exec(ctx, `INSERT INTO payment_methods(id,organizer_id,kind,provider,number,holder,note) VALUES('wallet-ruang','ruang','wallet','Wallet Contoh','080000000000','DEMO - JANGAN TRANSFER','Data fiktif.')`); e != nil {
		return e
	}
	loc, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Now().In(loc)
	soon := now.Add(time.Hour).Truncate(time.Minute)
	past := now.AddDate(0, 0, -30)
	music := time.Date(now.Year(), now.Month(), now.Day()+7, 16, 0, 0, 0, loc)
	festival := time.Date(now.Year(), now.Month(), now.Day()+14, 15, 0, 0, 0, loc)
	type event struct {
		id, slug, oid, title, category, venue, description, flyer, layout string
		lineup                                                            []Performer
		start                                                             time.Time
		price                                                             int64
	}
	es := []event{{"e1", "di-balik-layar", "ruang", "Di Balik Layar", "Teater", "Ruang Pertunjukan Tengah", "Sebuah pertunjukan tentang rumah, jarak, dan hal-hal yang belum sempat kita ucapkan. Pertemuan tiga orang dalam satu malam membuka kembali cerita yang mereka simpan bertahun-tahun.", "/assets/rehearsal.webp", "/assets/venue-layout.svg", []Performer{{"Alya Pramesti", "Aktor · Ratih", "/assets/lineup/artist-woman-01.webp"}, {"Bagas Wiratama", "Aktor · Bima", "/assets/lineup/artist-man-01.webp"}, {"Dimas Kurnia", "Sutradara", "/assets/lineup/artist-man-02.webp"}, {"Nara Putri", "Penata & pemain musik", "/assets/lineup/artist-woman-02.webp"}}, soon, 45000}, {"e2", "sore-mendengar", "nada", "Sore Mendengar", "Musik", "Halaman Nada", "Sore untuk mendengarkan karya musisi lokal dan cerita di balik lagu mereka. Pertunjukan akustik dengan ruang untuk ngobrol setelah set berakhir.", "", "", []Performer{{"Senja Sebentar", "Band · set akustik", "/assets/lineup/band-01.webp"}, {"Laras & Kawan", "Band · folk", "/assets/lineup/band-02.webp"}}, music, 35000}, {"e3", "festival-panggung-kecil", "sudut", "Festival Panggung Kecil", "Teater", "Ruang Seni Sudut", "Tiga pertunjukan, satu ruang pertemuan. Pilih sesi yang ingin kamu tonton, lalu ikut percakapan dengan para pembuatnya.", "/assets/rehearsal.webp", "", []Performer{{"Teater Langkah", "Rumah yang Berjalan", "/assets/lineup/band-01.webp"}, {"Kolektif Esok", "Surat untuk Esok", "/assets/lineup/band-02.webp"}, {"Panggung Sudut", "Percakapan Terakhir", "/assets/lineup/band-01.webp"}}, festival, 30000}, {"e4", "rumah-yang-kita-bawa", "ruang", "Rumah yang Kita Bawa", "Teater", "Ruang Pertunjukan Tengah", "Pertunjukan tentang pulang dan ingatan. Terima kasih sudah menjadi bagian dari pertemuan ini. Obrolan tentang karya masih terbuka di ruang komunitas.", "/assets/rehearsal.webp", "", []Performer{{"Rani Mahesa", "Aktor", "/assets/lineup/artist-woman-02.webp"}, {"Dimas Kurnia", "Sutradara", "/assets/lineup/artist-man-02.webp"}}, past, 40000}}
	for _, v := range es {
		duration, age, label, minutes := "90 menit", "15+", "Pertunjukan malam", 90
		if v.id == "e2" {
			duration, age, label, minutes = "150 menit", "Semua usia", "Sesi sore", 150
		}
		if v.id == "e3" {
			duration, label = "Per sesi 60–90 menit", "Rumah yang Berjalan"
		}
		if v.id == "e4" {
			duration, minutes = "80 menit", 80
		}
		lineup, _ := json.Marshal(v.lineup)
		if _, e = tx.Exec(ctx, `INSERT INTO events(id,slug,organizer_id,title,category,city,venue,description,flyer_url,layout_url,lineup,published,duration,age) VALUES($1,$2,$3,$4,$5,'Karawang',$6,$7,$8,$9,$10,true,$11,$12)`, v.id, v.slug, v.oid, v.title, v.category, v.venue, v.description, v.flyer, v.layout, string(lineup), duration, age); e != nil {
			return e
		}
		if _, e = tx.Exec(ctx, `INSERT INTO event_sessions(id,event_id,label,starts_at,ends_at,price,capacity) VALUES($1,$2,$6,$3,$4,$5,60)`, "s-"+v.id, v.id, v.start, v.start.Add(time.Duration(minutes)*time.Minute), v.price, label); e != nil {
			return e
		}
	}
	for i, label := range []string{"Surat untuk Esok", "Percakapan Terakhir"} {
		start := festival.Add(time.Duration(i+1) * 2 * time.Hour)
		if i == 1 {
			start = start.Add(30 * time.Minute)
		}
		if _, e = tx.Exec(ctx, `INSERT INTO event_sessions(id,event_id,label,starts_at,ends_at,price,capacity) VALUES($1,'e3',$2,$3,$4,$5,40)`, []string{"s-e3b", "s-e3c"}[i], label, start, start.Add(90*time.Minute), 35000+i*5000); e != nil {
			return e
		}
	}
	posts := [][8]string{{"p1", "organizer-demo", "ruang", "Sebelum lampu panggung menyala.", "Malam tadi kami mencoba adegan terakhir dengan cara yang berbeda. Ternyata, satu jeda kecil bisa mengubah rasa seluruh cerita.\n\nKalau nonton teater, kalian lebih suka akhir yang tuntas atau yang bikin kepikiran sampai pulang?", "/assets/rehearsal.webp", "", ""}, {"p2", "customer-raka", "", "Datang nonton teater sendirian, pernah?", "Lagi pengin nonton pentas berikutnya, tapi teman-teman belum bisa ikut. Ada yang punya pengalaman nonton sendiri? Penasaran suasananya dan apa yang perlu disiapkan.", "", "", ""}, {"p3", "organizer-nada", "nada", "Lagu-lagunya selesai, obrolannya belum.", "Terima kasih sudah datang ke sesi dengar kemarin. Lagu mana yang masih kalian putar di kepala? Ceritain di sini, kami pengin tahu.", "", "", ""}, {"p4", "organizer-demo", "ruang", "Bawa sedikit cerita pulang.", "Kaus dan tote dari Rumah yang Kita Bawa masih tersedia. Hasil penjualannya membantu proses produksi pentas berikutnya. Yuk lihat koleksinya.", "/assets/merchandise.webp", "https://example.com/merch", "Lihat merchandise"}}
	for i, p := range posts {
		var oid any
		if p[2] != "" {
			oid = p[2]
		}
		if _, e = tx.Exec(ctx, `INSERT INTO posts(id,account_id,organizer_id,title,body,image_url,link_url,link_label,city,created_at) VALUES($1,$2,$3,$4,$5,$6,$7,$8,'Karawang',$9)`, p[0], p[1], oid, p[3], p[4], p[5], p[6], p[7], now.Add(-time.Duration(i+1)*6*time.Hour)); e != nil {
			return e
		}
	}
	for _, m := range [][3]string{{"p1", "event", "e1"}, {"p2", "organizer", "sudut"}, {"p4", "event", "e4"}} {
		if _, e = tx.Exec(ctx, `INSERT INTO post_mentions(post_id,kind,target_id) VALUES($1,$2,$3)`, m[0], m[1], m[2]); e != nil {
			return e
		}
	}
	for _, oid := range []string{"ruang", "nada"} {
		if _, e = tx.Exec(ctx, `INSERT INTO follows(account_id,organizer_id) VALUES('customer-demo',$1)`, oid); e != nil {
			return e
		}
	}
	if _, e = tx.Exec(ctx, `INSERT INTO votes(post_id,account_id) VALUES('p1','customer-rani'),('p1','customer-bagas'),('p2','customer-rani'); INSERT INTO comments(id,post_id,account_id,body) VALUES('c1','p1','customer-alya','Tim akhir yang bikin kepikiran! Rasanya jadi ada obrolan yang bisa dibawa pulang.'); INSERT INTO comments(id,post_id,account_id,parent_id,body) VALUES('c2','p1','organizer-demo','c1','Wah, menarik. Justru jeda setelah pentas itu yang lagi kami pikirkan juga.'); INSERT INTO comments(id,post_id,account_id,body) VALUES('c3','p2','customer-dimas','Pernah, dan ternyata enak bisa fokus. Biasanya habis pentas juga ada yang ngobrol di luar venue.')`); e != nil {
		return e
	}
	for _, p := range [][3]string{{"m1", "Kaus Rumah yang Kita Bawa", "Kaus hitam berbahan katun untuk menemani hari-hari di luar panggung. Merchandise dari pertunjukan Rumah yang Kita Bawa."}, {"m2", "Tote Cerita Pulang", "Tote kanvas untuk membawa buku, catatan latihan, dan cerita. Tetap tersedia setelah pertunjukan selesai."}} {
		if _, e = tx.Exec(ctx, `INSERT INTO products(id,organizer_id,event_id,name,description,price,image_url,variants,purchase_url) VALUES($1,'ruang','e4',$2,$3,$4,'/assets/merchandise.webp',$5,'https://example.com/merch')`, p[0], p[1], p[2], map[string]int{"m1": 125000, "m2": 75000}[p[0]], map[string]string{"m1": `["S","M","L","XL"]`, "m2": `["Satu ukuran"]`}[p[0]]); e != nil {
			return e
		}
	}
	if _, e = tx.Exec(ctx, `INSERT INTO orders(id,account_id,session_id,quantity,total,status,idempotency_key) VALUES('order-past','customer-rani','s-e4',1,40000,'approved','seed-past'); INSERT INTO order_history(order_id,status,note) VALUES('order-past','approved','Pesanan contoh untuk ulasan dengan bukti kehadiran.'); INSERT INTO reviews(id,event_id,account_id,body) VALUES('r1','e4','customer-rani','Bagian akhirnya masih kepikiran sampai sekarang. Pertunjukannya terasa dekat dengan pengalaman sendiri. Tempatnya nyaman, tapi beberapa dialog agak susah terdengar dari baris belakang.'),('r2','e4','customer-bagas','Ini pertama kali saya nonton teater. Panitianya membantu saat masuk dan diskusi setelah pentas bikin ceritanya makin kebayang. Pengin ikut pentas berikutnya.'); INSERT INTO review_replies(id,review_id,account_id,body) VALUES('rr1','r1','organizer-demo','Terima kasih sudah berbagi, Rani. Masukan soal suara kami catat untuk produksi berikutnya.')`); e != nil {
		return e
	}
	if _, e = tx.Exec(ctx, `INSERT INTO tickets(id,order_id,ordinal,token,checked_at,checked_by) VALUES('ticket-past','order-past',1,$1,$2,'organizer-demo')`, token(), past); e != nil {
		return e
	}
	if e = notify(ctx, tx, "customer-demo", "event", "Selamat datang di ruang komunitas", "Temukan proses kreatif dan pertunjukan di sekitarmu.", "/"); e != nil {
		return e
	}
	return tx.Commit(ctx)
}
