# litehtml-go

AI generated Go bindings for [litehtml](http://www.litehtml.com/), a lightweight HTML/CSS
rendering engine. litehtml parses HTML/CSS and performs layout, then calls back
into your code to draw text, backgrounds, borders, and images. This package
wraps that engine via cgo so you can drive it entirely from Go.

## Requirements

- Go 1.25+
- A C/C++ toolchain (`gcc` / `clang` with C++17 support)
No external shared libraries are needed -- the litehtml C++ source is vendored
directly in the `litehtml/` directory and compiled from source as part of the
normal `go build` process through cgo.

## Installation

```bash
go get github.com/andresbott/litehtml-go
```

That's it -- no git submodules or system libraries to install.

## Quick start

The core workflow is:

1. Implement the `DocumentContainer` interface (fonts, drawing, resources).
2. Create a `Document` from an HTML string.
3. Call `Render` to lay out at a given width.
4. Call `Draw` to trigger your container's drawing callbacks.

```go
package main

import (
    "fmt"
    litehtml "github.com/andresbott/litehtml-go"
)

func main() {
    container := NewMyContainer(800, 600) // your DocumentContainer impl

    doc, err := litehtml.NewDocument(
        "<html><body><h1>Hello!</h1></body></html>",
        container,
        "", // master CSS (empty = litehtml default)
        "", // user CSS
    )
    if err != nil {
        panic(err)
    }
    defer doc.Close()

    doc.Render(800)

    clip := litehtml.Position{X: 0, Y: 0, Width: 800, Height: 600}
    doc.Draw(0, 0, 0, &clip)

    fmt.Printf("Document size: %.0f x %.0f\n", doc.Width(), doc.Height())
}
```

## API overview

### Document

| Method | Description |
|--------|-------------|
| `NewDocument(html, container, masterCSS, userCSS)` | Parse HTML and return a `*Document`. |
| `doc.Render(maxWidth)` | Perform layout. Returns the actual width used. |
| `doc.Draw(hdc, x, y, clip)` | Trigger draw callbacks. Pass `nil` for clip to draw everything. |
| `doc.Width()` / `doc.Height()` | Document dimensions after `Render`. |
| `doc.OnMouseOver(...)` | Mouse move -- returns changed flag and redraw boxes. |
| `doc.OnLButtonDown(...)` / `doc.OnLButtonUp(...)` | Mouse button events. |
| `doc.OnMouseLeave()` | Mouse left the document area. |
| `doc.Close()` | Free C++ resources. Safe to call multiple times. Also called by the finalizer. |
| `MasterCSS()` | Returns litehtml's built-in default stylesheet. |

### DocumentContainer interface

You must implement all methods of `DocumentContainer`. The engine calls these
during parsing, layout, and drawing. The full list:

**Fonts and text**
- `CreateFont(descr FontDescription) (uintptr, FontMetrics)` -- create a font, return a handle and metrics.
- `DeleteFont(hFont uintptr)` -- release a font handle.
- `TextWidth(text string, hFont uintptr) float32` -- measure text width.
- `DrawText(hdc uintptr, text string, hFont uintptr, color WebColor, pos Position)` -- draw a text run.
- `TransformText(text string, tt TextTransform) string` -- apply `text-transform` (uppercase, lowercase, etc.).
- `PtToPx(pt float32) float32` -- convert CSS points to pixels.
- `GetDefaultFontSize() float32` -- default font size in pixels (typically 16).
- `GetDefaultFontName() string` -- default font family name.

**Drawing**
- `DrawSolidFill(hdc uintptr, layer BackgroundLayer, color WebColor)` -- fill a background rectangle.
- `DrawBorders(hdc uintptr, borders Borders, drawPos Position, root bool)` -- draw element borders.
- `DrawListMarker(hdc uintptr, marker ListMarker)` -- draw a list bullet or number.
- `DrawImage(hdc uintptr, layer BackgroundLayer, url, baseURL string)` -- draw a background/foreground image.
- `DrawLinearGradient(...)` / `DrawRadialGradient(...)` / `DrawConicGradient(...)` -- draw CSS gradients.
- `SetClip(pos Position, bdrRadius BorderRadiuses)` / `DelClip()` -- push/pop a clipping region.

**Resources and navigation**
- `LoadImage(src, baseurl string, redrawOnReady bool)` -- start loading an image.
- `GetImageSize(src, baseurl string) Size` -- return a previously loaded image's dimensions.
- `ImportCSS(url, baseurl string) (text, newBaseURL string)` -- fetch an external stylesheet.
- `Link(href, rel, mediaType string)` -- notification of a `<link>` element.

**Environment**
- `GetViewport() Position` -- return the viewport rectangle.
- `GetMediaFeatures() MediaFeatures` -- describe the display (size, color depth, resolution).
- `GetLanguage() (language, culture string)` -- document language (e.g. `"en"`, `"en-US"`).

**Events and UI**
- `SetCaption(caption string)` -- `<title>` content.
- `SetBaseURL(baseURL string)` -- `<base>` href.
- `OnAnchorClick(url string)` -- a link was clicked.
- `OnMouseEvent(event MouseEvent)` -- mouse enter/leave an element.
- `SetCursor(cursor string)` -- requested cursor style.
- `CreateElement(tagName string, attributes map[string]string) uintptr` -- optionally create a custom element (return 0 to use default).

See `litehtml_types.go` for all supporting types (`Position`, `WebColor`,
`FontMetrics`, `Borders`, `BackgroundLayer`, gradients, etc.).

## Examples

Three working examples are included under `examples/`:

### basic

Renders a simple HTML page to a PNG image using Go's bundled fonts
(`golang.org/x/image/font/gofont`).

```bash
go run ./examples/basic
# writes output.png
```

### rich

Renders a more complex page with headings, styled divs, lists, multiple font
sizes, colors, borders, and an embedded generated image.

```bash
go run ./examples/rich
# writes output.png
```

### dump

Prints every container callback invocation to stdout. Useful for understanding
what litehtml does internally during parse, layout, and draw phases.

```bash
go run ./examples/dump
```

## Development

### Running tests

```bash
make test
```

### Full verification (lint, tests, license check, benchmarks, coverage)

```bash
make verify
```

This requires:
- [golangci-lint](https://golangci-lint.run/)
- [go-licence-detector](https://github.com/elastic/go-licence-detector)

### Coverage

The coverage threshold is 70%. If you add new exported functions, add
corresponding tests. Run the coverage report:

```bash
make cover-report
```

## Architecture

```
litehtml-go/
  litehtml/          # litehtml C++ source (vendored)
  bridge.h           # flat C API over litehtml's C++ classes
  bridge.cpp         # C++ implementation of the bridge
  bridge_gumbo.c     # compiles litehtml's Gumbo HTML parser
  container.go       # Go <-> C callback trampolines, DocumentContainer interface
  litehtml.go        # Document type and cgo bindings
  litehtml_types.go  # Go types mirroring litehtml's C++ enums and structs
  litehtml_test.go   # tests
  examples/          # runnable examples
```

The binding works by:

1. `bridge.h` / `bridge.cpp` expose litehtml's C++ API as plain C functions
   and a struct of callback function pointers (`lh_container_callbacks`).
2. `container.go` registers Go `DocumentContainer` implementations in a handle
   map and provides `//export` trampoline functions that C calls back into.
3. `litehtml.go` wraps the C document lifecycle (`create`, `render`, `draw`,
   `destroy`) in a Go-friendly `Document` type with a finalizer for safety.

## License

This project provides Go bindings. The litehtml engine itself is licensed under
its own terms (see `litehtml/LICENSE`).
