// package main ...
package main

import (
	"os"

	"paepcke.de/codescore"
)

// func main ...
func main() {
	c := codescore.GetDefaultConfig()
	l := len(os.Args)
	var opt string
	switch {
	case l > 1:
		for i := 1; i < l; i++ {
			o := os.Args[i]
			switch {
			case o[0] == '-':
				switch {
				case o == "--file" || o == "-f":
					c.File = true
					opt += "[--file] "
				case o == "--file-full" || o == "-F":
					c.FileFull = true
					opt += "[--file-full] "
				case o == "--help" || o == "-h":
					out(_syntax)
					os.Exit(0)
				case o == "--verbose" || o == "-v":
					c.Verbose = true
					opt += "[--verbose] "
				case o == "--silent" || o == "-q":
					c.Silent = true
					opt += "[--silent] "
				case o == "--debug" || o == "-d":
					c.Debug = true
					opt += "[--debug] "
				case o == "--enable-hidden-files" || o == "-e":
					c.SkipHidden = false
					opt += "[--enable-hidden-files] "
				case o == "--score-only" || o == "-s":
					c.ScoreOnly = true
					opt += "[--score-only] "
				case o == "--goo" || o == "-g":
					c.Goo = true
					opt += "[--goo] "
				case o == "--exclude" || o == "-e":
					i++
					switch {
					case i < l:
						c.Exclude = append(c.Exclude, os.Args[i])
						opt += "[--exclude " + os.Args[i] + "] "
					default:
						errExit("exclude switch value missing")
					}
				default:
					errExit("unkown commandline switch [" + o + "]")
				}
			case o == ".", o == "*":
				if c.Path != _empty {
					errExit("more than one [file|directory] path specified")
				}
				var err error
				if c.Path, err = os.Getwd(); err != nil {
					errExit("unable to validate current directory path")
				}
			case isFile(o):
				if c.Path != _empty {
					errExit(" more than one [file|directory] path specified")
				}
				c.Path = o
			case isDir(o):
				if c.Path != _empty {
					errExit("more than one [file|directory] path specified")
				}
				c.Path = o
			default:
				errExit("invalid path or option [" + o + "]")
			}
		}
		out(c.Start())
	default:
		out(_syntax)
	}
}

const (
	_syntax string = "syntax: codescore [options] <file|directory>\n\n--file [-f]\n\t\tcreate .go.codescore info file\n\n--file-full [-F]\n\t\tcreate .go.codescore info file and dump all details into the file\n\n--score-only [-s]\n\t\tprint only the score to stdout\n\n--enable-hidden-files [-e]\n\t\tenable scanning hidden files and directories\n\n--exclude [-e]\n\t\texclude all directories matching any of the keywords\n\t\tthis option can be specified several times\n\n--verbose [-v]\n--silent [-q]\n--debug [-d]\n--help [-h]\n"
	_empty  string = ""
)

//
// LITTLE GENERIC HELPER SECTION
//

const (
	_modeDir uint32 = 1 << (32 - 1 - 0)
)

// out ...
func out(msg string) {
	os.Stdout.Write([]byte(msg))
}

// errExit ...
func errExit(msg string) {
	out("[error] " + msg)
	os.Exit(1)
}

// isDir ...
func isDir(filename string) bool {
	fi, err := os.Stat(filename)
	if err != nil {
		return false
	}
	return uint32(fi.Mode())&_modeDir != 0
}

// isFile ...
func isFile(filename string) bool {
	fi, err := os.Lstat(filename)
	if err != nil {
		return false
	}
	return fi.Mode().IsRegular()
}
