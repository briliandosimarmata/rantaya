package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func page(r *http.Request) (int, int) {
	n, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if n < 1 {
		n = 1
	}
	if n > 1000 {
		n = 1000
	}
	return 20, (n - 1) * 20
}
func (a *App) organizers(w http.ResponseWriter, r *http.Request) error {
	q := r.URL.Query()
	l, o := page(r)
	items, e := a.list(r.Context(), `SELECT doc FROM organizer_view WHERE ($1='' OR doc->>'city'=$1) AND ($2='' OR doc->>'name' ILIKE $3 OR doc->>'description' ILIKE $3) ORDER BY (doc->>'followers')::int DESC,id LIMIT $4 OFFSET $5`, q.Get("city"), q.Get("q"), "%"+q.Get("q")+"%", l, o)
	if e != nil {
		return e
	}
	send(w, 200, items)
	return nil
}
func (a *App) organizer(w http.ResponseWriter, r *http.Request) error {
	doc, e := a.one(r.Context(), `SELECT doc FROM organizer_view WHERE id=$1 OR slug=$1`, r.PathValue("id"))
	if e != nil {
		return e
	}
	send(w, 200, doc)
	return nil
}
func (a *App) editOrganizer(w http.ResponseWriter, r *http.Request) error {
	u, e := a.actor(r, "organizer")
	if e != nil {
		return e
	}
	var in struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		About       string `json:"about"`
		City        string `json:"city"`
		Category    string `json:"category"`
		Avatar      string `json:"avatar_url"`
		Cover       string `json:"cover_url"`
	}
	if e = decode(r, &in); e != nil {
		return e
	}
	if !textOK(in.Name, 2, 80) || !textOK(in.Description, 10, 500) || !textOK(in.About, 0, 5000) || !textOK(in.City, 2, 60) || !validCategory(in.Category) {
		return bad("Lengkapi nama, deskripsi, kota, dan kategori ruang.")
	}
	if !a.mediaOK(r.Context(), u["id"].(string), in.Avatar) || !a.mediaOK(r.Context(), u["id"].(string), in.Cover) {
		return bad("Gambar ruang tidak valid.")
	}
	_, e = a.DB.Exec(r.Context(), `UPDATE organizers SET name=$2,description=$3,about=$4,city=$5,category=$6,avatar_url=$7,cover_url=$8 WHERE account_id=$1`, u["id"], in.Name, in.Description, in.About, in.City, in.Category, in.Avatar, in.Cover)
	if e != nil {
		return e
	}
	doc, e := a.one(r.Context(), `SELECT doc FROM organizer_view WHERE account_id=$1`, u["id"])
	if e != nil {
		return e
	}
	send(w, 200, doc)
	return nil
}
func validCategory(s string) bool {
	return s == "Teater" || s == "Musik" || s == "Komedi" || s == "Lainnya"
}
func (a *App) follow(w http.ResponseWriter, r *http.Request) error {
	u, e := a.me(r)
	if e != nil {
		return e
	}
	var in struct {
		Active bool `json:"active"`
	}
	if e = decode(r, &in); e != nil {
		return e
	}
	id := r.PathValue("id")
	var exists bool
	e = a.DB.QueryRow(r.Context(), `SELECT EXISTS(SELECT 1 FROM organizers WHERE id=$1)`, id).Scan(&exists)
	if e != nil {
		return e
	}
	if !exists {
		return missing()
	}
	if in.Active {
		_, e = a.DB.Exec(r.Context(), `INSERT INTO follows(account_id,organizer_id) VALUES($1,$2) ON CONFLICT DO NOTHING`, u["id"], id)
	} else {
		_, e = a.DB.Exec(r.Context(), `DELETE FROM follows WHERE account_id=$1 AND organizer_id=$2`, u["id"], id)
	}
	if e != nil {
		return e
	}
	send(w, 200, map[string]bool{"active": in.Active})
	return nil
}
func (a *App) events(w http.ResponseWriter, r *http.Request) error {
	q := r.URL.Query()
	l, o := page(r)
	from, to := q.Get("from"), q.Get("to")
	for _, d := range []string{from, to} {
		if d != "" {
			if _, e := time.Parse("2006-01-02", d); e != nil {
				return bad("Tanggal tidak valid.")
			}
		}
	}
	items, e := a.list(r.Context(), `SELECT v.doc FROM event_view v WHERE (v.published OR v.organizer_id=COALESCE((SELECT id FROM organizers WHERE account_id=$1),'')) AND ($2='' OR v.city=$2) AND ($3='' OR v.category=$3) AND ($4='' OR v.organizer_id=$4) AND ($5='' OR v.doc->>'title' ILIKE $6 OR v.doc->>'description' ILIKE $6) AND ($7='all' OR ($7='past' AND (v.doc->>'ends_at')::timestamptz<now()) OR ($7 IN ('','upcoming') AND (v.doc->>'ends_at')::timestamptz>=now())) AND EXISTS(SELECT 1 FROM event_sessions s WHERE s.event_id=v.id AND ($8='' OR (s.starts_at AT TIME ZONE 'Asia/Jakarta')::date>=NULLIF($8,'')::date) AND ($9='' OR (s.starts_at AT TIME ZONE 'Asia/Jakarta')::date<=NULLIF($9,'')::date)) ORDER BY (v.doc->>'starts_at')::timestamptz,v.id LIMIT $10 OFFSET $11`, a.uid(r), q.Get("city"), q.Get("category"), q.Get("organizer"), q.Get("q"), "%"+q.Get("q")+"%", q.Get("period"), from, to, l, o)
	if e != nil {
		return e
	}
	send(w, 200, items)
	return nil
}
func (a *App) event(w http.ResponseWriter, r *http.Request) error {
	doc, e := a.one(r.Context(), `SELECT doc FROM event_view WHERE (id=$1 OR slug=$1) AND (published OR organizer_id=(SELECT id FROM organizers WHERE account_id=$2))`, r.PathValue("id"), a.uid(r))
	if e != nil {
		return e
	}
	send(w, 200, doc)
	return nil
}

