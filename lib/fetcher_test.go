// Package transit provides functions for fetching and processing public transit
// information.
package transit

import (
	"os"
	"testing"
)

func validateSize (p string) bool {
	fi, e := os.Stat(p);
	if e != nil {
		return false
	}
	// get the size
	return fi.Size()>0
}

func TestDownload(t *testing.T) {
	targetPath := "./test/golang.doc"
	downloadErr := Download("https://golang.org/doc/", targetPath, validateSize)
	if downloadErr != nil {
		t.Errorf("Download returned error: %v", downloadErr)
	}

	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		t.Errorf("%v does not exists.", targetPath)
	}
}
