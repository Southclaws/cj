//go:build !embed_web

package web

import "embed"

var Dist embed.FS

const Embedded = false
