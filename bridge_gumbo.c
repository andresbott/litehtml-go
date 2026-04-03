/*
 * Unity build for gumbo HTML parser.
 * cgo compiles .c files in the package directory; this file #includes
 * all gumbo sources from the vendored litehtml submodule.
 */

#include "litehtml/src/gumbo/attribute.c"
#include "litehtml/src/gumbo/char_ref.c"
#include "litehtml/src/gumbo/error.c"
#include "litehtml/src/gumbo/parser.c"
#include "litehtml/src/gumbo/string_buffer.c"
#include "litehtml/src/gumbo/string_piece.c"
#include "litehtml/src/gumbo/tag.c"
#include "litehtml/src/gumbo/tokenizer.c"
#include "litehtml/src/gumbo/utf8.c"
#include "litehtml/src/gumbo/util.c"
#include "litehtml/src/gumbo/vector.c"