type Performer struct {
	Name  string `json:"name"`
	Role  string `json:"role"`
	Photo string `json:"photo"`
}
type SessionInput struct {
	ID       string    `json:"id"`
	Label    string    `json:"label"`
	Starts   time.Time `json:"starts_at"`
	Ends     time.Time `json:"ends_at"`
	Price    int64     `json:"price"`
	Capacity int       `json:"capacity"`
}
type EventInput struct {
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Category    string         `json:"category"`
	City        string         `json:"city"`
	Venue       string         `json:"venue"`
	Address     string         `json:"address"`
	Maps        string         `json:"maps_url"`
	Duration    string         `json:"duration"`
	Language    string         `json:"language"`
	Age         string         `json:"age"`
	Flyer       string         `json:"flyer_url"`
	Trailer     string         `json:"trailer_url"`
	Layout      string         `json:"layout_url"`
	Lineup      []Performer    `json:"lineup"`
	Published   bool           `json:"published"`
	Sessions    []SessionInput `json:"sessions"`
}

func (a *App) saveEvent(w http.ResponseWriter, r *http.Request) error {
	u, e := a.actor(r, "organizer")
	if e != nil {
		return e
	}
	uid := u["id"].(string)
	oid, _ := u["organizer_id"].(string)
	var in EventInput
	if e = decode(r, &in); e != nil {
		return e
	}
	if !textOK(in.Title, 3, 120) || !textOK(in.Description, 10, 10000) || !validCategory(in.Category) || !textOK(in.City, 2, 60) || !textOK(in.Venue, 2, 150) || len(in.Sessions) < 1 || len(in.Sessions) > 30 || len(in.Lineup) > 40 || !textOK(in.Address, 0, 400) || !textOK(in.Duration, 0, 80) || !textOK(in.Language, 0, 80) || !textOK(in.Age, 0, 60) {
		return bad("Lengkapi informasi event dan minimal satu sesi.")
	}
	if !mapsOK(in.Maps) || !safeURL(in.Trailer) {
		return bad("Tautan Maps atau trailer tidak valid.")
	}
	for _, v := range []string{in.Flyer, in.Layout} {
		if !a.mediaOK(r.Context(), uid, v) {
			return bad("Media event tidak valid.")
		}
	}
	for _, p := range in.Lineup {
		if !textOK(p.Name, 1, 80) || !textOK(p.Role, 1, 80) || !a.mediaOK(r.Context(), uid, p.Photo) {
			return bad("Nama, peran, atau foto pengisi tidak valid.")
		}
	}
	paid := false
	for _, s := range in.Sessions {
		if !textOK(s.Label, 2, 100) || s.Starts.IsZero() || !s.Ends.After(s.Starts) || s.Price < 0 || s.Price > 100000000 || s.Capacity < 1 || s.Capacity > 100000 {
			return bad("Jadwal, harga, atau kapasitas sesi tidak valid.")
		}
		paid = paid || s.Price > 0
	}
	tx, e := a.DB.Begin(r.Context())
	if e != nil {
		return e
	}
	defer tx.Rollback(r.Context())
	if in.Published {
		var complete, methods bool
		e = tx.QueryRow(r.Context(), `SELECT length(description)>=10,EXISTS(SELECT 1 FROM payment_methods WHERE organizer_id=$1 AND enabled) FROM organizers WHERE id=$1`, oid).Scan(&complete, &methods)
		if e != nil {
			return e
		}
		if !complete || paid && !methods {
			return bad("Lengkapi profil ruang dan rekening pembayaran sebelum menerbitkan event berbayar.")
		}
	}
	id := r.PathValue("id")
	oldPublished := false
	if id != "" {
		var owner string
		e = tx.QueryRow(r.Context(), `SELECT organizer_id,published FROM events WHERE id=$1 FOR UPDATE`, id).Scan(&owner, &oldPublished)
		if e != nil {
			return e
		}
		if owner != oid {
			return forbidden()
		}
	} else {
		id = ID()
	}
	lineup, _ := json.Marshal(in.Lineup)
	if in.Lineup == nil {
		lineup = []byte("[]")
	}
	_, e = tx.Exec(r.Context(), `INSERT INTO events(id,slug,organizer_id,title,description,category,city,venue,address,maps_url,duration,language,age,flyer_url,trailer_url,layout_url,lineup,published) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18) ON CONFLICT(id) DO UPDATE SET title=excluded.title,description=excluded.description,category=excluded.category,city=excluded.city,venue=excluded.venue,address=excluded.address,maps_url=excluded.maps_url,duration=excluded.duration,language=excluded.language,age=excluded.age,flyer_url=excluded.flyer_url,trailer_url=excluded.trailer_url,layout_url=excluded.layout_url,lineup=excluded.lineup,published=excluded.published,updated_at=now()`, id, slug(in.Title)+"-"+id[:8], oid, in.Title, in.Description, in.Category, in.City, in.Venue, in.Address, in.Maps, in.Duration, in.Language, in.Age, in.Flyer, in.Trailer, in.Layout, string(lineup), in.Published)
	if e != nil {
		return e
	}
	keep := []string{}
	seen := map[string]bool{}
	for _, s := range in.Sessions {
		if s.ID == "" {
			s.ID = ID()
		} else {
			if seen[s.ID] {
				return bad("Sesi terduplikasi.")
			}
			var parent string
			var start, end time.Time
			var price, reserved int64
			var booked bool
			e = tx.QueryRow(r.Context(), `SELECT event_id,starts_at,ends_at,price,EXISTS(SELECT 1 FROM orders WHERE session_id=$1),COALESCE((SELECT sum(quantity) FROM orders WHERE session_id=$1 AND (status IN ('approved','awaiting_review','correction_requested') OR status='awaiting_payment' AND expires_at>now())),0) FROM event_sessions WHERE id=$1 FOR UPDATE`, s.ID).Scan(&parent, &start, &end, &price, &booked, &reserved)
			if e != nil {
				return e
			}
			if parent != id {
				return forbidden()
			}
			if booked && (!s.Starts.Equal(start) || !s.Ends.Equal(end) || s.Price != price) {
				return conflict("Jadwal dan harga sesi yang sudah memiliki pesanan tidak dapat diubah.")
			}
			if int64(s.Capacity) < reserved {
				return conflict("Kapasitas tidak boleh lebih kecil dari tiket yang dipesan.")
			}
		}
		seen[s.ID] = true
		keep = append(keep, s.ID)
		_, e = tx.Exec(r.Context(), `INSERT INTO event_sessions(id,event_id,label,starts_at,ends_at,price,capacity) VALUES($1,$2,$3,$4,$5,$6,$7) ON CONFLICT(id) DO UPDATE SET label=excluded.label,starts_at=excluded.starts_at,ends_at=excluded.ends_at,price=excluded.price,capacity=excluded.capacity`, s.ID, id, s.Label, s.Starts, s.Ends, s.Price, s.Capacity)
		if e != nil {
			return e
		}
	}
	var removals int
	e = tx.QueryRow(r.Context(), `SELECT count(*) FROM event_sessions s WHERE event_id=$1 AND NOT(id=ANY($2::text[])) AND EXISTS(SELECT 1 FROM orders WHERE session_id=s.id)`, id, keep).Scan(&removals)
	if e != nil {
		return e
	}
	if removals > 0 {
		return conflict("Sesi yang memiliki pesanan tidak dapat dihapus.")
	}
	if _, e = tx.Exec(r.Context(), `DELETE FROM event_sessions WHERE event_id=$1 AND NOT(id=ANY($2::text[]))`, id, keep); e != nil {
		return e
	}
	if in.Published && !oldPublished {
		var savedSlug string
		if e = tx.QueryRow(r.Context(), `SELECT slug FROM events WHERE id=$1`, id).Scan(&savedSlug); e != nil {
			return e
		}
		_, e = tx.Exec(r.Context(), `INSERT INTO notifications(id,account_id,kind,title,body,url) SELECT $2||f.account_id,f.account_id,'event','Event baru dari ruang yang kamu ikuti',$3,'/event/'||$4 FROM follows f JOIN accounts a ON a.id=f.account_id WHERE f.organizer_id=$1 AND COALESCE((a.preferences->>'events')::boolean,true)`, oid, ID(), in.Title, savedSlug)
		if e != nil {
			return e
		}
	}
	if e = tx.Commit(r.Context()); e != nil {
		return e
	}
	doc, e := a.one(r.Context(), `SELECT doc FROM event_view WHERE id=$1`, id)
	if e != nil {
		return e
	}
	send(w, 200, doc)
	return nil
}

