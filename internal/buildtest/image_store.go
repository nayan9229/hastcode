package buildtest

import (
	"encoding/xml"
	"io"
	"net/http"
	"net/http/httptest"
	"slices"
	"strconv"
	"strings"
	"sync"
	"testing"

	"hastcode.com/site/builder"
	"hastcode.com/site/internal/imagestore"
)

// ImageStore is a fake Image Store: an httptest server speaking the four S3
// operations the builder and the prune command use (HeadObject, PutObject,
// ListObjectsV2, DeleteObjects) for one bucket, path-style
// (/<bucket>/<key>). It holds objects in memory and records every request.
// A request with the wrong bucket gets 404 NoSuchBucket; one not signed
// with the fake's access key gets 403.
//
//	store := buildtest.NewImageStore(t)
//	out := buildtest.MustBuild(t, store.Configure(buildtest.Config()), tree)
//	puts := store.Requests(http.MethodPut)
type ImageStore struct {
	// Bucket is the bucket the fake serves; Configure sets it as
	// images.bucket.
	Bucket string
	// BaseURL is the public URL Configure sets as images.base_url, with
	// a trailing slash.
	BaseURL string
	// AccessKeyID is the access key Env hands the builder; every request
	// must be signed with it.
	AccessKeyID string

	srv *httptest.Server

	mu       sync.Mutex
	objects  map[string]StoredObject
	requests []StoreRequest
	failures map[string]int // method + " " + key -> status
}

// StoredObject is an object in the fake store.
type StoredObject struct {
	Data         []byte
	ContentType  string
	CacheControl string
}

// StoreRequest is one request the fake received. Op is the S3 operation:
// "HeadObject", "PutObject", "ListObjectsV2", "DeleteObjects" (or
// "Unknown"). Key is the object key for HeadObject and PutObject, ""
// otherwise. Header is the request's headers; Body is the request body
// for PutObject and DeleteObjects.
type StoreRequest struct {
	Method string
	Op     string
	Key    string
	Query  string
	Header http.Header
	Body   []byte
}

// NewImageStore starts an empty fake Image Store, closed when the test ends.
func NewImageStore(tb testing.TB) *ImageStore {
	tb.Helper()
	s := &ImageStore{
		Bucket:      "hastcode-images-test",
		BaseURL:     "https://images.hastcode.test/",
		AccessKeyID: "test-access-key",
		objects:     map[string]StoredObject{},
		failures:    map[string]int{},
	}
	s.srv = httptest.NewServer(http.HandlerFunc(s.serve))
	tb.Cleanup(s.srv.Close)
	return s
}

// Endpoint is the fake's S3 endpoint URL (Env sets it as R2_ENDPOINT).
func (s *ImageStore) Endpoint() string { return s.srv.URL }

// Env is the environment a CI build would have: the R2 credentials and the
// endpoint override pointing at the fake.
func (s *ImageStore) Env() map[string]string {
	return map[string]string{
		imagestore.EnvAccountID:       "test-account",
		imagestore.EnvAccessKeyID:     s.AccessKeyID,
		imagestore.EnvSecretAccessKey: "test-secret",
		imagestore.EnvEndpoint:        s.srv.URL,
	}
}

// Configure returns cfg switched to the r2 store against this fake:
// images.store r2, base_url BaseURL, bucket Bucket, and Getenv reading Env.
func (s *ImageStore) Configure(cfg builder.Config) builder.Config {
	return s.ConfigureEnv(cfg, s.Env())
}

// ConfigureEnv is Configure with an explicit environment (drop a variable
// to test missing credentials).
func (s *ImageStore) ConfigureEnv(cfg builder.Config, env map[string]string) builder.Config {
	cfg.Images.Store = builder.StoreR2
	cfg.Images.BaseURL = s.BaseURL
	cfg.Images.Bucket = s.Bucket
	cfg.Getenv = func(k string) string { return env[k] }
	return cfg
}

// Seed puts an object in the store without recording a request.
func (s *ImageStore) Seed(key string, obj StoredObject) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.objects[key] = obj
}

// Object returns the stored object at key.
func (s *ImageStore) Object(key string) (StoredObject, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	o, ok := s.objects[key]
	return o, ok
}

// Keys are the stored object keys, sorted.
func (s *ImageStore) Keys() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	keys := make([]string, 0, len(s.objects))
	for k := range s.objects {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	return keys
}

// Requests are the requests received so far with the given HTTP method
// ("" for all), in arrival order.
func (s *ImageStore) Requests(method string) []StoreRequest {
	s.mu.Lock()
	defer s.mu.Unlock()
	var out []StoreRequest
	for _, r := range s.requests {
		if method == "" || r.Method == method {
			out = append(out, r)
		}
	}
	return out
}

// RequestKeys are the sorted keys of the requests with the given method.
func (s *ImageStore) RequestKeys(method string) []string {
	var keys []string
	for _, r := range s.Requests(method) {
		keys = append(keys, r.Key)
	}
	slices.Sort(keys)
	return keys
}

// ResetRequests forgets the recorded requests (the objects stay).
func (s *ImageStore) ResetRequests() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.requests = nil
}

// Fail makes every request with method for key answer status with an S3
// XML error body (and store nothing).
func (s *ImageStore) Fail(method, key string, status int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.failures[method+" "+key] = status
}

