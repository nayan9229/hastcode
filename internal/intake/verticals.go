package intake

import (
	"slices"

	"hastcode.com/site/builder"
)

// Vertical is a Vertical Intake files Leads under (site.yaml "verticals").
type Vertical struct {
	Slug string
	Name string
}

// LoadVerticals reads the site config (site.yaml) and returns its
// Verticals in display order: the ones a Source may name, in the order the
// queue groups Leads.
func LoadVerticals(siteConfig string) ([]Vertical, error) {
	cfg, err := builder.LoadConfig(siteConfig)
	if err != nil {
		return nil, err
	}
	vs := slices.Clone(cfg.Verticals)
	slices.SortStableFunc(vs, func(a, b builder.Vertical) int { return a.Order - b.Order })
	out := make([]Vertical, len(vs))
	for i, v := range vs {
		out[i] = Vertical{Slug: v.Slug, Name: v.Name}
	}
	return out, nil
}
