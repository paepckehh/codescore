//go:build windows

// walker windows does not implement (yet) skiping hardline
// files(inode seen based on windows 'fileid') and no support
// for filesystem boundary crossing checks
package codescore

import (
	"os"
)

const (
	_modeDir     uint32 = 1 << (32 - 1 - 0)
	_modeSymlink uint32 = 1 << (32 - 1 - 4)
)

func (c *Config) fastWalker() {
	exclude := false
	if len(c.Exclude) > 0 {
		exclude = true
	}
	skipme := false
	for i := 0; i < 1; i++ {
		go func() {
			for path := range channelDir {
				list, err := os.ReadDir(path)
				if err != nil {
					errOut("[unable to read directory] [" + path + "] [" + err.Error() + "]")
					walk.Done()
					continue
				}
				for _, item := range list {
					fname := item.Name()
					if c.SkipHidden {
						if fname[0] == '.' {
							continue // skip hidden files if requested
						}
					}
					ftype := uint32(item.Type())
					name := path + "/" + fname
					if ftype&_modeSymlink != 0 {
						continue // skip symlinks
					}
					if ftype&_modeDir != 0 {
						if exclude {
							skipme = false
							for _, exclude := range c.Exclude {
								if fname == exclude {
									skipme = true
									break // skip exclude list
								}
							}
							if skipme {
								continue
							}
						}
						walk.Add(1)
						channelDir <- name
						continue
					}
					l := len(fname)
					if l > 3 {
						if fname[l-3:] == ".go" {
							channelGoFiles <- name
							continue
						}
					}
				}
				walk.Done()
			}
		}()
	}
}
