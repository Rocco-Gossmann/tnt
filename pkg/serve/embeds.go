package serve

import "embed"

//go:embed views/index.html views/main.css views/htmx.js
var views embed.FS
