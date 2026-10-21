package vercel

import (
	"encoding/json"
	"fmt"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

// NewRouter returns an http.Handler that serves a Build Output directory
// locally the way Vercel's router does, for the subset of config.json the
// builder writes: routes are read from <outDir>/config.json, never
// duplicated in Go. functions maps a site path (e.g. "/api/subscribe") to
// the handler Vercel would run there; functions and static/ files form the
// filesystem phase.
//
// Routing, per https://vercel.com/docs/build-output-api/configuration#routes:
// routes before {"handle":"filesystem"} run first, in order. A matching
// route with "continue" adds its headers and routing goes on; any other
// matching route ends the phase: with a Location header it answers with
// its status (a redirect), otherwise its dest (or the request path) is
// looked up in the filesystem and served with the route's status. On a
// filesystem miss, the routes after {"handle":"filesystem"} run the same
// way. src is a case-insensitive regular expression unless caseSensitive.
//
// It is a development aid, not a Vercel emulator: a config using anything
// else (other handle phases, has/missing, a regex Go cannot compile) is an
// error so the local check never silently disagrees with production.
func NewRouter(outDir string, functions map[string]http.Handler) (http.Handler, error) {
	data, err := os.ReadFile(filepath.Join(outDir, "config.json"))
	if err != nil {
		return nil, err
	}
	var cfg struct {
		Config
		Routes []json.RawMessage `json:"routes"`
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("config.json: %w", err)
	}
	rt := &router{static: filepath.Join(outDir, "static"), functions: functions}
	phase := &rt.main
	for i, raw := range cfg.Routes {
		var r Route
		dec := json.NewDecoder(strings.NewReader(string(raw)))
		dec.DisallowUnknownFields()
		if err := dec.Decode(&r); err != nil {
			return nil, fmt.Errorf("config.json routes[%d]: unsupported by the local router: %w", i, err)
		}
		if r.Handle != "" {
			if r.Handle != HandleFilesystem || phase != &rt.main || r.Src != "" {
				return nil, fmt.Errorf("config.json routes[%d]: handle %q unsupported by the local router", i, r.Handle)
			}
			phase = &rt.miss
			continue
		}
		cr, err := compileRoute(r)
		if err != nil {
			return nil, fmt.Errorf("config.json routes[%d]: %w", i, err)
		}
		*phase = append(*phase, cr)
	}
	return rt, nil
}

type compiledRoute struct {
	Route
	re *regexp.Regexp
}

func compileRoute(r Route) (compiledRoute, error) {
	if r.Continue && (r.Dest != "" || r.Status != 0) {
		return compiledRoute{}, fmt.Errorf("continue with dest or status unsupported by the local router")
	}
	src := r.Src
	if !strings.HasPrefix(src, "^") {
		src = "^" + src
	}
	if !strings.HasSuffix(src, "$") {
		src += "$"
	}
	if !r.CaseSensitive {
		src = "(?i)" + src
	}
	re, err := regexp.Compile(src)
	if err != nil {
		return compiledRoute{}, fmt.Errorf("src %q: %w", r.Src, err)
	}
	return compiledRoute{Route: r, re: re}, nil
}

type router struct {
	static    string
	functions map[string]http.Handler
	main      []compiledRoute
	miss      []compiledRoute
}

func (rt *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	p := req.URL.Path
	if done := rt.phase(w, req, rt.main, &p); done {
		return
	}
	if rt.filesystem(w, req, p, 0) {
		return
	}
	if done := rt.phase(w, req, rt.miss, &p); done {
		return
	}
	http.NotFound(w, req)
}

// phase applies routes to *p. It returns true once a response is written;
// otherwise *p is the path the next phase looks up.
func (rt *router) phase(w http.ResponseWriter, req *http.Request, routes []compiledRoute, p *string) bool {
	for _, r := range routes {
		m := r.re.FindStringSubmatchIndex(*p)
		if m == nil {
			continue
		}
		for k, v := range r.Headers {
			w.Header().Set(k, string(r.re.ExpandString(nil, v, *p, m)))
		}
		if r.Continue {
			continue
		}
		if _, redirect := r.Headers["Location"]; redirect && r.Status != 0 {
			w.WriteHeader(r.Status)
			return true
		}
		if r.Dest != "" {
			*p = string(r.re.ExpandString(nil, r.Dest, *p, m))
		}
		if r.Status != 0 || r.Dest != "" {
			if rt.filesystem(w, req, *p, r.Status) {
				return true
			}
		}
		return false
	}
	return false
}

// filesystem serves a function or a static file at p with status (0: the
// handler's own / 200). It reports whether p existed.
func (rt *router) filesystem(w http.ResponseWriter, req *http.Request, p string, status int) bool {
	if h, ok := rt.functions[p]; ok {
		h.ServeHTTP(w, req)
		return true
	}
	if strings.Contains(p, "..") || !strings.HasPrefix(p, "/") {
		return false
	}
	file := filepath.Join(rt.static, filepath.FromSlash(p[1:]))
	if strings.HasSuffix(p, "/") {
		file = filepath.Join(file, "index.html")
	}
	info, err := os.Stat(file)
	if err != nil || info.IsDir() {
		return false
	}
	data, err := os.ReadFile(file)
	if err != nil {
		return false
	}
	ctype := mime.TypeByExtension(path.Ext(file))
	if ctype == "" {
		ctype = http.DetectContentType(data)
	}
	w.Header().Set("Content-Type", ctype)
	if status == 0 {
		status = http.StatusOK
	}
	w.WriteHeader(status)
	if req.Method != http.MethodHead {
		_, _ = w.Write(data)
	}
	return true
}
