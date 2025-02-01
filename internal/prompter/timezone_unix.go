//go:build unix && !ios && !android

package prompter

import (
	"os"
	"slices"
	"strings"

	"github.com/jsnfwlr/keyper-cli/internal/files"
)

var zoneDirs = []string{
	"/usr/share/zoneinfo/",
	"/usr/share/lib/zoneinfo/",
	"/usr/lib/locale/TZ/",
	"/etc/zoneinfo",
}

var zoneDir string

var exclude = []string{
	"SystemV",
	"US",
	"posix",
	"right",
}

func Timezones() (regions []string, cityMap map[string][]string) {
	r := []string{}
	c := map[string][]string{}
	for _, zoneDir = range zoneDirs {
		if files.Exists(zoneDir, false) {
			rfs, _ := os.ReadDir(zoneDir)
			if len(rfs) == 0 {
				continue
			}
			for _, rf := range rfs {
				if rf.IsDir() && !slices.Contains(exclude, rf.Name()) {
					r = append(r, rf.Name())
					rc := []string{}
					cfs, _ := os.ReadDir(zoneDir + rf.Name())
					for _, cf := range cfs {
						if cf.Name() != strings.ToUpper(cf.Name()[:1])+cf.Name()[1:] {
							continue
						}
						rc = append(rc, cf.Name())
					}
					c[rf.Name()] = rc
				}
			}

		}
	}
	return r, c
}

// func ReadFile(path string) {
// 	files, _ := os.ReadDir(zoneDir + path)
// 	for _, f := range files {
// 		if f.Name() != strings.ToUpper(f.Name()[:1])+f.Name()[1:] {
// 			continue
// 		}
// 		if f.IsDir() {
// 			ReadFile(path + "/" + f.Name())
// 		} else {
// 			return (path + "/" + f.Name())[1:]
// 		}
// 	}
// }
