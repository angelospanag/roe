// Package migrations embeds the goose SQL migration files into the binary.
package migrations

import "embed"

//go:embed *.sql
var FS embed.FS
