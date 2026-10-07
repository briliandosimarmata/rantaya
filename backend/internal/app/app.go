package app

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed schema.sql reminders.sql checkout_continuity.sql
var schema embed.FS

type Config struct {
	Addr, DatabaseURL, AppURL, UploadDir, GoogleID, GoogleSecret, GoogleRedirect string
	AllowedOrigins                                                               []string
	Demo, Seed, Secure, Production                                               bool
	HoldMinutes                                                                  int
}

func Env(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
func Configuration() Config {
	n, _ := strconv.Atoi(Env("HOLD_MINUTES", "30"))
	if n < 1 || n > 120 {
		n = 30
	}
	production := Env("APP_ENV", "development") == "production"
	return Config{Addr: Env("API_ADDR", "127.0.0.1:8080"), DatabaseURL: Env("DATABASE_URL", "postgres://ruang:ruang_local@localhost:5432/ruang?sslmode=disable"), AppURL: strings.TrimRight(Env("APP_URL", "http://localhost:5173"), "/"), AllowedOrigins: configuredOrigins(os.Getenv("APP_ALLOWED_ORIGINS")), UploadDir: Env("UPLOAD_DIR", "./var/uploads"), GoogleID: os.Getenv("GOOGLE_CLIENT_ID"), GoogleSecret: os.Getenv("GOOGLE_CLIENT_SECRET"), GoogleRedirect: Env("GOOGLE_REDIRECT_URI", "http://localhost:5173/api/auth/google/callback"), Demo: Env("DEMO_LOGIN", "false") == "true", Seed: Env("SEED_DEMO", "false") == "true", Secure: production || Env("COOKIE_SECURE", "false") == "true", Production: production, HoldMinutes: n}
}

func configuredOrigins(value string) []string {
	var origins []string
	for _, origin := range strings.Split(value, ",") {
		if origin = strings.TrimSpace(origin); origin != "" {
			origins = append(origins, origin)
		}
	}
	return origins
}

func validOrigin(origin string) bool {
	u, err := url.Parse(origin)
	return err == nil && (u.Scheme == "http" || u.Scheme == "https") && u.Hostname() != "" &&
		u.User == nil && !strings.Contains(u.Host, "*") && origin == u.Scheme+"://"+u.Host
}

func (c Config) allowsOrigin(origin string) bool {
	if origin == "" {
		return false
	}
	if origin == c.AppURL {
		return true
	}
	for _, allowed := range c.AllowedOrigins {
		if origin == allowed {
			return true
		}
	}
	return false
}

type App struct {
	DB     *pgxpool.Pool
	Config Config
	client *http.Client
	mu     sync.Mutex
	limits map[string]limit
}
type limit struct {
	start time.Time
	n     int
}
type Fault struct {
	Status        int
	Code, Message string
}

func (f *Fault) Error() string      { return f.Message }
func bad(message string) error      { return &Fault{400, "invalid_input", message} }
func forbidden() error              { return &Fault{403, "forbidden", "Akun ini tidak memiliki akses."} }
func conflict(message string) error { return &Fault{409, "conflict", message} }
func missing() error                { return &Fault{404, "not_found", "Data tidak ditemukan."} }
func ID() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(err)
	}
	return hex.EncodeToString(b[:])
}
func token() string {
	var b [32]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(err)
	}
	return hex.EncodeToString(b[:])
}
func hash(s string) string { h := sha256.Sum256([]byte(s)); return hex.EncodeToString(h[:]) }
func decode(r *http.Request, v any) error {
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(v); err != nil {
		return bad("Format permintaan tidak valid.")
	}
	var extra any
	if err := d.Decode(&extra); err != io.EOF {
		return bad("Hanya satu objek JSON yang diizinkan.")
	}
	return nil
}
func send(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
func serve(fn func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			var f *Fault
			if errors.Is(err, pgx.ErrNoRows) {
				f = &Fault{404, "not_found", "Data tidak ditemukan."}
			} else if !errors.As(err, &f) {
				slog.Error("request failed", "path", r.URL.Path, "error", err)
				f = &Fault{500, "internal_error", "Terjadi kendala. Coba lagi."}
			}
			send(w, f.Status, map[string]any{"error": map[string]string{"code": f.Code, "message": f.Message}})
		}
	}
}

