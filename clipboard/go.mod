module github.com/phoenix-tui/phoenix/clipboard

go 1.25.1

require (
	github.com/google/uuid v1.6.0
	github.com/phoenix-tui/phoenix/tea v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.11.1
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/phoenix-tui/phoenix/terminal v0.1.0-beta.3 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/sys v0.37.0 // indirect
	golang.org/x/term v0.36.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/phoenix-tui/phoenix/terminal => ../terminal

replace github.com/phoenix-tui/phoenix/tea => ../tea

replace github.com/phoenix-tui/phoenix/testing => ../testing
