package app

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDemoSecureConfiguration(t *testing.T) {
	for _, tc := range []struct {
		name, environment, cookieSecure, demo, seed string
		secure, production                          bool
	}{
		{"local demo", "development", "false", "true", "true", false, false},
		{"HTTPS demo", "demo", "true", "true", "true", true, false},
		{"production cannot disable secure cookie", "production", "false", "false", "false", true, true},
		{"production rejects demo login", "production", "true", "true", "false", true, true},
		{"production rejects demo seed", "production", "false", "false", "true", true, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("APP_ENV", tc.environment)
			t.Setenv("COOKIE_SECURE", tc.cookieSecure)
			t.Setenv("DEMO_LOGIN", tc.demo)
			t.Setenv("SEED_DEMO", tc.seed)
			c := Configuration()
			if c.Secure != tc.secure || c.Production != tc.production {
				t.Fatal("cookie security must be independent of demo mode; production always requires secure cookies")
			}
			if c.Production && (c.Demo || c.Seed) {
				if _, err := New(context.Background(), c); err == nil || !strings.Contains(err.Error(), "must be false in production") {
					t.Fatal("production must reject demo login and seed before database startup")
				}
			}
		})
	}
}

func TestConfiguredOrigins(t *testing.T) {
	t.Setenv("APP_URL", "http://192.0.2.10:5173")
	t.Setenv("APP_ALLOWED_ORIGINS", " http://localhost:5173, ,http://127.0.0.1:5173 ")
	c := Configuration()
	for _, origin := range []string{c.AppURL, "http://localhost:5173", "http://127.0.0.1:5173"} {
		if !c.allowsOrigin(origin) {
			t.Errorf("configured origin rejected: %s", origin)
		}
	}
	for _, origin := range []string{"", "null", "http://localhost:5174", "https://localhost:5173", "http://localhost.evil.test:5173"} {
		if c.allowsOrigin(origin) {
			t.Errorf("unconfigured origin accepted: %s", origin)
		}
	}
	t.Setenv("APP_ALLOWED_ORIGINS", "")
	if Configuration().allowsOrigin("http://localhost:5173") {
		t.Fatal("localhost must require explicit configuration when APP_URL is different")
	}
}

func TestOriginConfigurationValidation(t *testing.T) {
	for _, origin := range []string{"*", "http://*.example.test", "null", "javascript:alert(1)", "https://user:password@example.test", "https://example.test/path", "https://example.test?query=1", "https://example.test#fragment"} {
		for _, additional := range []bool{false, true} {
			c := Config{AppURL: origin}
			if additional {
				c.AppURL = "https://example.test"
				c.AllowedOrigins = []string{origin}
			}
			if _, err := New(context.Background(), c); err == nil || !strings.Contains(err.Error(), "origin") {
				t.Errorf("invalid origin configuration not rejected before DB startup: %q", origin)
			}
		}
	}
}

