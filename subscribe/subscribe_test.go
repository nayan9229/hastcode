package subscribe_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"

	"hastcode.com/site/subscribe"
)

const testAPIKey = "test-buttondown-key-0123456789"

// fakeButtondown stands in for Buttondown's subscriber API and records
// every request it receives.
type fakeButtondown struct {
	*httptest.Server
	status int    // status to answer with (default 201)
	body   string // response body

	mu       sync.Mutex
	received []receivedSubscriber
}

type receivedSubscriber struct {
	Method        string
	Path          string
	Authorization string
	EmailAddress  string   `json:"email_address"`
	Tags          []string `json:"tags"`
	IPAddress     string   `json:"ip_address"`
}

func newFakeButtondown(t *testing.T) *fakeButtondown {
	t.Helper()
	f := &fakeButtondown{status: http.StatusCreated, body: `{"id":"sub_1"}`}
	f.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var got receivedSubscriber
		_ = json.NewDecoder(r.Body).Decode(&got)
		got.Method, got.Path, got.Authorization = r.Method, r.URL.Path, r.Header.Get("Authorization")
		f.mu.Lock()
		f.received = append(f.received, got)
		status, body := f.status, f.body
		f.mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(f.Close)
	return f
}

func (f *fakeButtondown) respond(status int, body string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.status, f.body = status, body
}

func (f *fakeButtondown) requests() []receivedSubscriber {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]receivedSubscriber(nil), f.received...)
}

func newHandler(f *fakeButtondown) http.Handler {
	return subscribe.New(subscribe.Config{APIKey: testAPIKey, BaseURL: f.URL})
}

// form is one subscribe request as the newsletter form (or its script) sends it.
type form struct {
	email, source, honeypot string
	ip                      string
	acceptJSON              bool
}

func post(h http.Handler, f form) *httptest.ResponseRecorder {
	v := url.Values{}
	v.Set("email", f.email)
	v.Set("source", f.source)
	v.Set("website", f.honeypot)
	req := httptest.NewRequest(http.MethodPost, "/api/subscribe", strings.NewReader(v.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if f.acceptJSON {
		req.Header.Set("Accept", "application/json")
	} else {
		req.Header.Set("Accept", "text/html,application/xhtml+xml,*/*;q=0.8")
	}
	ip := f.ip
	if ip == "" {
		ip = "203.0.113.7"
	}
	req.Header.Set("X-Forwarded-For", ip)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func TestSubscribe_ValidFormPostRedirectsToSubscribedAndForwardsToButtondown(t *testing.T) {
	bd := newFakeButtondown(t)

	rec := post(newHandler(bd), form{email: "reader@example.com", source: "/distributed-systems/raft-leader-election/"})

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("status = %d, want 303; body %q", rec.Code, rec.Body.String())
	}
	if got := rec.Header().Get("Location"); got != "/subscribed/" {
		t.Errorf("Location = %q, want /subscribed/", got)
	}
	reqs := bd.requests()
	if len(reqs) != 1 {
		t.Fatalf("Buttondown received %d requests, want 1", len(reqs))
	}
	got := reqs[0]
	if got.Method != http.MethodPost || got.Path != "/v1/subscribers" {
		t.Errorf("Buttondown request = %s %s, want POST /v1/subscribers", got.Method, got.Path)
	}
	if got.Authorization != "Token "+testAPIKey {
		t.Errorf("Authorization = %q, want Token <key>", got.Authorization)
	}
	if got.EmailAddress != "reader@example.com" {
		t.Errorf("email_address = %q", got.EmailAddress)
	}
	if len(got.Tags) != 1 || got.Tags[0] != "/distributed-systems/raft-leader-election/" {
		t.Errorf("tags = %q, want the source page path", got.Tags)
	}
}

// jsonResult is the body the newsletter script reads.
type jsonResult struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

func decodeJSON(t *testing.T, rec *httptest.ResponseRecorder) jsonResult {
	t.Helper()
	if ct := rec.Header().Get("Content-Type"); !strings.HasPrefix(ct, "application/json") {
		t.Fatalf("Content-Type = %q, want application/json; body %q", ct, rec.Body.String())
	}
	var res jsonResult
	if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
		t.Fatalf("decode JSON body %q: %v", rec.Body.String(), err)
	}
	return res
}

func TestSubscribe_ValidScriptedPostGetsJSONSuccess(t *testing.T) {
	bd := newFakeButtondown(t)

	rec := post(newHandler(bd), form{email: "reader@example.com", source: "/briefs/postgres-18/", acceptJSON: true})

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body %q", rec.Code, rec.Body.String())
	}
	res := decodeJSON(t, rec)
	if !res.OK || res.Message == "" {
		t.Errorf("JSON = %+v, want ok with a message", res)
	}
	reqs := bd.requests()
	if len(reqs) != 1 || reqs[0].EmailAddress != "reader@example.com" || len(reqs[0].Tags) != 1 || reqs[0].Tags[0] != "/briefs/postgres-18/" {
		t.Errorf("Buttondown received %+v, want one subscriber tagged with the source", reqs)
	}
}

