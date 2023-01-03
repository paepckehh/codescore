// package codescore
package codescore

// import
import (
	"io/fs"
	"os"
	"strconv"
	"syscall"
)

//
// Display IO
//

// out ...
func out(msg string) {
	os.Stdout.Write([]byte(msg))
}

//
// Error Display IO
//

// errOut ...
func errOut(msg string) {
	out("[error] " + msg + "\n")
}

// errExit ...
func errExit(msg string) {
	errOut(msg)
	os.Exit(1)
}

//
// Little Display IO Helper
//

// itoa
func itoa(in int) string { return strconv.Itoa(in) }

// itoaFixed
func itoaFixed(in int) string {
	s := itoa(in)
	l := len(s)
	switch l {
	case 3:
		return s
	case 2:
		return "0" + s
	case 1:
		return "00" + s
	}
	errExit("internal error itoa score")
	return "" // unreachable
}

//
// File I/O
//

// const
const (
	_slashfwd = "/"
	_file     = ".codescore."
	_fileGoo  = ".goo.codescore."
)

// writeScoreFile ...
func (c *Config) writeScoreFile(sum result) {
	file := _slashfwd + _file
	switch {
	case c.Goo:
		file = _slashfwd + _fileGoo
		c.cleanScoreFiles()
		err := os.WriteFile(c.Path+file+itoaFixed(sum.score), []byte(itoaFixed(sum.score)), 0o644)
		if err != nil {
			errOut("unable to write: " + err.Error())
		}
	case c.File:
		c.cleanScoreFiles()
		err := os.WriteFile(c.Path+file+itoaFixed(sum.score), []byte(itoaFixed(sum.score)), 0o644)
		if err != nil {
			errOut("unable to write: " + err.Error())
		}
	case c.FileFull:
		report := []byte(c.scoreReport(sum))
		if len(sum.details) > 0 {
			report = append(report, []byte("\n"+sum.details)...)
		}
		c.cleanScoreFiles()
		err := os.WriteFile(c.Path+file+itoaFixed(sum.score), report, 0o644)
		if err != nil {
			errOut("unable to write: " + err.Error())
		}
	}
}

// cleanScoreFiles ...
func (c *Config) cleanScoreFiles() {
	list := readDir(c.Path)
	if c.Goo {
		for _, filename := range list {
			name := filename.Name()
			switch {
			case name[0] != '.':
				return // its a sorted list
			case len(name) != 18:
				continue
			case name[:15] == _fileGoo:
				syscall.Unlink(c.Path + "/" + name)
			}
		}
		return
	}
	for _, filename := range list {
		name := filename.Name()
		switch {
		case name[0] != '.':
			return
		case len(name) != 14:
			continue
		case name[:11] == _file:
			syscall.Unlink(c.Path + "/" + name)
		}
	}
}

// readDir ...
func readDir(path string) []fs.DirEntry {
	list, err := os.ReadDir(path)
	if err != nil {
		errExit("unable to list directory [" + path + "]")
	}
	return list
}
