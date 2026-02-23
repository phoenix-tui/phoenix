module github.com/phoenix-tui/phoenix/clipboard

go 1.25.1

require (
	github.com/google/uuid v1.6.0
	github.com/phoenix-tui/phoenix/tea v0.2.4
	github.com/stretchr/testify v1.11.1
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/phoenix-tui/phoenix/terminal v0.2.4 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/sys v0.40.0 // indirect
	golang.org/x/term v0.39.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/phoenix-tui/phoenix/terminal => ../terminal

replace github.com/phoenix-tui/phoenix/tea => ../tea

replace github.com/phoenix-tui/phoenix/testing => ../testing
