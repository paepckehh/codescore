//go:build !windows

package codescore

import (
	"os"
	"syscall"
)

const (
	_modeDir     uint32 = 1 << (32 - 1 - 0)
	_modeSymlink uint32 = 1 << (32 - 1 - 4)
)

func (c *Config) fastWalker() {
	fi, err := os.Stat(c.Path)
	if err != nil {
		errExit("[stat deviceID root dir] [" + c.Path + "] [" + err.Error() + "]")
	}
	d, ok := fi.Sys().(*syscall.Stat_t)
	if !ok {
		errExit("[stat deviceID root dir] [" + c.Path + "]")
	}
	exclude := false
	if len(c.Exclude) > 0 {
		exclude = true
	}
	skipme := false
	rootNodeDeviceID := uint64(d.Dev)
	for i := 0; i < 1; i++ {
		go func() {
			inodeSeen := make(map[uint64]struct{})
			for path := range channelDir {
				list, err := os.ReadDir(path)
				if err != nil {
					errOut("[unable to read directory] [" + path + "] [" + err.Error() + "]")
					walk.Done()
					continue
				}
				for _, item := range list {
					fi, _ := item.Info()
					fname := item.Name()
					if c.SkipHidden {
						if fname[0] == '.' {
							continue // skip hidden files if requested
						}
					}
					ftype := uint32(item.Type())
					name := path + "/" + fname
					inode := fi.Sys().(*syscall.Stat_t).Ino
					if _, ok := inodeSeen[inode]; ok {
						continue // skip inode if we seen it already
					}
					inodeSeen[inode] = struct{}{}
					if ftype&_modeSymlink != 0 {
						continue // skip symlinks
					}
					if ftype&_modeDir != 0 {
						st, _ := fi.Sys().(*syscall.Stat_t)
						if uint64(st.Dev) != rootNodeDeviceID {
							continue // skip dirtargets outside fs boundary
						}
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
