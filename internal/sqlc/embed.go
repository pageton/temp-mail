// Package sqlc contains the SQLite schema embedded into the binary.
package sqlc

import _ "embed"

//go:embed schema.sql
var Schema string
