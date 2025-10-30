package service

import (
	"bytes"
	"fmt"
	"html"
	"regexp"
	"strings"

	"github.com/phoenix-tui/phoenix/clipboard/internal/domain/value"
)

// RichTextCodec provides HTML and RTF encoding/decoding functionality.
// It's a domain service that handles rich text format conversions.
type RichTextCodec struct{}

// NewRichTextCodec creates a new RichTextCodec instance.
func NewRichTextCodec() *RichTextCodec {
	return &RichTextCodec{}
}

// EncodeHTML encodes text with styles into HTML format.
// Supports basic formatting: bold, italic, underline, and colors.
func (c *RichTextCodec) EncodeHTML(text string, styles value.TextStyles) (string, error) {
	if text == "" {
		return "", nil
	}

	// Escape HTML special characters
	encoded := html.EscapeString(text)

	// If no styles, return as-is
	if styles.IsPlain() {
		return encoded, nil
	}

	// Apply styles
	if styles.Bold {
		encoded = fmt.Sprintf("<strong>%s</strong>", encoded)
	}
	if styles.Italic {
		encoded = fmt.Sprintf("<em>%s</em>", encoded)
	}
	if styles.Underline {
		encoded = fmt.Sprintf("<u>%s</u>", encoded)
	}
	if styles.Color != "" {
		encoded = fmt.Sprintf(`<span style="color:%s">%s</span>`, styles.Color, encoded)
	}

	return encoded, nil
}

// DecodeHTML decodes HTML content to plain text and extracts styles.
// This is a simplified decoder that handles basic tags only.
// Returns the plain text, detected styles, and any error.
func (c *RichTextCodec) DecodeHTML(htmlContent string) (string, value.TextStyles, error) {
	if htmlContent == "" {
		return "", value.NewTextStyles(), nil
	}

	styles := value.NewTextStyles()
	text := htmlContent

	// Detect bold
	if strings.Contains(text, "<strong>") || strings.Contains(text, "<b>") {
		styles = styles.WithBold(true)
	}

	// Detect italic
	if strings.Contains(text, "<em>") || strings.Contains(text, "<i>") {
		styles = styles.WithItalic(true)
	}

	// Detect underline
	if strings.Contains(text, "<u>") {
		styles = styles.WithUnderline(true)
	}

	// Detect color from style attribute
	colorRegex := regexp.MustCompile(`style="color:\s*(#[0-9A-Fa-f]{6})"`)
	if matches := colorRegex.FindStringSubmatch(text); len(matches) > 1 {
		color := matches[1]
		var err error
		styles, err = styles.WithColor(color)
		if err != nil {
			return "", styles, fmt.Errorf("invalid color in HTML: %w", err)
		}
	}

	// Strip all HTML tags
	plainText, err := c.StripHTMLTags(text)
	if err != nil {
		return "", styles, err
	}

	return plainText, styles, nil
}

// StripHTMLTags removes all HTML tags from the input, leaving only plain text.
// This provides security by removing potentially dangerous HTML.
func (c *RichTextCodec) StripHTMLTags(htmlContent string) (string, error) {
	if htmlContent == "" {
		return "", nil
	}

	// Remove all HTML tags using regex
	tagRegex := regexp.MustCompile(`<[^>]*>`)
	stripped := tagRegex.ReplaceAllString(htmlContent, "")

	// Unescape HTML entities
	unescaped := html.UnescapeString(stripped)

	return unescaped, nil
}

// EncodeRTF encodes text with styles into RTF format.
// Supports basic formatting: bold, italic, underline, and colors.
func (c *RichTextCodec) EncodeRTF(text string, styles value.TextStyles) (string, error) {
	if text == "" {
		return "", nil
	}

	var buf bytes.Buffer

	// RTF header
	buf.WriteString("{\\rtf1\\ansi\\deff0\n")

	// Color table (if color is specified)
	if styles.Color != "" {
		r, g, b, err := value.ParseColorRGB(styles.Color)
		if err != nil {
			return "", fmt.Errorf("invalid color: %w", err)
		}
		buf.WriteString(fmt.Sprintf("{\\colortbl ;\\red%d\\green%d\\blue%d;}\n", r, g, b))
	}

	// Apply formatting
	if styles.Bold {
		buf.WriteString("\\b ")
	}
	if styles.Italic {
		buf.WriteString("\\i ")
	}
	if styles.Underline {
		buf.WriteString("\\ul ")
	}
	if styles.Color != "" {
		buf.WriteString("\\cf1 ")
	}

	// Escape RTF special characters
	escapedText := c.escapeRTF(text)
	buf.WriteString(escapedText)

	// Close formatting tags
	if styles.Bold {
		buf.WriteString("\\b0 ")
	}
	if styles.Italic {
		buf.WriteString("\\i0 ")
	}
	if styles.Underline {
		buf.WriteString("\\ul0 ")
	}
	if styles.Color != "" {
		buf.WriteString("\\cf0 ")
	}

	// RTF footer
	buf.WriteString("\n}")

	return buf.String(), nil
}

