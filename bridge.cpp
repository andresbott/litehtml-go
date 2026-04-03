/*
 * litehtml-go C++ bridge.
 *
 * Unity build: all litehtml sources are compiled as part of this file.
 * Then the go_container class and extern "C" API are defined.
 */

/* ── Unity build: litehtml C++ sources ── */

/*
 * Unity build ordering notes:
 * - css_parser.cpp must come before css_selector.cpp (template specialization order)
 * - internal.h defines `#define in /` (operator overload trick). Three files
 *   use this macro: css_selector.cpp, html_tag.cpp, style.cpp. They are grouped
 *   together and `#undef in` is placed after the last one. All other files must
 *   come before or after this group.
 */

/* Phase 1: files that do NOT use the `in` macro */
#include "litehtml/src/codepoint.cpp"
#include "litehtml/src/css_length.cpp"
#include "litehtml/src/css_tokenizer.cpp"
#include "litehtml/src/css_parser.cpp"
#include "litehtml/src/document.cpp"
#include "litehtml/src/document_container.cpp"
#include "litehtml/src/el_anchor.cpp"
#include "litehtml/src/el_base.cpp"
#include "litehtml/src/el_before_after.cpp"
#include "litehtml/src/el_body.cpp"
#include "litehtml/src/el_break.cpp"
#include "litehtml/src/el_cdata.cpp"
#include "litehtml/src/el_comment.cpp"
#include "litehtml/src/el_div.cpp"
#include "litehtml/src/element.cpp"
#include "litehtml/src/el_font.cpp"
#include "litehtml/src/el_image.cpp"
#include "litehtml/src/el_link.cpp"
#include "litehtml/src/el_para.cpp"
#include "litehtml/src/el_script.cpp"
#include "litehtml/src/el_space.cpp"
#include "litehtml/src/el_style.cpp"
#include "litehtml/src/el_table.cpp"
#include "litehtml/src/el_td.cpp"
#include "litehtml/src/el_text.cpp"
#include "litehtml/src/el_title.cpp"
#include "litehtml/src/el_tr.cpp"
#include "litehtml/src/encodings.cpp"
#undef out
#undef inout
#undef countof
#include "litehtml/src/html.cpp"
#include "litehtml/src/html_microsyntaxes.cpp"
#include "litehtml/src/iterators.cpp"
#include "litehtml/src/media_query.cpp"
#include "litehtml/src/stylesheet.cpp"
#include "litehtml/src/table.cpp"
#include "litehtml/src/tstring_view.cpp"
#include "litehtml/src/url.cpp"
#include "litehtml/src/url_path.cpp"
#include "litehtml/src/utf8_strings.cpp"
#include "litehtml/src/web_color.cpp"
#include "litehtml/src/num_cvt.cpp"
#include "litehtml/src/strtod.cpp"
#include "litehtml/src/string_id.cpp"
#include "litehtml/src/css_properties.cpp"
#include "litehtml/src/line_box.cpp"
#include "litehtml/src/css_borders.cpp"
#include "litehtml/src/render_item.cpp"
#include "litehtml/src/render_block_context.cpp"
#include "litehtml/src/render_block.cpp"
#include "litehtml/src/render_inline_context.cpp"
#include "litehtml/src/render_table.cpp"
#include "litehtml/src/render_flex.cpp"
#include "litehtml/src/render_image.cpp"
#include "litehtml/src/formatting_context.cpp"
#include "litehtml/src/flex_item.cpp"
#include "litehtml/src/flex_line.cpp"
#include "litehtml/src/background.cpp"
#include "litehtml/src/gradient.cpp"

/* Phase 2: files that use the `in` macro from internal.h */
#include "litehtml/src/css_selector.cpp"
#include "litehtml/src/html_tag.cpp"
#include "litehtml/src/style.cpp"
#undef in

/* ── Bridge code ── */

#include "bridge.h"
#include "litehtml/include/litehtml/document.h"
#include "litehtml/include/litehtml/document_container.h"
#include "litehtml/include/litehtml/master_css.h"
#include <string>
#include <cstring>
#include <vector>

/* ── Helper: convert C++ types to C flat structs and back ── */

static lh_position to_lh_position(const litehtml::position& p) {
    return {p.x, p.y, p.width, p.height};
}

static litehtml::position from_lh_position(const lh_position& p) {
    return litehtml::position(p.x, p.y, p.width, p.height);
}

