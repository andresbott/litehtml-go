package litehtml

/*
#include "bridge.h"
#include <stdlib.h>
*/
import "C"
import (
	"sync"
	"unsafe"
)

// DocumentContainer is the callback interface that consumers implement
// to provide font, drawing, and resource loading services to litehtml.
type DocumentContainer interface {
	CreateFont(descr FontDescription) (uintptr, FontMetrics)
	DeleteFont(hFont uintptr)
	TextWidth(text string, hFont uintptr) float32
	DrawText(hdc uintptr, text string, hFont uintptr, color WebColor, pos Position)
	PtToPx(pt float32) float32
	GetDefaultFontSize() float32
	GetDefaultFontName() string
	DrawListMarker(hdc uintptr, marker ListMarker)
	LoadImage(src, baseurl string, redrawOnReady bool)
	GetImageSize(src, baseurl string) Size
	DrawImage(hdc uintptr, layer BackgroundLayer, url, baseURL string)
	DrawSolidFill(hdc uintptr, layer BackgroundLayer, color WebColor)
	DrawLinearGradient(hdc uintptr, layer BackgroundLayer, gradient LinearGradient)
	DrawRadialGradient(hdc uintptr, layer BackgroundLayer, gradient RadialGradient)
	DrawConicGradient(hdc uintptr, layer BackgroundLayer, gradient ConicGradient)
	DrawBorders(hdc uintptr, borders Borders, drawPos Position, root bool)
	SetCaption(caption string)
	SetBaseURL(baseURL string)
	Link(href, rel, mediaType string)
	OnAnchorClick(url string)
	OnMouseEvent(event MouseEvent)
	SetCursor(cursor string)
	TransformText(text string, tt TextTransform) string
	ImportCSS(url, baseurl string) (text string, newBaseURL string)
	SetClip(pos Position, bdrRadius BorderRadiuses)
	DelClip()
	GetViewport() Position
	CreateElement(tagName string, attributes map[string]string) uintptr
	GetMediaFeatures() MediaFeatures
	GetLanguage() (language, culture string)
}

// ── Handle registry ──

var (
	containerMu  sync.Mutex
	containerMap = make(map[uintptr]DocumentContainer)
	nextHandle   uintptr
)

func registerContainer(c DocumentContainer) uintptr {
	containerMu.Lock()
	defer containerMu.Unlock()
	nextHandle++
	containerMap[nextHandle] = c
	return nextHandle
}

func lookupContainer(handle uintptr) DocumentContainer {
	containerMu.Lock()
	defer containerMu.Unlock()
	return containerMap[handle]
}

func unregisterContainer(handle uintptr) {
	containerMu.Lock()
	defer containerMu.Unlock()
	delete(containerMap, handle)
}

// ── Helper conversions from C to Go types ──

func positionFromC(p C.lh_position) Position {
	return Position{float32(p.x), float32(p.y), float32(p.width), float32(p.height)}
}

func webColorFromC(c C.lh_web_color) WebColor {
	return WebColor{uint8(c.red), uint8(c.green), uint8(c.blue), uint8(c.alpha)}
}

func borderFromC(b C.lh_border) Border {
	return Border{float32(b.width), BorderStyle(b.style), webColorFromC(b.color)}
}

func borderRadiusesFromC(r C.lh_border_radiuses) BorderRadiuses {
	return BorderRadiuses{
		float32(r.top_left_x), float32(r.top_left_y),
		float32(r.top_right_x), float32(r.top_right_y),
		float32(r.bottom_right_x), float32(r.bottom_right_y),
		float32(r.bottom_left_x), float32(r.bottom_left_y),
	}
}

func bordersFromC(b *C.lh_borders) Borders {
	return Borders{
		Left:   borderFromC(b.left),
		Top:    borderFromC(b.top),
		Right:  borderFromC(b.right),
		Bottom: borderFromC(b.bottom),
		Radius: borderRadiusesFromC(b.radius),
	}
}

func backgroundLayerFromC(l *C.lh_background_layer) BackgroundLayer {
	return BackgroundLayer{
		BorderBox:    positionFromC(l.border_box),
		BorderRadius: borderRadiusesFromC(l.border_radius),
		ClipBox:      positionFromC(l.clip_box),
		OriginBox:    positionFromC(l.origin_box),
		Attachment:   BackgroundAttachment(l.attachment),
		Repeat:       BackgroundRepeat(l.repeat),
		IsRoot:       l.is_root != 0,
	}
}

