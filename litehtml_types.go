package litehtml

// Position represents a rectangle with x, y, width, height.
type Position struct {
	X, Y, Width, Height float32
}

// Size represents a width/height pair.
type Size struct {
	Width, Height float32
}

// PointF represents a 2D floating-point coordinate.
type PointF struct {
	X, Y float32
}

// WebColor represents an RGBA color.
type WebColor struct {
	Red, Green, Blue, Alpha uint8
}

// FontStyle mirrors litehtml::font_style.
type FontStyle int

const (
	FontStyleNormal FontStyle = iota
	FontStyleItalic
)

// TextDecorationStyle mirrors litehtml::text_decoration_style.
type TextDecorationStyle int

const (
	TextDecorationStyleSolid TextDecorationStyle = iota
	TextDecorationStyleDouble
	TextDecorationStyleDotted
	TextDecorationStyleDashed
	TextDecorationStyleWavy
)

// TextTransform mirrors litehtml::text_transform.
type TextTransform int

const (
	TextTransformNone TextTransform = iota
	TextTransformCapitalize
	TextTransformUppercase
	TextTransformLowercase
)

// BorderStyle mirrors litehtml::border_style.
type BorderStyle int

const (
	BorderStyleNone BorderStyle = iota
	BorderStyleHidden
	BorderStyleDotted
	BorderStyleDashed
	BorderStyleSolid
	BorderStyleDouble
	BorderStyleGroove
	BorderStyleRidge
	BorderStyleInset
	BorderStyleOutset
)

// BackgroundAttachment mirrors litehtml::background_attachment.
type BackgroundAttachment int

const (
	BackgroundAttachmentScroll BackgroundAttachment = iota
	BackgroundAttachmentFixed
)

// BackgroundRepeat mirrors litehtml::background_repeat.
type BackgroundRepeat int

const (
	BackgroundRepeatRepeat BackgroundRepeat = iota
	BackgroundRepeatRepeatX
	BackgroundRepeatRepeatY
	BackgroundRepeatNoRepeat
)

// ListStyleType mirrors litehtml::list_style_type.
type ListStyleType int

const (
	ListStyleTypeNone ListStyleType = iota
	ListStyleTypeCircle
	ListStyleTypeDisc
	ListStyleTypeSquare
	ListStyleTypeArmenian
	ListStyleTypeCjkIdeographic
	ListStyleTypeDecimal
	ListStyleTypeDecimalLeadingZero
	ListStyleTypeGeorgian
	ListStyleTypeHebrew
	ListStyleTypeHiragana
	ListStyleTypeHiraganaIroha
	ListStyleTypeKatakana
	ListStyleTypeKatakanaIroha
	ListStyleTypeLowerAlpha
	ListStyleTypeLowerGreek
	ListStyleTypeLowerLatin
	ListStyleTypeLowerRoman
	ListStyleTypeUpperAlpha
	ListStyleTypeUpperLatin
	ListStyleTypeUpperRoman
)

// MediaType mirrors litehtml::media_type.
type MediaType int

const (
	MediaTypeUnknown MediaType = iota
	MediaTypeAll
	MediaTypePrint
	MediaTypeScreen
)

// MouseEvent mirrors litehtml::mouse_event.
type MouseEvent int

const (
	MouseEventEnter MouseEvent = iota
	MouseEventLeave
)

// FontMetrics holds font measurement data returned by CreateFont.
type FontMetrics struct {
	FontSize   float32
	Height     float32
	Ascent     float32
	Descent    float32
	XHeight    float32
	ChWidth    float32
	DrawSpaces bool
	SubShift   float32
	SuperShift float32
}

// FontDescription describes a font request from the layout engine.
type FontDescription struct {
	Family              string
	Size                float32
	Style               FontStyle
	Weight              int
	DecorationLine      int
	DecorationThickness float32
	DecorationStyle     TextDecorationStyle
	DecorationColor     WebColor
	EmphasisStyle       string
	EmphasisColor       WebColor
	EmphasisPosition    int
}

// Border describes one side of a border.
type Border struct {
	Width float32
	Style BorderStyle
	Color WebColor
}

// BorderRadiuses holds the 4-corner border radius values.
type BorderRadiuses struct {
	TopLeftX, TopLeftY         float32
	TopRightX, TopRightY       float32
	BottomRightX, BottomRightY float32
	BottomLeftX, BottomLeftY   float32
}

// Borders holds all four border sides plus radiuses.
type Borders struct {
	Left, Top, Right, Bottom Border
	Radius                   BorderRadiuses
}

// BackgroundLayer holds geometry for a background draw call.
type BackgroundLayer struct {
	BorderBox    Position
	BorderRadius BorderRadiuses
	ClipBox      Position
	OriginBox    Position
	Attachment   BackgroundAttachment
	Repeat       BackgroundRepeat
	IsRoot       bool
}

// ColorPoint is a single color stop in a gradient.
type ColorPoint struct {
	Offset float32
	Color  WebColor
}

// LinearGradient holds data for a linear gradient draw call.
type LinearGradient struct {
	Start, End  PointF
	ColorPoints []ColorPoint
}

// RadialGradient holds data for a radial gradient draw call.
type RadialGradient struct {
	Position    PointF
	Radius      PointF
	ColorPoints []ColorPoint
}

// ConicGradient holds data for a conic gradient draw call.
type ConicGradient struct {
	Position    PointF
	Angle       float32
	Radius      float32
	ColorPoints []ColorPoint
}

// ListMarker holds data for a list marker draw call.
type ListMarker struct {
	Image      string
	BaseURL    string
	MarkerType ListStyleType
	Color      WebColor
	Pos        Position
	Index      int
	Font       uintptr
}

// MediaFeatures describes the media environment.
type MediaFeatures struct {
	Type         MediaType
	Width        float32
	Height       float32
	DeviceWidth  float32
	DeviceHeight float32
	Color        int
	ColorIndex   int
	Monochrome   int
	Resolution   float32
}
