package transit

import (
	"os"
	"path/filepath"
	"hash/crc32"
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

func CRC32 (s string) int {
	crc32q := crc32.MakeTable(0xD5828281)
	return int(crc32.Checksum([]byte(s), crc32q))
}

func IsFileSizeGreaterThanZero (p string) bool {
	fi, e := os.Stat(p);
	if e != nil {
		return false
	}
	// get the size
	return fi.Size()>0
}