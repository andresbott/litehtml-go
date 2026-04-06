package litehtml

// The litehtml/ directory contains vendored C/C++ sources with no .go files.
// Without this embed directive, `go mod vendor` skips the entire tree because
// it only copies directories it recognizes as Go packages or that are referenced
// by //go:embed. The cgo bridge files (bridge.cpp, bridge_gumbo.c) #include
// sources from litehtml/src/ and headers from litehtml/include/, so the tree
// must be present for the build to succeed.
//
// This embed is unused at runtime; it exists solely to ensure `go mod vendor`
// copies the C/C++ source tree into the consumer's vendor directory.

import "embed"

//go:embed all:litehtml/src all:litehtml/include
var _ embed.FS
