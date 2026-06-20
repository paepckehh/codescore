package main

// Injected by linker flags during build.
var (
	version = "v0.1.xx-dev"
	commit  = "dev-unknown"
	date    = "dev-unknown"
)

// Version returns the full version string (e.g. "codescore version v0.1.48 (e00b8ea-dirty) built 2025-06-20").
func Version() string {
	return "codescore: " + version + " (commit:" + commit + ", built:" + date + ")"
}
