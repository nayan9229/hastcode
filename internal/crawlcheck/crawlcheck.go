// Package crawlcheck is `make crawl-check`: a read-only check, run by the
// Author against live Production, that search and answer crawlers can
// still reach the site. Two things outside the build can shut one out
// unnoticed: Cloudflare's managed robots.txt block, prepended to our file,
// and an edge rule answering a crawler's user agent with 403
// (docs/research/indexing-submission-channels.md §0).
//
// It only sends GET requests, never follows a redirect, and never talks to
// anything but the three addresses it checks: no IndexNow, no WebSub hub,
// no search engine API. It sends crawler user-agent strings from the
// Author's own IP address, so it shows what Cloudflare does with a user
// agent, not what a verified crawler sees from its own addresses. It runs
// locally only, never in GitHub Actions (CI's addresses would take a
// different path through Cloudflare's edge).
package crawlcheck

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"slices"
	"strings"
	"time"
)

// DefaultTimeout bounds each request when Checker.Timeout is zero.
const DefaultTimeout = 10 * time.Second

// UserAgent is sent with every request that is not made as a crawler
// (robots.txt, the redirect checks), so the site's logs show who asked.
const UserAgent = "HastCode-CrawlCheck/1 (+https://hastcode.com/contact/)"

// maxRobotsBytes bounds the robots.txt body read (Google reads 500 KiB).
const maxRobotsBytes = 512 << 10

// The note's runbook steps the failures point to.
const (
	note       = "docs/research/indexing-submission-channels.md"
	fixBlocked = "runbook step 1 in " + note
	fixHTTPS   = "turn on Always Use HTTPS (runbook step 2 in " + note + ")"
)

// allowGroups are the explicit allow groups the build writes
// (builder/step_write_robots.go).
var allowGroups = []string{"Googlebot", "Bingbot", "OAI-SearchBot", "PerplexityBot"}

// sitemapPaths are the two Sitemap: lines the build writes.
var sitemapPaths = []string{"/sitemap.xml", "/news-sitemap.xml"}

// accessPaths are what each crawler asks for.
var accessPaths = []string{"/", "/sitemap.xml"}

// Checker checks one site. New derives the three addresses from the
// site_url; each can be overridden (the tests point all three at local
// stand-ins).
type Checker struct {
	// Site is the Production origin, e.g. https://hastcode.com.
	Site string
	// HTTP is the plain-http origin, which must redirect to Site + "/".
	HTTP string
	// WWW is the www origin, which must redirect to Site.
	WWW string
	// Crawlers defaults to Crawlers.
	Crawlers []Crawler
	// Client sends every request; nil is a plain client. Its redirect
	// policy is replaced: the check never follows one.
	Client *http.Client
	// Timeout bounds each request; zero is DefaultTimeout.
	Timeout time.Duration
	// Out receives the report, one line per check and the summary;
	// nil is os.Stdout.
	Out io.Writer
}

// New returns a Checker for the site at siteURL (site.yaml site_url): the
// origin itself, its http:// twin and its www. twin.
func New(siteURL string) (Checker, error) {
	u, err := url.Parse(strings.TrimSpace(siteURL))
	if err != nil || u.Scheme != "https" || u.Host == "" {
		return Checker{}, fmt.Errorf("site_url %q is not an https:// address", siteURL)
	}
	return Checker{
		Site: "https://" + u.Host,
		HTTP: "http://" + u.Host,
		WWW:  "https://www." + u.Host,
	}, nil
}

// Result is one check's outcome.
type Result struct {
	OK bool
	// Line is the plain-words line printed after PASS or FAIL.
	Line string
}

// Report is every check's outcome, in the order run.
type Report struct {
	Results []Result
}

// Failed is the number of failed checks.
func (r Report) Failed() int {
	n := 0
	for _, res := range r.Results {
		if !res.OK {
			n++
		}
	}
	return n
}

