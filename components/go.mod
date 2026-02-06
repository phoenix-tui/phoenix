module github.com/phoenix-tui/phoenix/components

go 1.25.1

require (
	github.com/phoenix-tui/phoenix/tea v0.2.0
	github.com/rivo/uniseg v0.4.7
)

require (
	github.com/phoenix-tui/phoenix/core v0.2.0 // indirect
	github.com/unilibs/uniwidth v0.2.0 // indirect
)

require (
	github.com/phoenix-tui/phoenix/style v0.2.0
	github.com/phoenix-tui/phoenix/terminal v0.2.0 // indirect
	golang.org/x/sys v0.40.0 // indirect
	golang.org/x/term v0.39.0 // indirect
)

replace github.com/phoenix-tui/phoenix/tea => ../tea

// Local development
replace github.com/phoenix-tui/phoenix/terminal => ../terminal

// Local development
replace github.com/phoenix-tui/phoenix/testing => ../testing

replace github.com/phoenix-tui/phoenix/style => ../style

replace github.com/phoenix-tui/phoenix/core => ../core
