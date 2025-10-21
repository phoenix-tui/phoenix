module github.com/phoenix-tui/phoenix/tea

go 1.25.1

require (
	github.com/phoenix-tui/phoenix/terminal v0.1.0-beta.2
	github.com/phoenix-tui/phoenix/testing v0.1.0-beta.2
	github.com/stretchr/testify v1.11.1
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/sys v0.37.0 // indirect
	golang.org/x/term v0.36.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

// Local development: use terminal from feature/tea-exec-process branch
replace github.com/phoenix-tui/phoenix/terminal => ../terminal

// Local development: use testing from feature/tea-exec-process branch
replace github.com/phoenix-tui/phoenix/testing => ../testing
