# litehtml-go

Go bindings for [litehtml](https://github.com/nicehash/litehtml), a lightweight HTML/CSS rendering engine. Uses cgo to bridge Go with vendored C++ litehtml source. The engine calls back into a user-provided `DocumentContainer` interface for drawing.

## Build & Test

```bash
make test             # run tests
make lint             # golangci-lint (must be installed)
make benchmark        # run benchmarks
make coverage         # check coverage (70% threshold)
make verify           # all checks: lint, test, license-check, benchmark, coverage
```

Requires: Go 1.25+, C++17 compiler (gcc/clang), cgo enabled.

## Architecture

- Single package `litehtml`
- `bridge.h` / `bridge.cpp` / `bridge_gumbo.c` — flat C API over C++ classes
- `container.go` — `DocumentContainer` interface + Go↔C callback trampolines (handle-based)
- `litehtml.go` — `Document` type and cgo bindings
- `litehtml_types.go` — Go types mirroring C++ enums/structs
- `litehtml/` — vendored C++ litehtml source (do not edit directly, see update procedure below)
- `vendor_compat.go` — `//go:embed` directive that forces `go mod vendor` to copy the C/C++ tree (see below)

## Conventions

- Go naming: CamelCase types mirroring C++ names (e.g., `FontStyle`, `WebColor`, `Position`)
- Resource management: explicit `Close()` methods with finalizers
- Linting: golangci-lint with gocyclo (max 20), nestif (max 5), dupl, gosec
- Tests excluded from nestif and dupl lint rules

## go mod vendor and the litehtml/ directory

### The problem

`go mod vendor` only copies directories that are part of the Go import graph (i.e., Go packages transitively imported by the consumer). The `litehtml/` tree contains only C/C++ sources with zero `.go` files, so `go mod vendor` skips it entirely. However, `bridge.cpp` and `bridge_gumbo.c` in the root `#include` files from `litehtml/src/` and `litehtml/include/`, so vendored builds fail with missing file errors.

This only affects `go mod vendor` users. The module cache (`go env GOMODCACHE`) downloads the full module archive and works fine.

### The fix

`vendor_compat.go` in the root package contains:

```go
//go:embed all:litehtml/src all:litehtml/include
var litehtmlSources embed.FS
```

This `//go:embed` directive tells the Go toolchain that these directories are needed, which forces `go mod vendor` to copy them into the consumer's `vendor/` directory. The embedded data is unused at runtime — it exists solely for vendor compatibility.

The `all:` prefix is required to include files starting with `.` or `_` (some litehtml sources may have such names in the future).

Only `litehtml/src/` and `litehtml/include/` are embedded because those are the only directories referenced by the cgo build (via `#include` in bridge files and `-I` flags in `litehtml.go`). Other directories like `litehtml/containers/`, `litehtml/doc/`, `litehtml/support/` are not needed for building.

### Why not dummy .go files?

Adding `vendor_hint.go` files to each subdirectory does NOT work because:

1. `go mod vendor` only copies packages in the import graph — having a `.go` file makes a directory a "package" but doesn't put it in the import graph
2. Blank-importing the sub-packages from the root would put them in the graph, but then cgo tries to compile the C/C++ files in those directories as separate packages, causing "C++ source files not allowed when not using cgo or SWIG" errors
3. Adding `import "C"` to the dummy files would cause duplicate symbol errors (same sources compiled both in the sub-package and via `#include` in the root bridge files)

### Updating the vendored litehtml C++ library

When updating the upstream litehtml source in the `litehtml/` directory:

1. **Replace the `litehtml/` directory** with the new upstream source
2. **Check if directory structure changed**: if upstream added/removed/renamed directories under `litehtml/src/` or `litehtml/include/`, no changes to `vendor_compat.go` are needed (the embed globs `all:litehtml/src` and `all:litehtml/include` are recursive)
3. **Check if source files were added/removed in `litehtml/src/`**: update the `#include` lists in `bridge.cpp` (for `.cpp` files) and `bridge_gumbo.c` (for `.c` files in `litehtml/src/gumbo/`)
4. **Check if headers were added/removed in `litehtml/include/`**: update the `-I` flags in `litehtml.go` cgo directives if new include directories were introduced
5. **Check if new top-level directories are needed**: if the build requires files from directories OTHER than `litehtml/src/` and `litehtml/include/` (e.g., a new `litehtml/third_party/`), add them to the `//go:embed` directive in `vendor_compat.go`
6. **Verify**: run `make verify` and test a vendored consumer build:
   ```bash
   # From a separate test project that imports this library:
   go mod vendor && go build -mod=vendor ./...
   ```
