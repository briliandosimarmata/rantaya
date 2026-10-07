package app

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"net/http"
	"strings"
	"time"
)

const orderSelect = `SELECT to_jsonb(x)-'idempotency_key'||jsonb_build_object('event',ev.doc,'session',to_jsonb(s),'customer',jsonb_build_object('name',a.name,'email',a.email),'history',COALESCE((SELECT jsonb_agg(to_jsonb(h) ORDER BY h.id) FROM order_history h WHERE h.order_id=x.id),'[]'),'tickets',COALESCE((SELECT jsonb_agg(to_jsonb(t)-'token'||jsonb_build_object('qr','ruang:ticket:'||t.token) ORDER BY t.ordinal) FROM tickets t WHERE t.order_id=x.id),'[]')) FROM orders x JOIN event_sessions s ON s.id=x.session_id JOIN event_view ev ON ev.id=s.event_id JOIN accounts a ON a.id=x.account_id`

func (a *App) orders(w http.ResponseWriter, r *http.Request) error {
	u, e := a.me(r)
	if e != nil {
		return e
	}
	if e = a.Expire(r.Context()); e != nil {
		return e
	}
	q := orderSelect + ` WHERE x.account_id=$1 ORDER BY x.created_at DESC LIMIT 100`
	if u["role"] == "organizer" {
		q = orderSelect + ` WHERE ev.organizer_id=$1 ORDER BY x.created_at DESC LIMIT 100`
		items, e := a.list(r.Context(), q, u["organizer_id"])
		if e != nil {
			return e
		}
		send(w, 200, items)
		return nil
	}
	items, e := a.list(r.Context(), q, u["id"])
	if e != nil {
		return e
	}
	send(w, 200, items)
	return nil
}
func (a *App) order(w http.ResponseWriter, r *http.Request) error {
	u, e := a.me(r)
	if e != nil {
		return e
	}
	if e = a.Expire(r.Context()); e != nil {
		return e
	}
	doc, e := a.one(r.Context(), orderSelect+` WHERE x.id=$1 AND (x.account_id=$2 OR ev.organizer_id=COALESCE((SELECT id FROM organizers WHERE account_id=$2),''))`, r.PathValue("id"), u["id"])
	if e != nil {
		return e
	}
	send(w, 200, doc)
	return nil
}
func (a *App) createOrder(w http.ResponseWriter, r *http.Request) error {
	u, e := a.actor(r, "customer")
	if e != nil {
		return e
	}
	var in struct {
		Session  string `json:"session_id"`
		Quantity int    `json:"quantity"`
		Key      string `json:"idempotency_key"`
	}
	if e = decode(r, &in); e != nil {
		return e
	}
	if in.Quantity < 1 || in.Quantity > 6 || !textOK(in.Key, 8, 100) {
		return bad("Pilih sesi, jumlah 1–6 tiket, dan kunci permintaan.")
	}
	ctx := r.Context()
	tx, e := a.DB.Begin(ctx)
	if e != nil {
		return e
	}
	defer tx.Rollback(ctx)
	var eventID, oid, title string
	var starts time.Time
	var price int64
	var capacity int
	var published bool
	e = tx.QueryRow(ctx, `SELECT s.event_id,e.organizer_id,e.title,s.starts_at,s.price,s.capacity,e.published FROM event_sessions s JOIN events e ON e.id=s.event_id WHERE s.id=$1 FOR UPDATE OF s`, in.Session).Scan(&eventID, &oid, &title, &starts, &price, &capacity, &published)
	if e != nil {
		return e
	}
	var existing string
	e = tx.QueryRow(ctx, `SELECT id FROM orders WHERE account_id=$1 AND idempotency_key=$2`, u["id"], in.Key).Scan(&existing)
	if e == nil {
		if e = tx.Commit(ctx); e != nil {
			return e
		}
		r.SetPathValue("id", existing)
		return a.order(w, r)
	}
	if !errors.Is(e, pgx.ErrNoRows) {
		return e
	}
	if !published || !starts.After(time.Now()) {
		return bad("Sesi sudah dimulai atau belum diterbitkan.")
	}
	var reserved int
	e = tx.QueryRow(ctx, `SELECT COALESCE(sum(quantity),0) FROM orders WHERE session_id=$1 AND (status IN ('awaiting_review','correction_requested','approved') OR status='awaiting_payment' AND expires_at>now())`, in.Session).Scan(&reserved)
	if e != nil {
		return e
	}
	if in.Quantity > capacity-reserved {
		return conflict("Kuota tiket tidak mencukupi. Pilih sesi atau jumlah lain.")
	}
	if price > 0 {
		var enabled bool
		e = tx.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM payment_methods WHERE organizer_id=$1 AND enabled)`, oid).Scan(&enabled)
		if e != nil {
			return e
		}
		if !enabled {
			return bad("Pengelola belum menyediakan metode pembayaran.")
		}
	}
	id := ID()
	status := "awaiting_payment"
	var expiry any = time.Now().Add(time.Duration(a.Config.HoldMinutes) * time.Minute)
	if price == 0 {
		status = "approved"
		expiry = nil
	}
	_, e = tx.Exec(ctx, `INSERT INTO orders(id,account_id,session_id,quantity,total,status,expires_at,idempotency_key) VALUES($1,$2,$3,$4,$5,$6,$7,$8)`, id, u["id"], in.Session, in.Quantity, price*int64(in.Quantity), status, expiry, in.Key)
	if e != nil {
		return e
	}
	if _, e = tx.Exec(ctx, `INSERT INTO order_history(order_id,status,actor_id,note) VALUES($1,$2,$3,'Pesanan dibuat.')`, id, status, u["id"]); e != nil {
		return e
	}
	if price == 0 {
		if e = issueTickets(ctx, tx, id, in.Quantity); e != nil {
			return e
		}
	}
	label := "Pesanan dibuat, lanjutkan pembayaran"
	if price == 0 {
		label = "Tiket gratis sudah terbit"
	}
	if e = notify(ctx, tx, u["id"].(string), "order", label, title, "/transaksi/"+id); e != nil {
		return e
	}
	if e = tx.Commit(ctx); e != nil {
		return e
	}
	r.SetPathValue("id", id)
	return a.order(w, r)
}
func issueTickets(ctx context.Context, tx pgx.Tx, order string, qty int) error {
	for i := 1; i <= qty; i++ {
		if _, e := tx.Exec(ctx, `INSERT INTO tickets(id,order_id,ordinal,token) VALUES($1,$2,$3,$4) ON CONFLICT(order_id,ordinal) DO NOTHING`, ID(), order, i, token()); e != nil {
			return e
		}
	}
	return nil
}
func (a *App) selectPayment(w http.ResponseWriter, r *http.Request) error {
	u, e := a.actor(r, "customer")
	if e != nil {
		return e
	}
	var in struct {
		Method string `json:"payment_method_id"`
	}
	if e = decode(r, &in); e != nil {
		return e
	}
	tag, e := a.DB.Exec(r.Context(), `UPDATE orders x SET payment_method_id=p.id,payment_snapshot=to_jsonb(p)-'enabled',updated_at=now() FROM payment_methods p,event_sessions s,events ev WHERE x.id=$1 AND x.account_id=$2 AND x.status='awaiting_payment' AND x.expires_at>now() AND p.id=$3 AND p.enabled AND s.id=x.session_id AND ev.id=s.event_id AND ev.organizer_id=p.organizer_id`, r.PathValue("id"), u["id"], in.Method)
	if e != nil {
		return e
	}
	if tag.RowsAffected() == 0 {
		return conflict("Metode pembayaran tidak tersedia atau pesanan tidak bisa diubah.")
	}
	return a.order(w, r)
}
func (a *App) submitProof(w http.ResponseWriter, r *http.Request) error {
	u, e := a.actor(r, "customer")
	if e != nil {
		return e
	}
	var in struct {
		Upload string `json:"upload_id"`
	}
	if e = decode(r, &in); e != nil {
		return e
	}
	ctx := r.Context()
	tx, e := a.DB.Begin(ctx)
	if e != nil {
		return e
	}
	defer tx.Rollback(ctx)
	var owner, status, orgOwner, existingProof string
	var method *string
	var expiry *time.Time
	e = tx.QueryRow(ctx, `SELECT x.account_id,x.status,o.account_id,x.payment_method_id,x.expires_at,COALESCE(x.proof_id,'') FROM orders x JOIN event_sessions s ON s.id=x.session_id JOIN events ev ON ev.id=s.event_id JOIN organizers o ON o.id=ev.organizer_id WHERE x.id=$1 FOR UPDATE OF x`, r.PathValue("id")).Scan(&owner, &status, &orgOwner, &method, &expiry, &existingProof)
	if e != nil {
		return e
	}
	if owner != u["id"] {
		return forbidden()
	}
	if status == "awaiting_review" && existingProof == in.Upload {
		if e = tx.Commit(ctx); e != nil {
			return e
		}
		return a.order(w, r)
	}
	if status != "awaiting_payment" && status != "correction_requested" {
		return conflict("Pesanan tidak menerima bukti pembayaran.")
	}
	if method == nil {
		return bad("Pilih rekening atau wallet terlebih dahulu.")
	}
	if status == "awaiting_payment" && (expiry == nil || !expiry.After(time.Now())) {
		return conflict("Waktu pembayaran berakhir. Jika sudah transfer, hubungi pengelola.")
	}
	var valid bool
	e = tx.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM uploads WHERE id=$1 AND account_id=$2 AND purpose='proof')`, in.Upload, u["id"]).Scan(&valid)
	if e != nil {
		return e
	}
	if !valid {
		return bad("Bukti transfer tidak valid.")
	}
	_, e = tx.Exec(ctx, `UPDATE orders SET status='awaiting_review',proof_id=$2,note='',expires_at=NULL,updated_at=now() WHERE id=$1`, r.PathValue("id"), in.Upload)
	if e != nil {
		return e
	}
	if _, e = tx.Exec(ctx, `INSERT INTO order_history(order_id,status,actor_id,note) VALUES($1,'awaiting_review',$2,'Bukti transfer diunggah.')`, r.PathValue("id"), u["id"]); e != nil {
		return e
	}
	if e = notify(ctx, tx, owner, "order", "Bukti transfer terkirim", "Menunggu pengelola memeriksa pembayaran.", "/transaksi/"+r.PathValue("id")); e != nil {
		return e
	}
	if e = notify(ctx, tx, orgOwner, "payment", "Bukti transfer perlu diperiksa", "Periksa mutasi rekening sebelum menyetujui.", "/kelola/pembayaran/"+r.PathValue("id")); e != nil {
		return e
	}
	if e = tx.Commit(ctx); e != nil {
		return e
	}
	return a.order(w, r)
}
func (a *App) cancelOrder(w http.ResponseWriter, r *http.Request) error {
	u, e := a.actor(r, "customer")
	if e != nil {
		return e
	}
	var in struct {
		Unpaid bool `json:"confirm_unpaid"`
	}
	if e = decode(r, &in); e != nil {
		return e
	}
	if !in.Unpaid {
		return bad("Konfirmasi bahwa kamu belum mentransfer pembayaran.")
	}
	ctx := r.Context()
	tx, e := a.DB.Begin(ctx)
	if e != nil {
		return e
	}
	defer tx.Rollback(ctx)
	tag, e := tx.Exec(ctx, `UPDATE orders SET status='cancelled',updated_at=now() WHERE id=$1 AND account_id=$2 AND status='awaiting_payment' AND proof_id IS NULL`, r.PathValue("id"), u["id"])
	if e != nil {
		return e
	}
	if tag.RowsAffected() == 0 {
		return conflict("Pesanan hanya bisa dibatalkan sebelum bukti transfer diunggah.")
	}
	if _, e = tx.Exec(ctx, `INSERT INTO order_history(order_id,status,actor_id,note) VALUES($1,'cancelled',$2,'Penonton mengonfirmasi belum transfer.')`, r.PathValue("id"), u["id"]); e != nil {
		return e
	}
	if e = notify(ctx, tx, u["id"].(string), "order", "Pesanan dibatalkan", "Kuota tiket dikembalikan.", "/transaksi/"+r.PathValue("id")); e != nil {
		return e
	}
	if e = tx.Commit(ctx); e != nil {
		return e
	}
	return a.order(w, r)
}
func (a *App) reviewOrder(w http.ResponseWriter, r *http.Request) error {
	u, e := a.actor(r, "organizer")
	if e != nil {
		return e
	}
	var in struct {
		Action   string `json:"action"`
		Note     string `json:"note"`
		Received bool   `json:"received_funds"`
	}
	if e = decode(r, &in); e != nil {
		return e
	}
	if in.Action != "approve" && in.Action != "correction" {
		return bad("Tindakan tidak valid.")
	}
	if in.Action == "approve" && !in.Received {
		return bad("Konfirmasi dana benar-benar sudah masuk ke rekening.")
	}
	if in.Action == "correction" && !textOK(in.Note, 10, 500) {
		return bad("Tuliskan alasan perbaikan bukti transfer.")
	}
	ctx := r.Context()
	tx, e := a.DB.Begin(ctx)
	if e != nil {
		return e
	}
	defer tx.Rollback(ctx)
	var owner, status, oid string
	var qty int
	e = tx.QueryRow(ctx, `SELECT x.account_id,x.status,ev.organizer_id,x.quantity FROM orders x JOIN event_sessions s ON s.id=x.session_id JOIN events ev ON ev.id=s.event_id WHERE x.id=$1 FOR UPDATE OF x`, r.PathValue("id")).Scan(&owner, &status, &oid, &qty)
	if e != nil {
		return e
	}
	if oid != u["organizer_id"] {
		return forbidden()
	}
	if status == "approved" && in.Action == "approve" {
		if e = tx.Commit(ctx); e != nil {
			return e
		}
		return a.order(w, r)
	}
	if status != "awaiting_review" {
		return conflict("Pesanan sudah diperiksa atau belum memiliki bukti transfer.")
	}
	status = "approved"
	title := "Pembayaran disetujui, tiket sudah terbit"
	if in.Action == "correction" {
		status = "correction_requested"
		title = "Bukti transfer perlu diperbaiki"
	} else {
		if e = issueTickets(ctx, tx, r.PathValue("id"), qty); e != nil {
			return e
		}
	}
	_, e = tx.Exec(ctx, `UPDATE orders SET status=$2,note=$3,updated_at=now() WHERE id=$1`, r.PathValue("id"), status, in.Note)
	if e != nil {
		return e
	}
	if _, e = tx.Exec(ctx, `INSERT INTO order_history(order_id,status,actor_id,note) VALUES($1,$2,$3,$4)`, r.PathValue("id"), status, u["id"], in.Note); e != nil {
		return e
	}
	if e = notify(ctx, tx, owner, "order", title, in.Note, "/transaksi/"+r.PathValue("id")); e != nil {
		return e
	}
	if e = tx.Commit(ctx); e != nil {
		return e
	}
	return a.order(w, r)
}
func (a *App) tickets(w http.ResponseWriter, r *http.Request) error {
	u, e := a.actor(r, "customer")
	if e != nil {
		return e
	}
	items, e := a.list(r.Context(), `SELECT to_jsonb(t)-'token'||jsonb_build_object('qr','ruang:ticket:'||t.token,'event',ev.doc,'session',to_jsonb(s)) FROM tickets t JOIN orders x ON x.id=t.order_id JOIN event_sessions s ON s.id=x.session_id JOIN event_view ev ON ev.id=s.event_id WHERE x.account_id=$1 AND x.status='approved' ORDER BY s.starts_at DESC,t.ordinal LIMIT 200`, u["id"])
	if e != nil {
		return e
	}
	send(w, 200, items)
	return nil
}
func (a *App) checkIn(w http.ResponseWriter, r *http.Request) error {
	u, e := a.actor(r, "organizer")
	if e != nil {
		return e
	}
	var in struct {
		Code  string `json:"code"`
		Event string `json:"event_id"`
	}
	if e = decode(r, &in); e != nil {
		return e
	}
	code := strings.TrimSpace(strings.TrimPrefix(in.Code, "ruang:ticket:"))
	if len(code) != 64 {
		return bad("Kode tiket tidak valid.")
	}
	ctx := r.Context()
	tx, e := a.DB.Begin(ctx)
	if e != nil {
		return e
	}
	defer tx.Rollback(ctx)
	var id, oid, eventID, name, status string
	var checked *time.Time
	var starts, ends time.Time
	e = tx.QueryRow(ctx, `SELECT t.id,ev.organizer_id,ev.id,a.name,t.checked_at,s.starts_at,s.ends_at,x.status FROM tickets t JOIN orders x ON x.id=t.order_id JOIN accounts a ON a.id=x.account_id JOIN event_sessions s ON s.id=x.session_id JOIN events ev ON ev.id=s.event_id WHERE t.token=$1 FOR UPDATE OF t`, code).Scan(&id, &oid, &eventID, &name, &checked, &starts, &ends, &status)
	if e != nil {
		return missing()
	}
	if oid != u["organizer_id"] {
		return forbidden()
	}
	if eventID != in.Event {
		return bad("Tiket berasal dari event yang berbeda.")
	}
	if checked != nil {
		return conflict("Tiket sudah digunakan pada " + checked.In(time.FixedZone("WIB", 7*3600)).Format("02 Jan 2006 15:04") + " WIB.")
	}
	if status != "approved" {
		return conflict("Pembayaran tiket belum disetujui.")
	}
	now := time.Now()
	if now.Before(starts.Add(-2*time.Hour)) || now.After(ends.Add(2*time.Hour)) {
		return conflict("Check-in dibuka dua jam sebelum sesi sampai dua jam setelah selesai.")
	}
	_, e = tx.Exec(ctx, `UPDATE tickets SET checked_at=now(),checked_by=$2 WHERE id=$1`, id, u["id"])
	if e != nil {
		return e
	}
	if e = tx.Commit(ctx); e != nil {
		return e
	}
	send(w, 200, map[string]any{"valid": true, "ticket_id": id, "name": name, "message": "Tiket valid. Selamat menikmati pertunjukan."})
	return nil
}
func (a *App) paymentMethods(w http.ResponseWriter, r *http.Request) error {
	oid := r.URL.Query().Get("organizer")
	u, e := a.me(r)
	own := ""
	if e == nil && u["role"] == "organizer" {
		own = u["organizer_id"].(string)
	}
	if oid == "" {
		oid = own
	}
	if oid == "" {
		return bad("Pilih pengelola.")
	}
	items, e := a.list(r.Context(), `SELECT to_jsonb(p) FROM payment_methods p WHERE organizer_id=$1 AND (enabled OR organizer_id=$2) ORDER BY provider,id`, oid, own)
	if e != nil {
		return e
	}
	send(w, 200, items)
	return nil
}
func (a *App) savePaymentMethod(w http.ResponseWriter, r *http.Request) error {
	u, e := a.actor(r, "organizer")
	if e != nil {
		return e
	}
	var in struct {
		Kind     string `json:"kind"`
		Provider string `json:"provider"`
		Number   string `json:"number"`
		Holder   string `json:"holder"`
		Note     string `json:"note"`
		Enabled  bool   `json:"enabled"`
	}
	if e = decode(r, &in); e != nil {
		return e
	}
	if (in.Kind != "bank" && in.Kind != "wallet") || !textOK(in.Provider, 2, 60) || !textOK(in.Number, 3, 80) || !textOK(in.Holder, 2, 100) || !textOK(in.Note, 0, 500) {
		return bad("Lengkapi bank/wallet, nomor, dan nama pemilik.")
	}
	id := r.PathValue("id")
	if id == "" {
		id = ID()
		_, e = a.DB.Exec(r.Context(), `INSERT INTO payment_methods(id,organizer_id,kind,provider,number,holder,note,enabled) VALUES($1,$2,$3,$4,$5,$6,$7,$8)`, id, u["organizer_id"], in.Kind, in.Provider, in.Number, in.Holder, in.Note, in.Enabled)
	} else {
		var tag pgconn.CommandTag
		tag, e = a.DB.Exec(r.Context(), `UPDATE payment_methods SET kind=$3,provider=$4,number=$5,holder=$6,note=$7,enabled=$8 WHERE id=$1 AND organizer_id=$2`, id, u["organizer_id"], in.Kind, in.Provider, in.Number, in.Holder, in.Note, in.Enabled)
		if e == nil && tag.RowsAffected() == 0 {
			return forbidden()
		}
	}
	if e != nil {
		return e
	}
	doc, e := a.one(r.Context(), `SELECT to_jsonb(p) FROM payment_methods p WHERE id=$1`, id)
	if e != nil {
		return e
	}
	send(w, 200, doc)
	return nil
}
