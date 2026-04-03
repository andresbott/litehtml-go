// Example: dump all litehtml container callbacks
//
// This program parses a rich HTML document, renders it, calls Draw, and
// prints every container callback invocation to stdout. It is useful as a
// debugging and learning tool to see exactly what litehtml does internally
// during parsing, layout, and drawing.
package main

import (
	"fmt"
	"strings"

	litehtml "github.com/andresbott/litehtml-go"
)

// loggingContainer implements DocumentContainer by printing every callback
// invocation with its arguments to stdout.
type loggingContainer struct {
	fontCounter uintptr // simple counter to assign font handles
}

func (c *loggingContainer) CreateFont(descr litehtml.FontDescription) (uintptr, litehtml.FontMetrics) {
	c.fontCounter++
	h := c.fontCounter
	style := "normal"
	if descr.Style == litehtml.FontStyleItalic {
		style = "italic"
	}
	fmt.Printf("  CreateFont        font=%d family=%q size=%.1f weight=%d style=%s\n",
		h, descr.Family, descr.Size, descr.Weight, style)
	return h, litehtml.FontMetrics{
		FontSize:   descr.Size,
		Height:     descr.Size * 1.2,
		Ascent:     descr.Size * 0.8,
		Descent:    descr.Size * 0.2,
		XHeight:    descr.Size * 0.5,
		ChWidth:    descr.Size * 0.6,
		DrawSpaces: true,
	}
}

func (c *loggingContainer) DeleteFont(hFont uintptr) {
	fmt.Printf("  DeleteFont        font=%d\n", hFont)
}

func (c *loggingContainer) TextWidth(text string, hFont uintptr) float32 {
	w := float32(len(text)) * 8.0
	// TextWidth is called very frequently during layout; omit from log
	// to keep output readable. Uncomment the next line to see every call:
	// fmt.Printf("  TextWidth         font=%d text=%q -> %.0f\n", hFont, text, w)
	return w
}

func (c *loggingContainer) DrawText(hdc uintptr, text string, hFont uintptr, color litehtml.WebColor, pos litehtml.Position) {
	fmt.Printf("  DrawText          font=%d pos=(%.0f,%.0f %.0fx%.0f) color=rgba(%d,%d,%d,%d) text=%q\n",
		hFont, pos.X, pos.Y, pos.Width, pos.Height,
		color.Red, color.Green, color.Blue, color.Alpha, text)
}

func (c *loggingContainer) PtToPx(pt float32) float32 {
	return pt * 96.0 / 72.0
}

func (c *loggingContainer) GetDefaultFontSize() float32 {
	return 16
}

func (c *loggingContainer) GetDefaultFontName() string {
	return "serif"
}

func (c *loggingContainer) DrawListMarker(hdc uintptr, marker litehtml.ListMarker) {
	fmt.Printf("  DrawListMarker    type=%d pos=(%.0f,%.0f %.0fx%.0f) index=%d color=rgba(%d,%d,%d,%d)\n",
		marker.MarkerType, marker.Pos.X, marker.Pos.Y, marker.Pos.Width, marker.Pos.Height,
		marker.Index, marker.Color.Red, marker.Color.Green, marker.Color.Blue, marker.Color.Alpha)
}

func (c *loggingContainer) LoadImage(src, baseurl string, redrawOnReady bool) {
	fmt.Printf("  LoadImage         src=%q baseurl=%q redraw=%v\n", src, baseurl, redrawOnReady)
}

func (c *loggingContainer) GetImageSize(src, baseurl string) litehtml.Size {
	fmt.Printf("  GetImageSize      src=%q baseurl=%q\n", src, baseurl)
	return litehtml.Size{}
}

func (c *loggingContainer) DrawImage(hdc uintptr, layer litehtml.BackgroundLayer, url, baseURL string) {
	fmt.Printf("  DrawImage         url=%q borderBox=(%.0f,%.0f %.0fx%.0f)\n",
		url, layer.BorderBox.X, layer.BorderBox.Y, layer.BorderBox.Width, layer.BorderBox.Height)
}

func (c *loggingContainer) DrawSolidFill(hdc uintptr, layer litehtml.BackgroundLayer, color litehtml.WebColor) {
	fmt.Printf("  DrawSolidFill     borderBox=(%.0f,%.0f %.0fx%.0f) color=rgba(%d,%d,%d,%d) root=%v\n",
		layer.BorderBox.X, layer.BorderBox.Y, layer.BorderBox.Width, layer.BorderBox.Height,
		color.Red, color.Green, color.Blue, color.Alpha, layer.IsRoot)
}

func (c *loggingContainer) DrawLinearGradient(hdc uintptr, layer litehtml.BackgroundLayer, gradient litehtml.LinearGradient) {
	fmt.Printf("  DrawLinearGradient start=(%.1f,%.1f) end=(%.1f,%.1f) stops=%d\n",
		gradient.Start.X, gradient.Start.Y, gradient.End.X, gradient.End.Y, len(gradient.ColorPoints))
}

func (c *loggingContainer) DrawRadialGradient(hdc uintptr, layer litehtml.BackgroundLayer, gradient litehtml.RadialGradient) {
	fmt.Printf("  DrawRadialGradient pos=(%.1f,%.1f) radius=(%.1f,%.1f) stops=%d\n",
		gradient.Position.X, gradient.Position.Y, gradient.Radius.X, gradient.Radius.Y, len(gradient.ColorPoints))
}