static lh_web_color to_lh_web_color(const litehtml::web_color& c) {
    return {c.red, c.green, c.blue, c.alpha};
}

static litehtml::web_color from_lh_web_color(const lh_web_color& c) {
    return litehtml::web_color(c.red, c.green, c.blue, c.alpha);
}

static lh_border_radiuses to_lh_border_radiuses(const litehtml::border_radiuses& r) {
    return {
        r.top_left_x, r.top_left_y,
        r.top_right_x, r.top_right_y,
        r.bottom_right_x, r.bottom_right_y,
        r.bottom_left_x, r.bottom_left_y
    };
}

static lh_border to_lh_border(const litehtml::border& b) {
    return {b.width, (int)b.style, to_lh_web_color(b.color)};
}

static lh_borders to_lh_borders(const litehtml::borders& b) {
    lh_borders r;
    r.left = to_lh_border(b.left);
    r.top = to_lh_border(b.top);
    r.right = to_lh_border(b.right);
    r.bottom = to_lh_border(b.bottom);
    r.radius = to_lh_border_radiuses(b.radius);
    return r;
}

static lh_background_layer to_lh_background_layer(const litehtml::background_layer& l) {
    lh_background_layer r;
    r.border_box = to_lh_position(l.border_box);
    r.border_radius = to_lh_border_radiuses(l.border_radius);
    r.clip_box = to_lh_position(l.clip_box);
    r.origin_box = to_lh_position(l.origin_box);
    r.attachment = (int)l.attachment;
    r.repeat = (int)l.repeat;
    r.is_root = l.is_root ? 1 : 0;
    return r;
}

static lh_font_description to_lh_font_description(const litehtml::font_description& d) {
    lh_font_description r;
    r.family = d.family.c_str();
    r.size = d.size;
    r.style = (int)d.style;
    r.weight = d.weight;
    r.decoration_line = d.decoration_line;
    r.decoration_thickness = d.decoration_thickness.val();
    r.decoration_style = (int)d.decoration_style;
    r.decoration_color = to_lh_web_color(d.decoration_color);
    r.emphasis_style = d.emphasis_style.c_str();
    r.emphasis_color = to_lh_web_color(d.emphasis_color);
    r.emphasis_position = d.emphasis_position;
    return r;
}

static std::vector<lh_color_point> to_lh_color_points(
    const std::vector<litehtml::background_layer::color_point>& pts) {
    std::vector<lh_color_point> out;
    out.reserve(pts.size());
    for (const auto& pt : pts) {
        out.push_back({pt.offset, to_lh_web_color(pt.color)});
    }
    return out;
}

/* ── go_container: bridges C++ virtual calls to C function pointers ── */

class go_container : public litehtml::document_container {
    uintptr_t m_go_handle;
    lh_container_callbacks m_cb;
    std::string m_default_font_name_buf;
    std::string m_transform_text_buf;

public:
    go_container(uintptr_t go_handle, lh_container_callbacks* cb)
        : m_go_handle(go_handle), m_cb(*cb) {}

    litehtml::uint_ptr create_font(const litehtml::font_description& descr,
                                    const litehtml::document* /*doc*/,
                                    litehtml::font_metrics* fm) override {
        lh_font_description cd = to_lh_font_description(descr);
        lh_font_metrics cfm = {};
        uintptr_t result = m_cb.create_font(m_go_handle, &cd, &cfm);
        if (fm) {
            fm->font_size = cfm.font_size;
            fm->height = cfm.height;
            fm->ascent = cfm.ascent;
            fm->descent = cfm.descent;
            fm->x_height = cfm.x_height;
            fm->ch_width = cfm.ch_width;
            fm->draw_spaces = cfm.draw_spaces != 0;
            fm->sub_shift = cfm.sub_shift;
            fm->super_shift = cfm.super_shift;
        }
        return result;
    }

    void delete_font(litehtml::uint_ptr hFont) override {
        m_cb.delete_font(m_go_handle, hFont);
    }

    litehtml::pixel_t text_width(const char* text, litehtml::uint_ptr hFont) override {
        return m_cb.text_width(m_go_handle, text, hFont);
    }