// Use pgx transactions directly; no HTTP framework or ORM.
func (a *App) list(ctx context.Context, q string, args ...any) ([]json.RawMessage, error) {
	rows, err := a.DB.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []json.RawMessage{}
	for rows.Next() {
		var b []byte
		if err = rows.Scan(&b); err != nil {
			return nil, err
		}
		out = append(out, json.RawMessage(b))
	}
	return out, rows.Err()
}
func (a *App) one(ctx context.Context, q string, args ...any) (json.RawMessage, error) {
	var b []byte
	err := a.DB.QueryRow(ctx, q, args...).Scan(&b)
	return json.RawMessage(b), err
}
func (a *App) me(r *http.Request) (map[string]any, error) {
	t := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	if t == "" {
		if c, err := r.Cookie("ruang_session"); err == nil {
			t = c.Value
		}
	}
	if t == "" {
		return nil, &Fault{401, "unauthenticated", "Silakan masuk terlebih dahulu."}
	}
	b, err := a.one(r.Context(), `SELECT to_jsonb(a)-'google_sub' || jsonb_build_object('organizer_id',o.id,'organizer_slug',o.slug) FROM accounts a JOIN auth_sessions s ON s.account_id=a.id LEFT JOIN organizers o ON o.account_id=a.id WHERE s.token_hash=$1 AND s.expires_at>now()`, hash(t))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, &Fault{401, "unauthenticated", "Sesi sudah berakhir. Silakan masuk lagi."}
	}
	if err != nil {
		return nil, err
	}
	var m map[string]any
	err = json.Unmarshal(b, &m)
	return m, err
}
func (a *App) actor(r *http.Request, role string) (map[string]any, error) {
	u, err := a.me(r)
	if err != nil {
		return nil, err
	}
	if role != "" && u["role"] != role {
		return nil, forbidden()
	}
	return u, nil
}
func (a *App) uid(r *http.Request) string {
	u, err := a.me(r)
	if err != nil {
		return ""
	}
	return u["id"].(string)
}
func textOK(s string, min, max int) bool {
	n := len([]rune(strings.TrimSpace(s)))
	return n >= min && n <= max
}
func safeURL(s string) bool {
	if s == "" {
		return true
	}
	u, e := url.Parse(s)
	return e == nil && u.User == nil && u.Host != "" && (u.Scheme == "https" || u.Scheme == "http")
}
func mapsOK(s string) bool {
	if s == "" {
		return true
	}
	u, e := url.Parse(s)
	if e != nil || u.Scheme != "https" || u.User != nil {
		return false
	}
	h := strings.ToLower(u.Hostname())
	return h == "maps.app.goo.gl" || h == "maps.google.com" || h == "goo.gl" && strings.HasPrefix(u.Path, "/maps") || (h == "google.com" || h == "www.google.com" || h == "google.co.id" || h == "www.google.co.id") && strings.HasPrefix(u.Path, "/maps")
}
func slug(s string) string {
	var b strings.Builder
	dash := false
	for _, r := range strings.ToLower(s) {
		if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' {
			b.WriteRune(r)
			dash = false
		} else if !dash && b.Len() > 0 {
			b.WriteByte('-')
			dash = true
		}
	}
	v := strings.Trim(b.String(), "-")
	if v == "" {
		v = "ruang"
	}
	return v
}
func notify(ctx context.Context, tx pgx.Tx, uid, kind, title, body, path string) error {
	_, err := tx.Exec(ctx, `INSERT INTO notifications(id,account_id,kind,title,body,url) SELECT $1,$2,$3,$4,$5,$6 FROM accounts WHERE id=$2 AND ($3<>'reply' OR COALESCE((preferences->>'replies')::boolean,true))`, ID(), uid, kind, title, body, path)
	return err
}
func New(ctx context.Context, c Config) (*App, error) {
	if c.Production && (c.Demo || c.Seed) {
		return nil, fmt.Errorf("DEMO_LOGIN and SEED_DEMO must be false in production")
	}
	if !validOrigin(c.AppURL) {
		return nil, fmt.Errorf("APP_URL must be an exact HTTP(S) origin without credentials, path or wildcard")
	}
	for _, origin := range c.AllowedOrigins {
		if !validOrigin(origin) {
			return nil, fmt.Errorf("APP_ALLOWED_ORIGINS must contain exact HTTP(S) origins without credentials, paths or wildcards")
		}
	}
	pconf, err := pgxpool.ParseConfig(c.DatabaseURL)
	if err != nil {
		return nil, err
	}
	if pconf.MaxConns > 8 {
		pconf.MaxConns = 8
	}
	pconf.MinConns = 0
	pconf.MaxConnIdleTime = 5 * time.Minute
	db, err := pgxpool.NewWithConfig(ctx, pconf)
	if err != nil {
		return nil, err
	}
	a := &App{DB: db, Config: c, client: &http.Client{Timeout: 15 * time.Second}, limits: map[string]limit{}}
	if err = db.Ping(ctx); err != nil {
		db.Close()
		return nil, err
	}
	if err = a.Migrate(ctx); err != nil {
		db.Close()
		return nil, err
	}
	if err = os.MkdirAll(c.UploadDir, 0700); err != nil {
		db.Close()
		return nil, err
	}
	if c.Seed {
		err = a.Seed(ctx)
		if err != nil {
			db.Close()
			return nil, err
		}
	}
	return a, nil
}
func (a *App) Migrate(ctx context.Context) error {
	tx, e := a.DB.Begin(ctx)
	if e != nil {
		return e
	}
	defer tx.Rollback(ctx)
	if _, e = tx.Exec(ctx, `SELECT pg_advisory_xact_lock(7465821)`); e != nil {
		return e
	}
	if _, e = tx.Exec(ctx, `CREATE TABLE IF NOT EXISTS schema_migrations(version integer PRIMARY KEY,applied_at timestamptz NOT NULL DEFAULT now())`); e != nil {
		return e
	}
	var exists bool
	e = tx.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version=1)`).Scan(&exists)
	if e != nil {
		return e
	}
	if !exists {
		b, _ := schema.ReadFile("schema.sql")
		if _, e = tx.Exec(ctx, string(b)); e != nil {
			return e
		}
		if _, e = tx.Exec(ctx, `INSERT INTO schema_migrations(version) VALUES(1)`); e != nil {
			return e
		}
	}
	if e = tx.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version=2)`).Scan(&exists); e != nil {
		return e
	}
	if !exists {
		b, _ := schema.ReadFile("reminders.sql")
		if _, e = tx.Exec(ctx, string(b)); e != nil {
			return e
		}
		if _, e = tx.Exec(ctx, `INSERT INTO schema_migrations(version) VALUES(2)`); e != nil {
			return e
		}
	}
	if e = tx.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version=3)`).Scan(&exists); e != nil {
		return e
	}
	if !exists {
		b, err := schema.ReadFile("checkout_continuity.sql")
		if err != nil {
			return err
		}
		if _, e = tx.Exec(ctx, string(b)); e != nil {
			return e
		}
		if _, e = tx.Exec(ctx, `INSERT INTO schema_migrations(version) VALUES(3)`); e != nil {
			return e
		}
	}
	return tx.Commit(ctx)
}
func (a *App) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", serve(a.health))
	mux.HandleFunc("GET /api/auth/config", serve(a.authConfig))
	mux.HandleFunc("GET /api/auth/me", serve(a.getMe))
	mux.HandleFunc("POST /api/auth/demo", serve(a.demoLogin))
	mux.HandleFunc("POST /api/auth/logout", serve(a.logout))
	mux.HandleFunc("GET /api/auth/google/start", serve(a.googleStart))
	mux.HandleFunc("GET /api/auth/google/callback", serve(a.googleCallback))
	mux.HandleFunc("PATCH /api/profile", serve(a.profile))
	mux.HandleFunc("GET /api/organizers", serve(a.organizers))
	mux.HandleFunc("GET /api/organizers/{id}", serve(a.organizer))
	mux.HandleFunc("PATCH /api/organizer", serve(a.editOrganizer))
	mux.HandleFunc("POST /api/organizers/{id}/follow", serve(a.follow))
	mux.HandleFunc("GET /api/events", serve(a.events))
	mux.HandleFunc("GET /api/events/{id}", serve(a.event))
	mux.HandleFunc("POST /api/events", serve(a.saveEvent))
	mux.HandleFunc("PUT /api/events/{id}", serve(a.saveEvent))
	mux.HandleFunc("GET /api/posts", serve(a.posts))
	mux.HandleFunc("GET /api/posts/{id}", serve(a.post))
	mux.HandleFunc("POST /api/posts", serve(a.createPost))
	mux.HandleFunc("DELETE /api/posts/{id}", serve(a.deletePost))
	mux.HandleFunc("PATCH /api/posts/{id}/pin", serve(a.pinPost))
	mux.HandleFunc("POST /api/posts/{id}/vote", serve(a.vote))
	mux.HandleFunc("GET /api/posts/{id}/comments", serve(a.comments))
	mux.HandleFunc("POST /api/posts/{id}/comments", serve(a.addComment))
	mux.HandleFunc("GET /api/mentions", serve(a.mentions))
	mux.HandleFunc("GET /api/bookmarks", serve(a.bookmarks))
	mux.HandleFunc("POST /api/bookmarks", serve(a.bookmark))
	mux.HandleFunc("GET /api/reviews", serve(a.reviews))
	mux.HandleFunc("POST /api/reviews", serve(a.saveReview))
	mux.HandleFunc("POST /api/reviews/{id}/reply", serve(a.replyReview))
	mux.HandleFunc("POST /api/reviews/{id}/helpful", serve(a.helpfulReview))
	mux.HandleFunc("GET /api/products", serve(a.products))
	mux.HandleFunc("GET /api/products/{id}", serve(a.product))
	mux.HandleFunc("POST /api/products", serve(a.saveProduct))
	mux.HandleFunc("PUT /api/products/{id}", serve(a.saveProduct))
	mux.HandleFunc("GET /api/payment-methods", serve(a.paymentMethods))
	mux.HandleFunc("POST /api/payment-methods", serve(a.savePaymentMethod))
	mux.HandleFunc("PUT /api/payment-methods/{id}", serve(a.savePaymentMethod))
	mux.HandleFunc("POST /api/uploads", serve(a.upload))
	mux.HandleFunc("GET /api/uploads/{id}", serve(a.download))
	mux.HandleFunc("GET /api/orders", serve(a.orders))
	mux.HandleFunc("POST /api/orders", serve(a.createOrder))
	mux.HandleFunc("GET /api/orders/{id}", serve(a.order))
	mux.HandleFunc("PATCH /api/orders/{id}/payment", serve(a.selectPayment))
	mux.HandleFunc("POST /api/orders/{id}/proof", serve(a.submitProof))
	mux.HandleFunc("POST /api/orders/{id}/cancel", serve(a.cancelOrder))
	mux.HandleFunc("POST /api/orders/{id}/review", serve(a.reviewOrder))
	mux.HandleFunc("GET /api/tickets", serve(a.tickets))
	mux.HandleFunc("POST /api/check-in", serve(a.checkIn))
	mux.HandleFunc("GET /api/notifications", serve(a.notifications))
	mux.HandleFunc("POST /api/notifications/read", serve(a.readNotifications))
	mux.HandleFunc("GET /api/dashboard", serve(a.dashboard))
	mux.HandleFunc("POST /api/activity/click", serve(a.activityClick))
	mux.HandleFunc("POST /api/reports", serve(a.report))
	mux.HandleFunc("GET /api/admin/reports", serve(a.reports))
	mux.HandleFunc("PATCH /api/admin/reports/{id}", serve(a.resolveReport))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		origin := r.Header.Get("Origin")
		allowedOrigin := a.Config.allowsOrigin(origin)
		if allowedOrigin {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Vary", "Origin")
		}
		if r.Method == "OPTIONS" {
			if !allowedOrigin {
				send(w, 403, map[string]any{"error": map[string]string{"message": "Origin tidak diizinkan."}})
				return
			}
			w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization,Idempotency-Key")
			w.WriteHeader(204)
			return
		}
		if r.Method != "GET" && r.Method != "HEAD" && (origin != "" && !allowedOrigin || r.Header.Get("Sec-Fetch-Site") == "cross-site") {
			send(w, 403, map[string]any{"error": map[string]string{"code": "origin_rejected", "message": "Origin tidak diizinkan."}})
			return
		}
		ip := r.RemoteAddr
		if i := strings.LastIndex(ip, ":"); i >= 0 {
			ip = ip[:i]
		}
		a.mu.Lock()
		l := a.limits[ip]
		if time.Since(l.start) > time.Minute {
			l = limit{start: time.Now()}
		}
		l.n++
		a.limits[ip] = l
		if len(a.limits) > 10000 {
			for k, v := range a.limits {
				if time.Since(v.start) > 2*time.Minute {
					delete(a.limits, k)
				}
			}
		}
		a.mu.Unlock()
		if l.n > 600 {
			w.Header().Set("Retry-After", "60")
			send(w, 429, map[string]any{"error": map[string]string{"message": "Terlalu banyak permintaan. Coba sebentar lagi."}})
			return
		}
		r.Body = http.MaxBytesReader(w, r.Body, 10<<20)
		ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
		defer cancel()
		mux.ServeHTTP(w, r.WithContext(ctx))
	})
}
func (a *App) health(w http.ResponseWriter, r *http.Request) error {
	if err := a.DB.Ping(r.Context()); err != nil {
		return err
	}
	send(w, 200, map[string]any{"status": "ok", "service": "ruang-go-api"})
	return nil
}
func (a *App) Expire(ctx context.Context) error {
	tx, e := a.DB.Begin(ctx)
	if e != nil {
		return e
	}
	defer tx.Rollback(ctx)
	rows, e := tx.Query(ctx, `UPDATE orders SET status='expired',updated_at=now() WHERE status='awaiting_payment' AND expires_at<=now() RETURNING id,account_id`)
	if e != nil {
		return e
	}
	type pair struct{ id, uid string }
	items := []pair{}
	for rows.Next() {
		var p pair
		if e = rows.Scan(&p.id, &p.uid); e != nil {
			rows.Close()
			return e
		}
		items = append(items, p)
	}
	rows.Close()
	if e = rows.Err(); e != nil {
		return e
	}
	for _, p := range items {
		if _, e = tx.Exec(ctx, `INSERT INTO order_history(order_id,status,note) VALUES($1,'expired','Batas unggah bukti berakhir.')`, p.id); e != nil {
			return e
		}
		if e = notify(ctx, tx, p.uid, "order", "Waktu pembayaran berakhir", "Jika sudah transfer, hubungi pengelola melalui detail pesanan.", "/transaksi/"+p.id); e != nil {
			return e
		}
	}
	return tx.Commit(ctx)
}

// Remind once per account/session, and prune expired authentication state.
func (a *App) Remind(ctx context.Context) error {
	tx, e := a.DB.Begin(ctx)
	if e != nil {
		return e
	}
	defer tx.Rollback(ctx)
	_, e = tx.Exec(ctx, `WITH due AS (
 INSERT INTO event_reminders(account_id,session_id)
 SELECT DISTINCT o.account_id,s.id FROM orders o JOIN event_sessions s ON s.id=o.session_id JOIN accounts a ON a.id=o.account_id
 WHERE o.status='approved' AND s.starts_at>now() AND s.starts_at<=now()+interval '24 hours'
 AND COALESCE((a.preferences->>'reminders')::boolean,true)
 ON CONFLICT DO NOTHING RETURNING account_id,session_id
 ) INSERT INTO notifications(id,account_id,kind,title,body,url)
 SELECT 'reminder-'||d.account_id||'-'||d.session_id,d.account_id,'reminder','Pertunjukanmu segera dimulai',e.title||' · '||s.label,'/tiket'
 FROM due d JOIN event_sessions s ON s.id=d.session_id JOIN events e ON e.id=s.event_id`)
	if e != nil {
		return e
	}
	_, e = tx.Exec(ctx, `WITH due AS (
 INSERT INTO payment_reminders(order_id)
 SELECT id FROM orders WHERE status='awaiting_payment' AND expires_at>now() AND expires_at<=now()+interval '5 minutes'
 ON CONFLICT DO NOTHING RETURNING order_id
 ) INSERT INTO notifications(id,account_id,kind,title,body,url)
 SELECT 'payment-reminder-'||o.id,o.account_id,'order','Batas pembayaran sebentar lagi','Lanjutkan pembayaran atau batalkan jika belum transfer.','/transaksi/'||o.id
 FROM due d JOIN orders o ON o.id=d.order_id WHERE o.status='awaiting_payment'`)
	if e != nil {
		return e
	}
	if _, e = tx.Exec(ctx, `DELETE FROM auth_sessions WHERE expires_at<=now()`); e != nil {
		return e
	}
	if _, e = tx.Exec(ctx, `DELETE FROM oauth_states WHERE expires_at<=now()`); e != nil {
		return e
	}
	return tx.Commit(ctx)
}