func colorPointsFromC(pts *C.lh_color_point, count C.int) []ColorPoint {
	if count == 0 || pts == nil {
		return nil
	}
	n := int(count)
	slice := unsafe.Slice(pts, n)
	result := make([]ColorPoint, n)
	for i := 0; i < n; i++ {
		result[i] = ColorPoint{
			Offset: float32(slice[i].offset),
			Color:  webColorFromC(slice[i].color),
		}
	}
	return result
}

func fontDescriptionFromC(d *C.lh_font_description) FontDescription {
	return FontDescription{
		Family:              C.GoString(d.family),
		Size:                float32(d.size),
		Style:               FontStyle(d.style),
		Weight:              int(d.weight),
		DecorationLine:      int(d.decoration_line),
		DecorationThickness: float32(d.decoration_thickness),
		DecorationStyle:     TextDecorationStyle(d.decoration_style),
		DecorationColor:     webColorFromC(d.decoration_color),
		EmphasisStyle:       C.GoString(d.emphasis_style),
		EmphasisColor:       webColorFromC(d.emphasis_color),
		EmphasisPosition:    int(d.emphasis_position),
	}
}

// ── Callback trampoline functions (called from C) ──

//export goCreateFont
func goCreateFont(goHandle C.uintptr_t, descr *C.lh_font_description, fm *C.lh_font_metrics) C.uintptr_t {
	c := lookupContainer(uintptr(goHandle))
	if c == nil {
		return 0
	}
	hFont, metrics := c.CreateFont(fontDescriptionFromC(descr))
	fm.font_size = C.float(metrics.FontSize)
	fm.height = C.float(metrics.Height)
	fm.ascent = C.float(metrics.Ascent)
	fm.descent = C.float(metrics.Descent)
	fm.x_height = C.float(metrics.XHeight)
	fm.ch_width = C.float(metrics.ChWidth)
	if metrics.DrawSpaces {
		fm.draw_spaces = 1
	} else {
		fm.draw_spaces = 0
	}
	fm.sub_shift = C.float(metrics.SubShift)
	fm.super_shift = C.float(metrics.SuperShift)
	return C.uintptr_t(hFont)
}

//export goDeleteFont
func goDeleteFont(goHandle C.uintptr_t, hFont C.uintptr_t) {
	if c := lookupContainer(uintptr(goHandle)); c != nil {
		c.DeleteFont(uintptr(hFont))
	}
}

//export goTextWidth
func goTextWidth(goHandle C.uintptr_t, text *C.char, hFont C.uintptr_t) C.float {
	if c := lookupContainer(uintptr(goHandle)); c != nil {
		return C.float(c.TextWidth(C.GoString(text), uintptr(hFont)))
	}
	return 0
}

//export goDrawText
func goDrawText(goHandle C.uintptr_t, hdc C.uintptr_t, text *C.char, hFont C.uintptr_t, color C.lh_web_color, pos C.lh_position) {
	if c := lookupContainer(uintptr(goHandle)); c != nil {
		c.DrawText(uintptr(hdc), C.GoString(text), uintptr(hFont), webColorFromC(color), positionFromC(pos))
	}
}

//export goPtToPx
func goPtToPx(goHandle C.uintptr_t, pt C.float) C.float {
	if c := lookupContainer(uintptr(goHandle)); c != nil {
		return C.float(c.PtToPx(float32(pt)))
	}
	return 0
}

//export goGetDefaultFontSize
func goGetDefaultFontSize(goHandle C.uintptr_t) C.float {
	if c := lookupContainer(uintptr(goHandle)); c != nil {
		return C.float(c.GetDefaultFontSize())
	}
	return 16
}

// goDefaultFontNameBuf is used to hold the return value of goGetDefaultFontName
// so the C side gets a valid pointer. We use a global because the C++ side
// copies it into its own buffer immediately.
var goDefaultFontNameBuf *C.char

//export goGetDefaultFontName
func goGetDefaultFontName(goHandle C.uintptr_t) *C.char {
	if c := lookupContainer(uintptr(goHandle)); c != nil {
		name := c.GetDefaultFontName()
		if goDefaultFontNameBuf != nil {
			C.free(unsafe.Pointer(goDefaultFontNameBuf))
		}
		goDefaultFontNameBuf = C.CString(name)
		return goDefaultFontNameBuf
	}
	return nil
}

