package value

import (
	"fmt"
	"regexp"
	"strings"
)

// TextStyles represents text formatting options for rich text content.
// This is a value object that encapsulates styling information.
type TextStyles struct {
	Bold      bool
	Italic    bool
	Underline bool
	Color     string // hex color like "#FF0000" (RGB format)
}

// NewTextStyles creates a new TextStyles value object with default values (no styling).
func NewTextStyles() TextStyles {
	return TextStyles{
		Bold:      false,
		Italic:    false,
		Underline: false,
		Color:     "",
	}
}

// NewTextStylesWithColor creates a new TextStyles with the specified color.
func NewTextStylesWithColor(color string) (TextStyles, error) {
	if color != "" {
		if err := validateColor(color); err != nil {
			return TextStyles{}, err
		}
	}
	return TextStyles{
		Bold:      false,
		Italic:    false,
		Underline: false,
		Color:     color,
	}, nil
}

// WithBold returns a new TextStyles with bold enabled or disabled.
func (t TextStyles) WithBold(bold bool) TextStyles {
	return TextStyles{
		Bold:      bold,
		Italic:    t.Italic,
		Underline: t.Underline,
		Color:     t.Color,
	}
}

// WithItalic returns a new TextStyles with italic enabled or disabled.
func (t TextStyles) WithItalic(italic bool) TextStyles {
	return TextStyles{
		Bold:      t.Bold,
		Italic:    italic,
		Underline: t.Underline,
		Color:     t.Color,
	}
}

// WithUnderline returns a new TextStyles with underline enabled or disabled.
func (t TextStyles) WithUnderline(underline bool) TextStyles {
	return TextStyles{
		Bold:      t.Bold,
		Italic:    t.Italic,
		Underline: underline,
		Color:     t.Color,
	}
}

// WithColor returns a new TextStyles with the specified color.
// Color should be in hex format: "#RRGGBB" (e.g., "#FF0000" for red).
// Pass empty string to clear the color.
func (t TextStyles) WithColor(color string) (TextStyles, error) {
	if color != "" {
		if err := validateColor(color); err != nil {
			return t, err
		}
	}
	return TextStyles{
		Bold:      t.Bold,
		Italic:    t.Italic,
		Underline: t.Underline,
		Color:     color,
	}, nil
}

// IsPlain returns true if the TextStyles has no formatting applied.
func (t TextStyles) IsPlain() bool {
	return !t.Bold && !t.Italic && !t.Underline && t.Color == ""
}

// Equals compares two TextStyles for equality.
func (t TextStyles) Equals(other TextStyles) bool {
	return t.Bold == other.Bold &&
		t.Italic == other.Italic &&
		t.Underline == other.Underline &&
		strings.EqualFold(t.Color, other.Color) // Case-insensitive color comparison
}

// String returns a human-readable representation of the TextStyles.
func (t TextStyles) String() string {
	if t.IsPlain() {
		return "plain"
	}

	var styles []string
	if t.Bold {
		styles = append(styles, "bold")
	}
	if t.Italic {
		styles = append(styles, "italic")
	}
	if t.Underline {
		styles = append(styles, "underline")
	}
	if t.Color != "" {
		styles = append(styles, fmt.Sprintf("color:%s", t.Color))
	}

	return strings.Join(styles, ",")
}

// validateColor validates that the color is in hex format: #RRGGBB
func validateColor(color string) error {
	// Color must be in format #RRGGBB
	matched, err := regexp.MatchString(`^#[0-9A-Fa-f]{6}$`, color)
	if err != nil {
		return fmt.Errorf("failed to validate color: %w", err)
	}
	if !matched {
		return fmt.Errorf("invalid color format: %s (expected #RRGGBB)", color)
	}
	return nil
}

// ParseColorRGB parses a color string and returns RGB components (0-255).
// The color must be in hex format: #RRGGBB.
func ParseColorRGB(color string) (r, g, b int, err error) {
	if err := validateColor(color); err != nil {
		return 0, 0, 0, err
	}

	// Parse hex values
	var rgb int
	_, err = fmt.Sscanf(color, "#%06x", &rgb)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("failed to parse color: %w", err)
	}

	r = (rgb >> 16) & 0xFF
	g = (rgb >> 8) & 0xFF
	b = rgb & 0xFF

	return r, g, b, nil
}

// FormatColorRGB formats RGB components (0-255) into a hex color string (#RRGGBB).
func FormatColorRGB(r, g, b int) (string, error) {
	// Validate RGB values
	if r < 0 || r > 255 || g < 0 || g > 255 || b < 0 || b > 255 {
		return "", fmt.Errorf("RGB values must be in range 0-255")
	}

	return fmt.Sprintf("#%02X%02X%02X", r, g, b), nil
}
