package service

import (
	"strings"
	"testing"

	"github.com/phoenix-tui/phoenix/clipboard/internal/domain/value"
)

func TestNewRichTextCodec(t *testing.T) {
	codec := NewRichTextCodec()
	if codec == nil {
		t.Error("Expected non-nil codec")
	}
}

// HTML Tests

func TestEncodeHTML_Plain(t *testing.T) {
	codec := NewRichTextCodec()
	styles := value.NewTextStyles()

	result, err := codec.EncodeHTML("Hello, World!", styles)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expected := "Hello, World!"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestEncodeHTML_Bold(t *testing.T) {
	codec := NewRichTextCodec()
	styles := value.NewTextStyles().WithBold(true)

	result, err := codec.EncodeHTML("Bold text", styles)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expected := "<strong>Bold text</strong>"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestEncodeHTML_Italic(t *testing.T) {
	codec := NewRichTextCodec()
	styles := value.NewTextStyles().WithItalic(true)

	result, err := codec.EncodeHTML("Italic text", styles)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expected := "<em>Italic text</em>"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestEncodeHTML_Underline(t *testing.T) {
	codec := NewRichTextCodec()
	styles := value.NewTextStyles().WithUnderline(true)

	result, err := codec.EncodeHTML("Underline text", styles)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expected := "<u>Underline text</u>"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestEncodeHTML_Color(t *testing.T) {
	codec := NewRichTextCodec()
	styles, _ := value.NewTextStyles().WithColor("#FF0000")

	result, err := codec.EncodeHTML("Red text", styles)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expected := `<span style="color:#FF0000">Red text</span>`
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestEncodeHTML_AllStyles(t *testing.T) {
	codec := NewRichTextCodec()
	styles := value.NewTextStyles().WithBold(true).WithItalic(true).WithUnderline(true)
	styles, _ = styles.WithColor("#FF0000")

	result, err := codec.EncodeHTML("Styled text", styles)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Should nest tags: strong > em > u > span
	expected := `<span style="color:#FF0000"><u><em><strong>Styled text</strong></em></u></span>`
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestEncodeHTML_EscapeSpecialChars(t *testing.T) {
	codec := NewRichTextCodec()
	styles := value.NewTextStyles()

	result, err := codec.EncodeHTML("<script>alert('xss')</script>", styles)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Should escape HTML special characters
	if strings.Contains(result, "<script>") {
		t.Error("Failed to escape HTML special characters")
	}
	expected := "&lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt;"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestEncodeHTML_Empty(t *testing.T) {
	codec := NewRichTextCodec()
	styles := value.NewTextStyles()

	result, err := codec.EncodeHTML("", styles)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result != "" {
		t.Errorf("Expected empty string, got %q", result)
	}
}

func TestDecodeHTML_Plain(t *testing.T) {
	codec := NewRichTextCodec()

	text, styles, err := codec.DecodeHTML("Hello, World!")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if text != "Hello, World!" {
		t.Errorf("Expected 'Hello, World!', got %q", text)
	}
	if !styles.IsPlain() {
		t.Error("Expected plain styles")
	}
}

func TestDecodeHTML_Bold(t *testing.T) {
	codec := NewRichTextCodec()

	tests := []string{
		"<strong>Bold text</strong>",
		"<b>Bold text</b>",
	}

	for _, html := range tests {
		text, styles, err := codec.DecodeHTML(html)
		if err != nil {
			t.Fatalf("Unexpected error for %q: %v", html, err)
		}

		if text != "Bold text" {
			t.Errorf("Expected 'Bold text', got %q", text)
		}
		if !styles.Bold {
			t.Error("Expected bold style to be detected")
		}
	}
}

func TestDecodeHTML_Italic(t *testing.T) {
	codec := NewRichTextCodec()

	tests := []string{
		"<em>Italic text</em>",
		"<i>Italic text</i>",
	}

	for _, html := range tests {
		text, styles, err := codec.DecodeHTML(html)
		if err != nil {
			t.Fatalf("Unexpected error for %q: %v", html, err)
		}

		if text != "Italic text" {
			t.Errorf("Expected 'Italic text', got %q", text)
		}
		if !styles.Italic {
			t.Error("Expected italic style to be detected")
		}
	}
}

func TestDecodeHTML_Underline(t *testing.T) {
	codec := NewRichTextCodec()

	text, styles, err := codec.DecodeHTML("<u>Underline text</u>")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if text != "Underline text" {
		t.Errorf("Expected 'Underline text', got %q", text)
	}
	if !styles.Underline {
		t.Error("Expected underline style to be detected")
	}
}

func TestDecodeHTML_Color(t *testing.T) {
	codec := NewRichTextCodec()

	text, styles, err := codec.DecodeHTML(`<span style="color:#FF0000">Red text</span>`)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if text != "Red text" {
		t.Errorf("Expected 'Red text', got %q", text)
	}
	if styles.Color == "" {
		t.Error("Expected color to be detected")
	}
	expectedColor := "#FF0000"
	if !strings.EqualFold(styles.Color, expectedColor) {
		t.Errorf("Expected color %q, got %q", expectedColor, styles.Color)
	}
}

func TestDecodeHTML_AllStyles(t *testing.T) {
	codec := NewRichTextCodec()

	html := `<span style="color:#FF0000"><u><em><strong>Styled text</strong></em></u></span>`
	text, styles, err := codec.DecodeHTML(html)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if text != "Styled text" {
		t.Errorf("Expected 'Styled text', got %q", text)
	}
	if !styles.Bold {
		t.Error("Expected bold to be detected")
	}
	if !styles.Italic {
		t.Error("Expected italic to be detected")
	}
	if !styles.Underline {
		t.Error("Expected underline to be detected")
	}
	if styles.Color == "" {
		t.Error("Expected color to be detected")
	}
}

func TestDecodeHTML_Empty(t *testing.T) {
	codec := NewRichTextCodec()

	text, styles, err := codec.DecodeHTML("")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if text != "" {
		t.Errorf("Expected empty string, got %q", text)
	}
	if !styles.IsPlain() {
		t.Error("Expected plain styles")
	}
}

func TestStripHTMLTags(t *testing.T) {
	codec := NewRichTextCodec()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Plain text", "Hello", "Hello"},
		{"Simple tag", "<p>Hello</p>", "Hello"},
		{"Multiple tags", "<strong><em>Hello</em></strong>", "Hello"},
		{"Attributes", `<span style="color:red">Hello</span>`, "Hello"},
		{"Nested tags", "<div><p><strong>Hello</strong></p></div>", "Hello"},
		{"HTML entities", "&lt;script&gt;", "<script>"},
		{"Empty", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := codec.StripHTMLTags(tt.input)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

// RTF Tests

func TestEncodeRTF_Plain(t *testing.T) {
	codec := NewRichTextCodec()
	styles := value.NewTextStyles()

	result, err := codec.EncodeRTF("Hello, World!", styles)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Should have RTF header and footer
	if !strings.HasPrefix(result, "{\\rtf1\\ansi\\deff0") {
		t.Error("Expected RTF header")
	}
	if !strings.HasSuffix(result, "}") {
		t.Error("Expected RTF footer")
	}
	if !strings.Contains(result, "Hello, World!") {
		t.Error("Expected text in RTF")
	}
}

func TestEncodeRTF_Bold(t *testing.T) {
	codec := NewRichTextCodec()
	styles := value.NewTextStyles().WithBold(true)

	result, err := codec.EncodeRTF("Bold text", styles)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !strings.Contains(result, "\\b ") {
		t.Error("Expected bold code in RTF")
	}
	if !strings.Contains(result, "\\b0 ") {
		t.Error("Expected bold end code in RTF")
	}
}

func TestEncodeRTF_Italic(t *testing.T) {
	codec := NewRichTextCodec()
	styles := value.NewTextStyles().WithItalic(true)

	result, err := codec.EncodeRTF("Italic text", styles)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !strings.Contains(result, "\\i ") {
		t.Error("Expected italic code in RTF")
	}
	if !strings.Contains(result, "\\i0 ") {
		t.Error("Expected italic end code in RTF")
	}
}

func TestEncodeRTF_Underline(t *testing.T) {
	codec := NewRichTextCodec()
	styles := value.NewTextStyles().WithUnderline(true)

	result, err := codec.EncodeRTF("Underline text", styles)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !strings.Contains(result, "\\ul ") {
		t.Error("Expected underline code in RTF")
	}
	if !strings.Contains(result, "\\ul0 ") {
		t.Error("Expected underline end code in RTF")
	}
}

func TestEncodeRTF_Color(t *testing.T) {
	codec := NewRichTextCodec()
	styles, _ := value.NewTextStyles().WithColor("#FF0000")

	result, err := codec.EncodeRTF("Red text", styles)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Should have color table
	if !strings.Contains(result, "\\colortbl") {
		t.Error("Expected color table in RTF")
	}
	if !strings.Contains(result, "\\red255\\green0\\blue0") {
		t.Error("Expected RGB values in color table")
	}
	if !strings.Contains(result, "\\cf1 ") {
		t.Error("Expected color reference in RTF")
	}
}

func TestEncodeRTF_AllStyles(t *testing.T) {
	codec := NewRichTextCodec()
	styles := value.NewTextStyles().WithBold(true).WithItalic(true).WithUnderline(true)
	styles, _ = styles.WithColor("#FF0000")

	result, err := codec.EncodeRTF("Styled text", styles)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Check all style codes
	if !strings.Contains(result, "\\b ") {
		t.Error("Expected bold code")
	}
	if !strings.Contains(result, "\\i ") {
		t.Error("Expected italic code")
	}
	if !strings.Contains(result, "\\ul ") {
		t.Error("Expected underline code")
	}
	if !strings.Contains(result, "\\colortbl") {
		t.Error("Expected color table")
	}
}

func TestEncodeRTF_EscapeSpecialChars(t *testing.T) {
	codec := NewRichTextCodec()
	styles := value.NewTextStyles()

	tests := []struct {
		name     string
		input    string
		contains string
	}{
		{"Backslash", "test\\value", "\\\\"},
		{"Left brace", "test{value", "\\{"},
		{"Right brace", "test}value", "\\}"},
		{"Newline", "line1\nline2", "\\line"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := codec.EncodeRTF(tt.input, styles)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if !strings.Contains(result, tt.contains) {
				t.Errorf("Expected %q to contain %q", result, tt.contains)
			}
		})
	}
}

func TestEncodeRTF_Empty(t *testing.T) {
	codec := NewRichTextCodec()
	styles := value.NewTextStyles()

	result, err := codec.EncodeRTF("", styles)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result != "" {
		t.Errorf("Expected empty string, got %q", result)
	}
}

func TestDecodeRTF_Plain(t *testing.T) {
	codec := NewRichTextCodec()

	rtf := "{\\rtf1\\ansi\\deff0\nHello, World!\n}"
	text, styles, err := codec.DecodeRTF(rtf)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if text != "Hello, World!" {
		t.Errorf("Expected 'Hello, World!', got %q", text)
	}
	if !styles.IsPlain() {
		t.Error("Expected plain styles")
	}
}

func TestDecodeRTF_Bold(t *testing.T) {
	codec := NewRichTextCodec()

	rtf := "{\\rtf1\\ansi\\deff0\n\\b Bold text\\b0 \n}"
	text, styles, err := codec.DecodeRTF(rtf)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if text != "Bold text" {
		t.Errorf("Expected 'Bold text', got %q", text)
	}
	if !styles.Bold {
		t.Error("Expected bold to be detected")
	}
}

func TestDecodeRTF_Color(t *testing.T) {
	codec := NewRichTextCodec()

	rtf := "{\\rtf1\\ansi\\deff0\n{\\colortbl ;\\red255\\green0\\blue0;}\n\\cf1 Red text\\cf0 \n}"
	text, styles, err := codec.DecodeRTF(rtf)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if text != "Red text" {
		t.Errorf("Expected 'Red text', got %q", text)
	}
	if styles.Color == "" {
		t.Error("Expected color to be detected")
	}
	expectedColor := "#FF0000"
	if !strings.EqualFold(styles.Color, expectedColor) {
		t.Errorf("Expected color %q, got %q", expectedColor, styles.Color)
	}
}

func TestDecodeRTF_InvalidHeader(t *testing.T) {
	codec := NewRichTextCodec()

	_, _, err := codec.DecodeRTF("Not valid RTF")
	if err == nil {
		t.Error("Expected error for invalid RTF")
	}
}

func TestDecodeRTF_Empty(t *testing.T) {
	codec := NewRichTextCodec()

	text, styles, err := codec.DecodeRTF("")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if text != "" {
		t.Errorf("Expected empty string, got %q", text)
	}
	if !styles.IsPlain() {
		t.Error("Expected plain styles")
	}
}

func TestStripRTFFormatting(t *testing.T) {
	codec := NewRichTextCodec()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Plain", "{\\rtf1\\ansi\\deff0\nHello\n}", "Hello"},
		{"Bold", "{\\rtf1\\ansi\\deff0\n\\b Bold\\b0 \n}", "Bold"},
		{"Color table", "{\\rtf1\\ansi\\deff0\n{\\colortbl ;\\red255\\green0\\blue0;}\nText\n}", "Text"},
		{"Empty", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := codec.StripRTFFormatting(tt.input)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

// Format Conversion Tests

func TestHTMLToRTF(t *testing.T) {
	codec := NewRichTextCodec()

	html := "<strong>Bold text</strong>"
	rtf, err := codec.HTMLToRTF(html)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Should be valid RTF
	if !strings.HasPrefix(rtf, "{\\rtf1\\ansi\\deff0") {
		t.Error("Expected RTF header")
	}
	if !strings.Contains(rtf, "\\b ") {
		t.Error("Expected bold code in RTF")
	}
	if !strings.Contains(rtf, "Bold text") {
		t.Error("Expected text in RTF")
	}
}

func TestRTFToHTML(t *testing.T) {
	codec := NewRichTextCodec()

	rtf := "{\\rtf1\\ansi\\deff0\n\\b Bold text\\b0 \n}"
	html, err := codec.RTFToHTML(rtf)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Should be valid HTML
	if !strings.Contains(html, "<strong>") {
		t.Error("Expected strong tag in HTML")
	}
	if !strings.Contains(html, "Bold text") {
		t.Error("Expected text in HTML")
	}
}

func TestHTMLToRTF_Empty(t *testing.T) {
	codec := NewRichTextCodec()

	rtf, err := codec.HTMLToRTF("")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if rtf != "" {
		t.Errorf("Expected empty string, got %q", rtf)
	}
}

func TestRTFToHTML_Empty(t *testing.T) {
	codec := NewRichTextCodec()

	html, err := codec.RTFToHTML("")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if html != "" {
		t.Errorf("Expected empty string, got %q", html)
	}
}

// Round Trip Tests

func TestHTMLRoundTrip(t *testing.T) {
	codec := NewRichTextCodec()

	tests := []struct {
		name   string
		text   string
		styles value.TextStyles
	}{
		{"Plain", "Hello", value.NewTextStyles()},
		{"Bold", "Bold text", value.NewTextStyles().WithBold(true)},
		{"Italic", "Italic text", value.NewTextStyles().WithItalic(true)},
		{"Underline", "Underline text", value.NewTextStyles().WithUnderline(true)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Encode to HTML
			html, err := codec.EncodeHTML(tt.text, tt.styles)
			if err != nil {
				t.Fatalf("Encode failed: %v", err)
			}

			// Decode back
			text, styles, err := codec.DecodeHTML(html)
			if err != nil {
				t.Fatalf("Decode failed: %v", err)
			}

			// Verify text
			if text != tt.text {
				t.Errorf("Text mismatch: expected %q, got %q", tt.text, text)
			}

			// Verify styles
			if tt.styles.Bold != styles.Bold {
				t.Errorf("Bold mismatch: expected %v, got %v", tt.styles.Bold, styles.Bold)
			}
			if tt.styles.Italic != styles.Italic {
				t.Errorf("Italic mismatch: expected %v, got %v", tt.styles.Italic, styles.Italic)
			}
			if tt.styles.Underline != styles.Underline {
				t.Errorf("Underline mismatch: expected %v, got %v", tt.styles.Underline, styles.Underline)
			}
		})
	}
}

func TestRTFRoundTrip(t *testing.T) {
	codec := NewRichTextCodec()

	tests := []struct {
		name   string
		text   string
		styles value.TextStyles
	}{
		{"Plain", "Hello", value.NewTextStyles()},
		{"Bold", "Bold text", value.NewTextStyles().WithBold(true)},
		{"Italic", "Italic text", value.NewTextStyles().WithItalic(true)},
		{"Underline", "Underline text", value.NewTextStyles().WithUnderline(true)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Encode to RTF
			rtf, err := codec.EncodeRTF(tt.text, tt.styles)
			if err != nil {
				t.Fatalf("Encode failed: %v", err)
			}

			// Decode back
			text, styles, err := codec.DecodeRTF(rtf)
			if err != nil {
				t.Fatalf("Decode failed: %v", err)
			}

			// Verify text
			if text != tt.text {
				t.Errorf("Text mismatch: expected %q, got %q", tt.text, text)
			}

			// Verify styles
			if tt.styles.Bold != styles.Bold {
				t.Errorf("Bold mismatch: expected %v, got %v", tt.styles.Bold, styles.Bold)
			}
			if tt.styles.Italic != styles.Italic {
				t.Errorf("Italic mismatch: expected %v, got %v", tt.styles.Italic, styles.Italic)
			}
			if tt.styles.Underline != styles.Underline {
				t.Errorf("Underline mismatch: expected %v, got %v", tt.styles.Underline, styles.Underline)
			}
		})
	}
}
