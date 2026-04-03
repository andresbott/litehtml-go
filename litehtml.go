package litehtml

/*
#cgo CXXFLAGS: -std=c++17 -I${SRCDIR}/litehtml/include -I${SRCDIR}/litehtml/include/litehtml -I${SRCDIR}/litehtml/src -I${SRCDIR}/litehtml/src/gumbo/include
#cgo CFLAGS: -I${SRCDIR}/litehtml/src/gumbo/include -I${SRCDIR}/litehtml/src/gumbo/include/gumbo
#cgo LDFLAGS: -lstdc++ -lm
#include "bridge.h"
#include <stdlib.h>

// Function pointer assignment helpers.
// cgo cannot take the address of a Go function directly for struct fields,
// so we use C helper functions.

extern uintptr_t goCreateFont(uintptr_t, lh_font_description*, lh_font_metrics*);
extern void      goDeleteFont(uintptr_t, uintptr_t);
extern float     goTextWidth(uintptr_t, const char*, uintptr_t);
extern void      goDrawText(uintptr_t, uintptr_t, const char*, uintptr_t, lh_web_color, lh_position);
extern float     goPtToPx(uintptr_t, float);
extern float     goGetDefaultFontSize(uintptr_t);
extern const char* goGetDefaultFontName(uintptr_t);
extern void      goDrawListMarker(uintptr_t, uintptr_t, lh_list_marker*);
extern void      goLoadImage(uintptr_t, const char*, const char*, int);
extern void      goGetImageSize(uintptr_t, const char*, const char*, lh_size*);
extern void      goDrawImage(uintptr_t, uintptr_t, lh_background_layer*, const char*, const char*);
extern void      goDrawSolidFill(uintptr_t, uintptr_t, lh_background_layer*, lh_web_color);
extern void      goDrawLinearGradient(uintptr_t, uintptr_t, lh_background_layer*, lh_linear_gradient*);
extern void      goDrawRadialGradient(uintptr_t, uintptr_t, lh_background_layer*, lh_radial_gradient*);
extern void      goDrawConicGradient(uintptr_t, uintptr_t, lh_background_layer*, lh_conic_gradient*);
extern void      goDrawBorders(uintptr_t, uintptr_t, lh_borders*, lh_position, int);
extern void      goSetCaption(uintptr_t, const char*);
extern void      goSetBaseURL(uintptr_t, const char*);
extern void      goLink(uintptr_t, const char*, const char*, const char*);
extern void      goOnAnchorClick(uintptr_t, const char*);
extern void      goOnMouseEvent(uintptr_t, int);
extern void      goSetCursor(uintptr_t, const char*);
extern const char* goTransformText(uintptr_t, const char*, int);
extern void      goImportCSS(uintptr_t, const char*, const char*, lh_import_css_result*);
extern void      goSetClip(uintptr_t, lh_position, lh_border_radiuses);
extern void      goDelClip(uintptr_t);
extern void      goGetViewport(uintptr_t, lh_position*);
extern uintptr_t goCreateElement(uintptr_t, const char*);
extern void      goGetMediaFeatures(uintptr_t, lh_media_features*);
extern void      goGetLanguage(uintptr_t, lh_language_result*);

static lh_container_callbacks make_callbacks() {
    lh_container_callbacks cb;
    cb.create_font = goCreateFont;
    cb.delete_font = goDeleteFont;
    cb.text_width = goTextWidth;
    cb.draw_text = goDrawText;
    cb.pt_to_px = goPtToPx;
    cb.get_default_font_size = goGetDefaultFontSize;
    cb.get_default_font_name = goGetDefaultFontName;
    cb.draw_list_marker = goDrawListMarker;
    cb.load_image = goLoadImage;
    cb.get_image_size = goGetImageSize;
    cb.draw_image = goDrawImage;
    cb.draw_solid_fill = goDrawSolidFill;
    cb.draw_linear_gradient = goDrawLinearGradient;
    cb.draw_radial_gradient = goDrawRadialGradient;
    cb.draw_conic_gradient = goDrawConicGradient;
    cb.draw_borders = goDrawBorders;
    cb.set_caption = goSetCaption;
    cb.set_base_url = goSetBaseURL;
    cb.link = goLink;
    cb.on_anchor_click = goOnAnchorClick;
    cb.on_mouse_event = goOnMouseEvent;
    cb.set_cursor = goSetCursor;
    cb.transform_text = goTransformText;
    cb.import_css = goImportCSS;
    cb.set_clip = goSetClip;
    cb.del_clip = goDelClip;
    cb.get_viewport = goGetViewport;
    cb.create_element = goCreateElement;
    cb.get_media_features = goGetMediaFeatures;
    cb.get_language = goGetLanguage;
    return cb;
}
*/
import "C"
import (
	"errors"
	"runtime"
	"unsafe"
)

// Document represents a parsed and renderable HTML document.
type Document struct {
	handle  C.uintptr_t
	cHandle uintptr
}

// MasterCSS returns the default master CSS stylesheet built into litehtml.
func MasterCSS() string {
	return C.GoString(C.lh_master_css())
}

