module github.com/phoenix-tui/phoenix/examples/cobra-cli

go 1.25.1

replace github.com/phoenix-tui/phoenix/components => ../../components

replace github.com/phoenix-tui/phoenix/style => ../../style

replace github.com/phoenix-tui/phoenix/tea => ../../tea

replace github.com/phoenix-tui/phoenix/core => ../../core

replace github.com/phoenix-tui/phoenix/terminal => ../../terminal

replace github.com/phoenix-tui/phoenix/testing => ../../testing

require (
	github.com/phoenix-tui/phoenix/components v0.2.4
	github.com/phoenix-tui/phoenix/style v0.2.4
	github.com/phoenix-tui/phoenix/tea v0.2.4
	github.com/spf13/cobra v1.10.2
)

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/phoenix-tui/phoenix/core v0.2.4 // indirect
	github.com/phoenix-tui/phoenix/terminal v0.2.4 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/spf13/pflag v1.0.10 // indirect
	github.com/unilibs/uniwidth v0.2.0 // indirect
	golang.org/x/sys v0.40.0 // indirect
	golang.org/x/term v0.39.0 // indirect
)
