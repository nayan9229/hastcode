// Command crawlcheck is `make crawl-check`: a read-only check that search
// and answer crawlers can still reach live Production (robots.txt, each
// crawler's user agent, the http and www redirects). All logic is in
// internal/crawlcheck; main is flags and exit codes:
//
//	go run ./cmd/crawlcheck [-config site.yaml] [-site URL] [-http URL] [-www URL] [-timeout 10s]
//
// The site comes from -config's site_url; -site, -http and -www override
// the three addresses derived from it. Local only, never in GitHub Actions.
//
//	0  every check passed
//	1  a check failed (a network error or timeout included)
//	2  usage: bad flags, unreadable site config
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"

	"hastcode.com/site/builder"
	"hastcode.com/site/internal/crawlcheck"
)

func main() {
	config := flag.String("config", "site.yaml", "site config; its site_url is the Production site checked")
	site := flag.String("site", "", "override the site address (default: site_url)")
	httpURL := flag.String("http", "", "override the http:// address that must redirect to the site (default: http:// + the site's host)")
	www := flag.String("www", "", "override the www address that must redirect to the site (default: https://www. + the site's host)")
	timeout := flag.Duration("timeout", crawlcheck.DefaultTimeout, "bound on each request")
	flag.Parse()
	if flag.NArg() > 0 {
		fmt.Fprintf(os.Stderr, "crawlcheck: unexpected argument %q\n", flag.Arg(0))
		os.Exit(2)
	}
	cfg, err := builder.LoadConfig(*config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "crawlcheck: %v\n", err)
		os.Exit(2)
	}
	c, err := crawlcheck.New(cfg.SiteURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "crawlcheck: %s: %v\n", *config, err)
		os.Exit(2)
	}
	for _, o := range []struct {
		flag string
		into *string
	}{{*site, &c.Site}, {*httpURL, &c.HTTP}, {*www, &c.WWW}} {
		if o.flag != "" {
			*o.into = o.flag
		}
	}
	c.Timeout = *timeout
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	code := c.Run(ctx).ExitCode()
	stop()
	os.Exit(code)
}
