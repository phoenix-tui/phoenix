module github.com/phoenix-tui/phoenix/components

go 1.25.1

require (
	github.com/phoenix-tui/phoenix/tea v0.1.0-beta.3
	github.com/rivo/uniseg v0.4.7
)

require (
	github.com/phoenix-tui/phoenix/terminal v0.1.0-beta.3 // indirect
	golang.org/x/sys v0.37.0 // indirect
	golang.org/x/term v0.36.0 // indirect
)

replace github.com/phoenix-tui/phoenix/tea => ../tea

// Local development
replace github.com/phoenix-tui/phoenix/terminal => ../terminal

// Local development
replace github.com/phoenix-tui/phoenix/testing => ../testing
