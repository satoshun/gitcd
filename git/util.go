package git

import (
	"os"
	"path"
)

func Exists(p string) bool {
	if _, err := os.Stat(path.Join(p, ".git")); err == nil {
		return true
	}

	return false
}
