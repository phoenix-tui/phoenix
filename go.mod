module github.com/phoenix-tui/phoenix

go 1.25

// Phoenix TUI Framework - Root module
// This is an umbrella module that provides convenient access to all Phoenix libraries.
//
// Individual libraries can be imported directly:
// - github.com/phoenix-tui/phoenix/clipboard
// - github.com/phoenix-tui/phoenix/components
// - github.com/phoenix-tui/phoenix/core
// - github.com/phoenix-tui/phoenix/layout
// - github.com/phoenix-tui/phoenix/mouse
// - github.com/phoenix-tui/phoenix/render
// - github.com/phoenix-tui/phoenix/style
// - github.com/phoenix-tui/phoenix/tea
// - github.com/phoenix-tui/phoenix/terminal
// - github.com/phoenix-tui/phoenix/testing
//
// Each library has its own go.mod and can be versioned independently.
// This root module uses replace directives to point to local subdirectories
// for development (similar to opentelemetry-go and kubernetes).

// Replace directives for local development
// These allow local development and testing of interdependent modules.
replace github.com/phoenix-tui/phoenix/clipboard => ./clipboard

replace github.com/phoenix-tui/phoenix/components => ./components

replace github.com/phoenix-tui/phoenix/core => ./core

replace github.com/phoenix-tui/phoenix/layout => ./layout

replace github.com/phoenix-tui/phoenix/mouse => ./mouse

replace github.com/phoenix-tui/phoenix/render => ./render

replace github.com/phoenix-tui/phoenix/style => ./style

replace github.com/phoenix-tui/phoenix/tea => ./tea

replace github.com/phoenix-tui/phoenix/terminal => ./terminal

replace github.com/phoenix-tui/phoenix/testing => ./testing
