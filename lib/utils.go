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

// CreateFile creates a file in path with the supplied content.
func CreateFile(path, content string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}

	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()

	_, err = f.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}

// Exists returns true if the path exists, false otherwise.
func Exists(path string) bool {
	result := false
	if _, err := os.Stat(path); err == nil {
		result = true
	}
	return result
}