//export goDrawListMarker
func goDrawListMarker(goHandle C.uintptr_t, hdc C.uintptr_t, marker *C.lh_list_marker) {
	if c := lookupContainer(uintptr(goHandle)); c != nil {
		m := ListMarker{
			Image:      C.GoString(marker.image),
			BaseURL:    C.GoString(marker.baseurl),
			MarkerType: ListStyleType(marker.marker_type),
			Color:      webColorFromC(marker.color),
			Pos:        positionFromC(marker.pos),
			Index:      int(marker.index),
			Font:       uintptr(marker.font),
		}
		c.DrawListMarker(uintptr(hdc), m)
	}
}

//export goLoadImage
func goLoadImage(goHandle C.uintptr_t, src *C.char, baseurl *C.char, redrawOnReady C.int) {
	if c := lookupContainer(uintptr(goHandle)); c != nil {
		c.LoadImage(C.GoString(src), C.GoString(baseurl), redrawOnReady != 0)
	}
}

//export goGetImageSize
func goGetImageSize(goHandle C.uintptr_t, src *C.char, baseurl *C.char, sz *C.lh_size) {
	if c := lookupContainer(uintptr(goHandle)); c != nil {
		s := c.GetImageSize(C.GoString(src), C.GoString(baseurl))
		sz.width = C.float(s.Width)
		sz.height = C.float(s.Height)
	}
}

//export goDrawImage
func goDrawImage(goHandle C.uintptr_t, hdc C.uintptr_t, layer *C.lh_background_layer, url *C.char, baseURL *C.char) {
	if c := lookupContainer(uintptr(goHandle)); c != nil {
		c.DrawImage(uintptr(hdc), backgroundLayerFromC(layer), C.GoString(url), C.GoString(baseURL))
	}
}

//export goDrawSolidFill
func goDrawSolidFill(goHandle C.uintptr_t, hdc C.uintptr_t, layer *C.lh_background_layer, color C.lh_web_color) {
	if c := lookupContainer(uintptr(goHandle)); c != nil {
		c.DrawSolidFill(uintptr(hdc), backgroundLayerFromC(layer), webColorFromC(color))
	}
}

//export goDrawLinearGradient
func goDrawLinearGradient(goHandle C.uintptr_t, hdc C.uintptr_t, layer *C.lh_background_layer, gradient *C.lh_linear_gradient) {
	if c := lookupContainer(uintptr(goHandle)); c != nil {
		g := LinearGradient{
			Start:       PointF{float32(gradient.start.x), float32(gradient.start.y)},
			End:         PointF{float32(gradient.end.x), float32(gradient.end.y)},
			ColorPoints: colorPointsFromC(gradient.color_points, gradient.color_points_count),
		}
		c.DrawLinearGradient(uintptr(hdc), backgroundLayerFromC(layer), g)
	}
}

//export goDrawRadialGradient
func goDrawRadialGradient(goHandle C.uintptr_t, hdc C.uintptr_t, layer *C.lh_background_layer, gradient *C.lh_radial_gradient) {
	if c := lookupContainer(uintptr(goHandle)); c != nil {
		g := RadialGradient{
			Position:    PointF{float32(gradient.position.x), float32(gradient.position.y)},
			Radius:      PointF{float32(gradient.radius.x), float32(gradient.radius.y)},
			ColorPoints: colorPointsFromC(gradient.color_points, gradient.color_points_count),
		}
		c.DrawRadialGradient(uintptr(hdc), backgroundLayerFromC(layer), g)
	}
}

//export goDrawConicGradient
func goDrawConicGradient(goHandle C.uintptr_t, hdc C.uintptr_t, layer *C.lh_background_layer, gradient *C.lh_conic_gradient) {
	if c := lookupContainer(uintptr(goHandle)); c != nil {
		g := ConicGradient{
			Position:    PointF{float32(gradient.position.x), float32(gradient.position.y)},
			Angle:       float32(gradient.angle),
			Radius:      float32(gradient.radius),
			ColorPoints: colorPointsFromC(gradient.color_points, gradient.color_points_count),
		}
		c.DrawConicGradient(uintptr(hdc), backgroundLayerFromC(layer), g)
	}
}

//export goDrawBorders
func goDrawBorders(goHandle C.uintptr_t, hdc C.uintptr_t, borders *C.lh_borders, drawPos C.lh_position, root C.int) {
	if c := lookupContainer(uintptr(goHandle)); c != nil {
		c.DrawBorders(uintptr(hdc), bordersFromC(borders), positionFromC(drawPos), root != 0)
	}
}

//export goSetCaption
func goSetCaption(goHandle C.uintptr_t, caption *C.char) {
	if c := lookupContainer(uintptr(goHandle)); c != nil {
		c.SetCaption(C.GoString(caption))
	}
}

