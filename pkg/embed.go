package pkg

import "embed"

//go:embed views/index.html
//go:embed views/main.css
//go:embed views/htmx.js
var Views embed.FS

