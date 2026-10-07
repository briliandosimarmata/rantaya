package app

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5"
	"net/http"
	"net/url"
	"strings"
)

func (a *App) authConfig(w http.ResponseWriter, r *http.Request) error {
	send(w, 200, map[string]any{"demo": a.Config.Demo, "google": a.Config.GoogleID != "" && a.Config.GoogleSecret != "", "cities": []string{"Karawang", "Jakarta"}})
	return nil
}
func (a *App) getMe(w http.ResponseWriter, r *http.Request) error {
	u, e := a.me(r)
	if e != nil {
		var f *Fault
		if errors.As(e, &f) && f.Status == 401 {
			send(w, 200, map[string]any{"user": nil, "follows": []string{}, "bookmarks": []any{}, "votes": []string{}, "helpful_reviews": []string{}, "unread": 0})
			return nil
		}
		return e
	}
	uid := u["id"].(string)
	follows, e := a.list(r.Context(), `SELECT to_jsonb(organizer_id) FROM follows WHERE account_id=$1`, uid)
	if e != nil {
		return e
	}
	bookmarks, e := a.list(r.Context(), `SELECT jsonb_build_object('kind',kind,'id',target_id) FROM bookmarks WHERE account_id=$1`, uid)
	if e != nil {
		return e
	}
	votes, e := a.list(r.Context(), `SELECT to_jsonb(post_id) FROM votes WHERE account_id=$1`, uid)
	if e != nil {
		return e
	}
	helpful, e := a.list(r.Context(), `SELECT to_jsonb(review_id) FROM review_helpful WHERE account_id=$1`, uid)
	if e != nil {
		return e
	}
	var unread int
	e = a.DB.QueryRow(r.Context(), `SELECT count(*) FROM notifications WHERE account_id=$1 AND read_at IS NULL`, uid).Scan(&unread)
	if e != nil {
		return e
	}
	send(w, 200, map[string]any{"user": u, "follows": follows, "bookmarks": bookmarks, "votes": votes, "helpful_reviews": helpful, "unread": unread})
	return nil
}
func (a *App) session(w http.ResponseWriter, r *http.Request, uid string) error {
	t := token()
	_, e := a.DB.Exec(r.Context(), `INSERT INTO auth_sessions(token_hash,account_id,expires_at) VALUES($1,$2,now()+interval '14 days')`, hash(t), uid)
	if e != nil {
		return e
	}
	http.SetCookie(w, &http.Cookie{Name: "ruang_session", Value: t, Path: "/", HttpOnly: true, Secure: a.Config.Secure, SameSite: http.SameSiteLaxMode, MaxAge: 14 * 86400})
	return nil
}
func (a *App) demoLogin(w http.ResponseWriter, r *http.Request) error {
	if !a.Config.Demo {
		return forbidden()
	}
	var in struct {
		Role string `json:"role"`
	}
	if e := decode(r, &in); e != nil {
		return e
	}
	ids := map[string]string{"customer": "customer-demo", "organizer": "organizer-demo", "admin": "admin-demo"}
	uid, ok := ids[in.Role]
	if !ok {
		return bad("Jenis akun tidak valid.")
	}
	if e := a.session(w, r, uid); e != nil {
		return e
	}
	send(w, 200, map[string]string{"message": "Berhasil masuk."})
	return nil
}
func (a *App) logout(w http.ResponseWriter, r *http.Request) error {
	if c, e := r.Cookie("ruang_session"); e == nil {
		if _, e = a.DB.Exec(r.Context(), `DELETE FROM auth_sessions WHERE token_hash=$1`, hash(c.Value)); e != nil {
			return e
		}
	}
	http.SetCookie(w, &http.Cookie{Name: "ruang_session", Value: "", Path: "/", HttpOnly: true, Secure: a.Config.Secure, SameSite: http.SameSiteLaxMode, MaxAge: -1})
	send(w, 200, map[string]bool{"ok": true})
	return nil
}
func (a *App) profile(w http.ResponseWriter, r *http.Request) error {
	u, e := a.me(r)
	if e != nil {
		return e
	}
	var in struct {
		Name           string          `json:"name"`
		City           string          `json:"city"`
		Bio            string          `json:"bio"`
		AvatarURL      string          `json:"avatar_url"`
		Interests      []string        `json:"interests"`
		Preferences    map[string]bool `json:"preferences"`
		OnboardingDone bool            `json:"onboarding_done"`
	}
	if e = decode(r, &in); e != nil {
		return e
	}
	if !textOK(in.Name, 2, 80) || !textOK(in.City, 2, 60) || !textOK(in.Bio, 0, 500) || len(in.Interests) > 10 {
		return bad("Nama, kota, atau bio tidak valid.")
	}
	if !a.mediaOK(r.Context(), u["id"].(string), in.AvatarURL) {
		return bad("Gambar profil tidak valid.")
	}
	b, _ := json.Marshal(in.Interests)
	prefs, _ := json.Marshal(in.Preferences)
	_, e = a.DB.Exec(r.Context(), `UPDATE accounts SET name=$2,city=$3,bio=$4,avatar_url=$5,interests=$6,preferences=$7,onboarding_done=$8 WHERE id=$1`, u["id"], strings.TrimSpace(in.Name), in.City, in.Bio, in.AvatarURL, string(b), string(prefs), in.OnboardingDone)
	if e != nil {
		return e
	}
	return a.getMe(w, r)
}
func (a *App) googleStart(w http.ResponseWriter, r *http.Request) error {
	if a.Config.GoogleID == "" || a.Config.GoogleSecret == "" {
		return &Fault{503, "google_unconfigured", "Login Google belum dikonfigurasi."}
	}
	role := r.URL.Query().Get("role")
	if role != "customer" && role != "organizer" {
		return bad("Jenis akun tidak valid.")
	}
	state := token()
	verifier := token()
	digest := sha256Bytes(verifier)
	_, e := a.DB.Exec(r.Context(), `INSERT INTO oauth_states(id,role,verifier,next_path,expires_at) VALUES($1,$2,$3,$4,now()+interval '10 minutes')`, hash(state), role, verifier, safeNext(r.URL.Query().Get("next")))
	if e != nil {
		return e
	}
	http.SetCookie(w, &http.Cookie{Name: "ruang_oauth", Value: state, Path: "/api/auth/google", HttpOnly: true, Secure: a.Config.Secure, SameSite: http.SameSiteLaxMode, MaxAge: 600})
	v := url.Values{"client_id": {a.Config.GoogleID}, "redirect_uri": {a.Config.GoogleRedirect}, "response_type": {"code"}, "scope": {"openid email profile"}, "state": {state}, "code_challenge": {base64.RawURLEncoding.EncodeToString(digest)}, "code_challenge_method": {"S256"}, "prompt": {"select_account"}}
	http.Redirect(w, r, "https://accounts.google.com/o/oauth2/v2/auth?"+v.Encode(), 302)
	return nil
}
func sha256Bytes(s string) []byte        { b, _ := hexDecode(hash(s)); return b }
func hexDecode(s string) ([]byte, error) { return hex.DecodeString(s) }
func (a *App) googleCallback(w http.ResponseWriter, r *http.Request) error {
	state := r.URL.Query().Get("state")
	c, e := r.Cookie("ruang_oauth")
	if e != nil || state == "" || c.Value != state {
		return bad("Sesi login Google tidak valid. Mulai lagi.")
	}
	var role, verifier, next string
	e = a.DB.QueryRow(r.Context(), `DELETE FROM oauth_states WHERE id=$1 AND expires_at>now() RETURNING role,verifier,next_path`, hash(state)).Scan(&role, &verifier, &next)
	if e != nil {
		return bad("Sesi login Google sudah berakhir.")
	}
	http.SetCookie(w, &http.Cookie{Name: "ruang_oauth", Path: "/api/auth/google", MaxAge: -1, HttpOnly: true, Secure: a.Config.Secure, SameSite: http.SameSiteLaxMode})
	if r.URL.Query().Get("error") != "" {
		http.Redirect(w, r, a.Config.AppURL+"/masuk/"+role+"?error=google_cancelled", 303)
		return nil
	}
	code := r.URL.Query().Get("code")
	if code == "" {
		return bad("Kode Google tidak ditemukan.")
	}
	values := url.Values{"client_id": {a.Config.GoogleID}, "client_secret": {a.Config.GoogleSecret}, "code": {code}, "code_verifier": {verifier}, "grant_type": {"authorization_code"}, "redirect_uri": {a.Config.GoogleRedirect}}
	req, e := http.NewRequestWithContext(r.Context(), "POST", "https://oauth2.googleapis.com/token", strings.NewReader(values.Encode()))
	if e != nil {
		return e
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, e := a.client.Do(req)
	if e != nil {
		return e
	}
	defer resp.Body.Close()
	var creds struct {
		AccessToken string `json:"access_token"`
	}
	if resp.StatusCode != 200 || json.NewDecoder(resp.Body).Decode(&creds) != nil || creds.AccessToken == "" {
		return bad("Google gagal memverifikasi login. Coba lagi.")
	}
	req, e = http.NewRequestWithContext(r.Context(), "GET", "https://openidconnect.googleapis.com/v1/userinfo", nil)
	if e != nil {
		return e
	}
	req.Header.Set("Authorization", "Bearer "+creds.AccessToken)
	info, e := a.client.Do(req)
	if e != nil {
		return e
	}
	defer info.Body.Close()
	var profile struct {
		Sub, Email, Name, Picture string
		Verified                  bool `json:"email_verified"`
	}
	if info.StatusCode != 200 || json.NewDecoder(info.Body).Decode(&profile) != nil || !profile.Verified || profile.Sub == "" {
		return bad("Akun Google belum terverifikasi.")
	}
	tx, e := a.DB.Begin(r.Context())
	if e != nil {
		return e
	}
	defer tx.Rollback(r.Context())
	var uid string
	e = tx.QueryRow(r.Context(), `SELECT id FROM accounts WHERE google_sub=$1 AND role=$2`, profile.Sub, role).Scan(&uid)
	fresh := errors.Is(e, pgx.ErrNoRows)
	if e != nil && !fresh {
		return e
	}
	if fresh {
		uid = ID()
		name := profile.Name
		if !textOK(name, 2, 80) {
			name = "Teman Panggung"
		}
		_, e = tx.Exec(r.Context(), `INSERT INTO accounts(id,role,google_sub,email,name,avatar_url) VALUES($1,$2,$3,$4,$5,$6)`, uid, role, profile.Sub, profile.Email, name, profile.Picture)
		if e != nil {
			return e
		}
		if role == "organizer" {
			oid := ID()
			_, e = tx.Exec(r.Context(), `INSERT INTO organizers(id,account_id,slug,name) VALUES($1,$2,$3,$4)`, oid, uid, slug(name)+"-"+oid[:8], name)
			if e != nil {
				return e
			}
		}
	}
	if e = tx.Commit(r.Context()); e != nil {
		return e
	}
	if e = a.session(w, r, uid); e != nil {
		return e
	}
	dest := "/"
	if role == "organizer" {
		dest = "/kelola"
	}
	if next = safeNext(next); next != "" {
		dest = next
	}
	if fresh {
		dest = "/onboarding?next=" + url.QueryEscape(dest)
	}
	http.Redirect(w, r, a.Config.AppURL+dest, 303)
	return nil
}

func safeNext(value string) string {
	if len(value) > 2048 {
		return ""
	}
	u, err := url.Parse(value)
	if err != nil || u.IsAbs() || u.Host != "" || !strings.HasPrefix(u.Path, "/") || strings.HasPrefix(u.Path, "//") || strings.ContainsAny(u.Path, "\\\r\n") {
		return ""
	}
	return u.String()
}