// ExitCode is the command's exit status: 0 when every check passed, 1 on
// any failure (a network error or timeout included).
func (r Report) ExitCode() int {
	if r.Failed() > 0 {
		return 1
	}
	return 0
}

// Run runs every check, printing each line as it goes and the summary
// last.
func (c Checker) Run(ctx context.Context) Report {
	out := c.Out
	if out == nil {
		out = os.Stdout
	}
	site := strings.TrimSuffix(c.Site, "/")
	fmt.Fprintf(out, "crawl-check: read-only GETs to %s, %s and %s.\n", site, c.HTTP, c.WWW)
	fmt.Fprintln(out, "It sends crawler user agents from this machine's IP address, so it shows what Cloudflare does with each user agent, not what a verified crawler sees from its own addresses.")
	var rep Report
	add := func(ok bool, format string, a ...any) {
		res := Result{OK: ok, Line: fmt.Sprintf(format, a...)}
		rep.Results = append(rep.Results, res)
		word := "PASS"
		if !ok {
			word = "FAIL"
		}
		fmt.Fprintf(out, "%s %s\n", word, res.Line)
	}
	crawlers := c.Crawlers
	if crawlers == nil {
		crawlers = Crawlers
	}

	c.checkRobots(ctx, site, crawlers, add)
	for _, cr := range crawlers {
		c.checkAccess(ctx, site, cr, add)
	}
	c.checkRedirect(ctx, c.HTTP+"/", func(loc string) bool { return loc == site+"/" }, site+"/", fixHTTPS, add)
	c.checkRedirect(ctx, c.WWW+"/", func(loc string) bool { return sameOrigin(loc, site) }, site+"/",
		"restore the www-to-apex redirect rule in Cloudflare", add)

	if n := rep.Failed(); n > 0 {
		fmt.Fprintf(out, "crawl-check: %d of %d checks failed.\n", n, len(rep.Results))
	} else {
		fmt.Fprintf(out, "crawl-check: all %d checks passed.\n", len(rep.Results))
	}
	return rep
}

type addFunc func(ok bool, format string, a ...any)

func (c Checker) checkRobots(ctx context.Context, site string, crawlers []Crawler, add addFunc) {
	u := site + "/robots.txt"
	resp, err := c.get(ctx, u, "")
	if err != nil {
		add(false, "robots.txt: %s", err)
		return
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, maxRobotsBytes))
	resp.Body.Close()
	if err != nil {
		add(false, "robots.txt: reading %s: %v", u, err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		add(false, "robots.txt: %s answered %d, not 200", u, resp.StatusCode)
		return
	}
	r := parseRobots(string(body))
	for _, token := range allowGroups {
		if r.hasAllowAll(token) {
			add(true, "robots.txt has our %s allow group", token)
		} else {
			add(false, "robots.txt has lost our %s allow group (User-agent: %s / Allow: /); the build writes it (builder/step_write_robots.go), so something in front of the site changed the file", token, token)
		}
	}
	for _, p := range sitemapPaths {
		want := site + p
		if slices.Contains(r.sitemaps, want) {
			add(true, "robots.txt has Sitemap: %s", want)
		} else {
			add(false, "robots.txt has no Sitemap: %s line; the build writes it (builder/step_write_robots.go)", want)
		}
	}
	for _, cr := range crawlers {
		rules, followed := r.rulesFor(cr)
		group := "the * groups"
		if followed != "*" {
			group = "the " + followed + " group"
		}
		if d, ok := decide(rules, "/"); ok && !d.allow {
			add(false, "robots.txt stops %s crawling / (%s, line %d, in %s). If that line is in Cloudflare's managed block, change it in Cloudflare (%s); if it is ours, see builder/step_write_robots.go", cr.Name, d, d.line, group, fixBlocked)
			continue
		}
		add(true, "robots.txt lets %s crawl / (it follows %s)", cr.Name, group)
	}
}

