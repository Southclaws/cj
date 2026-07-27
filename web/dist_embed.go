//go:build embed_web

package web

import "embed"

//go:embed all:dist
var Dist embed.FS

const Embedded = true
