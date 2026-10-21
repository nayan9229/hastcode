# 05: Explainer head and structured data

**What to build:** Every Explainer page carries what Google, Discover and answer engines need to index and cite it: an `Article` JSON-LD block (headline, description, image in 16:9/4:3/1:1, datePublished and dateModified with timezone matching the visible dates, author as `Person` with url and sameAs, publisher as `Organization` named "Hast Code", mainEntityOfPage, articleSection = Vertical name, keywords = tags), a `BreadcrumbList` (Home → Vertical → Explainer), a canonical URL with trailing slash, Open Graph and Twitter card tags using the hero, `max-image-preview:large` robots meta, `theme-color` for both colour schemes, and `lang="en"`. No `FAQPage`, `HowTo`, `Speakable`, or `TechArticle`.

**Blocked by:** 04 (Hero image pipeline)

**Status:** resolved

- [x] Build-seam test decodes the JSON-LD and asserts every required field, all three image ratio URLs, and that dates are RFC 3339 with offset and equal to the frontmatter values
- [x] Visible dates on the page are the same instants as the JSON-LD dates
- [x] `BreadcrumbList` has three items with correct names and absolute URLs
- [x] Canonical, og:image (≥1200 px variant), og:type article, twitter:card summary_large_image present
- [x] Robots meta contains `max-image-preview:large`; no `noindex` anywhere in production output
- [x] `dateModified` equals `reviewed` from frontmatter, never the build time
