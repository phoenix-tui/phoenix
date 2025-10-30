module github.com/phoenix-tui/phoenix/examples/context-menu

go 1.25.1

replace github.com/phoenix-tui/phoenix/mouse => ../../mouse

replace github.com/phoenix-tui/phoenix/style => ../../style

replace github.com/phoenix-tui/phoenix/core => ../../core

replace github.com/phoenix-tui/phoenix/tea => ../../tea

replace github.com/phoenix-tui/phoenix/terminal => ../../terminal

replace github.com/phoenix-tui/phoenix/testing => ../../testing

require (
	github.com/phoenix-tui/phoenix/mouse v0.0.0
	github.com/phoenix-tui/phoenix/style v0.0.0
	github.com/phoenix-tui/phoenix/tea v0.0.0
)

require (
	github.com/phoenix-tui/phoenix/core v0.1.0-beta.1 // indirect
	github.com/phoenix-tui/phoenix/terminal v0.1.0-beta.3 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/unilibs/uniwidth v0.1.0-beta // indirect
	golang.org/x/sys v0.37.0 // indirect
	golang.org/x/term v0.36.0 // indirect
)