    void draw_text(litehtml::uint_ptr hdc, const char* text, litehtml::uint_ptr hFont,
                   litehtml::web_color color, const litehtml::position& pos) override {
        m_cb.draw_text(m_go_handle, hdc, text, hFont, to_lh_web_color(color), to_lh_position(pos));
    }

    litehtml::pixel_t pt_to_px(float pt) const override {
        return m_cb.pt_to_px(m_go_handle, pt);
    }

    litehtml::pixel_t get_default_font_size() const override {
        return m_cb.get_default_font_size(m_go_handle);
    }

    const char* get_default_font_name() const override {
        const char* name = m_cb.get_default_font_name(m_go_handle);
        if (name) {
            const_cast<go_container*>(this)->m_default_font_name_buf = name;
            return m_default_font_name_buf.c_str();
        }
        return "";
    }

    void draw_list_marker(litehtml::uint_ptr hdc, const litehtml::list_marker& marker) override {
        lh_list_marker cm;
        cm.image = marker.image.c_str();
        cm.baseurl = marker.baseurl ? marker.baseurl : "";
        cm.marker_type = (int)marker.marker_type;
        cm.color = to_lh_web_color(marker.color);
        cm.pos = to_lh_position(marker.pos);
        cm.index = marker.index;
        cm.font = marker.font;
        m_cb.draw_list_marker(m_go_handle, hdc, &cm);
    }

    void load_image(const char* src, const char* baseurl, bool redraw_on_ready) override {
        m_cb.load_image(m_go_handle, src, baseurl ? baseurl : "", redraw_on_ready ? 1 : 0);
    }

    void get_image_size(const char* src, const char* baseurl, litehtml::size& sz) override {
        lh_size csz = {0, 0};
        m_cb.get_image_size(m_go_handle, src, baseurl ? baseurl : "", &csz);
        sz.width = csz.width;
        sz.height = csz.height;
    }

    void draw_image(litehtml::uint_ptr hdc, const litehtml::background_layer& layer,
                    const std::string& url, const std::string& base_url) override {
        lh_background_layer cl = to_lh_background_layer(layer);
        m_cb.draw_image(m_go_handle, hdc, &cl, url.c_str(), base_url.c_str());
    }

    void draw_solid_fill(litehtml::uint_ptr hdc, const litehtml::background_layer& layer,
                         const litehtml::web_color& color) override {
        lh_background_layer cl = to_lh_background_layer(layer);
        m_cb.draw_solid_fill(m_go_handle, hdc, &cl, to_lh_web_color(color));
    }

    void draw_linear_gradient(litehtml::uint_ptr hdc, const litehtml::background_layer& layer,
                              const litehtml::background_layer::linear_gradient& gradient) override {
        lh_background_layer cl = to_lh_background_layer(layer);
        auto pts = to_lh_color_points(gradient.color_points);
        lh_linear_gradient cg;
        cg.start = {gradient.start.x, gradient.start.y};
        cg.end = {gradient.end.x, gradient.end.y};
        cg.color_points = pts.data();
        cg.color_points_count = (int)pts.size();
        m_cb.draw_linear_gradient(m_go_handle, hdc, &cl, &cg);
    }

    void draw_radial_gradient(litehtml::uint_ptr hdc, const litehtml::background_layer& layer,
                              const litehtml::background_layer::radial_gradient& gradient) override {
        lh_background_layer cl = to_lh_background_layer(layer);
        auto pts = to_lh_color_points(gradient.color_points);
        lh_radial_gradient cg;
        cg.position = {gradient.position.x, gradient.position.y};
        cg.radius = {gradient.radius.x, gradient.radius.y};
        cg.color_points = pts.data();
        cg.color_points_count = (int)pts.size();
        m_cb.draw_radial_gradient(m_go_handle, hdc, &cl, &cg);
    }

    void draw_conic_gradient(litehtml::uint_ptr hdc, const litehtml::background_layer& layer,
                             const litehtml::background_layer::conic_gradient& gradient) override {
        lh_background_layer cl = to_lh_background_layer(layer);
        auto pts = to_lh_color_points(gradient.color_points);
        lh_conic_gradient cg;
        cg.position = {gradient.position.x, gradient.position.y};
        cg.angle = gradient.angle;
        cg.radius = gradient.radius;
        cg.color_points = pts.data();
        cg.color_points_count = (int)pts.size();
        m_cb.draw_conic_gradient(m_go_handle, hdc, &cl, &cg);
    }