// hasAllowAll reports whether a group naming token allows /.
func (r robots) hasAllowAll(token string) bool {
	for _, g := range r.groupsFor(token) {
		for _, ru := range g.rules {
			if ru.allow && ru.pattern == "/" {
				return true
			}
		}
	}
	return false
}

func (c Checker) checkAccess(ctx context.Context, site string, cr Crawler, add addFunc) {
	var blocked, other, failed []string
	for _, p := range accessPaths {
		resp, err := c.get(ctx, site+p, cr.UserAgent)
		if err != nil {
			failed = append(failed, err.Error())
			continue
		}
		_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 1<<20))
		resp.Body.Close()
		switch resp.StatusCode {
		case http.StatusOK:
		case http.StatusForbidden:
			blocked = append(blocked, p)
		default:
			other = append(other, fmt.Sprintf("%s answered %d", p, resp.StatusCode))
		}
	}
	paths := strings.Join(accessPaths, " and ")
	switch {
	case len(blocked) > 0:
		add(false, "Cloudflare is blocking %s: 403 on %s with its user agent. Fix: %s", cr.Name, strings.Join(blocked, " and "), fixBlocked)
	case len(other) > 0:
		add(false, "%s does not get 200: %s", cr.Name, strings.Join(other, ", "))
	case len(failed) > 0:
		add(false, "%s: %s", cr.Name, failed[0])
	default:
		add(true, "%s gets 200 on %s", cr.Name, paths)
	}
}

// checkRedirect checks that u answers 301 or 308 with a location good
// accepts; want is the location named in a failure.
func (c Checker) checkRedirect(ctx context.Context, u string, good func(string) bool, want, fix string, add addFunc) {
	resp, err := c.get(ctx, u, "")
	if err != nil {
		add(false, "redirect check: %s", err)
		return
	}
	resp.Body.Close()
	loc := resp.Header.Get("Location")
	switch {
	case resp.StatusCode != http.StatusMovedPermanently && resp.StatusCode != http.StatusPermanentRedirect:
		add(false, "%s answers %d instead of redirecting to %s: %s", u, resp.StatusCode, want, fix)
	case !good(loc):
		add(false, "%s redirects to %q, not %s: %s", u, loc, want, fix)
	default:
		add(true, "%s redirects (%d) to %s", u, resp.StatusCode, loc)
	}
}

// get sends one GET without following redirects; a network error or
// timeout comes back in plain words.
func (c Checker) get(ctx context.Context, u, userAgent string) (*http.Response, error) {
	timeout := c.Timeout
	if timeout == 0 {
		timeout = DefaultTimeout
	}
	client := http.Client{}
	if c.Client != nil {
		client = *c.Client
	}
	client.Timeout = timeout
	client.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, fmt.Errorf("bad address %s: %v", u, err)
	}
	if userAgent == "" {
		userAgent = UserAgent
	}
	req.Header.Set("User-Agent", userAgent)
	resp, err := client.Do(req)
	if err != nil {
		return nil, plainError(u, timeout, err)
	}
	return resp, nil
}

func plainError(u string, timeout time.Duration, err error) error {
	var ne net.Error
	if errors.As(err, &ne) && ne.Timeout() {
		return fmt.Errorf("could not reach %s: no answer within %s", u, timeout)
	}
	var dns *net.DNSError
	if errors.As(err, &dns) {
		return fmt.Errorf("could not reach %s: the name %s does not resolve (offline?)", u, dns.Name)
	}
	var op *net.OpError
	if errors.As(err, &op) {
		return fmt.Errorf("could not reach %s: %v", u, op.Err)
	}
	return fmt.Errorf("could not reach %s: %v", u, err)
}

// sameOrigin reports whether loc is on origin (scheme and host).
func sameOrigin(loc, origin string) bool {
	l, err1 := url.Parse(loc)
	o, err2 := url.Parse(origin)
	return err1 == nil && err2 == nil && l.Scheme == o.Scheme && strings.EqualFold(l.Host, o.Host)
}
