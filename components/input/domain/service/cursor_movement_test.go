package service

import "testing"

func TestCursorMovementService_MoveLeft(t *testing.T) {
	svc := NewCursorMovementService()

	tests := []struct {
		name       string
		content    string
		currentPos int
		want       int
	}{
		{"simple text", "hello", 3, 2},
		{"at start", "hello", 0, 0},
		{"at end", "hello", 5, 4},
		{"emoji", "👋hello", 2, 1},
		{"emoji at cursor", "hello👋world", 6, 5},
		{"cjk text", "你好世界", 2, 1},
		{"combining chars", "é", 1, 0}, // é as e + combining acute
		{"empty string", "", 0, 0},
		{"single char", "a", 1, 0},
		{"negative pos clamped", "hello", -5, 0},
		{"beyond end clamped", "hello", 10, 5}, // Beyond end returns max (no movement)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := svc.MoveLeft(tt.content, tt.currentPos)
			if got != tt.want {
				t.Errorf("MoveLeft(%q, %d) = %d, want %d", tt.content, tt.currentPos, got, tt.want)
			}
		})
	}
}

func TestCursorMovementService_MoveRight(t *testing.T) {
	svc := NewCursorMovementService()

	tests := []struct {
		name       string
		content    string
		currentPos int
		want       int
	}{
		{"simple text", "hello", 2, 3},
		{"at start", "hello", 0, 1},
		{"at end", "hello", 5, 5},
		{"before end", "hello", 4, 5},
		{"emoji", "👋hello", 0, 1},
		{"emoji after cursor", "hello👋world", 5, 6},
		{"cjk text", "你好世界", 1, 2},
		{"combining chars", "é", 0, 1},
		{"empty string", "", 0, 0},
		{"single char", "a", 0, 1},
		{"beyond end stays", "hello", 10, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := svc.MoveRight(tt.content, tt.currentPos)
			if got != tt.want {
				t.Errorf("MoveRight(%q, %d) = %d, want %d", tt.content, tt.currentPos, got, tt.want)
			}
		})
	}
}

func TestCursorMovementService_GraphemeCount(t *testing.T) {
	svc := NewCursorMovementService()

	tests := []struct {
		name    string
		content string
		want    int
	}{
		{"empty", "", 0},
		{"ascii", "hello", 5},
		{"emoji", "👋👋👋", 3},
		{"emoji mixed", "hello👋world", 11},
		{"cjk", "你好世界", 4},
		{"combining", "é", 1},           // e + combining acute = 1 grapheme
		{"complex emoji", "👨‍👩‍👧‍👦", 1}, // family emoji = 1 grapheme
		{"mixed script", "Hello世界👋", 8},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := svc.GraphemeCount(tt.content)
			if got != tt.want {
				t.Errorf("GraphemeCount(%q) = %d, want %d", tt.content, got, tt.want)
			}
		})
	}
}

func TestCursorMovementService_SplitAtCursor(t *testing.T) {
	svc := NewCursorMovementService()

	tests := []struct {
		name       string
		content    string
		pos        int
		wantBefore string
		wantAt     string
		wantAfter  string
	}{
		{
			name:       "middle of simple text",
			content:    "hello",
			pos:        2,
			wantBefore: "he",
			wantAt:     "l",
			wantAfter:  "lo",
		},
		{
			name:       "at start",
			content:    "hello",
			pos:        0,
			wantBefore: "",
			wantAt:     "h",
			wantAfter:  "ello",
		},
		{
			name:       "at end",
			content:    "hello",
			pos:        5,
			wantBefore: "hello",
			wantAt:     "",
			wantAfter:  "",
		},
		{
			name:       "before emoji",
			content:    "hello👋world",
			pos:        5,
			wantBefore: "hello",
			wantAt:     "👋",
			wantAfter:  "world",
		},
		{
			name:       "after emoji",
			content:    "hello👋world",
			pos:        6,
			wantBefore: "hello👋",
			wantAt:     "w",
			wantAfter:  "orld",
		},
		{
			name:       "cjk text",
			content:    "你好世界",
			pos:        2,
			wantBefore: "你好",
			wantAt:     "世",
			wantAfter:  "界",
		},
		{
			name:       "empty string",
			content:    "",
			pos:        0,
			wantBefore: "",
			wantAt:     "",
			wantAfter:  "",
		},
		{
			name:       "negative position",
			content:    "hello",
			pos:        -1,
			wantBefore: "",
			wantAt:     "h",
			wantAfter:  "ello",
		},
		{
			name:       "beyond end",
			content:    "hello",
			pos:        10,
			wantBefore: "hello",
			wantAt:     "",
			wantAfter:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			before, at, after := svc.SplitAtCursor(tt.content, tt.pos)
			if before != tt.wantBefore {
				t.Errorf("before = %q, want %q", before, tt.wantBefore)
			}
			if at != tt.wantAt {
				t.Errorf("at = %q, want %q", at, tt.wantAt)
			}
			if after != tt.wantAfter {
				t.Errorf("after = %q, want %q", after, tt.wantAfter)
			}
		})
	}
}

