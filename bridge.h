#ifndef LITEHTML_GO_BRIDGE_H
#define LITEHTML_GO_BRIDGE_H

#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

/* ── Flat C structs for the cgo boundary ── */

typedef struct {
    float x, y, width, height;
} lh_position;

typedef struct {
    float width, height;
} lh_size;

typedef struct {
    float x, y;
} lh_pointf;

typedef struct {
    uint8_t red, green, blue, alpha;
} lh_web_color;

typedef struct {
    float font_size;
    float height;
    float ascent;
    float descent;
    float x_height;
    float ch_width;
    int   draw_spaces;
    float sub_shift;
    float super_shift;
} lh_font_metrics;

typedef struct {
    const char* family;
    float       size;
    int         style;
    int         weight;
    int         decoration_line;
    float       decoration_thickness;
    int         decoration_style;
    lh_web_color decoration_color;
    const char* emphasis_style;
    lh_web_color emphasis_color;
    int         emphasis_position;
} lh_font_description;

typedef struct {
    float        width;
    int          style;
    lh_web_color color;
} lh_border;

typedef struct {
    float top_left_x,     top_left_y;
    float top_right_x,    top_right_y;
    float bottom_right_x, bottom_right_y;
    float bottom_left_x,  bottom_left_y;
} lh_border_radiuses;

typedef struct {
    lh_border         left, top, right, bottom;
    lh_border_radiuses radius;
} lh_borders;

typedef struct {
    lh_position        border_box;
    lh_border_radiuses border_radius;
    lh_position        clip_box;
    lh_position        origin_box;
    int                attachment;
    int                repeat;
    int                is_root;
} lh_background_layer;

typedef struct {
    float        offset;
    lh_web_color color;
} lh_color_point;

typedef struct {
    lh_pointf      start;
    lh_pointf      end;
    lh_color_point* color_points;
    int             color_points_count;
} lh_linear_gradient;

typedef struct {
    lh_pointf      position;
    lh_pointf      radius;
    lh_color_point* color_points;
    int             color_points_count;
} lh_radial_gradient;

typedef struct {
    lh_pointf      position;
    float          angle;
    float          radius;
    lh_color_point* color_points;
    int             color_points_count;
} lh_conic_gradient;

typedef struct {
    const char*  image;
    const char*  baseurl;
    int          marker_type;
    lh_web_color color;
    lh_position  pos;
    int          index;
    uintptr_t    font;
} lh_list_marker;

typedef struct {
    int   _type;
    float width;
    float height;
    float device_width;
    float device_height;
    int   color;
    int   color_index;
    int   monochrome;
    float resolution;
} lh_media_features;

/* ── Callback result structs ── */

typedef struct {
    char* text;
    char* baseurl;
} lh_import_css_result;

typedef struct {
    char* language;
    char* culture;
} lh_language_result;

/* ── Container callbacks ── */

typedef struct {
    uintptr_t   (*create_font)(uintptr_t go_handle, lh_font_description* descr, lh_font_metrics* fm);
    void        (*delete_font)(uintptr_t go_handle, uintptr_t hFont);
    float       (*text_width)(uintptr_t go_handle, const char* text, uintptr_t hFont);
    void        (*draw_text)(uintptr_t go_handle, uintptr_t hdc, const char* text, uintptr_t hFont, lh_web_color color, lh_position pos);
    float       (*pt_to_px)(uintptr_t go_handle, float pt);
    float       (*get_default_font_size)(uintptr_t go_handle);
    const char* (*get_default_font_name)(uintptr_t go_handle);
    void        (*draw_list_marker)(uintptr_t go_handle, uintptr_t hdc, lh_list_marker* marker);
    void        (*load_image)(uintptr_t go_handle, const char* src, const char* baseurl, int redraw_on_ready);
    void        (*get_image_size)(uintptr_t go_handle, const char* src, const char* baseurl, lh_size* sz);
    void        (*draw_image)(uintptr_t go_handle, uintptr_t hdc, lh_background_layer* layer, const char* url, const char* base_url);
    void        (*draw_solid_fill)(uintptr_t go_handle, uintptr_t hdc, lh_background_layer* layer, lh_web_color color);
    void        (*draw_linear_gradient)(uintptr_t go_handle, uintptr_t hdc, lh_background_layer* layer, lh_linear_gradient* gradient);
    void        (*draw_radial_gradient)(uintptr_t go_handle, uintptr_t hdc, lh_background_layer* layer, lh_radial_gradient* gradient);
    void        (*draw_conic_gradient)(uintptr_t go_handle, uintptr_t hdc, lh_background_layer* layer, lh_conic_gradient* gradient);
    void        (*draw_borders)(uintptr_t go_handle, uintptr_t hdc, lh_borders* borders, lh_position draw_pos, int root);
    void        (*set_caption)(uintptr_t go_handle, const char* caption);
    void        (*set_base_url)(uintptr_t go_handle, const char* base_url);
    void        (*link)(uintptr_t go_handle, const char* href, const char* rel, const char* type);
    void        (*on_anchor_click)(uintptr_t go_handle, const char* url);
    void        (*on_mouse_event)(uintptr_t go_handle, int event);
    void        (*set_cursor)(uintptr_t go_handle, const char* cursor);
    const char* (*transform_text)(uintptr_t go_handle, const char* text, int tt);
    void        (*import_css)(uintptr_t go_handle, const char* url, const char* baseurl, lh_import_css_result* result);
    void        (*set_clip)(uintptr_t go_handle, lh_position pos, lh_border_radiuses bdr_radius);
    void        (*del_clip)(uintptr_t go_handle);
    void        (*get_viewport)(uintptr_t go_handle, lh_position* viewport);
    uintptr_t   (*create_element)(uintptr_t go_handle, const char* tag_name);
    void        (*get_media_features)(uintptr_t go_handle, lh_media_features* media);
    void        (*get_language)(uintptr_t go_handle, lh_language_result* result);
} lh_container_callbacks;

/* ── Document API ── */

uintptr_t lh_document_create_from_string(
    const char* html,
    uintptr_t go_handle,
    lh_container_callbacks* cb,
    const char* master_css,
    const char* user_css);

void  lh_document_destroy(uintptr_t doc_handle);
float lh_document_render(uintptr_t doc_handle, float max_width);
void  lh_document_draw(uintptr_t doc_handle, uintptr_t hdc, float x, float y, lh_position* clip);
float lh_document_width(uintptr_t doc_handle);
float lh_document_height(uintptr_t doc_handle);

int lh_document_on_mouse_over(uintptr_t doc_handle, float x, float y, float client_x, float client_y, lh_position* redraw_boxes, int max_boxes, int* redraw_count);
int lh_document_on_lbutton_down(uintptr_t doc_handle, float x, float y, float client_x, float client_y, lh_position* redraw_boxes, int max_boxes, int* redraw_count);
int lh_document_on_lbutton_up(uintptr_t doc_handle, float x, float y, float client_x, float client_y, lh_position* redraw_boxes, int max_boxes, int* redraw_count);
int lh_document_on_mouse_leave(uintptr_t doc_handle, lh_position* redraw_boxes, int max_boxes, int* redraw_count);

/* Get the default master CSS string from litehtml */
const char* lh_master_css();

#ifdef __cplusplus
}
#endif

#endif /* LITEHTML_GO_BRIDGE_H */
