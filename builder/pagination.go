package builder

import (
	"strconv"

	"hastcode.com/site/internal/views"
)

// Paginated listings (Briefs index, Vertical hubs): page 1 lives at the
// listing's base URL, page n at {base}page/{n}/ (spec "URLs").

// listingPageSize is how many items a paginated listing shows per page.
const listingPageSize = 20

// listingPage is one page of a paginated listing.
type listingPage[T any] struct {
	Items []T
	// URL is the page's site-relative URL.
	URL string
	// Pagination is the page's position and its prev/next links, for
	// views.PaginationNav.
	Pagination views.Pagination
}

// paginate splits items, already in display order, into pages of size
// items under base (a site-relative URL with a trailing slash, e.g.
// "/briefs/"). There is always at least one page, so a listing with no
// items still has its base URL; no other page is ever empty.
func paginate[T any](items []T, size int, base string) []listingPage[T] {
	total := max(1, (len(items)+size-1)/size)
	pages := make([]listingPage[T], total)
	for i := range pages {
		n := i + 1
		p := views.Pagination{Current: n, Total: total}
		if n > 1 {
			p.Prev = pageURL(base, n-1)
		}
		if n < total {
			p.Next = pageURL(base, n+1)
		}
		pages[i] = listingPage[T]{
			Items:      items[i*size : min(len(items), n*size)],
			URL:        pageURL(base, n),
			Pagination: p,
		}
	}
	return pages
}

// pageURL is page n of the listing at base: base itself for page 1.
func pageURL(base string, n int) string {
	if n == 1 {
		return base
	}
	return base + "page/" + strconv.Itoa(n) + "/"
}
