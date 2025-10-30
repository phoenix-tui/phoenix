module github.com/phoenix-tui/phoenix/examples/wheel-scroll

go 1.25.3

replace github.com/phoenix-tui/phoenix/components => ../../components

replace github.com/phoenix-tui/phoenix/tea => ../../tea

replace github.com/phoenix-tui/phoenix/terminal => ../../terminal

replace github.com/phoenix-tui/phoenix/core => ../../core

replace github.com/phoenix-tui/phoenix/style => ../../style

require (
	github.com/phoenix-tui/phoenix/components v0.1.0-beta.4
	github.com/phoenix-tui/phoenix/tea v0.1.0-beta.4
)
