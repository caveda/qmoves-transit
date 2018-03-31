// Package transit provides functions for fetching and processing public transit
// information.
package transit

import (
	"os"
	"testing"
)

func TestDownload(t *testing.T) {
	targetPath := "./test/golang.doc"
	downloadErr := Download(targetPath, "https://golang.org/doc/")
	if downloadErr != nil {
		t.Errorf("Download returned error: %v", downloadErr)
	}

	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		t.Errorf("%v does not exists.", targetPath)
	}
}
