package serve

import "embed"

//go:embed views/*.html views/main.css views/htmx.js
var views embed.FS
