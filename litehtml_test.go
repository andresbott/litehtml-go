package litehtml

import (
	"strings"
	"sync"
	"testing"
)

// testContainer is a minimal DocumentContainer implementation for testing.
type testContainer struct {
	mu              sync.Mutex
	createFontCalls int
	drawTextCalls   int
	viewportWidth   float32
	viewportHeight  float32
}

func newTestContainer(w, h float32) *testContainer {
	return &testContainer{viewportWidth: w, viewportHeight: h}
}

func (c *testContainer) CreateFont(descr FontDescription) (uintptr, FontMetrics) {
	c.mu.Lock()
	c.createFontCalls++
	c.mu.Unlock()
	return 1, FontMetrics{
		FontSize:   descr.Size,
		Height:     descr.Size * 1.2,
		Ascent:     descr.Size * 0.8,
		Descent:    descr.Size * 0.2,
		XHeight:    descr.Size * 0.5,
		ChWidth:    descr.Size * 0.6,
		DrawSpaces: true,
	}
}

func (c *testContainer) DeleteFont(hFont uintptr) {}

func (c *testContainer) TextWidth(text string, hFont uintptr) float32 {
	return float32(len(text)) * 8.0
}

func (c *testContainer) DrawText(hdc uintptr, text string, hFont uintptr, color WebColor, pos Position) {
	c.mu.Lock()
	c.drawTextCalls++
	c.mu.Unlock()
}

func (c *testContainer) PtToPx(pt float32) float32 {
	return pt * 96.0 / 72.0
}

func (c *testContainer) GetDefaultFontSize() float32 {
	return 16
}

func (c *testContainer) GetDefaultFontName() string {
	return "serif"
}

func (c *testContainer) DrawListMarker(hdc uintptr, marker ListMarker) {}

func (c *testContainer) LoadImage(src, baseurl string, redrawOnReady bool) {}

func (c *testContainer) GetImageSize(src, baseurl string) Size {
	return Size{}
}

func (c *testContainer) DrawImage(hdc uintptr, layer BackgroundLayer, url, baseURL string) {}

func (c *testContainer) DrawSolidFill(hdc uintptr, layer BackgroundLayer, color WebColor) {}

func (c *testContainer) DrawLinearGradient(hdc uintptr, layer BackgroundLayer, gradient LinearGradient) {
}

func (c *testContainer) DrawRadialGradient(hdc uintptr, layer BackgroundLayer, gradient RadialGradient) {
}

func (c *testContainer) DrawConicGradient(hdc uintptr, layer BackgroundLayer, gradient ConicGradient) {
}

func (c *testContainer) DrawBorders(hdc uintptr, borders Borders, drawPos Position, root bool) {}

func (c *testContainer) SetCaption(caption string) {}

func (c *testContainer) SetBaseURL(baseURL string) {}

func (c *testContainer) Link(href, rel, mediaType string) {}

func (c *testContainer) OnAnchorClick(url string) {}

func (c *testContainer) OnMouseEvent(event MouseEvent) {}

func (c *testContainer) SetCursor(cursor string) {}

func (c *testContainer) TransformText(text string, tt TextTransform) string {
	switch tt {
	case TextTransformUppercase:
		return strings.ToUpper(text)
	case TextTransformLowercase:
		return strings.ToLower(text)
	default:
		return text
	}
}

func (c *testContainer) ImportCSS(url, baseurl string) (string, string) {
	return "", baseurl
}

func (c *testContainer) SetClip(pos Position, bdrRadius BorderRadiuses) {}

func (c *testContainer) DelClip() {}

func (c *testContainer) GetViewport() Position {
	return Position{0, 0, c.viewportWidth, c.viewportHeight}
}

func (c *testContainer) CreateElement(tagName string, attributes map[string]string) uintptr {
	return 0
}

func (c *testContainer) GetMediaFeatures() MediaFeatures {
	return MediaFeatures{
		Type:         MediaTypeScreen,
		Width:        c.viewportWidth,
		Height:       c.viewportHeight,
		DeviceWidth:  c.viewportWidth,
		DeviceHeight: c.viewportHeight,
		Color:        8,
		Resolution:   96,
	}
}

func (c *testContainer) GetLanguage() (string, string) {
	return "en", "en-US"
}

