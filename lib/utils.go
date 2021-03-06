package transit

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"hash/crc32"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
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

// CRC32 calculates the CRC32 of the given string
func CRC32(s string) int {
	crc32q := crc32.MakeTable(0xD5828281)
	return int(crc32.Checksum([]byte(s), crc32q))
}

// IsFileSizeGreaterThanZero returns true if the file
// p is greater than zero.
func IsFileSizeGreaterThanZero(p string) bool {
	fi, e := os.Stat(p)
	if e != nil {
		return false
	}
	// get the size
	return fi.Size() > 0
}

// ApplyRegexAllSubmatch compile and apply the supplied regular expression over
// content and returns the list of matches
func ApplyRegexAllSubmatch(content string, regex string) ([][]string, error) {
	// Compile expression
	r, err := regexp.Compile(regex)
	if err != nil {
		log.Printf("Error compiling schedule regex %v. Error: %v ", regex, err)
		return nil, err
	}

	matches := r.FindAllStringSubmatch(content, -1)
	if matches == nil {
		return nil, errors.New("No matches found inside the content")
	}

	return matches, nil
}

// GetEnvVariableValueBool returns true if the variable is defined
// and the value can be mapped to true. False otherwise.
func GetEnvVariableValueBool(v string) bool {
	result := false
	value := os.Getenv(v)
	b, err := strconv.ParseBool(value)
	if err == nil {
		result = b
	}
	return result
}

// MD5 returns the MD5 checksum of the input as string
func MD5(s string) string {
	hash := md5.New()
	io.WriteString(hash, s)
	return hex.EncodeToString(hash.Sum(nil))
}
