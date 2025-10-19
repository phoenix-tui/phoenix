package model

import (
	"fmt"

	"github.com/phoenix-tui/phoenix/style/domain/value"
)

// Style is a rich domain model that encapsulates all styling concerns for terminal UI elements.
// It follows DDD principles: data + behavior, immutability, and fluent API design.
//
// Style composes multiple value objects:
//   - Color (foreground, background, border color)
//   - Border (border type and sides)
//   - Padding (inner spacing)
//   - Margin (outer spacing)
//   - Size (width/height constraints)
//   - Alignment (horizontal and vertical)
//   - Text decorations (bold, italic, underline, strikethrough)
//
// All setter methods return a new Style instance (immutability).
// Methods are chainable for fluent API usage.
//
// Example:
//
//	style := model.NewStyle().
//	    Foreground(value.RGB(255, 255, 255)).
//	    Background(value.RGB(0, 0, 255)).
//	    Bold(true).
//	    Padding(value.NewPadding(1, 2, 1, 2)).
//	    Border(value.RoundedBorder()).
//	    BorderColor(value.RGB(128, 128, 128))
type Style struct {
	// Color properties
	foreground *value.Color
	background *value.Color

	// Border properties
	border       *value.Border
	borderColor  *value.Color
	borderTop    bool
	borderBottom bool
	borderLeft   bool
	borderRight  bool

	// Spacing properties
	padding *value.Padding
	margin  *value.Margin

	// Size constraints
	size *value.Size

	// Alignment
	alignment *value.Alignment

	// Text properties
	bold          bool
	italic        bool
	underline     bool
	strikethrough bool

	// Terminal capability (for color adaptation)
	terminalCapability value.TerminalCapability
}

// NewStyle creates a new Style with default values.
// Default values:
//   - No colors set (uses terminal defaults)
//   - No border
//   - No padding/margin
//   - No size constraints
//   - No alignment (content-dependent)
//   - No text decorations
//   - TrueColor terminal capability
func NewStyle() Style {
	return Style{
		foreground:         nil,
		background:         nil,
		border:             nil,
		borderColor:        nil,
		borderTop:          false,
		borderBottom:       false,
		borderLeft:         false,
		borderRight:        false,
		padding:            nil,
		margin:             nil,
		size:               nil,
		alignment:          nil,
		bold:               false,
		italic:             false,
		underline:          false,
		strikethrough:      false,
		terminalCapability: value.TrueColor,
	}
}

// Color methods (fluent)

// Foreground sets the foreground (text) color.
// Returns a new Style instance (immutability).
func (s Style) Foreground(c value.Color) Style {
	s.foreground = &c
	return s
}

// Background sets the background color.
// Returns a new Style instance (immutability).
func (s Style) Background(c value.Color) Style {
	s.background = &c
	return s
}

// Border methods (fluent)

// Border sets the border style.
// By default, enables all sides. Use BorderTop/Bottom/Left/Right to selectively enable sides.
// Returns a new Style instance (immutability).
func (s Style) Border(b value.Border) Style {
	s.border = &b
	// Default: enable all sides when border is set
	if s.border != nil && !s.borderTop && !s.borderBottom && !s.borderLeft && !s.borderRight {
		s.borderTop = true
		s.borderBottom = true
		s.borderLeft = true
		s.borderRight = true
	}
	return s
}

// BorderColor sets the border color.
// Returns a new Style instance (immutability).
func (s Style) BorderColor(c value.Color) Style {
	s.borderColor = &c
	return s
}

// BorderTop enables or disables the top border.
// Returns a new Style instance (immutability).
func (s Style) BorderTop(enabled bool) Style {
	s.borderTop = enabled
	return s
}

// BorderBottom enables or disables the bottom border.
// Returns a new Style instance (immutability).
func (s Style) BorderBottom(enabled bool) Style {
	s.borderBottom = enabled
	return s
}

// BorderLeft enables or disables the left border.
// Returns a new Style instance (immutability).
func (s Style) BorderLeft(enabled bool) Style {
	s.borderLeft = enabled
	return s
}

// BorderRight enables or disables the right border.
// Returns a new Style instance (immutability).
func (s Style) BorderRight(enabled bool) Style {
	s.borderRight = enabled
	return s
}

// Spacing methods (fluent)

// Padding sets the padding (inner spacing).
// Returns a new Style instance (immutability).
func (s Style) Padding(p value.Padding) Style {
	s.padding = &p
	return s
}