func TestSubscribe_InvalidEmailIsRejectedWithAMessage(t *testing.T) {
	for _, email := range []string{"", "not-an-email", "reader@", "@example.com", "Reader <reader@example.com>", "reader@example", "a b@example.com"} {
		t.Run(email, func(t *testing.T) {
			bd := newFakeButtondown(t)
			h := newHandler(bd)

			rec := post(h, form{email: email, source: "/x/", acceptJSON: true})
			if rec.Code != http.StatusBadRequest {
				t.Fatalf("JSON: status = %d, want 400", rec.Code)
			}
			if res := decodeJSON(t, rec); res.OK || !strings.Contains(res.Message, "valid email") {
				t.Errorf("JSON = %+v, want not ok with a 'valid email' message", res)
			}

			rec = post(h, form{email: email, source: "/x/"})
			if rec.Code != http.StatusBadRequest {
				t.Fatalf("form: status = %d, want 400", rec.Code)
			}
			if body := rec.Body.String(); !strings.Contains(body, "valid email") || !strings.Contains(body, `href="/x/"`) {
				t.Errorf("form body = %q, want the message and a link back to the source page", body)
			}

			if n := len(bd.requests()); n != 0 {
				t.Errorf("Buttondown received %d requests, want 0", n)
			}
		})
	}
}

func TestSubscribe_HoneypotHitLooksLikeSuccessButIsNotForwarded(t *testing.T) {
	bd := newFakeButtondown(t)
	h := newHandler(bd)
	bot := form{email: "bot@example.com", source: "/x/", honeypot: "https://spam.example"}

	if rec := post(h, bot); rec.Code != http.StatusSeeOther || rec.Header().Get("Location") != "/subscribed/" {
		t.Errorf("form: status = %d Location = %q, want 303 /subscribed/", rec.Code, rec.Header().Get("Location"))
	}
	bot.acceptJSON = true
	rec := post(h, bot)
	if rec.Code != http.StatusOK || !decodeJSON(t, rec).OK {
		t.Errorf("JSON: status = %d body %q, want 200 ok", rec.Code, rec.Body.String())
	}
	if n := len(bd.requests()); n != 0 {
		t.Errorf("Buttondown received %d requests, want 0", n)
	}
}

func TestSubscribe_RateLimitTripsAfterNRequestsFromOneIP(t *testing.T) {
	bd := newFakeButtondown(t)
	h := subscribe.New(subscribe.Config{APIKey: testAPIKey, BaseURL: bd.URL, RateLimit: 3, RateWindow: time.Hour})
	reader := form{email: "reader@example.com", source: "/x/", ip: "198.51.100.1", acceptJSON: true}

	for i := 1; i <= 3; i++ {
		if rec := post(h, reader); rec.Code != http.StatusOK {
			t.Fatalf("request %d: status = %d, want 200", i, rec.Code)
		}
	}
	rec := post(h, reader)
	if rec.Code != http.StatusTooManyRequests {
		t.Fatalf("request 4: status = %d, want 429", rec.Code)
	}
	if res := decodeJSON(t, rec); res.OK || res.Message == "" {
		t.Errorf("request 4: JSON = %+v, want not ok with a message", res)
	}
	reader.acceptJSON = false
	if rec := post(h, reader); rec.Code != http.StatusTooManyRequests || !strings.Contains(rec.Body.String(), "Too many") {
		t.Errorf("request 5 (form): status = %d body %q, want 429 with a message", rec.Code, rec.Body.String())
	}
	if n := len(bd.requests()); n != 3 {
		t.Errorf("Buttondown received %d requests, want 3", n)
	}

	other := form{email: "other@example.com", source: "/x/", ip: "198.51.100.2", acceptJSON: true}
	if rec := post(h, other); rec.Code != http.StatusOK {
		t.Errorf("another IP: status = %d, want 200", rec.Code)
	}
}

func TestSubscribe_UpstreamFailureIsABadGatewayWithAFriendlyMessage(t *testing.T) {
	bd := newFakeButtondown(t)
	bd.respond(http.StatusServiceUnavailable, `{"detail":"internal trace: key `+testAPIKey+`"}`)
	h := newHandler(bd)

	rec := post(h, form{email: "reader@example.com", source: "/x/", acceptJSON: true})
	if rec.Code != http.StatusBadGateway {
		t.Fatalf("JSON: status = %d, want 502", rec.Code)
	}
	if res := decodeJSON(t, rec); res.OK || !strings.Contains(res.Message, "try again") {
		t.Errorf("JSON = %+v, want not ok with a 'try again' message", res)
	}
	requireNoSecrets(t, rec)

	rec = post(h, form{email: "reader@example.com", source: "/x/"})
	if rec.Code != http.StatusBadGateway || !strings.Contains(rec.Body.String(), "try again") {
		t.Errorf("form: status = %d body %q, want 502 with a 'try again' message", rec.Code, rec.Body.String())
	}
	requireNoSecrets(t, rec)
}

