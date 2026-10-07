package app

import (
	"context"
	"encoding/json"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func (a *App) mediaOK(ctx context.Context, uid, value string) bool {
	if value == "" {
		return true
	}
	if strings.HasPrefix(value, "/assets/") {
		return !strings.Contains(value, "..") && !strings.ContainsAny(value, "?#\\")
	}
	if strings.HasPrefix(value, "/api/uploads/") {
		id := strings.TrimPrefix(value, "/api/uploads/")
		var ok bool
		e := a.DB.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM uploads WHERE id=$1 AND account_id=$2 AND purpose='media')`, id, uid).Scan(&ok)
		return e == nil && ok
	}
	return strings.HasPrefix(value, "https://") && safeURL(value)
}
func (a *App) upload(w http.ResponseWriter, r *http.Request) error {
	u, e := a.me(r)
	if e != nil {
		return e
	}
	if e = r.ParseMultipartForm(2 << 20); e != nil {
		return bad("Gambar maksimal 8 MB.")
	}
	if r.MultipartForm != nil {
		defer r.MultipartForm.RemoveAll()
	}
	purpose := r.FormValue("purpose")
	if purpose != "media" && purpose != "proof" {
		return bad("Tujuan unggahan tidak valid.")
	}
	file, header, e := r.FormFile("file")
	if e != nil {
		return bad("Pilih gambar untuk diunggah.")
	}
	defer file.Close()
	if header.Size <= 0 || header.Size > 8<<20 {
		return bad("Gambar maksimal 8 MB.")
	}
	prefix := make([]byte, 512)
	n, e := file.Read(prefix)
	if e != nil && e != io.EOF {
		return e
	}
	mime := http.DetectContentType(prefix[:n])
	ext := map[string]string{"image/jpeg": ".jpg", "image/png": ".png", "image/webp": ".webp"}[mime]
	if ext == "" {
		return bad("Gunakan gambar JPG, PNG, atau WebP. SVG tidak diterima.")
	}
	if _, e = file.Seek(0, io.SeekStart); e != nil {
		return e
	}
	if mime != "image/webp" {
		cfg, _, err := image.DecodeConfig(file)
		if err != nil || cfg.Width < 1 || cfg.Height < 1 || int64(cfg.Width)*int64(cfg.Height) > 40000000 {
			return bad("Gambar tidak valid atau resolusinya terlalu besar.")
		}
		if _, e = file.Seek(0, io.SeekStart); e != nil {
			return e
		}
	}
	id := ID()
	filename := id + ext
	f, e := os.OpenFile(filepath.Join(a.Config.UploadDir, filename), os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if e != nil {
		return e
	}
	size, e := io.Copy(f, io.LimitReader(file, (8<<20)+1))
	closeErr := f.Close()
	if e != nil || closeErr != nil || size > 8<<20 {
		_ = os.Remove(filepath.Join(a.Config.UploadDir, filename))
		return bad("Gagal menyimpan gambar.")
	}
	_, e = a.DB.Exec(r.Context(), `INSERT INTO uploads(id,account_id,purpose,filename,content_type,size) VALUES($1,$2,$3,$4,$5,$6)`, id, u["id"], purpose, filename, mime, size)
	if e != nil {
		_ = os.Remove(filepath.Join(a.Config.UploadDir, filename))
		return e
	}
	send(w, 201, map[string]any{"id": id, "url": "/api/uploads/" + id, "size": size})
	return nil
}
func (a *App) download(w http.ResponseWriter, r *http.Request) error {
	var owner, purpose, filename, mime string
	e := a.DB.QueryRow(r.Context(), `SELECT account_id,purpose,filename,content_type FROM uploads WHERE id=$1`, r.PathValue("id")).Scan(&owner, &purpose, &filename, &mime)
	if e != nil {
		return e
	}
	if purpose == "proof" {
		u, e := a.me(r)
		if e != nil {
			return e
		}
		if u["id"] != owner {
			var allowed bool
			e = a.DB.QueryRow(r.Context(), `SELECT EXISTS(SELECT 1 FROM orders x JOIN event_sessions s ON s.id=x.session_id JOIN events ev ON ev.id=s.event_id JOIN organizers o ON o.id=ev.organizer_id WHERE x.proof_id=$1 AND o.account_id=$2)`, r.PathValue("id"), u["id"]).Scan(&allowed)
			if e != nil {
				return e
			}
			if !allowed {
				return forbidden()
			}
		}
		w.Header().Set("Cache-Control", "private, no-store")
	} else {
		w.Header().Set("Cache-Control", "public, max-age=86400, immutable")
	}
	w.Header().Set("Content-Type", mime)
	http.ServeFile(w, r, filepath.Join(a.Config.UploadDir, filepath.Base(filename)))
	return nil
}
func (a *App) products(w http.ResponseWriter, r *http.Request) error {
	q := r.URL.Query()
	items, e := a.list(r.Context(), `SELECT doc FROM product_view WHERE ($1='' OR organizer_id=$1) AND ($2='' OR doc->>'name' ILIKE $3) ORDER BY doc->>'created_at' DESC LIMIT 100`, q.Get("organizer"), q.Get("q"), "%"+q.Get("q")+"%")
	if e != nil {
		return e
	}
	send(w, 200, items)
	return nil
}
func (a *App) product(w http.ResponseWriter, r *http.Request) error {
	doc, e := a.one(r.Context(), `SELECT doc FROM product_view WHERE id=$1`, r.PathValue("id"))
	if e != nil {
		return e
	}
	send(w, 200, doc)
	return nil
}
func (a *App) saveProduct(w http.ResponseWriter, r *http.Request) error {
	u, e := a.actor(r, "organizer")
	if e != nil {
		return e
	}
	var in struct {
		Name         string   `json:"name"`
		Description  string   `json:"description"`
		Price        int64    `json:"price"`
		Image        string   `json:"image_url"`
		Variants     []string `json:"variants"`
		Availability string   `json:"availability"`
		URL          string   `json:"purchase_url"`
		Event        string   `json:"event_id"`
	}
	if e = decode(r, &in); e != nil {
		return e
	}
	if !textOK(in.Name, 2, 100) || !textOK(in.Description, 10, 5000) || in.Price < 0 || in.Price > 100000000 || len(in.Variants) < 1 || len(in.Variants) > 20 || !safeURL(in.URL) || !textOK(in.Availability, 2, 80) || !a.mediaOK(r.Context(), u["id"].(string), in.Image) {
		return bad("Informasi merchandise tidak valid.")
	}
	for _, v := range in.Variants {
		if !textOK(v, 1, 80) {
			return bad("Varian tidak valid.")
		}
	}
	var event any
	if in.Event != "" {
		var owned bool
		e = a.DB.QueryRow(r.Context(), `SELECT EXISTS(SELECT 1 FROM events WHERE id=$1 AND organizer_id=$2)`, in.Event, u["organizer_id"]).Scan(&owned)
		if e != nil {
			return e
		}
		if !owned {
			return forbidden()
		}
		event = in.Event
	}
	variants, _ := json.Marshal(in.Variants)
	id := r.PathValue("id")
	if id == "" {
		id = ID()
		_, e = a.DB.Exec(r.Context(), `INSERT INTO products(id,organizer_id,event_id,name,description,price,image_url,variants,availability,purchase_url) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`, id, u["organizer_id"], event, in.Name, in.Description, in.Price, in.Image, string(variants), in.Availability, in.URL)
	} else {
		tag, err := a.DB.Exec(r.Context(), `UPDATE products SET event_id=$3,name=$4,description=$5,price=$6,image_url=$7,variants=$8,availability=$9,purchase_url=$10 WHERE id=$1 AND organizer_id=$2`, id, u["organizer_id"], event, in.Name, in.Description, in.Price, in.Image, string(variants), in.Availability, in.URL)
		e = err
		if e == nil && tag.RowsAffected() == 0 {
			return forbidden()
		}
	}
	if e != nil {
		return e
	}
	doc, e := a.one(r.Context(), `SELECT doc FROM product_view WHERE id=$1`, id)
	if e != nil {
		return e
	}
	send(w, 200, doc)
	return nil
}
func (a *App) notifications(w http.ResponseWriter, r *http.Request) error {
	u, e := a.me(r)
	if e != nil {
		return e
	}
	items, e := a.list(r.Context(), `SELECT to_jsonb(n) FROM notifications n WHERE account_id=$1 ORDER BY created_at DESC LIMIT 100`, u["id"])
	if e != nil {
		return e
	}
	send(w, 200, items)
	return nil
}
func (a *App) readNotifications(w http.ResponseWriter, r *http.Request) error {
	u, e := a.me(r)
	if e != nil {
		return e
	}
	var in struct {
		ID string `json:"id"`
	}
	if e = decode(r, &in); e != nil {
		return e
	}
	_, e = a.DB.Exec(r.Context(), `UPDATE notifications SET read_at=now() WHERE account_id=$1 AND read_at IS NULL AND ($2='' OR id=$2)`, u["id"], in.ID)
	if e != nil {
		return e
	}
	send(w, 200, map[string]bool{"ok": true})
	return nil
}
func (a *App) dashboard(w http.ResponseWriter, r *http.Request) error {
	u, e := a.actor(r, "organizer")
	if e != nil {
		return e
	}
	doc, e := a.one(r.Context(), `SELECT jsonb_build_object(
 'followers',(SELECT count(*) FROM follows WHERE organizer_id=$1),
 'posts',(SELECT count(*) FROM posts WHERE organizer_id=$1 AND NOT hidden),
 'events',(SELECT count(*) FROM events WHERE organizer_id=$1),
 'interactions',(SELECT count(*) FROM votes v JOIN posts p ON p.id=v.post_id WHERE p.organizer_id=$1 AND NOT p.hidden)+(SELECT count(*) FROM comments c JOIN posts p ON p.id=c.post_id WHERE p.organizer_id=$1 AND NOT p.hidden),
 'ticket_clicks',(SELECT count(*) FROM activity_clicks WHERE organizer_id=$1 AND kind='ticket'),
 'merch_clicks',(SELECT count(*) FROM activity_clicks WHERE organizer_id=$1 AND kind='merch'),
 'activity',(SELECT jsonb_agg(jsonb_build_object('date',d.day::date,'count',
   (SELECT count(*) FROM posts p WHERE p.organizer_id=$1 AND NOT hidden AND (p.created_at AT TIME ZONE 'Asia/Jakarta')::date=d.day::date)+
   (SELECT count(*) FROM votes v JOIN posts p ON p.id=v.post_id WHERE p.organizer_id=$1 AND NOT p.hidden AND (v.created_at AT TIME ZONE 'Asia/Jakarta')::date=d.day::date)+
   (SELECT count(*) FROM comments c JOIN posts p ON p.id=c.post_id WHERE p.organizer_id=$1 AND NOT p.hidden AND (c.created_at AT TIME ZONE 'Asia/Jakarta')::date=d.day::date)
 ) ORDER BY d.day) FROM generate_series((now() AT TIME ZONE 'Asia/Jakarta')::date-6,(now() AT TIME ZONE 'Asia/Jakarta')::date,interval '1 day') AS d(day)),
 'pending_payments',(SELECT count(*) FROM orders x JOIN event_sessions s ON s.id=x.session_id JOIN events ev ON ev.id=s.event_id WHERE ev.organizer_id=$1 AND x.status='awaiting_review'),
 'tickets',(SELECT count(*) FROM tickets t JOIN orders x ON x.id=t.order_id JOIN event_sessions s ON s.id=x.session_id JOIN events ev ON ev.id=s.event_id WHERE ev.organizer_id=$1),
 'checked_in',(SELECT count(*) FROM tickets t JOIN orders x ON x.id=t.order_id JOIN event_sessions s ON s.id=x.session_id JOIN events ev ON ev.id=s.event_id WHERE ev.organizer_id=$1 AND t.checked_at IS NOT NULL),
 'revenue',(SELECT COALESCE(sum(x.total),0) FROM orders x JOIN event_sessions s ON s.id=x.session_id JOIN events ev ON ev.id=s.event_id WHERE ev.organizer_id=$1 AND x.status='approved'))`, u["organizer_id"])
	if e != nil {
		return e
	}
	send(w, 200, doc)
	return nil
}

// Counts CTA clicks, not sales. No viewer identity or IP is persisted.
func (a *App) activityClick(w http.ResponseWriter, r *http.Request) error {
	var in struct {
		Kind   string `json:"kind"`
		Target string `json:"target_id"`
	}
	if e := decode(r, &in); e != nil {
		return e
	}
	var result int64
	if in.Kind == "ticket" {
		tag, e := a.DB.Exec(r.Context(), `INSERT INTO activity_clicks(organizer_id,kind,event_id) SELECT organizer_id,'ticket',id FROM events WHERE (id=$1 OR slug=$1) AND published`, in.Target)
		if e != nil {
			return e
		}
		result = tag.RowsAffected()
	} else if in.Kind == "merch" {
		tag, e := a.DB.Exec(r.Context(), `INSERT INTO activity_clicks(organizer_id,kind,product_id) SELECT organizer_id,'merch',id FROM products WHERE id=$1 AND availability<>'Habis'`, in.Target)
		if e != nil {
			return e
		}
		result = tag.RowsAffected()
	} else {
		return bad("Jenis aktivitas tidak valid.")
	}
	if result == 0 {
		return missing()
	}
	send(w, 201, map[string]bool{"ok": true})
	return nil
}
func (a *App) report(w http.ResponseWriter, r *http.Request) error {
	u, e := a.me(r)
	if e != nil {
		return e
	}
	var in struct {
		Kind   string `json:"kind"`
		Target string `json:"target_id"`
		Reason string `json:"reason"`
	}
	if e = decode(r, &in); e != nil {
		return e
	}
	if !textOK(in.Reason, 10, 1000) || (in.Kind != "post" && in.Kind != "review") {
		return bad("Tuliskan alasan laporan, minimal 10 karakter.")
	}
	var exists bool
	table := "posts"
	if in.Kind == "review" {
		table = "reviews"
	}
	e = a.DB.QueryRow(r.Context(), fmt.Sprintf(`SELECT EXISTS(SELECT 1 FROM %s WHERE id=$1 AND NOT hidden)`, table), in.Target).Scan(&exists)
	if e != nil {
		return e
	}
	if !exists {
		return missing()
	}
	_, e = a.DB.Exec(r.Context(), `INSERT INTO reports(id,account_id,kind,target_id,reason) VALUES($1,$2,$3,$4,$5)`, ID(), u["id"], in.Kind, in.Target, in.Reason)
	if e != nil {
		return e
	}
	send(w, 201, map[string]bool{"ok": true})
	return nil
}
func (a *App) reports(w http.ResponseWriter, r *http.Request) error {
	if _, e := a.actor(r, "admin"); e != nil {
		return e
	}
	items, e := a.list(r.Context(), `SELECT to_jsonb(r)||jsonb_build_object('content',CASE WHEN r.kind='post' THEN(SELECT doc FROM post_view WHERE id=r.target_id) ELSE(SELECT doc FROM review_view WHERE id=r.target_id) END) FROM reports r ORDER BY r.created_at DESC LIMIT 200`)
	if e != nil {
		return e
	}
	send(w, 200, items)
	return nil
}
func (a *App) resolveReport(w http.ResponseWriter, r *http.Request) error {
	if _, e := a.actor(r, "admin"); e != nil {
		return e
	}
	var in struct {
		Status string `json:"status"`
	}
	if e := decode(r, &in); e != nil {
		return e
	}
	if in.Status != "reviewed" && in.Status != "hidden" {
		return bad("Status laporan tidak valid.")
	}
	tx, e := a.DB.Begin(r.Context())
	if e != nil {
		return e
	}
	defer tx.Rollback(r.Context())
	var kind, id string
	e = tx.QueryRow(r.Context(), `UPDATE reports SET status=$2 WHERE id=$1 RETURNING kind,target_id`, r.PathValue("id"), in.Status).Scan(&kind, &id)
	if e != nil {
		return e
	}
	if in.Status == "hidden" {
		table := "posts"
		if kind == "review" {
			table = "reviews"
		}
		if _, e = tx.Exec(r.Context(), fmt.Sprintf(`UPDATE %s SET hidden=true WHERE id=$1`, table), id); e != nil {
			return e
		}
	}
	if e = tx.Commit(r.Context()); e != nil {
		return e
	}
	send(w, 200, map[string]bool{"ok": true})
	return nil
}
