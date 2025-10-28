module github.com/phoenix-tui/phoenix/testing

go 1.25.1

require github.com/phoenix-tui/phoenix/terminal v0.1.0-beta.1

require (
	golang.org/x/sys v0.37.0 // indirect
	golang.org/x/term v0.36.0 // indirect
)

replace github.com/phoenix-tui/phoenix/terminal => ../terminal