// PaddingTop sets only the top padding.
// Returns a new Style instance (immutability).
func (s Style) PaddingTop(top int) Style {
	if s.padding == nil {
		s.padding = &value.Padding{}
	}
	updated := value.NewPadding(top, s.padding.Right(), s.padding.Bottom(), s.padding.Left())
	s.padding = &updated
	return s
}

// PaddingRight sets only the right padding.
// Returns a new Style instance (immutability).
func (s Style) PaddingRight(right int) Style {
	if s.padding == nil {
		s.padding = &value.Padding{}
	}
	updated := value.NewPadding(s.padding.Top(), right, s.padding.Bottom(), s.padding.Left())
	s.padding = &updated
	return s
}

// PaddingBottom sets only the bottom padding.
// Returns a new Style instance (immutability).
func (s Style) PaddingBottom(bottom int) Style {
	if s.padding == nil {
		s.padding = &value.Padding{}
	}
	updated := value.NewPadding(s.padding.Top(), s.padding.Right(), bottom, s.padding.Left())
	s.padding = &updated
	return s
}

// PaddingLeft sets only the left padding.
// Returns a new Style instance (immutability).
func (s Style) PaddingLeft(left int) Style {
	if s.padding == nil {
		s.padding = &value.Padding{}
	}
	updated := value.NewPadding(s.padding.Top(), s.padding.Right(), s.padding.Bottom(), left)
	s.padding = &updated
	return s
}

// PaddingAll sets the same padding for all sides.
// Returns a new Style instance (immutability).
func (s Style) PaddingAll(all int) Style {
	padding := value.NewPadding(all, all, all, all)
	s.padding = &padding
	return s
}

// PaddingHorizontal sets the same padding for left and right sides.
// Returns a new Style instance (immutability).
func (s Style) PaddingHorizontal(horizontal int) Style {
	if s.padding == nil {
		s.padding = &value.Padding{}
	}
	updated := value.NewPadding(s.padding.Top(), horizontal, s.padding.Bottom(), horizontal)
	s.padding = &updated
	return s
}

// PaddingVertical sets the same padding for top and bottom sides.
// Returns a new Style instance (immutability).
func (s Style) PaddingVertical(vertical int) Style {
	if s.padding == nil {
		s.padding = &value.Padding{}
	}
	updated := value.NewPadding(vertical, s.padding.Right(), vertical, s.padding.Left())
	s.padding = &updated
	return s
}

// Margin sets the margin (outer spacing).
// Returns a new Style instance (immutability).
func (s Style) Margin(m value.Margin) Style {
	s.margin = &m
	return s
}

// MarginTop sets only the top margin.
// Returns a new Style instance (immutability).
func (s Style) MarginTop(top int) Style {
	if s.margin == nil {
		s.margin = &value.Margin{}
	}
	updated := value.NewMargin(top, s.margin.Right(), s.margin.Bottom(), s.margin.Left())
	s.margin = &updated
	return s
}

// MarginRight sets only the right margin.
// Returns a new Style instance (immutability).
func (s Style) MarginRight(right int) Style {
	if s.margin == nil {
		s.margin = &value.Margin{}
	}
	updated := value.NewMargin(s.margin.Top(), right, s.margin.Bottom(), s.margin.Left())
	s.margin = &updated
	return s
}

// MarginBottom sets only the bottom margin.
// Returns a new Style instance (immutability).
func (s Style) MarginBottom(bottom int) Style {
	if s.margin == nil {
		s.margin = &value.Margin{}
	}
	updated := value.NewMargin(s.margin.Top(), s.margin.Right(), bottom, s.margin.Left())
	s.margin = &updated
	return s
}

// MarginLeft sets only the left margin.
// Returns a new Style instance (immutability).
func (s Style) MarginLeft(left int) Style {
	if s.margin == nil {
		s.margin = &value.Margin{}
	}
	updated := value.NewMargin(s.margin.Top(), s.margin.Right(), s.margin.Bottom(), left)
	s.margin = &updated
	return s
}

// MarginAll sets the same margin for all sides.
// Returns a new Style instance (immutability).
func (s Style) MarginAll(all int) Style {
	margin := value.NewMargin(all, all, all, all)
	s.margin = &margin
	return s
}