func (c *loggingContainer) DrawConicGradient(hdc uintptr, layer litehtml.BackgroundLayer, gradient litehtml.ConicGradient) {
	fmt.Printf("  DrawConicGradient pos=(%.1f,%.1f) angle=%.1f stops=%d\n",
		gradient.Position.X, gradient.Position.Y, gradient.Angle, len(gradient.ColorPoints))
}

func (c *loggingContainer) DrawBorders(hdc uintptr, borders litehtml.Borders, drawPos litehtml.Position, root bool) {
	fmt.Printf("  DrawBorders       pos=(%.0f,%.0f %.0fx%.0f) root=%v left=%.0f top=%.0f right=%.0f bottom=%.0f\n",
		drawPos.X, drawPos.Y, drawPos.Width, drawPos.Height, root,
		borders.Left.Width, borders.Top.Width, borders.Right.Width, borders.Bottom.Width)
}

func (c *loggingContainer) SetCaption(caption string) {
	fmt.Printf("  SetCaption        %q\n", caption)
}

func (c *loggingContainer) SetBaseURL(baseURL string) {
	fmt.Printf("  SetBaseURL        %q\n", baseURL)
}

func (c *loggingContainer) Link(href, rel, mediaType string) {
	fmt.Printf("  Link              href=%q rel=%q type=%q\n", href, rel, mediaType)
}

func (c *loggingContainer) OnAnchorClick(url string) {
	fmt.Printf("  OnAnchorClick     %q\n", url)
}

func (c *loggingContainer) OnMouseEvent(event litehtml.MouseEvent) {
	fmt.Printf("  OnMouseEvent      %d\n", event)
}

func (c *loggingContainer) SetCursor(cursor string) {
	fmt.Printf("  SetCursor         %q\n", cursor)
}

func (c *loggingContainer) TransformText(text string, tt litehtml.TextTransform) string {
	switch tt {
	case litehtml.TextTransformUppercase:
		return strings.ToUpper(text)
	case litehtml.TextTransformLowercase:
		return strings.ToLower(text)
	default:
		return text
	}
}

func (c *loggingContainer) ImportCSS(url, baseurl string) (string, string) {
	fmt.Printf("  ImportCSS         url=%q baseurl=%q\n", url, baseurl)
	return "", baseurl
}

func (c *loggingContainer) SetClip(pos litehtml.Position, bdrRadius litehtml.BorderRadiuses) {
	fmt.Printf("  SetClip           pos=(%.0f,%.0f %.0fx%.0f)\n", pos.X, pos.Y, pos.Width, pos.Height)
}

func (c *loggingContainer) DelClip() {
	fmt.Printf("  DelClip\n")
}

func (c *loggingContainer) GetViewport() litehtml.Position {
	return litehtml.Position{X: 0, Y: 0, Width: 800, Height: 600}
}

func (c *loggingContainer) CreateElement(tagName string, attributes map[string]string) uintptr {
	fmt.Printf("  CreateElement     tag=%q\n", tagName)
	return 0
}

func (c *loggingContainer) GetMediaFeatures() litehtml.MediaFeatures {
	return litehtml.MediaFeatures{
		Type:         litehtml.MediaTypeScreen,
		Width:        800,
		Height:       600,
		DeviceWidth:  800,
		DeviceHeight: 600,
		Color:        8,
		Resolution:   96,
	}
}

func (c *loggingContainer) GetLanguage() (string, string) {
	return "en", "en-US"
}

func main() {
	// A richer HTML document to exercise more callbacks.
	html := `<html>
<head>
  <title>Callback Dump Example</title>
</head>
<body>
  <h1>Welcome to litehtml-go</h1>
  <p>This paragraph contains <strong>bold text</strong> and <em>italic text</em>,
     as well as a <a href="https://example.com">hyperlink</a>.</p>
  <p style="color: red; background-color: #eee; padding: 8px; border: 1px solid #ccc;">
     This paragraph has inline styles: red text, a grey background, padding, and a border.
  </p>
  <h2>A List</h2>
  <ul>
    <li>First item</li>
    <li>Second item</li>
    <li>Third item</li>
  </ul>
  <h2>Another Section</h2>
  <p>Final paragraph with <span style="font-size: 20px;">larger text</span> inline.</p>
</body>
</html>`

	container := &loggingContainer{}

	// --- Phase 1: Parse ---
	fmt.Println("=== PARSE (NewDocument) ===")
	doc, err := litehtml.NewDocument(html, container, "", "")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer doc.Close()

	// --- Phase 2: Render (layout) ---
	fmt.Println("\n=== RENDER (800px) ===")
	doc.Render(800)
	fmt.Printf("\n  Document size: %.0f x %.0f\n", doc.Width(), doc.Height())

	// --- Phase 3: Draw ---
	fmt.Println("\n=== DRAW ===")
	clip := litehtml.Position{X: 0, Y: 0, Width: 800, Height: doc.Height()}
	doc.Draw(0, 0, 0, &clip)

	fmt.Println("\n=== DONE ===")
}
