package views

import (
	_ "embed"
	"strings"
)

//go:embed search.js
var searchJS string

// searchScript is the search page's inline lazy loader (page_search.templ).
var searchScript = strings.TrimSuffix(minifyLines(searchJS), "\n")

// SearchPage is the view model for /search/.
type SearchPage struct {
	Page Page
	// Bundle is the site path the Pagefind bundle is served from
	// ("/pagefind/"); the loader fetches the UI and index from there.
	Bundle string
}