// trackingContainer extends testContainer to track all callback invocations.
type trackingContainer struct {
	testContainer
	drawSolidFillCalls  int
	drawBordersCalls    int
	drawListMarkerCalls int
	loadImageCalls      int
	getImageSizeCalls   int
	setCaptionCalls     int
	setBaseURLCalls     int
	linkCalls           int
	onAnchorClickCalls  int
	setClipCalls        int
	delClipCalls        int
	transformTextCalls  int
	importCSSCalls      int
	setCursorCalls      int
	ptToPxCalls         int
	drawImageCalls      int
	getLanguageCalls    int
	onMouseEventCalls   int
	lastCaption         string
}

func newTrackingContainer(w, h float32) *trackingContainer {
	return &trackingContainer{
		testContainer: testContainer{viewportWidth: w, viewportHeight: h},
	}
}

func (c *trackingContainer) DrawSolidFill(hdc uintptr, layer BackgroundLayer, color WebColor) {
	c.mu.Lock()
	c.drawSolidFillCalls++
	c.mu.Unlock()
}

func (c *trackingContainer) DrawBorders(hdc uintptr, borders Borders, drawPos Position, root bool) {
	c.mu.Lock()
	c.drawBordersCalls++
	c.mu.Unlock()
}

func (c *trackingContainer) DrawListMarker(hdc uintptr, marker ListMarker) {
	c.mu.Lock()
	c.drawListMarkerCalls++
	c.mu.Unlock()
}

func (c *trackingContainer) LoadImage(src, baseurl string, redrawOnReady bool) {
	c.mu.Lock()
	c.loadImageCalls++
	c.mu.Unlock()
}

func (c *trackingContainer) GetImageSize(src, baseurl string) Size {
	c.mu.Lock()
	c.getImageSizeCalls++
	c.mu.Unlock()
	return Size{Width: 100, Height: 50}
}

func (c *trackingContainer) DrawImage(hdc uintptr, layer BackgroundLayer, url, baseURL string) {
	c.mu.Lock()
	c.drawImageCalls++
	c.mu.Unlock()
}

func (c *trackingContainer) SetCaption(caption string) {
	c.mu.Lock()
	c.setCaptionCalls++
	c.lastCaption = caption
	c.mu.Unlock()
}

func (c *trackingContainer) SetBaseURL(baseURL string) {
	c.mu.Lock()
	c.setBaseURLCalls++
	c.mu.Unlock()
}

func (c *trackingContainer) Link(href, rel, mediaType string) {
	c.mu.Lock()
	c.linkCalls++
	c.mu.Unlock()
}

func (c *trackingContainer) OnAnchorClick(url string) {
	c.mu.Lock()
	c.onAnchorClickCalls++
	c.mu.Unlock()
}

func (c *trackingContainer) OnMouseEvent(event MouseEvent) {
	c.mu.Lock()
	c.onMouseEventCalls++
	c.mu.Unlock()
}

func (c *trackingContainer) SetCursor(cursor string) {
	c.mu.Lock()
	c.setCursorCalls++
	c.mu.Unlock()
}

func (c *trackingContainer) SetClip(pos Position, bdrRadius BorderRadiuses) {
	c.mu.Lock()
	c.setClipCalls++
	c.mu.Unlock()
}

func (c *trackingContainer) DelClip() {
	c.mu.Lock()
	c.delClipCalls++
	c.mu.Unlock()
}

func (c *trackingContainer) TransformText(text string, tt TextTransform) string {
	c.mu.Lock()
	c.transformTextCalls++
	c.mu.Unlock()
	switch tt {
	case TextTransformUppercase:
		return strings.ToUpper(text)
	case TextTransformLowercase:
		return strings.ToLower(text)
	default:
		return text
	}
}

func (c *trackingContainer) ImportCSS(url, baseurl string) (string, string) {
	c.mu.Lock()
	c.importCSSCalls++
	c.mu.Unlock()
	return "", baseurl
}

func (c *trackingContainer) PtToPx(pt float32) float32 {
	c.mu.Lock()
	c.ptToPxCalls++
	c.mu.Unlock()
	return pt * 96.0 / 72.0
}

func (c *trackingContainer) GetLanguage() (string, string) {
	c.mu.Lock()
	c.getLanguageCalls++
	c.mu.Unlock()
	return "en", "en-US"
}

// ── Tests ──

func TestBasicParse(t *testing.T) {
	container := newTestContainer(800, 600)
	doc, err := NewDocument("<html><body><p>Hello World</p></body></html>", container, "", "")
	if err != nil {
		t.Fatalf("NewDocument failed: %v", err)
	}
	defer doc.Close()

	doc.Render(800)

	w := doc.Width()
	h := doc.Height()
	if w <= 0 {
		t.Errorf("expected positive width, got %f", w)
	}
	if h <= 0 {
		t.Errorf("expected positive height, got %f", h)
	}
	t.Logf("Document size: %f x %f", w, h)
}