// NewDocument parses an HTML string and creates a Document.
// The container provides callbacks for fonts, drawing, and resources.
// Pass empty strings for masterCSS/userCSS to use litehtml defaults.
func NewDocument(html string, container DocumentContainer, masterCSS, userCSS string) (*Document, error) {
	cHandle := registerContainer(container)

	cHTML := C.CString(html)
	defer C.free(unsafe.Pointer(cHTML))

	var cMasterCSS *C.char
	if masterCSS != "" {
		cMasterCSS = C.CString(masterCSS)
		defer C.free(unsafe.Pointer(cMasterCSS))
	}

	var cUserCSS *C.char
	if userCSS != "" {
		cUserCSS = C.CString(userCSS)
		defer C.free(unsafe.Pointer(cUserCSS))
	}

	cb := C.make_callbacks()
	docHandle := C.lh_document_create_from_string(cHTML, C.uintptr_t(cHandle), &cb, cMasterCSS, cUserCSS)
	if docHandle == 0 {
		unregisterContainer(cHandle)
		return nil, errors.New("litehtml: failed to create document")
	}

	doc := &Document{
		handle:  docHandle,
		cHandle: cHandle,
	}
	runtime.SetFinalizer(doc, (*Document).Close)
	return doc, nil
}

// Render performs layout at the given maximum width. Returns the actual width used.
func (d *Document) Render(maxWidth float32) float32 {
	if d.handle == 0 {
		return 0
	}
	return float32(C.lh_document_render(d.handle, C.float(maxWidth)))
}

// Draw draws the document at position (x, y). Pass nil for clip to draw everything.
func (d *Document) Draw(hdc uintptr, x, y float32, clip *Position) {
	if d.handle == 0 {
		return
	}
	if clip != nil {
		cClip := C.lh_position{
			x: C.float(clip.X), y: C.float(clip.Y),
			width: C.float(clip.Width), height: C.float(clip.Height),
		}
		C.lh_document_draw(d.handle, C.uintptr_t(hdc), C.float(x), C.float(y), &cClip)
	} else {
		C.lh_document_draw(d.handle, C.uintptr_t(hdc), C.float(x), C.float(y), nil)
	}
}

// Width returns the document width after rendering.
func (d *Document) Width() float32 {
	if d.handle == 0 {
		return 0
	}
	return float32(C.lh_document_width(d.handle))
}

// Height returns the document height after rendering.
func (d *Document) Height() float32 {
	if d.handle == 0 {
		return 0
	}
	return float32(C.lh_document_height(d.handle))
}

const maxRedrawBoxes = 64

// OnMouseOver handles a mouse move event. Returns whether any elements changed and the list of boxes to redraw.
func (d *Document) OnMouseOver(x, y, clientX, clientY float32) (bool, []Position) {
	if d.handle == 0 {
		return false, nil
	}
	var boxes [maxRedrawBoxes]C.lh_position
	var count C.int
	result := C.lh_document_on_mouse_over(d.handle, C.float(x), C.float(y), C.float(clientX), C.float(clientY), &boxes[0], maxRedrawBoxes, &count)
	return result != 0, cPositionsToSlice(boxes[:], int(count))
}

// OnLButtonDown handles a left mouse button press.
func (d *Document) OnLButtonDown(x, y, clientX, clientY float32) (bool, []Position) {
	if d.handle == 0 {
		return false, nil
	}
	var boxes [maxRedrawBoxes]C.lh_position
	var count C.int
	result := C.lh_document_on_lbutton_down(d.handle, C.float(x), C.float(y), C.float(clientX), C.float(clientY), &boxes[0], maxRedrawBoxes, &count)
	return result != 0, cPositionsToSlice(boxes[:], int(count))
}

// OnLButtonUp handles a left mouse button release.
func (d *Document) OnLButtonUp(x, y, clientX, clientY float32) (bool, []Position) {
	if d.handle == 0 {
		return false, nil
	}
	var boxes [maxRedrawBoxes]C.lh_position
	var count C.int
	result := C.lh_document_on_lbutton_up(d.handle, C.float(x), C.float(y), C.float(clientX), C.float(clientY), &boxes[0], maxRedrawBoxes, &count)
	return result != 0, cPositionsToSlice(boxes[:], int(count))
}

// OnMouseLeave handles the mouse leaving the document area.
func (d *Document) OnMouseLeave() (bool, []Position) {
	if d.handle == 0 {
		return false, nil
	}
	var boxes [maxRedrawBoxes]C.lh_position
	var count C.int
	result := C.lh_document_on_mouse_leave(d.handle, &boxes[0], maxRedrawBoxes, &count)
	return result != 0, cPositionsToSlice(boxes[:], int(count))
}

// Close releases the C++ resources. Safe to call multiple times.
func (d *Document) Close() {
	if d.handle != 0 {
		C.lh_document_destroy(d.handle)
		d.handle = 0
		unregisterContainer(d.cHandle)
		runtime.SetFinalizer(d, nil)
	}
}

func cPositionsToSlice(boxes []C.lh_position, count int) []Position {
	if count <= 0 {
		return nil
	}
	if count > len(boxes) {
		count = len(boxes)
	}
	result := make([]Position, count)
	for i := 0; i < count; i++ {
		result[i] = positionFromC(boxes[i])
	}
	return result
}
