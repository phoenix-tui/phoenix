package service

// EnvironmentProvider is a port (interface) for reading environment variables.
// This keeps the domain layer pure by abstracting infrastructure concerns.
//
// Implementation lives in infrastructure/platform/environment.go.
type EnvironmentProvider interface {
	// Get returns environment variable value (empty string if not set).
	Get(key string) string

	// Platform returns the operating system ("linux", "darwin", "windows", etc).
	Platform() string
}