func TestCallbacksInvoked(t *testing.T) {
	container := newTestContainer(800, 600)
	doc, err := NewDocument("<html><body><p>Test</p></body></html>", container, "", "")
	if err != nil {
		t.Fatalf("NewDocument failed: %v", err)
	}
	defer doc.Close()

	doc.Render(800)

	container.mu.Lock()
	fontCalls := container.createFontCalls
	container.mu.Unlock()

	if fontCalls == 0 {
		t.Error("expected CreateFont to be called at least once")
	}
}

func TestDrawInvokesCallbacks(t *testing.T) {
	container := newTestContainer(800, 600)
	doc, err := NewDocument("<html><body><p>Draw me</p></body></html>", container, "", "")
	if err != nil {
		t.Fatalf("NewDocument failed: %v", err)
	}
	defer doc.Close()

	doc.Render(800)
	doc.Draw(0, 0, 0, nil)

	container.mu.Lock()
	drawCalls := container.drawTextCalls
	container.mu.Unlock()

	if drawCalls == 0 {
		t.Error("expected DrawText to be called at least once during Draw")
	}
}

func TestMultipleDocuments(t *testing.T) {
	c1 := newTestContainer(800, 600)
	c2 := newTestContainer(1024, 768)

	doc1, err := NewDocument("<html><body><p>Doc 1</p></body></html>", c1, "", "")
	if err != nil {
		t.Fatalf("NewDocument doc1 failed: %v", err)
	}
	defer doc1.Close()

	doc2, err := NewDocument("<html><body><h1>Doc 2</h1></body></html>", c2, "", "")
	if err != nil {
		t.Fatalf("NewDocument doc2 failed: %v", err)
	}
	defer doc2.Close()

	doc1.Render(800)
	doc2.Render(1024)

	if doc1.Width() <= 0 || doc1.Height() <= 0 {
		t.Error("doc1 should have positive dimensions")
	}
	if doc2.Width() <= 0 || doc2.Height() <= 0 {
		t.Error("doc2 should have positive dimensions")
	}
}

func TestCloseIsIdempotent(t *testing.T) {
	container := newTestContainer(800, 600)
	doc, err := NewDocument("<html><body>test</body></html>", container, "", "")
	if err != nil {
		t.Fatalf("NewDocument failed: %v", err)
	}
	doc.Close()
	doc.Close() // should not panic

	// Methods on closed doc should return zero values
	if doc.Width() != 0 {
		t.Error("Width on closed doc should return 0")
	}
}

func TestMasterCSS(t *testing.T) {
	css := MasterCSS()
	if css == "" {
		t.Error("MasterCSS should return non-empty string")
	}
}

func TestNewDocumentWithCustomCSS(t *testing.T) {
	container := newTestContainer(800, 600)
	masterCSS := MasterCSS()
	userCSS := "body { color: red; }"
	doc, err := NewDocument("<html><body><p>Hello</p></body></html>", container, masterCSS, userCSS)
	if err != nil {
		t.Fatalf("NewDocument with custom CSS failed: %v", err)
	}
	defer doc.Close()
	doc.Render(800)
	if doc.Width() <= 0 {
		t.Error("expected positive width")
	}
}

func TestDrawWithClip(t *testing.T) {
	container := newTestContainer(800, 600)
	doc, err := NewDocument("<html><body><p>Clip test</p></body></html>", container, "", "")
	if err != nil {
		t.Fatalf("NewDocument failed: %v", err)
	}
	defer doc.Close()
	doc.Render(800)
	clip := Position{X: 0, Y: 0, Width: 800, Height: 600}
	doc.Draw(0, 0, 0, &clip)
}

func TestMouseEvents(t *testing.T) {
	container := newTrackingContainer(800, 600)
	html := `<html><body><a href="https://example.com">Click me</a></body></html>`
	doc, err := NewDocument(html, container, "", "")
	if err != nil {
		t.Fatalf("NewDocument failed: %v", err)
	}
	defer doc.Close()
	doc.Render(800)

	changed, boxes := doc.OnMouseOver(50, 20, 50, 20)
	_ = changed
	_ = boxes

	changed, boxes = doc.OnLButtonDown(50, 20, 50, 20)
	_ = changed
	_ = boxes

	changed, boxes = doc.OnLButtonUp(50, 20, 50, 20)
	_ = changed
	_ = boxes

	changed, boxes = doc.OnMouseLeave()
	_ = changed
	_ = boxes
}