//export goSetBaseURL
func goSetBaseURL(goHandle C.uintptr_t, baseURL *C.char) {
	if c := lookupContainer(uintptr(goHandle)); c != nil {
		c.SetBaseURL(C.GoString(baseURL))
	}
}

//export goLink
func goLink(goHandle C.uintptr_t, href *C.char, rel *C.char, typ *C.char) {
	if c := lookupContainer(uintptr(goHandle)); c != nil {
		c.Link(C.GoString(href), C.GoString(rel), C.GoString(typ))
	}
}

//export goOnAnchorClick
func goOnAnchorClick(goHandle C.uintptr_t, url *C.char) {
	if c := lookupContainer(uintptr(goHandle)); c != nil {
		c.OnAnchorClick(C.GoString(url))
	}
}

//export goOnMouseEvent
func goOnMouseEvent(goHandle C.uintptr_t, event C.int) {
	if c := lookupContainer(uintptr(goHandle)); c != nil {
		c.OnMouseEvent(MouseEvent(event))
	}
}

//export goSetCursor
func goSetCursor(goHandle C.uintptr_t, cursor *C.char) {
	if c := lookupContainer(uintptr(goHandle)); c != nil {
		c.SetCursor(C.GoString(cursor))
	}
}

// goTransformTextBuf holds the return value so C gets a valid pointer.
var goTransformTextBuf *C.char

//export goTransformText
func goTransformText(goHandle C.uintptr_t, text *C.char, tt C.int) *C.char {
	if c := lookupContainer(uintptr(goHandle)); c != nil {
		result := c.TransformText(C.GoString(text), TextTransform(tt))
		if goTransformTextBuf != nil {
			C.free(unsafe.Pointer(goTransformTextBuf))
		}
		goTransformTextBuf = C.CString(result)
		return goTransformTextBuf
	}
	return text
}

//export goImportCSS
func goImportCSS(goHandle C.uintptr_t, url *C.char, baseurl *C.char, result *C.lh_import_css_result) {
	if c := lookupContainer(uintptr(goHandle)); c != nil {
		text, newBaseURL := c.ImportCSS(C.GoString(url), C.GoString(baseurl))
		result.text = C.CString(text)
		result.baseurl = C.CString(newBaseURL)
	}
}

//export goSetClip
func goSetClip(goHandle C.uintptr_t, pos C.lh_position, bdrRadius C.lh_border_radiuses) {
	if c := lookupContainer(uintptr(goHandle)); c != nil {
		c.SetClip(positionFromC(pos), borderRadiusesFromC(bdrRadius))
	}
}

//export goDelClip
func goDelClip(goHandle C.uintptr_t) {
	if c := lookupContainer(uintptr(goHandle)); c != nil {
		c.DelClip()
	}
}

//export goGetViewport
func goGetViewport(goHandle C.uintptr_t, viewport *C.lh_position) {
	if c := lookupContainer(uintptr(goHandle)); c != nil {
		v := c.GetViewport()
		viewport.x = C.float(v.X)
		viewport.y = C.float(v.Y)
		viewport.width = C.float(v.Width)
		viewport.height = C.float(v.Height)
	}
}

//export goCreateElement
func goCreateElement(goHandle C.uintptr_t, tagName *C.char) C.uintptr_t {
	if c := lookupContainer(uintptr(goHandle)); c != nil {
		return C.uintptr_t(c.CreateElement(C.GoString(tagName), nil))
	}
	return 0
}

//export goGetMediaFeatures
func goGetMediaFeatures(goHandle C.uintptr_t, media *C.lh_media_features) {
	if c := lookupContainer(uintptr(goHandle)); c != nil {
		mf := c.GetMediaFeatures()
		media._type = C.int(mf.Type)
		media.width = C.float(mf.Width)
		media.height = C.float(mf.Height)
		media.device_width = C.float(mf.DeviceWidth)
		media.device_height = C.float(mf.DeviceHeight)
		media.color = C.int(mf.Color)
		media.color_index = C.int(mf.ColorIndex)
		media.monochrome = C.int(mf.Monochrome)
		media.resolution = C.float(mf.Resolution)
	}
}

//export goGetLanguage
func goGetLanguage(goHandle C.uintptr_t, result *C.lh_language_result) {
	if c := lookupContainer(uintptr(goHandle)); c != nil {
		lang, culture := c.GetLanguage()
		result.language = C.CString(lang)
		result.culture = C.CString(culture)
	}
}
