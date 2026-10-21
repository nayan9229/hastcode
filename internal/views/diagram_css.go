package views

import _ "embed"

// diagramCSS styles Mermaid figures (partial_diagram.templ); it is part of
// the site stylesheet (buildStylesheet).
//
//go:embed diagram.css
var diagramCSS string