func TestMouseEventsOnClosedDoc(t *testing.T) {
	container := newTestContainer(800, 600)
	doc, err := NewDocument("<html><body>test</body></html>", container, "", "")
	if err != nil {
		t.Fatalf("NewDocument failed: %v", err)
	}
	doc.Close()

	changed, boxes := doc.OnMouseOver(0, 0, 0, 0)
	if changed || boxes != nil {
		t.Error("OnMouseOver on closed doc should return false, nil")
	}

	changed, boxes = doc.OnLButtonDown(0, 0, 0, 0)
	if changed || boxes != nil {
		t.Error("OnLButtonDown on closed doc should return false, nil")
	}

	changed, boxes = doc.OnLButtonUp(0, 0, 0, 0)
	if changed || boxes != nil {
		t.Error("OnLButtonUp on closed doc should return false, nil")
	}

	changed, boxes = doc.OnMouseLeave()
	if changed || boxes != nil {
		t.Error("OnMouseLeave on closed doc should return false, nil")
	}

	// Render and Draw on closed doc
	if doc.Render(800) != 0 {
		t.Error("Render on closed doc should return 0")
	}
	doc.Draw(0, 0, 0, nil) // should not panic

	if doc.Height() != 0 {
		t.Error("Height on closed doc should return 0")
	}
}

func TestRichHTMLCallbacks(t *testing.T) {
	container := newTrackingContainer(800, 600)
	html := `<html>
<head>
  <title>Test Title</title>
  <link rel="stylesheet" href="style.css" type="text/css">
</head>
<body>
  <p style="border: 2px solid red; background-color: #eee; padding: 10px;">
    Bordered paragraph
  </p>
  <ul>
    <li>Item one</li>
    <li>Item two</li>
  </ul>
  <img src="test.png" alt="test">
  <p style="text-transform: uppercase;">uppercase text</p>
</body>
</html>`

	doc, err := NewDocument(html, container, "", "")
	if err != nil {
		t.Fatalf("NewDocument failed: %v", err)
	}
	defer doc.Close()

	doc.Render(800)
	clip := Position{X: 0, Y: 0, Width: 800, Height: 600}
	doc.Draw(0, 0, 0, &clip)

	container.mu.Lock()
	defer container.mu.Unlock()

	if container.drawSolidFillCalls == 0 {
		t.Error("expected DrawSolidFill to be called")
	}
	if container.drawBordersCalls == 0 {
		t.Error("expected DrawBorders to be called")
	}
	if container.drawListMarkerCalls == 0 {
		t.Error("expected DrawListMarker to be called")
	}
	if container.loadImageCalls == 0 {
		t.Error("expected LoadImage to be called")
	}
	if container.getImageSizeCalls == 0 {
		t.Error("expected GetImageSize to be called")
	}
	if container.setCaptionCalls == 0 {
		t.Error("expected SetCaption to be called")
	}
	if container.lastCaption != "Test Title" {
		t.Errorf("expected caption 'Test Title', got %q", container.lastCaption)
	}
	if container.transformTextCalls == 0 {
		t.Error("expected TransformText to be called")
	}
	// PtToPx and GetLanguage may or may not be called depending on the HTML;
	// they are tested separately.
}

func TestSetClipAndDelClip(t *testing.T) {
	container := newTrackingContainer(800, 600)
	html := `<html><body>
<div style="overflow: hidden; width: 100px; height: 50px;">
  <p>This is clipped content that overflows the container div intentionally.</p>
</div>
</body></html>`

	doc, err := NewDocument(html, container, "", "")
	if err != nil {
		t.Fatalf("NewDocument failed: %v", err)
	}
	defer doc.Close()

	doc.Render(800)
	clip := Position{X: 0, Y: 0, Width: 800, Height: 600}
	doc.Draw(0, 0, 0, &clip)

	container.mu.Lock()
	defer container.mu.Unlock()

	if container.setClipCalls == 0 {
		t.Error("expected SetClip to be called for overflow:hidden")
	}
	if container.delClipCalls == 0 {
		t.Error("expected DelClip to be called")
	}
}

func TestImportCSS(t *testing.T) {
	container := newTrackingContainer(800, 600)
	html := `<html>
<head>
  <link rel="stylesheet" href="external.css">
</head>
<body><p>Test</p></body>
</html>`

	doc, err := NewDocument(html, container, "", "")
	if err != nil {
		t.Fatalf("NewDocument failed: %v", err)
	}
	defer doc.Close()
	doc.Render(800)

	container.mu.Lock()
	defer container.mu.Unlock()

	if container.importCSSCalls == 0 {
		t.Error("expected ImportCSS to be called for external stylesheet link")
	}
}
