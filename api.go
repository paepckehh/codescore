// package codescore provides a code quality score
package codescore

// Config ...
type Config struct {
	Path       string
	Exclude    []string
	File       bool
	FileFull   bool
	ScoreOnly  bool
	SkipHidden bool
	Verbose    bool
	Silent     bool
	Goo        bool
	Debug      bool
}

// GetDefaultConfig returns and Config object with sane defaults
func GetDefaultConfig() *Config {
	return &Config{
		SkipHidden: true,
	}
}

// Start calculates the codescore for all go source files in c.Path
func (c *Config) Start() string { return c.score() }
