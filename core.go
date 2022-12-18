// package codescore
package codescore

// import
import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/revivelib"
)

// const
const (
	_linefeed = '\n'
	chanBuff  = 1000
)

// var
var (
	bg, walk, collect sync.WaitGroup
	worker            = runtime.NumCPU()
	channelDir        = make(chan string, chanBuff*25)
	channelGoFiles    = make(chan string, chanBuff)
	channelResults    = make(chan result, chanBuff)
)

// type
type result struct {
	score, loc, rate int
	details          string
}

// score ...
func (c *Config) score() string {
	var (
		err error
		r   result
	)
	switch c.Path {
	case ".", "*", "":
		c.Path, err = os.Getwd()
		if err != nil {
			errExit("unable to get current working directory [" + c.Path + "] [" + err.Error() + "]")
		}
	}
	fi, err := os.Stat(c.Path)
	if err != nil {
		errExit("invalid target [directory|file] [" + c.Path + "] {" + err.Error() + "]")
	}
	switch mode := fi.Mode(); {
	case mode.IsRegular():
		r = scoreFile(c.Path, lintConfig())
	case mode.IsDir():
		r = c.scorePath()
	default:
		errExit("[invalid target] [" + c.Path + "]")
	}
	switch {
	case c.ScoreOnly:
		return itoa(r.score)
	case c.Verbose, c.Debug:
		return c.scoreReport(r) + "\n" + r.details
	}
	return c.scoreReport(r)
}

// scoreReport ...
func (c *Config) scoreReport(r result) string {
	if c.Goo {
		s := strings.Split(c.Path, _slashfwd)
		l := len(s)
		if l > 3 {
			i := 2
			for i < l-1 {
				c.Path = s[i] + _slashfwd
				i++
			}
			c.Path = s[i]
		}
	}
	return "[score:" + itoaFixed(r.score) + "/100] [loc:" + itoa(r.loc) + "] [err:" + itoa(r.rate) + "] -> [" + c.Path + "]"
}

// scorePath ...
func (c *Config) scorePath() result {
	walk.Add(1)
	channelDir <- c.Path
	go c.fastWalker()
	bg.Add(worker)
	go scoreWorker()
	collect.Add(1)
	sum := result{}
	var countFiles int
	go func() {
		for r := range channelResults {
			countFiles++
			sum.score += r.score
			sum.loc += r.loc
			sum.rate += r.rate
			sum.details += r.details
		}
		collect.Done()
	}()
	walk.Wait()
	close(channelDir)
	close(channelGoFiles)
	bg.Wait()
	close(channelResults)
	collect.Wait()
	if countFiles > 1 {
		sum.score = sum.score / countFiles
	}
	c.writeScoreFile(sum)
	return sum
}

// scoreWorker ...
func scoreWorker() {
	for i := 0; i < worker; i++ {
		go func() {
			config := lintConfig()
			for filename := range channelGoFiles {
				channelResults <- scoreFile(filename, config)
			}
			bg.Done()
		}()
	}
}

// scoreFile ...
func scoreFile(filename string, config *lint.Config) result {
	defer func() {
		if r := recover(); r != nil {
			out("[codereview] [external ast parser crash/panic] [skip] [" + filename + "]")
		}
	}()
	loc := countFileLoc(filename) // count first [fastio/cache]
	lin, err := revivelib.New(config, false, 0)
	if err != nil {
		errOut(err.Error())
	}
	packages := []*revivelib.LintPattern{}
	packages = append(packages, revivelib.Include(filename))
	failures, err := lin.Lint(packages...)
	if err != nil {
		errOut(err.Error())
	}
	var output string
	for failure := range failures {
		output += fmt.Sprintf("%v: %s\n", failure.Position.Start, failure.Failure)
	}
	errCount := countLines([]byte(output))
	return result{100 - (100 * errCount / loc), loc, errCount, output}
}

// countFileLoc ...
func countFileLoc(filename string) int {
	f, err := os.ReadFile(filename)
	if err != nil {
		errExit("invalid target [directory|file] [" + filename + "] {" + err.Error() + "]")
	}
	lines := countLines(f)
	if lines < 1 {
		return 1
	}
	return lines
}

// countLines ...
func countLines(in []byte) int {
	size, lines := len(in), 0
	for i := 0; i < size; i++ {
		if in[i] == _linefeed {
			lines++
		}
	}
	return lines
}