// DecodeRTF decodes RTF content to plain text and extracts styles.
// This is a simplified decoder that handles basic RTF tags only.
// Returns the plain text, detected styles, and any error.
func (c *RichTextCodec) DecodeRTF(rtfContent string) (string, value.TextStyles, error) {
	if rtfContent == "" {
		return "", value.NewTextStyles(), nil
	}

	// Check for RTF header
	if !strings.HasPrefix(strings.TrimSpace(rtfContent), "{\\rtf") {
		return "", value.NewTextStyles(), fmt.Errorf("invalid RTF format: missing header")
	}

	styles := value.NewTextStyles()

	// Detect bold
	if strings.Contains(rtfContent, "\\b ") {
		styles = styles.WithBold(true)
	}

	// Detect italic
	if strings.Contains(rtfContent, "\\i ") {
		styles = styles.WithItalic(true)
	}

	// Detect underline
	if strings.Contains(rtfContent, "\\ul ") {
		styles = styles.WithUnderline(true)
	}

	// Detect color from color table
	colorTableRegex := regexp.MustCompile(`\\colortbl\s*;\\red(\d+)\\green(\d+)\\blue(\d+);`)
	if matches := colorTableRegex.FindStringSubmatch(rtfContent); len(matches) > 3 {
		var r, g, b int
		fmt.Sscanf(matches[1], "%d", &r)
		fmt.Sscanf(matches[2], "%d", &g)
		fmt.Sscanf(matches[3], "%d", &b)

		color, err := value.FormatColorRGB(r, g, b)
		if err != nil {
			return "", styles, fmt.Errorf("invalid color in RTF: %w", err)
		}
		styles, err = styles.WithColor(color)
		if err != nil {
			return "", styles, fmt.Errorf("failed to set color: %w", err)
		}
	}

	// Strip RTF formatting
	plainText, err := c.StripRTFFormatting(rtfContent)
	if err != nil {
		return "", styles, err
	}

	return plainText, styles, nil
}

// StripRTFFormatting removes all RTF formatting from the input, leaving only plain text.
func (c *RichTextCodec) StripRTFFormatting(rtfContent string) (string, error) {
	if rtfContent == "" {
		return "", nil
	}

	// Remove RTF control words and groups
	text := rtfContent

	// Remove color table
	colorTableRegex := regexp.MustCompile(`\{\\colortbl[^}]*\}`)
	text = colorTableRegex.ReplaceAllString(text, "")

	// Remove RTF header
	text = strings.TrimPrefix(text, "{\\rtf1\\ansi\\deff0")
	text = strings.TrimPrefix(text, "\n")

	// Remove formatting codes
	text = strings.ReplaceAll(text, "\\b ", "")
	text = strings.ReplaceAll(text, "\\b0 ", "")
	text = strings.ReplaceAll(text, "\\i ", "")
	text = strings.ReplaceAll(text, "\\i0 ", "")
	text = strings.ReplaceAll(text, "\\ul ", "")
	text = strings.ReplaceAll(text, "\\ul0 ", "")
	text = regexp.MustCompile(`\\cf\d+ `).ReplaceAllString(text, "")

	// Remove remaining control sequences
	controlRegex := regexp.MustCompile(`\\[a-z]+\d*\s*`)
	text = controlRegex.ReplaceAllString(text, "")

	// Remove braces
	text = strings.ReplaceAll(text, "{", "")
	text = strings.ReplaceAll(text, "}", "")

	// Trim whitespace
	text = strings.TrimSpace(text)

	return text, nil
}

// HTMLToRTF converts HTML content to RTF format.
func (c *RichTextCodec) HTMLToRTF(htmlContent string) (string, error) {
	if htmlContent == "" {
		return "", nil
	}

	// Decode HTML to get text and styles
	text, styles, err := c.DecodeHTML(htmlContent)
	if err != nil {
		return "", fmt.Errorf("failed to decode HTML: %w", err)
	}

	// Encode to RTF
	rtf, err := c.EncodeRTF(text, styles)
	if err != nil {
		return "", fmt.Errorf("failed to encode RTF: %w", err)
	}

	return rtf, nil
}

// RTFToHTML converts RTF content to HTML format.
func (c *RichTextCodec) RTFToHTML(rtfContent string) (string, error) {
	if rtfContent == "" {
		return "", nil
	}

	// Decode RTF to get text and styles
	text, styles, err := c.DecodeRTF(rtfContent)
	if err != nil {
		return "", fmt.Errorf("failed to decode RTF: %w", err)
	}

	// Encode to HTML
	htmlContent, err := c.EncodeHTML(text, styles)
	if err != nil {
		return "", fmt.Errorf("failed to encode HTML: %w", err)
	}

	return htmlContent, nil
}

// escapeRTF escapes special characters for RTF format.
func (c *RichTextCodec) escapeRTF(text string) string {
	var buf bytes.Buffer

	for _, r := range text {
		switch r {
		case '\\':
			buf.WriteString("\\\\")
		case '{':
			buf.WriteString("\\{")
		case '}':
			buf.WriteString("\\}")
		case '\n':
			buf.WriteString("\\line\n")
		case '\r':
			// Skip carriage returns
			continue
		default:
			// For non-ASCII characters, use Unicode escape
			if r > 127 {
				buf.WriteString(fmt.Sprintf("\\u%d?", r))
			} else {
				buf.WriteRune(r)
			}
		}
	}

	return buf.String()
}
