module github.com/phoenix-tui/phoenix/render

go 1.25.1

require (
	github.com/rivo/uniseg v0.4.7
	github.com/stretchr/testify v1.11.1
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/phoenix-tui/phoenix/core => ../core
	github.com/phoenix-tui/phoenix/style => ../style
)