    void draw_borders(litehtml::uint_ptr hdc, const litehtml::borders& borders,
                      const litehtml::position& draw_pos, bool root) override {
        lh_borders cb = to_lh_borders(borders);
        m_cb.draw_borders(m_go_handle, hdc, &cb, to_lh_position(draw_pos), root ? 1 : 0);
    }

    void set_caption(const char* caption) override {
        m_cb.set_caption(m_go_handle, caption);
    }

    void set_base_url(const char* base_url) override {
        m_cb.set_base_url(m_go_handle, base_url);
    }

    void link(const std::shared_ptr<litehtml::document>& /*doc*/,
              const litehtml::element::ptr& el) override {
        // Extract href, rel, type attributes from the element
        // For now pass empty strings if attributes don't exist
        m_cb.link(m_go_handle, "", "", "");
    }

    void on_anchor_click(const char* url, const litehtml::element::ptr& /*el*/) override {
        m_cb.on_anchor_click(m_go_handle, url);
    }

    void on_mouse_event(const litehtml::element::ptr& /*el*/, litehtml::mouse_event event) override {
        m_cb.on_mouse_event(m_go_handle, (int)event);
    }

    void set_cursor(const char* cursor) override {
        m_cb.set_cursor(m_go_handle, cursor);
    }

    void transform_text(litehtml::string& text, litehtml::text_transform tt) override {
        const char* result = m_cb.transform_text(m_go_handle, text.c_str(), (int)tt);
        if (result) {
            text = result;
        }
    }

    void import_css(litehtml::string& text, const litehtml::string& url,
                    litehtml::string& baseurl) override {
        lh_import_css_result result = {nullptr, nullptr};
        m_cb.import_css(m_go_handle, url.c_str(), baseurl.c_str(), &result);
        if (result.text) {
            text = result.text;
            free(result.text);
        }
        if (result.baseurl) {
            baseurl = result.baseurl;
            free(result.baseurl);
        }
    }

    void set_clip(const litehtml::position& pos, const litehtml::border_radiuses& bdr_radius) override {
        m_cb.set_clip(m_go_handle, to_lh_position(pos), to_lh_border_radiuses(bdr_radius));
    }

    void del_clip() override {
        m_cb.del_clip(m_go_handle);
    }

    void get_viewport(litehtml::position& viewport) const override {
        lh_position cv = {};
        m_cb.get_viewport(m_go_handle, &cv);
        viewport = from_lh_position(cv);
    }

    litehtml::element::ptr create_element(const char* tag_name,
                                           const litehtml::string_map& /*attributes*/,
                                           const std::shared_ptr<litehtml::document>& /*doc*/) override {
        uintptr_t result = m_cb.create_element(m_go_handle, tag_name);
        (void)result; // For now, always return nullptr to use default elements
        return nullptr;
    }

    void get_media_features(litehtml::media_features& media) const override {
        lh_media_features cm = {};
        m_cb.get_media_features(m_go_handle, &cm);
        media.type = (litehtml::media_type)cm._type;
        media.width = cm.width;
        media.height = cm.height;
        media.device_width = cm.device_width;
        media.device_height = cm.device_height;
        media.color = cm.color;
        media.color_index = cm.color_index;
        media.monochrome = cm.monochrome;
        media.resolution = cm.resolution;
    }

    void get_language(litehtml::string& language, litehtml::string& culture) const override {
        lh_language_result result = {nullptr, nullptr};
        m_cb.get_language(m_go_handle, &result);
        if (result.language) {
            language = result.language;
            free(result.language);
        }
        if (result.culture) {
            culture = result.culture;
            free(result.culture);
        }
    }
};

/* ── Data holder for a document + its container ── */

struct doc_holder {
    go_container* container;
    litehtml::document::ptr doc;
};

/* ── Extern "C" API implementations ── */

