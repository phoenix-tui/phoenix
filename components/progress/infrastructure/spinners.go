package infrastructure

import "github.com/phoenix-tui/phoenix/components/progress/domain/value"

// Pre-defined spinner styles for common use cases.
// These are curated for visual appeal and cross-platform compatibility.

var (
	// SpinnerDots - Unicode Braille pattern dots (most popular)
	// ⠋ ⠙ ⠹ ⠸ ⠼ ⠴ ⠦ ⠧ ⠇ ⠏
	SpinnerDots = value.NewSpinnerStyle([]string{
		"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏",
	}, 10)

	// SpinnerLine - Classic ASCII line spinner
	// | / - \
	SpinnerLine = value.NewSpinnerStyle([]string{
		"|", "/", "-", "\\",
	}, 10)

	// SpinnerArrow - Rotating arrow
	// ← ↖ ↑ ↗ → ↘ ↓ ↙
	SpinnerArrow = value.NewSpinnerStyle([]string{
		"←", "↖", "↑", "↗", "→", "↘", "↓", "↙",
	}, 8)

	// SpinnerCircle - Rotating circle quarters
	// ◐ ◓ ◑ ◒
	SpinnerCircle = value.NewSpinnerStyle([]string{
		"◐", "◓", "◑", "◒",
	}, 8)

	// SpinnerBounce - Bouncing ball effect
	// ⠁ ⠂ ⠄ ⡀ ⢀ ⠠ ⠐ ⠈
	SpinnerBounce = value.NewSpinnerStyle([]string{
		"⠁", "⠂", "⠄", "⡀", "⢀", "⠠", "⠐", "⠈",
	}, 10)

	// SpinnerDotPulse - Pulsing dots
	// ⣾ ⣽ ⣻ ⢿ ⡿ ⣟ ⣯ ⣷
	SpinnerDotPulse = value.NewSpinnerStyle([]string{
		"⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷",
	}, 12)

	// SpinnerGrowVertical - Vertical growth
	// ▁ ▃ ▄ ▅ ▆ ▇ █ ▇ ▆ ▅ ▄ ▃
	SpinnerGrowVertical = value.NewSpinnerStyle([]string{
		"▁", "▃", "▄", "▅", "▆", "▇", "█", "▇", "▆", "▅", "▄", "▃",
	}, 8)

	// SpinnerGrowHorizontal - Horizontal growth
	// ▏ ▎ ▍ ▌ ▋ ▊ ▉ █ ▉ ▊ ▋ ▌ ▍ ▎
	SpinnerGrowHorizontal = value.NewSpinnerStyle([]string{
		"▏", "▎", "▍", "▌", "▋", "▊", "▉", "█", "▉", "▊", "▋", "▌", "▍", "▎",
	}, 8)

	// SpinnerBoxBounce - Box bouncing
	// ▖ ▘ ▝ ▗
	SpinnerBoxBounce = value.NewSpinnerStyle([]string{
		"▖", "▘", "▝", "▗",
	}, 10)

	// SpinnerSimpleDots - Simple ASCII dots
	// .  .. ...
	SpinnerSimpleDots = value.NewSpinnerStyle([]string{
		".  ", ".. ", "...",
	}, 6)

	// SpinnerClock - Clock rotation
	// 🕐 🕑 🕒 🕓 🕔 🕕 🕖 🕗 🕘 🕙 🕚 🕛
	SpinnerClock = value.NewSpinnerStyle([]string{
		"🕐", "🕑", "🕒", "🕓", "🕔", "🕕", "🕖", "🕗", "🕘", "🕙", "🕚", "🕛",
	}, 4)

	// SpinnerEarth - Spinning earth
	// 🌍 🌎 🌏
	SpinnerEarth = value.NewSpinnerStyle([]string{
		"🌍", "🌎", "🌏",
	}, 6)

	// SpinnerMoon - Moon phases
	// 🌑 🌒 🌓 🌔 🌕 🌖 🌗 🌘
	SpinnerMoon = value.NewSpinnerStyle([]string{
		"🌑", "🌒", "🌓", "🌔", "🌕", "🌖", "🌗", "🌘",
	}, 5)

	// SpinnerToggle - On/off toggle
	// ⊶ ⊷
	SpinnerToggle = value.NewSpinnerStyle([]string{
		"⊶", "⊷",
	}, 8)

	// SpinnerHamburger - Hamburger menu animation
	// ☱ ☲ ☴
	SpinnerHamburger = value.NewSpinnerStyle([]string{
		"☱", "☲", "☴",
	}, 8)
)

// styleRegistry maps style names to spinner styles.
var styleRegistry = map[string]*value.SpinnerStyle{
	"dots":            SpinnerDots,
	"line":            SpinnerLine,
	"arrow":           SpinnerArrow,
	"circle":          SpinnerCircle,
	"bounce":          SpinnerBounce,
	"dot-pulse":       SpinnerDotPulse,
	"grow-vertical":   SpinnerGrowVertical,
	"grow-horizontal": SpinnerGrowHorizontal,
	"box-bounce":      SpinnerBoxBounce,
	"simple-dots":     SpinnerSimpleDots,
	"clock":           SpinnerClock,
	"earth":           SpinnerEarth,
	"moon":            SpinnerMoon,
	"toggle":          SpinnerToggle,
	"hamburger":       SpinnerHamburger,
}

// GetSpinnerStyle retrieves a pre-defined spinner style by name.
// Returns SpinnerDots if the name is not found.
func GetSpinnerStyle(name string) *value.SpinnerStyle {
	style, ok := styleRegistry[name]
	if !ok {
		return SpinnerDots // Default fallback
	}
	return style
}

// AvailableStyles returns a list of all available spinner style names.
func AvailableStyles() []string {
	styles := make([]string, 0, len(styleRegistry))
	for name := range styleRegistry {
		styles = append(styles, name)
	}
	return styles
}
