# 03: Mermaid fences become inline SVG figures

**What to build:** The Author writes a Mermaid fence in an Explainer and the published page contains an inline SVG inside a `<figure>` with a caption taken from the fence's title line. The SVG uses `currentColor` and the site's CSS variables so it is legible in light and dark mode, keeps node labels as real text for crawlers, and carries explicit width/height or viewBox so it causes no layout shift. The renderer must run inside Vercel's build environment.

**Blocked by:** 01 (Walking skeleton)

**Status:** ready-for-agent

- [ ] Build-seam test: Explainer with two Mermaid fences → two `<figure>` elements each containing an `<svg>` and a `<figcaption>`
- [ ] Node label text from the fence is present as text nodes inside the SVG
- [ ] No hard-coded fill/stroke colours that would break dark mode; theme adaptation verified by inspecting attributes/styles
- [ ] A malformed Mermaid fence is a build error naming the page
- [ ] Rendering result cached by fence content hash so unchanged diagrams are not re-rendered
- [ ] Renderer choice and its install step documented in the repo README so the Vercel build command can run it