extern "C" {

uintptr_t lh_document_create_from_string(
    const char* html,
    uintptr_t go_handle,
    lh_container_callbacks* cb,
    const char* master_css_str,
    const char* user_css_str)
{
    auto* container = new go_container(go_handle, cb);
    std::string master = master_css_str ? master_css_str : litehtml::master_css;
    std::string user = user_css_str ? user_css_str : "";

    litehtml::document::ptr doc = litehtml::document::createFromString(html, container, master, user);
    if (!doc) {
        delete container;
        return 0;
    }

    auto* holder = new doc_holder{container, doc};
    return reinterpret_cast<uintptr_t>(holder);
}

void lh_document_destroy(uintptr_t doc_handle) {
    if (!doc_handle) return;
    auto* holder = reinterpret_cast<doc_holder*>(doc_handle);
    holder->doc.reset();
    delete holder->container;
    delete holder;
}

float lh_document_render(uintptr_t doc_handle, float max_width) {
    if (!doc_handle) return 0;
    auto* holder = reinterpret_cast<doc_holder*>(doc_handle);
    return holder->doc->render((litehtml::pixel_t)max_width);
}

void lh_document_draw(uintptr_t doc_handle, uintptr_t hdc, float x, float y, lh_position* clip) {
    if (!doc_handle) return;
    auto* holder = reinterpret_cast<doc_holder*>(doc_handle);
    if (clip) {
        litehtml::position p = from_lh_position(*clip);
        holder->doc->draw(hdc, (litehtml::pixel_t)x, (litehtml::pixel_t)y, &p);
    } else {
        holder->doc->draw(hdc, (litehtml::pixel_t)x, (litehtml::pixel_t)y, nullptr);
    }
}

float lh_document_width(uintptr_t doc_handle) {
    if (!doc_handle) return 0;
    return reinterpret_cast<doc_holder*>(doc_handle)->doc->width();
}

float lh_document_height(uintptr_t doc_handle) {
    if (!doc_handle) return 0;
    return reinterpret_cast<doc_holder*>(doc_handle)->doc->height();
}

static int mouse_event_helper(litehtml::document::ptr& doc,
    bool (litehtml::document::*fn)(litehtml::pixel_t, litehtml::pixel_t,
                                    litehtml::pixel_t, litehtml::pixel_t,
                                    litehtml::position::vector&),
    float x, float y, float cx, float cy,
    lh_position* boxes, int max_boxes, int* count)
{
    litehtml::position::vector redraw;
    bool result = (doc.get()->*fn)(x, y, cx, cy, redraw);
    int n = (int)redraw.size();
    if (n > max_boxes) n = max_boxes;
    for (int i = 0; i < n; i++) {
        boxes[i] = to_lh_position(redraw[i]);
    }
    *count = (int)redraw.size();
    return result ? 1 : 0;
}

int lh_document_on_mouse_over(uintptr_t doc_handle, float x, float y, float cx, float cy,
                               lh_position* boxes, int max_boxes, int* count) {
    if (!doc_handle) { *count = 0; return 0; }
    auto* h = reinterpret_cast<doc_holder*>(doc_handle);
    return mouse_event_helper(h->doc, &litehtml::document::on_mouse_over, x, y, cx, cy, boxes, max_boxes, count);
}

int lh_document_on_lbutton_down(uintptr_t doc_handle, float x, float y, float cx, float cy,
                                 lh_position* boxes, int max_boxes, int* count) {
    if (!doc_handle) { *count = 0; return 0; }
    auto* h = reinterpret_cast<doc_holder*>(doc_handle);
    return mouse_event_helper(h->doc, &litehtml::document::on_lbutton_down, x, y, cx, cy, boxes, max_boxes, count);
}

int lh_document_on_lbutton_up(uintptr_t doc_handle, float x, float y, float cx, float cy,
                               lh_position* boxes, int max_boxes, int* count) {
    if (!doc_handle) { *count = 0; return 0; }
    auto* h = reinterpret_cast<doc_holder*>(doc_handle);
    return mouse_event_helper(h->doc, &litehtml::document::on_lbutton_up, x, y, cx, cy, boxes, max_boxes, count);
}

int lh_document_on_mouse_leave(uintptr_t doc_handle, lh_position* boxes, int max_boxes, int* count) {
    if (!doc_handle) { *count = 0; return 0; }
    auto* h = reinterpret_cast<doc_holder*>(doc_handle);
    litehtml::position::vector redraw;
    bool result = h->doc->on_mouse_leave(redraw);
    int n = (int)redraw.size();
    if (n > max_boxes) n = max_boxes;
    for (int i = 0; i < n; i++) {
        boxes[i] = to_lh_position(redraw[i]);
    }
    *count = (int)redraw.size();
    return result ? 1 : 0;
}

const char* lh_master_css() {
    return litehtml::master_css;
}

} /* extern "C" */