func (s *ImageStore) serve(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	bucket, key, _ := strings.Cut(strings.TrimPrefix(r.URL.Path, "/"), "/")
	req := StoreRequest{Method: r.Method, Op: op(r, key), Key: key, Query: r.URL.RawQuery, Header: r.Header.Clone()}
	if req.Op == "PutObject" || req.Op == "DeleteObjects" {
		req.Body = body
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.requests = append(s.requests, req)

	if !strings.Contains(r.Header.Get("Authorization"), "Credential="+s.AccessKeyID+"/") {
		s3Error(w, r, http.StatusForbidden, "AccessDenied")
		return
	}
	if bucket != s.Bucket {
		s3Error(w, r, http.StatusNotFound, "NoSuchBucket")
		return
	}
	if status, ok := s.failures[r.Method+" "+key]; ok {
		s3Error(w, r, status, http.StatusText(status))
		return
	}

	switch req.Op {
	case "HeadObject":
		o, ok := s.objects[key]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", o.ContentType)
		w.Header().Set("Cache-Control", o.CacheControl)
		w.Header().Set("Content-Length", strconv.Itoa(len(o.Data)))
		w.Header().Set("ETag", `"fake"`)
	case "PutObject":
		s.objects[key] = StoredObject{Data: body, ContentType: r.Header.Get("Content-Type"), CacheControl: r.Header.Get("Cache-Control")}
		w.Header().Set("ETag", `"fake"`)
	case "ListObjectsV2":
		s.list(w, r)
	case "DeleteObjects":
		s.delete(w, body)
	default:
		s3Error(w, r, http.StatusNotImplemented, "NotImplemented")
	}
}

func op(r *http.Request, key string) string {
	q := r.URL.Query()
	switch {
	case r.Method == http.MethodHead && key != "":
		return "HeadObject"
	case r.Method == http.MethodPut && key != "":
		return "PutObject"
	case r.Method == http.MethodGet && key == "" && q.Get("list-type") == "2":
		return "ListObjectsV2"
	case r.Method == http.MethodPost && key == "" && q.Has("delete"):
		return "DeleteObjects"
	}
	return "Unknown"
}

// list answers ListObjectsV2: keys under prefix in lexical order, at most
// max-keys (default 1000) per page, continuing after continuation-token
// (the last key of the previous page).
func (s *ImageStore) list(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	prefix, after := q.Get("prefix"), q.Get("continuation-token")
	maxKeys := 1000
	if n, err := strconv.Atoi(q.Get("max-keys")); err == nil && n > 0 {
		maxKeys = n
	}
	var keys []string
	for k := range s.objects {
		if strings.HasPrefix(k, prefix) && k > after {
			keys = append(keys, k)
		}
	}
	slices.Sort(keys)
	type content struct {
		Key  string `xml:"Key"`
		Size int    `xml:"Size"`
	}
	res := struct {
		XMLName               xml.Name  `xml:"ListBucketResult"`
		Name                  string    `xml:"Name"`
		Prefix                string    `xml:"Prefix"`
		KeyCount              int       `xml:"KeyCount"`
		MaxKeys               int       `xml:"MaxKeys"`
		IsTruncated           bool      `xml:"IsTruncated"`
		Contents              []content `xml:"Contents"`
		ContinuationToken     string    `xml:"ContinuationToken,omitempty"`
		NextContinuationToken string    `xml:"NextContinuationToken,omitempty"`
	}{Name: s.Bucket, Prefix: prefix, MaxKeys: maxKeys, ContinuationToken: after}
	if len(keys) > maxKeys {
		keys = keys[:maxKeys]
		res.IsTruncated = true
		res.NextContinuationToken = keys[len(keys)-1]
	}
	for _, k := range keys {
		res.Contents = append(res.Contents, content{Key: k, Size: len(s.objects[k].Data)})
	}
	res.KeyCount = len(res.Contents)
	writeXML(w, res)
}

// delete answers DeleteObjects: removes each listed key (missing keys are
// reported deleted, as in S3).
func (s *ImageStore) delete(w http.ResponseWriter, body []byte) {
	var in struct {
		Objects []struct {
			Key string `xml:"Key"`
		} `xml:"Object"`
	}
	if err := xml.Unmarshal(body, &in); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	type deleted struct {
		Key string `xml:"Key"`
	}
	res := struct {
		XMLName xml.Name  `xml:"DeleteResult"`
		Deleted []deleted `xml:"Deleted"`
	}{}
	for _, o := range in.Objects {
		delete(s.objects, o.Key)
		res.Deleted = append(res.Deleted, deleted{Key: o.Key})
	}
	writeXML(w, res)
}

func writeXML(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/xml")
	_, _ = io.WriteString(w, xml.Header)
	_ = xml.NewEncoder(w).Encode(v)
}

// s3Error writes an S3 error response (no body for HEAD, as S3 does).
func s3Error(w http.ResponseWriter, r *http.Request, status int, code string) {
	if r.Method == http.MethodHead {
		w.WriteHeader(status)
		return
	}
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(status)
	_ = xml.NewEncoder(w).Encode(struct {
		XMLName xml.Name `xml:"Error"`
		Code    string   `xml:"Code"`
		Message string   `xml:"Message"`
	}{Code: code, Message: code})
}
