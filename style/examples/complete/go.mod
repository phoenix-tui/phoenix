module github.com/phoenix-tui/phoenix/style/examples/complete

go 1.25.1

replace github.com/phoenix-tui/phoenix/style => ../..

replace github.com/phoenix-tui/phoenix/core => ../../../core

require github.com/phoenix-tui/phoenix/style v0.0.0

require (
	github.com/phoenix-tui/phoenix/core v0.2.0 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/unilibs/uniwidth v0.2.0 // indirect
)
