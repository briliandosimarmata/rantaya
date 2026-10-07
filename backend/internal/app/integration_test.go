package app

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
	"time"
)

// Runs against an isolated TEST_DATABASE_URL. Never point this at a production database.
func TestCommerceIntegration(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("set TEST_DATABASE_URL to an isolated PostgreSQL database")
	}
	ctx := context.Background()
	a, e := New(ctx, Config{DatabaseURL: dsn, UploadDir: t.TempDir(), AppURL: "http://localhost:5173", HoldMinutes: 30})
	if e != nil {
		t.Fatal(e)
	}
	defer a.DB.Close()
	h := a.Handler()
	prefix := ID()
	customer, other, owner, wrongOwner, oid := prefix+"c", prefix+"x", prefix+"o", prefix+"w", prefix+"org"
	for _, u := range []struct{ id, role string }{{customer, "customer"}, {other, "customer"}, {owner, "organizer"}, {wrongOwner, "organizer"}} {
		if _, e = a.DB.Exec(ctx, `INSERT INTO accounts(id,role,email,name,onboarding_done) VALUES($1,$2,$1||'@example.test',$1,true)`, u.id, u.role); e != nil {
			t.Fatal(e)
		}
	}
	for _, u := range []struct{ id, account string }{{oid, owner}, {prefix + "wrongorg", wrongOwner}} {
		if _, e = a.DB.Exec(ctx, `INSERT INTO organizers(id,account_id,slug,name,description) VALUES($1,$2,$1,$1,'Komunitas seni untuk pengujian integrasi')`, u.id, u.account); e != nil {
			t.Fatal(e)
		}
	}
	tokens := map[string]string{}
	for _, uid := range []string{customer, other, owner, wrongOwner} {
		tokens[uid] = token()
		if _, e = a.DB.Exec(ctx, `INSERT INTO auth_sessions(token_hash,account_id,expires_at) VALUES($1,$2,now()+interval '1 hour')`, hash(tokens[uid]), uid); e != nil {
			t.Fatal(e)
		}
	}
	call := func(uid, method, path string, body any, want int) map[string]any {
		t.Helper()
		var payload []byte
		if body != nil {
			payload, _ = json.Marshal(body)
		}
		r := httptest.NewRequest(method, "/api"+path, bytes.NewReader(payload))
		if uid != "" {
			r.AddCookie(&http.Cookie{Name: "ruang_session", Value: tokens[uid]})
		}
		r.Header.Set("Origin", a.Config.AppURL)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		if w.Code != want {
			t.Fatalf("%s %s: got %d want %d: %s", method, path, w.Code, want, w.Body.String())
		}
		out := map[string]any{}
		_ = json.Unmarshal(w.Body.Bytes(), &out)
		return out
	}
	method := call(owner, "POST", "/payment-methods", map[string]any{"kind": "bank", "provider": "Bank TEST", "number": "0000000000", "holder": "TEST JANGAN TRANSFER", "note": "dummy", "enabled": true}, 200)
	mid := method["id"].(string)
	start := time.Now().Add(time.Hour)
	end := start.Add(time.Hour)
	event := call(owner, "POST", "/events", map[string]any{"title": "Integration event", "description": "Pertunjukan untuk pengujian otomatis saja.", "category": "Teater", "city": "Karawang", "venue": "Venue TEST", "lineup": []any{}, "published": true, "sessions": []any{map[string]any{"label": "Sesi TEST", "starts_at": start.Format(time.RFC3339), "ends_at": end.Format(time.RFC3339), "price": 50000, "capacity": 6}}}, 200)
	eid := event["id"].(string)
	sid := event["sessions"].([]any)[0].(map[string]any)["id"].(string)
	create := func(key string, qty int, want int) map[string]any {
		return call(customer, "POST", "/orders", map[string]any{"session_id": sid, "quantity": qty, "idempotency_key": prefix + key}, want)
	}
	order := create("first", 2, 200)
	id := order["id"].(string)
	if create("first", 2, 200)["id"] != id {
		t.Fatal("idempotency created duplicate order")
	}
	call(owner, "POST", "/orders", map[string]any{"session_id": sid, "quantity": 1, "idempotency_key": prefix + "role"}, 403)
	call(other, "GET", "/orders/"+id, nil, 404)
	call(wrongOwner, "GET", "/orders/"+id, nil, 404)
	call(customer, "PATCH", "/orders/"+id+"/payment", map[string]any{"payment_method_id": mid}, 200)
	call(owner, "PUT", "/payment-methods/"+mid, map[string]any{"kind": "bank", "provider": "Bank TEST", "number": "9999999999", "holder": "MASTER BARU", "note": "dummy", "enabled": true}, 200)
	snapshot := call(customer, "GET", "/orders/"+id, nil, 200)["payment_snapshot"].(map[string]any)
	if snapshot["number"] != "0000000000" {
		t.Fatal("payment snapshot changed with master")
	}
	uploadProof := func() string {
		var body bytes.Buffer
		writer := multipart.NewWriter(&body)
		_ = writer.WriteField("purpose", "proof")
		file, _ := writer.CreateFormFile("file", "proof.png")
		_ = png.Encode(file, image.NewRGBA(image.Rect(0, 0, 8, 8)))
		_ = writer.Close()
		r := httptest.NewRequest("POST", "/api/uploads", &body)
		r.Header.Set("Content-Type", writer.FormDataContentType())
		r.AddCookie(&http.Cookie{Name: "ruang_session", Value: tokens[customer]})
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		if w.Code != 201 {
			t.Fatal(w.Body.String())
		}
		var u map[string]any
		_ = json.Unmarshal(w.Body.Bytes(), &u)
		return u["id"].(string)
	}
	proof := uploadProof()
	call(other, "GET", "/uploads/"+proof, nil, 403)
	call("", "GET", "/uploads/"+proof, nil, 401)
	call(customer, "POST", "/orders/"+id+"/proof", map[string]any{"upload_id": proof}, 200)
	call(owner, "GET", "/uploads/"+proof, nil, 200)
	call(wrongOwner, "GET", "/uploads/"+proof, nil, 403)
	call(customer, "POST", "/orders/"+id+"/cancel", map[string]any{"confirm_unpaid": true}, 409)
	call(owner, "POST", "/orders/"+id+"/review", map[string]any{"action": "approve", "received_funds": false, "note": ""}, 400)
	call(owner, "POST", "/orders/"+id+"/review", map[string]any{"action": "correction", "received_funds": false, "note": "Bukti kurang jelas, mohon unggah ulang."}, 200)
	call(customer, "POST", "/orders/"+id+"/proof", map[string]any{"upload_id": uploadProof()}, 200)
	approved := call(owner, "POST", "/orders/"+id+"/review", map[string]any{"action": "approve", "received_funds": true, "note": ""}, 200)
	again := call(owner, "POST", "/orders/"+id+"/review", map[string]any{"action": "approve", "received_funds": true, "note": ""}, 200)
	tickets := approved["tickets"].([]any)
	if len(tickets) != 2 || len(again["tickets"].([]any)) != 2 {
		t.Fatal("one QR per admission; approval must be idempotent")
	}
	qr := tickets[0].(map[string]any)["qr"].(string)
	call(wrongOwner, "POST", "/check-in", map[string]any{"code": qr, "event_id": eid}, 403)
	call(owner, "POST", "/check-in", map[string]any{"code": qr, "event_id": eid + "wrong"}, 400)
	call(owner, "POST", "/check-in", map[string]any{"code": qr, "event_id": eid}, 200)
	call(owner, "POST", "/check-in", map[string]any{"code": qr, "event_id": eid}, 409)
	// Multiple gate requests for one QR must record exactly one admission.
	secondQR := tickets[1].(map[string]any)["qr"].(string)
	var gateWG sync.WaitGroup
	var gateMu sync.Mutex
	gateOK, gateConflict := 0, 0
	for i := 0; i < 8; i++ {
		gateWG.Add(1)
		go func() {
			defer gateWG.Done()
			b, _ := json.Marshal(map[string]any{"code": secondQR, "event_id": eid})
			r := httptest.NewRequest("POST", "/api/check-in", bytes.NewReader(b))
			r.AddCookie(&http.Cookie{Name: "ruang_session", Value: tokens[owner]})
			w := httptest.NewRecorder()
			h.ServeHTTP(w, r)
			gateMu.Lock()
			defer gateMu.Unlock()
			if w.Code == 200 {
				gateOK++
			} else if w.Code == 409 {
				gateConflict++
			}
		}()
	}
	gateWG.Wait()
	if gateOK != 1 || gateConflict != 7 {
		t.Fatalf("QR race: admitted=%d rejected=%d", gateOK, gateConflict)
	}
	cancelled := create("cancel", 1, 200)["id"].(string)
	call(customer, "POST", "/orders/"+cancelled+"/cancel", map[string]any{"confirm_unpaid": true}, 200)
	expired := create("expire", 1, 200)["id"].(string)
	if _, e = a.DB.Exec(ctx, `UPDATE orders SET expires_at=now()-interval '1 minute' WHERE id=$1`, expired); e != nil {
		t.Fatal(e)
	}
	if call(customer, "GET", "/orders/"+expired, nil, 200)["status"] != "expired" {
		t.Fatal("unpaid order did not expire")
	}
	// A payment deadline reminder is emitted once, even if the scheduler runs twice.
	deadline := create("deadline", 1, 200)["id"].(string)
	if _, e = a.DB.Exec(ctx, `UPDATE orders SET expires_at=now()+interval '3 minutes' WHERE id=$1`, deadline); e != nil {
		t.Fatal(e)
	}
	for i := 0; i < 2; i++ {
		if e = a.Remind(ctx); e != nil {
			t.Fatal(e)
		}
	}
	var deadlineCount int
	if e = a.DB.QueryRow(ctx, `SELECT count(*) FROM notifications WHERE account_id=$1 AND id='payment-reminder-'||$2`, customer, deadline).Scan(&deadlineCount); e != nil {
		t.Fatal(e)
	}
	if deadlineCount != 1 {
		t.Fatalf("payment deadline repeated: %d", deadlineCount)
	}
	call(customer, "POST", "/orders/"+deadline+"/cancel", map[string]any{"confirm_unpaid": true}, 200)
	// Click metrics stay with the owner of a published event; a bad target creates no record.
	call("", "POST", "/activity/click", map[string]any{"kind": "ticket", "target_id": eid}, 201)
	call("", "POST", "/activity/click", map[string]any{"kind": "ticket", "target_id": "missing"}, 404)
	stats := call(owner, "GET", "/dashboard", nil, 200)
	if stats["ticket_clicks"] != float64(1) || len(stats["activity"].([]any)) != 7 {
		t.Fatal("dashboard activity not scoped or incomplete")
	}
	if call(wrongOwner, "GET", "/dashboard", nil, 200)["ticket_clicks"] != float64(0) {
		t.Fatal("dashboard leaks another owner's metrics")
	}
	// Concurrent requests for the four remaining seats: exactly four succeed.
	var wg sync.WaitGroup
	var mu sync.Mutex
	ok := 0
	bad := []string{}
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			b, _ := json.Marshal(map[string]any{"session_id": sid, "quantity": 1, "idempotency_key": fmt.Sprintf("%sconcurrent%d", prefix, i)})
			r := httptest.NewRequest("POST", "/api/orders", bytes.NewReader(b))
			r.AddCookie(&http.Cookie{Name: "ruang_session", Value: tokens[other]})
			w := httptest.NewRecorder()
			h.ServeHTTP(w, r)
			mu.Lock()
			defer mu.Unlock()
			if w.Code == 200 {
				ok++
			} else if w.Code != 409 {
				bad = append(bad, w.Body.String())
			}
		}(i)
	}
	wg.Wait()
	if ok != 4 || len(bad) > 0 {
		t.Fatalf("capacity race: succeeded=%d errors=%v", ok, bad)
	}
	if e = a.Remind(ctx); e != nil {
		t.Fatal(e)
	}
	if e = a.Remind(ctx); e != nil {
		t.Fatal(e)
	}
	var count int
	_ = a.DB.QueryRow(ctx, `SELECT count(*) FROM notifications WHERE account_id=$1 AND kind='reminder'`, customer).Scan(&count)
	if count != 1 {
		t.Fatalf("reminder repeated: %d", count)
	}
	// Customer mentions appear in organizer's Mention feed; untrusted HTML remains plain text.
	post := call(customer, "POST", "/posts", map[string]any{"body": "Suka latihan di @Integration <script>test</script>", "city": "Karawang", "mentions": []any{map[string]any{"kind": "event", "id": eid}}}, 201)
	call(other, "DELETE", "/posts/"+post["id"].(string), nil, 403)
	// Category discovery includes audience posts mentioning either an event or its organizer.
	if _, e = a.DB.Exec(ctx, `UPDATE organizers SET category='Teater' WHERE id=$1`, oid); e != nil {
		t.Fatal(e)
	}
	mentioned := call(customer, "POST", "/posts", map[string]any{"body": "Obrolan dengan pengelola teater", "city": "Karawang", "mentions": []any{map[string]any{"kind": "organizer", "id": oid}}}, 201)
	feedContains := func(city, category, id string) bool {
		t.Helper()
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest("GET", "/api/posts?city="+city+"&category="+category, nil))
		if w.Code != 200 {
			t.Fatalf("category feed returned %d: %s", w.Code, w.Body.String())
		}
		var items []map[string]any
		if e := json.Unmarshal(w.Body.Bytes(), &items); e != nil {
			t.Fatal(e)
		}
		for _, item := range items {
			if item["id"] == id {
				return true
			}
		}
		return false
	}
	for _, id := range []string{post["id"].(string), mentioned["id"].(string)} {
		if !feedContains("Karawang", "Teater", id) || feedContains("Karawang", "Musik", id) || feedContains("Jakarta", "Teater", id) {
			t.Fatal("category feed must resolve mentions and remain scoped to the selected city")
		}
	}
	call(customer, "POST", "/reviews", map[string]any{"event_id": eid, "body": "Pengalaman menonton yang menyenangkan.", "attended": true}, 400)
	if _, e = a.DB.Exec(ctx, `UPDATE event_sessions SET starts_at=now()-interval '2 hours',ends_at=now()-interval '1 hour' WHERE id=$1`, sid); e != nil {
		t.Fatal(e)
	}
	var reviewJSON []byte
	review := call(customer, "POST", "/reviews", map[string]any{"event_id": eid, "body": "Pengalaman menonton yang menyenangkan.", "attended": true}, 200)
	_ = a.DB.QueryRow(ctx, `SELECT doc FROM review_view WHERE event_id=$1 AND account_id=$2`, eid, customer).Scan(&reviewJSON)
	_ = json.Unmarshal(reviewJSON, &review)
	if review["verified"] != true {
		t.Fatal("checked ticket not linked to review")
	}
	unverified := call(other, "POST", "/reviews", map[string]any{"event_id": eid, "body": "Hadir tetapi belum memiliki check-in tiket.", "attended": true}, 200)
	_ = a.DB.QueryRow(ctx, `SELECT doc FROM review_view WHERE event_id=$1 AND account_id=$2`, eid, other).Scan(&reviewJSON)
	_ = json.Unmarshal(reviewJSON, &unverified)
	if unverified["verified"] != false {
		t.Fatal("self-reported attendance must not receive verified badge")
	}
	t.Log("payment resume/snapshot/privacy/correction/approval/QR/capacity/cancel/expiry/reminder/review policies passed; prefix=" + strings.TrimSpace(prefix))
}
