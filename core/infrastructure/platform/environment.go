// Package platform provides platform-specific infrastructure implementations.
package platform

import "os"

// OsEnvironmentProvider implements service.EnvironmentProvider using os package.
// This is the infrastructure adapter that provides real environment variable access.
//
// Hexagonal Architecture:
//   - Domain layer defines EnvironmentProvider interface (port)
//   - Infrastructure layer provides OsEnvironmentProvider implementation (adapter)
//   - No domain code depends on this infrastructure implementation
//
// Usage:
//
//	env := platform.OsEnvironmentProvider{}
//	detector := service.NewCapabilitiesDetector(env)
//	caps := detector.Detect()
type OsEnvironmentProvider struct{}

// Get retrieves environment variable value using os.Getenv.
// Returns empty string if variable is not set.
func (e OsEnvironmentProvider) Get(key string) string {
	return os.Getenv(key)
}

// Lookup retrieves environment variable value with existence check using os.LookupEnv.
// Returns value and true if variable exists, empty string and false otherwise.
func (e OsEnvironmentProvider) Lookup(key string) (string, bool) {
	return os.LookupEnv(key)
}

// Platform returns the operating system platform ("linux", "darwin", "windows", etc).
func (e OsEnvironmentProvider) Platform() string {
	// Use build tags to return platform at compile time
	return platformString
}