// MarginHorizontal sets the same margin for left and right sides.
// Returns a new Style instance (immutability).
func (s Style) MarginHorizontal(horizontal int) Style {
	if s.margin == nil {
		s.margin = &value.Margin{}
	}
	updated := value.NewMargin(s.margin.Top(), horizontal, s.margin.Bottom(), horizontal)
	s.margin = &updated
	return s
}

// MarginVertical sets the same margin for top and bottom sides.
// Returns a new Style instance (immutability).
func (s Style) MarginVertical(vertical int) Style {
	if s.margin == nil {
		s.margin = &value.Margin{}
	}
	updated := value.NewMargin(vertical, s.margin.Right(), vertical, s.margin.Left())
	s.margin = &updated
	return s
}

// Size methods (fluent)

// Width sets the exact width.
// Returns a new Style instance (immutability).
func (s Style) Width(w int) Style {
	if s.size == nil {
		size := value.NewSize()
		s.size = &size
	}
	updated := s.size.SetWidth(w)
	s.size = &updated
	return s
}

// Height sets the exact height.
// Returns a new Style instance (immutability).
func (s Style) Height(h int) Style {
	if s.size == nil {
		size := value.NewSize()
		s.size = &size
	}
	updated := s.size.SetHeight(h)
	s.size = &updated
	return s
}

// MaxWidth sets the maximum width constraint.
// Returns a new Style instance (immutability).
func (s Style) MaxWidth(w int) Style {
	if s.size == nil {
		size := value.NewSize()
		s.size = &size
	}
	updated := s.size.SetMaxWidth(w)
	s.size = &updated
	return s
}

// MaxHeight sets the maximum height constraint.
// Returns a new Style instance (immutability).
func (s Style) MaxHeight(h int) Style {
	if s.size == nil {
		size := value.NewSize()
		s.size = &size
	}
	updated := s.size.SetMaxHeight(h)
	s.size = &updated
	return s
}

// MinWidth sets the minimum width constraint.
// Returns a new Style instance (immutability).
func (s Style) MinWidth(w int) Style {
	if s.size == nil {
		size := value.NewSize()
		s.size = &size
	}
	updated := s.size.SetMinWidth(w)
	s.size = &updated
	return s
}

// MinHeight sets the minimum height constraint.
// Returns a new Style instance (immutability).
func (s Style) MinHeight(h int) Style {
	if s.size == nil {
		size := value.NewSize()
		s.size = &size
	}
	updated := s.size.SetMinHeight(h)
	s.size = &updated
	return s
}

// Alignment methods (fluent)

// Align sets both horizontal and vertical alignment.
// Returns a new Style instance (immutability).
func (s Style) Align(a value.Alignment) Style {
	s.alignment = &a
	return s
}

// AlignHorizontal sets only the horizontal alignment.
// Returns a new Style instance (immutability).
func (s Style) AlignHorizontal(h value.HorizontalAlignment) Style {
	if s.alignment == nil {
		alignment := value.NewAlignment(h, value.AlignTop)
		s.alignment = &alignment
	} else {
		updated := value.NewAlignment(h, s.alignment.Vertical())
		s.alignment = &updated
	}
	return s
}

// AlignVertical sets only the vertical alignment.
// Returns a new Style instance (immutability).
func (s Style) AlignVertical(v value.VerticalAlignment) Style {
	if s.alignment == nil {
		alignment := value.NewAlignment(value.AlignLeft, v)
		s.alignment = &alignment
	} else {
		updated := value.NewAlignment(s.alignment.Horizontal(), v)
		s.alignment = &updated
	}
	return s
}

// Text decoration methods (fluent)

// Bold enables or disables bold text.
// Returns a new Style instance (immutability).
func (s Style) Bold(enabled bool) Style {
	s.bold = enabled
	return s
}

// Italic enables or disables italic text.
// Returns a new Style instance (immutability).
func (s Style) Italic(enabled bool) Style {
	s.italic = enabled
	return s
}

// Underline enables or disables underlined text.
// Returns a new Style instance (immutability).
func (s Style) Underline(enabled bool) Style {
	s.underline = enabled
	return s
}

// Strikethrough enables or disables strikethrough text.
// Returns a new Style instance (immutability).
func (s Style) Strikethrough(enabled bool) Style {
	s.strikethrough = enabled
	return s
}

// Terminal capability

// TerminalCapability sets the terminal capability for color adaptation.
// Returns a new Style instance (immutability).
func (s Style) TerminalCapability(tc value.TerminalCapability) Style {
	s.terminalCapability = tc
	return s
}