func TestCursorMovementService_ByteOffsetToGraphemeOffset(t *testing.T) {
	svc := NewCursorMovementService()

	tests := []struct {
		name       string
		content    string
		byteOffset int
		want       int
	}{
		{"ascii middle", "hello", 2, 2},
		{"ascii start", "hello", 0, 0},
		{"ascii end", "hello", 5, 5},
		{"emoji boundary", "hello👋", 5, 5},
		{"within emoji", "hello👋", 6, 5}, // Inside emoji bytes -> grapheme before
		{"after emoji", "hello👋world", 9, 6},
		{"cjk start", "你好", 0, 0},
		{"cjk second char", "你好", 3, 1},
		{"beyond end", "hello", 100, 5},
		{"negative", "hello", -1, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := svc.ByteOffsetToGraphemeOffset(tt.content, tt.byteOffset)
			if got != tt.want {
				t.Errorf("ByteOffsetToGraphemeOffset(%q, %d) = %d, want %d",
					tt.content, tt.byteOffset, got, tt.want)
			}
		})
	}
}

func TestCursorMovementService_GraphemeOffsetToByteOffset(t *testing.T) {
	svc := NewCursorMovementService()

	tests := []struct {
		name           string
		content        string
		graphemeOffset int
		want           int
	}{
		{"ascii middle", "hello", 2, 2},
		{"ascii start", "hello", 0, 0},
		{"ascii end", "hello", 5, 5},
		{"before emoji", "hello👋", 5, 5},
		{"after emoji", "hello👋world", 6, 9},
		{"cjk first", "你好", 0, 0},
		{"cjk second", "你好", 1, 3},
		{"cjk end", "你好", 2, 6},
		{"beyond end", "hello", 100, 5},
		{"negative", "hello", -1, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := svc.GraphemeOffsetToByteOffset(tt.content, tt.graphemeOffset)
			if got != tt.want {
				t.Errorf("GraphemeOffsetToByteOffset(%q, %d) = %d, want %d",
					tt.content, tt.graphemeOffset, got, tt.want)
			}
		})
	}
}

// Test round-trip conversion
func TestCursorMovementService_RoundTrip(t *testing.T) {
	svc := NewCursorMovementService()
	testStrings := []string{
		"hello",
		"hello👋world",
		"你好世界",
		"👨‍👩‍👧‍👦family",
		"é", // combining character
		"",
	}

	for _, content := range testStrings {
		t.Run(content, func(t *testing.T) {
			maxGraphemes := svc.GraphemeCount(content)
			for graphemeOffset := 0; graphemeOffset <= maxGraphemes; graphemeOffset++ {
				// Convert grapheme -> byte -> grapheme
				byteOffset := svc.GraphemeOffsetToByteOffset(content, graphemeOffset)
				backToGrapheme := svc.ByteOffsetToGraphemeOffset(content, byteOffset)

				if backToGrapheme != graphemeOffset {
					t.Errorf("round-trip failed for %q at grapheme %d: got %d",
						content, graphemeOffset, backToGrapheme)
				}
			}
		})
	}
}
