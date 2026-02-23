module github.com/phoenix-tui/phoenix/testing

go 1.25.1

require github.com/phoenix-tui/phoenix/terminal v0.2.4

require (
	golang.org/x/sys v0.40.0 // indirect
	golang.org/x/term v0.39.0 // indirect
)

replace github.com/phoenix-tui/phoenix/terminal => ../terminal