func TestSubscribe_UnreachableUpstreamIsABadGateway(t *testing.T) {
	bd := newFakeButtondown(t)
	h := newHandler(bd)
	bd.Close()

	rec := post(h, form{email: "reader@example.com", source: "/x/", acceptJSON: true})
	if rec.Code != http.StatusBadGateway {
		t.Fatalf("status = %d, want 502", rec.Code)
	}
	requireNoSecrets(t, rec)
}

func TestSubscribe_AlreadySubscribedReaderGetsTheSameSuccess(t *testing.T) {
	bd := newFakeButtondown(t)
	bd.respond(http.StatusBadRequest, `{"code":"subscriber_already_exists","detail":"A subscriber with this email already exists.","metadata":{"subscriber_id":"sub_9"}}`)
	h := newHandler(bd)

	rec := post(h, form{email: "reader@example.com", source: "/x/", acceptJSON: true})
	if rec.Code != http.StatusOK || !decodeJSON(t, rec).OK {
		t.Errorf("JSON: status = %d body %q, want 200 ok", rec.Code, rec.Body.String())
	}
	if rec := post(h, form{email: "reader@example.com", source: "/x/"}); rec.Code != http.StatusSeeOther {
		t.Errorf("form: status = %d, want 303", rec.Code)
	}
}

func TestSubscribe_EmailButtondownRejectsIsAnInvalidEmail(t *testing.T) {
	bd := newFakeButtondown(t)
	bd.respond(http.StatusBadRequest, `{"code":"email_invalid","detail":"The email address provided is not valid.","metadata":{"field":"email_address"}}`)

	rec := post(newHandler(bd), form{email: "reader@mailinator.example", source: "/x/", acceptJSON: true})
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
	if res := decodeJSON(t, rec); res.OK || !strings.Contains(res.Message, "valid email") {
		t.Errorf("JSON = %+v, want not ok with a 'valid email' message", res)
	}
}

func TestSubscribe_OnlyAcceptsPOST(t *testing.T) {
	bd := newFakeButtondown(t)
	rec := httptest.NewRecorder()
	newHandler(bd).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/subscribe?email=reader@example.com", nil))

	if rec.Code != http.StatusMethodNotAllowed || rec.Header().Get("Allow") != http.MethodPost {
		t.Errorf("GET: status = %d Allow = %q, want 405 with Allow: POST", rec.Code, rec.Header().Get("Allow"))
	}
	if n := len(bd.requests()); n != 0 {
		t.Errorf("Buttondown received %d requests, want 0", n)
	}
}

func TestSubscribe_MissingAPIKeyFailsWithoutCallingButtondown(t *testing.T) {
	bd := newFakeButtondown(t)
	h := subscribe.New(subscribe.Config{BaseURL: bd.URL})

	if rec := post(h, form{email: "reader@example.com", source: "/x/", acceptJSON: true}); rec.Code != http.StatusBadGateway {
		t.Errorf("status = %d, want 502", rec.Code)
	}
	if n := len(bd.requests()); n != 0 {
		t.Errorf("Buttondown received %d requests, want 0", n)
	}
}

func TestSubscribe_ConfigFromEnvironment(t *testing.T) {
	bd := newFakeButtondown(t)
	t.Setenv("BUTTONDOWN_API_KEY", "env-key-42")
	t.Setenv("BUTTONDOWN_BASE_URL", bd.URL)

	if rec := post(subscribe.New(subscribe.ConfigFromEnv()), form{email: "reader@example.com", source: "/x/"}); rec.Code != http.StatusSeeOther {
		t.Fatalf("status = %d, want 303", rec.Code)
	}
	if reqs := bd.requests(); len(reqs) != 1 || reqs[0].Authorization != "Token env-key-42" {
		t.Errorf("Buttondown received %+v, want one request with the env key", reqs)
	}
}

// requireNoSecrets fails if the response leaks the API key or upstream detail.
func requireNoSecrets(t *testing.T, rec *httptest.ResponseRecorder) {
	t.Helper()
	all := rec.Body.String()
	for k, vs := range rec.Header() {
		all += k + strings.Join(vs, "")
	}
	if strings.Contains(all, testAPIKey) || strings.Contains(all, "internal trace") {
		t.Errorf("response leaks the API key or upstream body: %q", all)
	}
}
