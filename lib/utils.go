package transit

import (
	"os"
	"path/filepath"
)

// CurrentPath returns the absolute current path.
func CurrentPath() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	return exPath
}
