package files

import (
	"os"
	"path/filepath"
	"strings"
)

func DirWalk(root, suffix string, depth, maxDepth int) ([]string, error) {
	depth++
	dirContents, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}
	files := []string{}
	for _, entry := range dirContents {
		if entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") {
			if depth < maxDepth {
				f, err := DirWalk(filepath.Join(root, entry.Name()), suffix, depth, maxDepth)
				if err != nil {
					return nil, err
				}
				files = append(files, f...)
			}
		} else if strings.HasSuffix(entry.Name(), suffix) {
			files = append(files, filepath.Join(root, entry.Name()))
		}
	}
	return files, nil
}

func DirectoryOf(path string) string {
	return filepath.Dir(path)
}