// Getters (for rendering)

// GetForeground returns the foreground color if set.
// Returns (color, true) if set, (zero value, false) otherwise.
func (s Style) GetForeground() (value.Color, bool) {
	if s.foreground == nil {
		return value.Color{}, false
	}
	return *s.foreground, true
}

// GetBackground returns the background color if set.
// Returns (color, true) if set, (zero value, false) otherwise.
func (s Style) GetBackground() (value.Color, bool) {
	if s.background == nil {
		return value.Color{}, false
	}
	return *s.background, true
}

// GetBorder returns the border if set.
// Returns (border, true) if set, (zero value, false) otherwise.
func (s Style) GetBorder() (value.Border, bool) {
	if s.border == nil {
		return value.Border{}, false
	}
	return *s.border, true
}

// GetBorderColor returns the border color if set.
// Returns (color, true) if set, (zero value, false) otherwise.
func (s Style) GetBorderColor() (value.Color, bool) {
	if s.borderColor == nil {
		return value.Color{}, false
	}
	return *s.borderColor, true
}

// GetBorderTop returns whether the top border is enabled.
func (s Style) GetBorderTop() bool {
	return s.borderTop && s.border != nil
}

// GetBorderBottom returns whether the bottom border is enabled.
func (s Style) GetBorderBottom() bool {
	return s.borderBottom && s.border != nil
}

// GetBorderLeft returns whether the left border is enabled.
func (s Style) GetBorderLeft() bool {
	return s.borderLeft && s.border != nil
}

// GetBorderRight returns whether the right border is enabled.
func (s Style) GetBorderRight() bool {
	return s.borderRight && s.border != nil
}

// GetPadding returns the padding if set.
// Returns (padding, true) if set, (zero value, false) otherwise.
func (s Style) GetPadding() (value.Padding, bool) {
	if s.padding == nil {
		return value.Padding{}, false
	}
	return *s.padding, true
}

// GetMargin returns the margin if set.
// Returns (margin, true) if set, (zero value, false) otherwise.
func (s Style) GetMargin() (value.Margin, bool) {
	if s.margin == nil {
		return value.Margin{}, false
	}
	return *s.margin, true
}

// GetSize returns the size constraints if set.
// Returns (size, true) if set, (zero value, false) otherwise.
func (s Style) GetSize() (value.Size, bool) {
	if s.size == nil {
		return value.Size{}, false
	}
	return *s.size, true
}

// GetAlignment returns the alignment if set.
// Returns (alignment, true) if set, (zero value, false) otherwise.
func (s Style) GetAlignment() (value.Alignment, bool) {
	if s.alignment == nil {
		return value.Alignment{}, false
	}
	return *s.alignment, true
}

// GetBold returns whether bold text is enabled.
func (s Style) GetBold() bool {
	return s.bold
}

// GetItalic returns whether italic text is enabled.
func (s Style) GetItalic() bool {
	return s.italic
}

// GetUnderline returns whether underline text is enabled.
func (s Style) GetUnderline() bool {
	return s.underline
}

// GetStrikethrough returns whether strikethrough text is enabled.
func (s Style) GetStrikethrough() bool {
	return s.strikethrough
}

// GetTerminalCapability returns the terminal capability.
func (s Style) GetTerminalCapability() value.TerminalCapability {
	return s.terminalCapability
}

// Validate validates the style for consistency.
// Returns an error if the style has invalid combinations.
func (s Style) Validate() error {
	// Validate size constraints
	if s.size != nil {
	}

	// Validate padding (must be non-negative)
	if s.padding != nil {
		if s.padding.Top() < 0 || s.padding.Right() < 0 ||
			s.padding.Bottom() < 0 || s.padding.Left() < 0 {
			return fmt.Errorf("padding values cannot be negative")
		}
	}

	// Validate margin (must be non-negative)
	if s.margin != nil {
		if s.margin.Top() < 0 || s.margin.Right() < 0 ||
			s.margin.Bottom() < 0 || s.margin.Left() < 0 {
			return fmt.Errorf("margin values cannot be negative")
		}
	}

	// Border validation: if any border side is enabled, border must be set
	if (s.borderTop || s.borderBottom || s.borderLeft || s.borderRight) && s.border == nil {
		return fmt.Errorf("border sides enabled but no border style set")
	}

	return nil
}
