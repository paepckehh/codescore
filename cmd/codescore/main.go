// package main ...
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"paepcke.de/codescore"
)

// usage holds the help text for the CLI.
var usage = `usage: codescore [options] <path>

Options:

     --file, -f               Create .go.codescore info file
     --file-full, -F          Create .go.codescore file with details
     --score-only, -s         Print only the score to stdout
     --enable-hidden-files, -e
                             Enable scanning hidden files and directories
     --exclude <dir>          Exclude directories matching <dir> (repeatable)
     --verbose, -v            Verbose output
     --silent, -q             Silent output
     --debug, -d              Debug output
     --goo, -g                Goo mode
     --version, -V            Print version information
     --help, -h               Show this help

Targets:
     <path>                   File or directory to score;
                              "." or "*" are aliases for current directory
                              No target defaults to current directory
Version:
` + Version() + "\n"

// excludeFlag implements flag.Value for repeatable --exclude flags.
type excludeFlag []string

func (e *excludeFlag) String() string { return strings.Join(*e, ",") }
func (e *excludeFlag) Set(v string) error {
	*e = append(*e, v)
	return nil
}

// expandArgs transforms os.Args so the flag package can parse them:
//   - Expands compound short flags (-vd → -v -d)
//   - Maps -h to --help
func expandArgs(args []string) []string {
	out := make([]string, 0, len(args)*2)
	for _, arg := range args {
		if len(arg) >= 2 && arg[0] == '-' && arg[1] != '-' {
			for i := 1; i < len(arg); i++ {
				out = append(out, "-"+string(arg[i]))
			}
		} else if arg == "-h" || arg == "--help" {
			out = append(out, "--help")
		} else {
			out = append(out, arg)
		}
	}
	return out
}

// run performs all CLI logic: flag parsing, validation, and scoring.
func run() (string, error) {
	c := codescore.GetDefaultConfig()

	var excludes excludeFlag

	flagSet := flag.NewFlagSet("codescore", flag.ContinueOnError)
	flagSet.SetOutput(io.Discard)
	flagSet.Usage = func() {}

	// Boolean flags — default is current value from GetDefaultConfig().
	flagSet.BoolVar(&c.File, "file", c.File, "create .go.codescore info file")
	flagSet.BoolVar(&c.File, "f", c.File, "")
	flagSet.BoolVar(&c.FileFull, "file-full", c.FileFull, "create .go.codescore file with details")
	flagSet.BoolVar(&c.FileFull, "F", c.FileFull, "")
	flagSet.BoolVar(&c.ScoreOnly, "score-only", c.ScoreOnly, "print only the score")
	flagSet.BoolVar(&c.ScoreOnly, "s", c.ScoreOnly, "")
	flagSet.BoolVar(&c.SkipHidden, "enable-hidden-files", c.SkipHidden, "enable scanning hidden files and directories")
	flagSet.BoolVar(&c.SkipHidden, "e", c.SkipHidden, "")
	flagSet.BoolVar(&c.Verbose, "verbose", c.Verbose, "verbose output")
	flagSet.BoolVar(&c.Verbose, "v", c.Verbose, "")
	flagSet.BoolVar(&c.Silent, "silent", c.Silent, "silent mode")
	flagSet.BoolVar(&c.Silent, "q", c.Silent, "")
	flagSet.BoolVar(&c.Debug, "debug", c.Debug, "debug output")
	flagSet.BoolVar(&c.Debug, "d", c.Debug, "")
	flagSet.BoolVar(&c.Goo, "goo", c.Goo, "goo mode")
	flagSet.BoolVar(&c.Goo, "g", c.Goo, "")

	versionFlag := flagSet.Bool("version", false, "print version information")
	_ = flagSet.Bool("V", false, "print version information")

	// --exclude is long-flag only; no short alias to avoid collision with -e.
	flagSet.Var(&excludes, "exclude", "exclude directories by keyword (repeatable)")

	help := flagSet.Bool("help", false, "show this help")

	if err := flagSet.Parse(expandArgs(os.Args[1:])); err != nil {
		return "", err
	}

	if *help {
		fmt.Print(usage)
		os.Exit(0)
	}

	if *versionFlag {
		out(Version())
		os.Exit(0)
	}

	positionals := flagSet.Args()
	if len(positionals) == 0 {
		fmt.Print(usage)
		os.Exit(0)
	}

	if len(positionals) > 1 {
		return "", fmt.Errorf("more than one [file|directory] path specified")
	}

	path := positionals[0]
	switch path {
	case ".", "*":
		var err error
		if path, err = os.Getwd(); err != nil {
			return "", fmt.Errorf("unable to get current working directory [%s] [%s]", path, err)
		}
	}

	c.Path = path
	c.Exclude = excludes

	// Validate that path is a regular file or directory.
	fi, err := os.Stat(path)
	if err != nil {
		return "", fmt.Errorf("invalid target [directory|file] [%s] {%s}", path, err)
	}
	if !fi.Mode().IsRegular() && !fi.Mode().IsDir() {
		return "", fmt.Errorf("[invalid target] [%s]", path)
	}

	return c.Start(), nil
}

func main() {
	// Pre-check for --help anywhere in args since flag.Parse stops at
	// the first non-flag positional argument.
	for _, arg := range os.Args[1:] {
		switch arg {
		case "--help", "-h":
			fmt.Print(usage)
			os.Exit(0)
		case "--version", "-V":
			out(Version())
			os.Exit(0)
		}
	}

	result, err := run()
	if err != nil {
		out("[error] " + err.Error())
		os.Exit(1)
	}
	out(result)
}

// ---- CLI helpers ----

// out writes a message to stdout.
func out(msg string) {
	os.Stdout.Write([]byte(msg))
}