func TestOriginMiddleware(t *testing.T) {
	a := &App{Config: Config{AppURL: "http://192.0.2.10:5173", AllowedOrigins: []string{"http://localhost:5173", "http://127.0.0.1:5173"}, Demo: true}, limits: map[string]limit{}}
	h := a.Handler()
	for _, tc := range []struct {
		name, method, origin, fetchSite string
		status                          int
	}{
		{"LAN login", "POST", a.Config.AppURL, "same-origin", 400},
		{"localhost login", "POST", "http://localhost:5173", "same-origin", 400},
		{"loopback login", "POST", "http://127.0.0.1:5173", "same-origin", 400},
		{"untrusted origin", "POST", "https://evil.test", "", 403},
		{"lookalike host", "POST", "http://localhost.evil.test:5173", "", 403},
		{"different port", "POST", "http://localhost:5174", "", 403},
		{"opaque origin", "POST", "null", "", 403},
		{"cross-site with trusted origin", "POST", "http://localhost:5173", "cross-site", 403},
		{"cross-site without origin", "POST", "", "cross-site", 403},
		{"native request without origin", "POST", "", "", 400},
		{"LAN preflight", "OPTIONS", a.Config.AppURL, "", 204},
		{"localhost preflight", "OPTIONS", "http://localhost:5173", "", 204},
		{"untrusted preflight", "OPTIONS", "https://evil.test", "", 403},
		{"missing preflight origin", "OPTIONS", "", "", 403},
	} {
		t.Run(tc.name, func(t *testing.T) {
			// Invalid role reaches input validation for trusted origins without needing a DB.
			r := httptest.NewRequest(tc.method, "/api/auth/demo", strings.NewReader(`{"role":"invalid"}`))
			r.Header.Set("Origin", tc.origin)
			r.Header.Set("Sec-Fetch-Site", tc.fetchSite)
			w := httptest.NewRecorder()
			h.ServeHTTP(w, r)
			if w.Code != tc.status {
				t.Fatalf("status %d, want %d", w.Code, tc.status)
			}
			if tc.status != 403 && tc.origin != "" {
				if w.Header().Get("Access-Control-Allow-Origin") != tc.origin || w.Header().Get("Access-Control-Allow-Credentials") != "true" {
					t.Fatal("CORS must return the exact configured origin with credentials")
				}
			}
			if tc.method == "POST" && tc.status == 403 {
				var body struct{ Error struct{ Code string } }
				if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil || body.Error.Code != "origin_rejected" {
					t.Fatal("untrusted mutation must be rejected by origin protection")
				}
			}
			if tc.status == 403 && !a.Config.allowsOrigin(tc.origin) && w.Header().Get("Access-Control-Allow-Origin") != "" {
				t.Fatal("untrusted origin received CORS permission")
			}
		})
	}
}

func TestURLPolicies(t *testing.T) {
	for _, v := range []string{"javascript:alert(1)", "data:text/html,test", "https://u:p@example.com"} {
		if safeURL(v) {
			t.Errorf("unsafe URL accepted: %s", v)
		}
	}
	for _, v := range []string{"", "https://example.com/tiket", "https://wa.me/628123"} {
		if !safeURL(v) {
			t.Errorf("valid URL rejected: %s", v)
		}
	}
	for _, v := range []string{"https://maps.app.goo.gl/abc", "https://www.google.com/maps?q=Karawang"} {
		if !mapsOK(v) {
			t.Errorf("Maps URL rejected: %s", v)
		}
	}
	for _, v := range []string{"https://google.com.evil.test/maps", "https://example.com/maps", "javascript:alert(1)"} {
		if mapsOK(v) {
			t.Errorf("invalid Maps URL accepted: %s", v)
		}
	}
}
func TestStrictJSON(t *testing.T) {
	for _, body := range []string{`{"title":"ok"} {}`, `{"title":"ok"} garbage`, `{"title":"ok","unexpected":1}`} {
		var dst struct {
			Title string `json:"title"`
		}
		r := httptest.NewRequest("POST", "/", strings.NewReader(body))
		if decode(r, &dst) == nil {
			t.Errorf("accepted invalid JSON %s", body)
		}
	}
	var dst struct {
		Title string `json:"title"`
	}
	if e := decode(httptest.NewRequest("POST", "/", bytes.NewBufferString(`{"title":"ok"}`)), &dst); e != nil {
		t.Fatal(e)
	}
}
func TestOpaqueIdentifiers(t *testing.T) {
	a, b := token(), token()
	if a == b || len(a) < 40 {
		t.Fatal("ticket/session tokens must be independent and unpredictable")
	}
	if hash(a) == a {
		t.Fatal("session token stored unhashed")
	}
}

func TestLoginReturnDestination(t *testing.T) {
	for _, value := range []string{"https://evil.test", "//evil.test", "/%2f%2fevil.test", "/%5cevil.test", "javascript:alert(1)", "/\r\nLocation: https://evil.test"} {
		if safeNext(value) != "" {
			t.Errorf("external return destination accepted: %q", value)
		}
	}
	for _, value := range []string{"/event/di-balik-layar/tiket?session=s-e1", "/transaksi/order-1", "/?tab=following"} {
		if safeNext(value) != value {
			t.Errorf("local return destination lost: %q", value)
		}
	}
}