func (a *App) posts(w http.ResponseWriter, r *http.Request) error {
	q := r.URL.Query()
	l, o := page(r)
	order := "p.pinned DESC,p.created_at DESC,p.id DESC"
	if q.Get("sort") == "popular" {
		order = "(p.doc->>'votes')::int DESC,p.created_at DESC,p.id DESC"
	}
	items, e := a.list(r.Context(), `SELECT p.doc FROM post_view p WHERE NOT hidden AND ($1='' OR p.city=$1) AND ($2='' OR p.doc->>'body' ILIKE $3 OR p.doc->>'title' ILIKE $3) AND ($4='' OR ($5='mentions' AND p.organizer_id IS NULL AND EXISTS(SELECT 1 FROM post_mentions m WHERE m.post_id=p.id AND (m.kind='organizer' AND m.target_id=$4 OR m.kind='event' AND m.target_id IN(SELECT id FROM events WHERE organizer_id=$4)))) OR ($5!='mentions' AND p.organizer_id=$4)) AND ($6='' OR EXISTS(SELECT 1 FROM post_mentions m WHERE m.post_id=p.id AND m.kind='event' AND m.target_id=$6)) AND ($7!='following' OR EXISTS(SELECT 1 FROM follows f WHERE f.account_id=$8 AND (f.organizer_id=p.organizer_id))) AND ($9='' OR p.organizer_id IN(SELECT id FROM organizers WHERE category=$9) OR EXISTS(SELECT 1 FROM post_mentions m WHERE m.post_id=p.id AND (m.kind='event' AND m.target_id IN(SELECT id FROM events WHERE category=$9) OR m.kind='organizer' AND m.target_id IN(SELECT id FROM organizers WHERE category=$9)))) ORDER BY `+order+` LIMIT $10 OFFSET $11`, q.Get("city"), q.Get("q"), "%"+q.Get("q")+"%", q.Get("organizer"), q.Get("mode"), q.Get("event"), q.Get("tab"), a.uid(r), q.Get("category"), l, o)
	if e != nil {
		return e
	}
	send(w, 200, items)
	return nil
}
func (a *App) post(w http.ResponseWriter, r *http.Request) error {
	doc, e := a.one(r.Context(), `SELECT doc FROM post_view WHERE id=$1 AND NOT hidden`, r.PathValue("id"))
	if e != nil {
		return e
	}
	send(w, 200, doc)
	return nil
}
func (a *App) createPost(w http.ResponseWriter, r *http.Request) error {
	u, e := a.me(r)
	if e != nil {
		return e
	}
	if u["role"] == "admin" {
		return forbidden()
	}
	var in struct {
		Title    string `json:"title"`
		Body     string `json:"body"`
		Image    string `json:"image_url"`
		City     string `json:"city"`
		Link     string `json:"link_url"`
		Label    string `json:"link_label"`
		Mentions []struct {
			Kind string `json:"kind"`
			ID   string `json:"id"`
		} `json:"mentions"`
	}
	if e = decode(r, &in); e != nil {
		return e
	}
	if !textOK(in.Body, 1, 5000) || !textOK(in.Title, 0, 120) || !textOK(in.City, 2, 60) || !textOK(in.Label, 0, 80) || len(in.Mentions) > 20 || !safeURL(in.Link) {
		return bad("Isi postingan atau tautan tidak valid.")
	}
	if in.Link != "" && u["role"] != "organizer" {
		return forbidden()
	}
	if !a.mediaOK(r.Context(), u["id"].(string), in.Image) {
		return bad("Gambar tidak valid.")
	}
	tx, e := a.DB.Begin(r.Context())
	if e != nil {
		return e
	}
	defer tx.Rollback(r.Context())
	id := ID()
	var oid any
	if u["role"] == "organizer" {
		oid = u["organizer_id"]
	}
	_, e = tx.Exec(r.Context(), `INSERT INTO posts(id,account_id,organizer_id,title,body,image_url,city,link_url,link_label) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9)`, id, u["id"], oid, in.Title, in.Body, in.Image, in.City, in.Link, in.Label)
	if e != nil {
		return e
	}
	notified := map[string]bool{}
	for _, m := range in.Mentions {
		if m.Kind != "event" && m.Kind != "organizer" {
			return bad("Mention tidak valid.")
		}
		var owner string
		if m.Kind == "event" {
			e = tx.QueryRow(r.Context(), `SELECT o.account_id FROM events e JOIN organizers o ON o.id=e.organizer_id WHERE e.id=$1 AND e.published`, m.ID).Scan(&owner)
		} else {
			e = tx.QueryRow(r.Context(), `SELECT account_id FROM organizers WHERE id=$1`, m.ID).Scan(&owner)
		}
		if e != nil {
			return bad("Tujuan mention tidak tersedia.")
		}
		_, e = tx.Exec(r.Context(), `INSERT INTO post_mentions(post_id,kind,target_id) VALUES($1,$2,$3) ON CONFLICT DO NOTHING`, id, m.Kind, m.ID)
		if e != nil {
			return e
		}
		if owner != u["id"] && !notified[owner] {
			if e = notify(r.Context(), tx, owner, "mention", "Ruangmu disebut dalam percakapan", in.Body, "/post/"+id); e != nil {
				return e
			}
			notified[owner] = true
		}
	}
	if e = tx.Commit(r.Context()); e != nil {
		return e
	}
	doc, e := a.one(r.Context(), `SELECT doc FROM post_view WHERE id=$1`, id)
	if e != nil {
		return e
	}
	send(w, 201, doc)
	return nil
}
func (a *App) deletePost(w http.ResponseWriter, r *http.Request) error {
	u, e := a.me(r)
	if e != nil {
		return e
	}
	tag, e := a.DB.Exec(r.Context(), `UPDATE posts SET hidden=true WHERE id=$1 AND account_id=$2`, r.PathValue("id"), u["id"])
	if e != nil {
		return e
	}
	if tag.RowsAffected() == 0 {
		return forbidden()
	}
	send(w, 200, map[string]bool{"ok": true})
	return nil
}
func (a *App) pinPost(w http.ResponseWriter, r *http.Request) error {
	u, e := a.actor(r, "organizer")
	if e != nil {
		return e
	}
	var in struct {
		Pinned bool `json:"pinned"`
	}
	if e = decode(r, &in); e != nil {
		return e
	}
	tag, e := a.DB.Exec(r.Context(), `UPDATE posts SET pinned=$3 WHERE id=$1 AND account_id=$2 AND NOT hidden`, r.PathValue("id"), u["id"], in.Pinned)
	if e != nil {
		return e
	}
	if tag.RowsAffected() == 0 {
		return forbidden()
	}
	send(w, 200, map[string]bool{"pinned": in.Pinned})
	return nil
}
func (a *App) vote(w http.ResponseWriter, r *http.Request) error {
	u, e := a.me(r)
	if e != nil {
		return e
	}
	var in struct {
		Active bool `json:"active"`
	}
	if e = decode(r, &in); e != nil {
		return e
	}
	id := r.PathValue("id")
	var exists bool
	e = a.DB.QueryRow(r.Context(), `SELECT EXISTS(SELECT 1 FROM posts WHERE id=$1 AND NOT hidden)`, id).Scan(&exists)
	if e != nil {
		return e
	}
	if !exists {
		return missing()
	}
	if in.Active {
		_, e = a.DB.Exec(r.Context(), `INSERT INTO votes(post_id,account_id) VALUES($1,$2) ON CONFLICT DO NOTHING`, id, u["id"])
	} else {
		_, e = a.DB.Exec(r.Context(), `DELETE FROM votes WHERE post_id=$1 AND account_id=$2`, id, u["id"])
	}
	if e != nil {
		return e
	}
	send(w, 200, map[string]bool{"active": in.Active})
	return nil
}
func (a *App) comments(w http.ResponseWriter, r *http.Request) error {
	items, e := a.list(r.Context(), `SELECT to_jsonb(c)||jsonb_build_object('author',COALESCE(o.name,a.name),'avatar_url',COALESCE(o.avatar_url,a.avatar_url),'official',o.id IS NOT NULL) FROM comments c JOIN accounts a ON a.id=c.account_id LEFT JOIN organizers o ON o.account_id=a.id JOIN posts p ON p.id=c.post_id WHERE c.post_id=$1 AND NOT p.hidden ORDER BY c.created_at,c.id LIMIT 200`, r.PathValue("id"))
	if e != nil {
		return e
	}
	send(w, 200, items)
	return nil
}
func (a *App) addComment(w http.ResponseWriter, r *http.Request) error {
	u, e := a.me(r)
	if e != nil {
		return e
	}
	var in struct {
		Body   string `json:"body"`
		Parent string `json:"parent_id"`
	}
	if e = decode(r, &in); e != nil {
		return e
	}
	if !textOK(in.Body, 1, 2000) {
		return bad("Komentar harus berisi 1–2000 karakter.")
	}
	tx, e := a.DB.Begin(r.Context())
	if e != nil {
		return e
	}
	defer tx.Rollback(r.Context())
	id := r.PathValue("id")
	var owner string
	e = tx.QueryRow(r.Context(), `SELECT account_id FROM posts WHERE id=$1 AND NOT hidden`, id).Scan(&owner)
	if e != nil {
		return e
	}
	var parent any
	if in.Parent != "" {
		var post, replyOwner string
		var grandparent *string
		e = tx.QueryRow(r.Context(), `SELECT post_id,account_id,parent_id FROM comments WHERE id=$1`, in.Parent).Scan(&post, &replyOwner, &grandparent)
		if e != nil || post != id || grandparent != nil {
			return bad("Balasan hanya satu tingkat pada komentar event yang sama.")
		}
		parent = in.Parent
		if replyOwner != u["id"] && replyOwner != owner {
			if e = notify(r.Context(), tx, replyOwner, "reply", "Ada balasan untuk komentarmu", in.Body, "/post/"+id); e != nil {
				return e
			}
		}
	}
	cid := ID()
	_, e = tx.Exec(r.Context(), `INSERT INTO comments(id,post_id,account_id,parent_id,body) VALUES($1,$2,$3,$4,$5)`, cid, id, u["id"], parent, in.Body)
	if e != nil {
		return e
	}
	if owner != u["id"] {
		if e = notify(r.Context(), tx, owner, "reply", "Ada komentar baru", in.Body, "/post/"+id); e != nil {
			return e
		}
	}
	if e = tx.Commit(r.Context()); e != nil {
		return e
	}
	send(w, 201, map[string]string{"id": cid})
	return nil
}
func (a *App) mentions(w http.ResponseWriter, r *http.Request) error {
	q := r.URL.Query().Get("q")
	if len(q) > 100 {
		return bad("Pencarian terlalu panjang.")
	}
	items, e := a.list(r.Context(), `SELECT doc FROM (SELECT jsonb_build_object('kind','organizer','id',id,'slug',slug,'name',name,'city',city,'label','Pengelola') doc,name FROM organizers WHERE name ILIKE $1 UNION ALL SELECT jsonb_build_object('kind','event','id',id,'slug',slug,'name',title,'city',city,'label','Event') doc,title FROM events WHERE published AND title ILIKE $1) s ORDER BY name LIMIT 12`, "%"+q+"%")
	if e != nil {
		return e
	}
	send(w, 200, items)
	return nil
}
func (a *App) bookmarks(w http.ResponseWriter, r *http.Request) error {
	u, e := a.me(r)
	if e != nil {
		return e
	}
	items, e := a.list(r.Context(), `SELECT jsonb_build_object('kind',b.kind,'item',CASE WHEN b.kind='post' THEN (SELECT doc FROM post_view WHERE id=b.target_id AND NOT hidden) ELSE (SELECT doc FROM event_view WHERE id=b.target_id AND published) END) FROM bookmarks b WHERE account_id=$1 ORDER BY created_at DESC LIMIT 200`, u["id"])
	if e != nil {
		return e
	}
	send(w, 200, items)
	return nil
}
func (a *App) bookmark(w http.ResponseWriter, r *http.Request) error {
	u, e := a.me(r)
	if e != nil {
		return e
	}
	var in struct {
		Kind   string `json:"kind"`
		ID     string `json:"id"`
		Active bool   `json:"active"`
	}
	if e = decode(r, &in); e != nil {
		return e
	}
	var exists bool
	if in.Kind == "event" {
		e = a.DB.QueryRow(r.Context(), `SELECT EXISTS(SELECT 1 FROM events WHERE id=$1 AND published)`, in.ID).Scan(&exists)
	} else if in.Kind == "post" {
		e = a.DB.QueryRow(r.Context(), `SELECT EXISTS(SELECT 1 FROM posts WHERE id=$1 AND NOT hidden)`, in.ID).Scan(&exists)
	} else {
		return bad("Jenis simpan tidak valid.")
	}
	if e != nil {
		return e
	}
	if !exists {
		return missing()
	}
	if in.Active {
		_, e = a.DB.Exec(r.Context(), `INSERT INTO bookmarks(account_id,kind,target_id) VALUES($1,$2,$3) ON CONFLICT DO NOTHING`, u["id"], in.Kind, in.ID)
	} else {
		_, e = a.DB.Exec(r.Context(), `DELETE FROM bookmarks WHERE account_id=$1 AND kind=$2 AND target_id=$3`, u["id"], in.Kind, in.ID)
	}
	if e != nil {
		return e
	}
	send(w, 200, map[string]bool{"active": in.Active})
	return nil
}
func (a *App) reviews(w http.ResponseWriter, r *http.Request) error {
	q := r.URL.Query()
	items, e := a.list(r.Context(), `SELECT doc || jsonb_build_object('event_starts_at',(SELECT min(starts_at) FROM event_sessions WHERE event_id=review_view.event_id),'organizer_name',(SELECT name FROM organizers WHERE id=review_view.organizer_id)) FROM review_view WHERE NOT hidden AND ($1='' OR event_id=$1) AND ($2='' OR organizer_id=$2) ORDER BY (doc->>'created_at')::timestamptz DESC LIMIT 100`, q.Get("event"), q.Get("organizer"))
	if e != nil {
		return e
	}
	send(w, 200, items)
	return nil
}
func (a *App) saveReview(w http.ResponseWriter, r *http.Request) error {
	u, e := a.actor(r, "customer")
	if e != nil {
		return e
	}
	var in struct {
		EventID  string `json:"event_id"`
		Body     string `json:"body"`
		Attended bool   `json:"attended"`
	}
	if e = decode(r, &in); e != nil {
		return e
	}
	if !textOK(in.Body, 10, 3000) || !in.Attended {
		return bad("Ceritakan pengalamanmu dan konfirmasi bahwa kamu hadir.")
	}
	var ended bool
	e = a.DB.QueryRow(r.Context(), `SELECT published AND (SELECT max(ends_at) FROM event_sessions WHERE event_id=e.id)<now() FROM events e WHERE e.id=$1`, in.EventID).Scan(&ended)
	if e != nil {
		return e
	}
	if !ended {
		return bad("Ulasan dibuka setelah event selesai.")
	}
	_, e = a.DB.Exec(r.Context(), `INSERT INTO reviews(id,event_id,account_id,body) VALUES($1,$2,$3,$4) ON CONFLICT(event_id,account_id) DO UPDATE SET body=excluded.body,updated_at=now()`, ID(), in.EventID, u["id"], in.Body)
	if e != nil {
		return e
	}
	send(w, 200, map[string]bool{"ok": true})
	return nil
}
func (a *App) replyReview(w http.ResponseWriter, r *http.Request) error {
	u, e := a.actor(r, "organizer")
	if e != nil {
		return e
	}
	var in struct {
		Body string `json:"body"`
	}
	if e = decode(r, &in); e != nil {
		return e
	}
	if !textOK(in.Body, 1, 2000) {
		return bad("Isi balasan tidak valid.")
	}
	tx, e := a.DB.Begin(r.Context())
	if e != nil {
		return e
	}
	defer tx.Rollback(r.Context())
	var oid, owner string
	e = tx.QueryRow(r.Context(), `SELECT organizer_id,account_id FROM review_view WHERE id=$1 AND NOT hidden`, r.PathValue("id")).Scan(&oid, &owner)
	if e != nil {
		return e
	}
	if oid != u["organizer_id"] {
		return forbidden()
	}
	_, e = tx.Exec(r.Context(), `INSERT INTO review_replies(id,review_id,account_id,body) VALUES($1,$2,$3,$4)`, ID(), r.PathValue("id"), u["id"], in.Body)
	if e != nil {
		return e
	}
	if e = notify(r.Context(), tx, owner, "reply", "Pengelola membalas ulasanmu", in.Body, "/ruang/"+fmt.Sprint(u["organizer_slug"])+"?tab=ulasan"); e != nil {
		return e
	}
	if e = tx.Commit(r.Context()); e != nil {
		return e
	}
	send(w, 201, map[string]bool{"ok": true})
	return nil
}
func (a *App) helpfulReview(w http.ResponseWriter, r *http.Request) error {
	u, e := a.me(r)
	if e != nil {
		return e
	}
	var in struct {
		Active bool `json:"active"`
	}
	if e = decode(r, &in); e != nil {
		return e
	}
	if in.Active {
		_, e = a.DB.Exec(r.Context(), `INSERT INTO review_helpful(review_id,account_id) SELECT id,$2 FROM reviews WHERE id=$1 AND NOT hidden ON CONFLICT DO NOTHING`, r.PathValue("id"), u["id"])
	} else {
		_, e = a.DB.Exec(r.Context(), `DELETE FROM review_helpful WHERE review_id=$1 AND account_id=$2`, r.PathValue("id"), u["id"])
	}
	if e != nil {
		return e
	}
	send(w, 200, map[string]bool{"ok": true})
	return nil
}
