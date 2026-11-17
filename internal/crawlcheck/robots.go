package crawlcheck

import (
	"strings"
)

// robots is a parsed robots.txt, read the RFC 9309 way: a group is one or
// more user-agent lines followed by its rules; a user-agent line after a
// rule starts the next group. Sitemap lines belong to no group. Other
// records (Cloudflare's Content-Signal) are ignored and neither end nor
// start a group.
type robots struct {
	groups   []group
	sitemaps []string
}

type group struct {
	agents []string // product tokens as written
	rules  []rule
}

type rule struct {
	allow   bool
	pattern string
	line    int // 1-based line in the file
}

func (r rule) String() string {
	if r.allow {
		return "Allow: " + r.pattern
	}
	return "Disallow: " + r.pattern
}

func parseRobots(body string) robots {
	var r robots
	var cur *group
	inRules := false
	for i, line := range strings.Split(body, "\n") {
		if hash := strings.IndexByte(line, '#'); hash >= 0 {
			line = line[:hash]
		}
		key, value, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		key = strings.ToLower(strings.TrimSpace(key))
		value = strings.TrimSpace(value)
		switch key {
		case "user-agent":
			if cur == nil || inRules {
				r.groups = append(r.groups, group{})
				cur = &r.groups[len(r.groups)-1]
				inRules = false
			}
			cur.agents = append(cur.agents, value)
		case "allow", "disallow":
			if cur == nil {
				continue // a rule before any user-agent line belongs to no group
			}
			inRules = true
			cur.rules = append(cur.rules, rule{allow: key == "allow", pattern: value, line: i + 1})
		case "sitemap":
			r.sitemaps = append(r.sitemaps, value)
		}
	}
	return r
}

// groupsFor returns the groups naming token, case-insensitively and as
// a whole token (Applebot-Extended is not Applebot).
func (r robots) groupsFor(token string) []group {
	var out []group
	for _, g := range r.groups {
		for _, a := range g.agents {
			if strings.EqualFold(a, token) {
				out = append(out, g)
				break
			}
		}
	}
	return out
}

// rulesFor returns the combined rules c follows: those of every group
// naming its first named token, else of every * group (RFC 9309 §2.2.1).
func (r robots) rulesFor(c Crawler) (rules []rule, followed string) {
	followed = "*"
	gs := r.groupsFor("*")
	for _, t := range c.Tokens {
		if named := r.groupsFor(t); len(named) > 0 {
			gs, followed = named, t
			break
		}
	}
	for _, g := range gs {
		rules = append(rules, g.rules...)
	}
	return rules, followed
}

// decide returns the rule that decides path for these rules: the longest
// matching pattern, Allow winning a tie (RFC 9309 §2.2.2). ok is false
// when no rule matches, which means allowed.
func decide(rules []rule, path string) (decider rule, ok bool) {
	for _, r := range rules {
		if r.pattern == "" || !matches(r.pattern, path) {
			continue
		}
		if !ok || len(r.pattern) > len(decider.pattern) ||
			(len(r.pattern) == len(decider.pattern) && r.allow && !decider.allow) {
			decider, ok = r, true
		}
	}
	return decider, ok
}

// matches reports whether pattern matches path from its start, with *
// matching any run of characters and a trailing $ anchoring the end.
func matches(pattern, path string) bool {
	anchored := strings.HasSuffix(pattern, "$")
	parts := strings.Split(strings.TrimSuffix(pattern, "$"), "*")
	if !strings.HasPrefix(path, parts[0]) {
		return false
	}
	rest := path[len(parts[0]):]
	if len(parts) == 1 {
		return !anchored || rest == ""
	}
	for i, part := range parts[1:] {
		if anchored && i == len(parts)-2 {
			return strings.HasSuffix(rest, part)
		}
		at := strings.Index(rest, part)
		if at < 0 {
			return false
		}
		rest = rest[at+len(part):]
	}
	return true
}
