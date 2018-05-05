// Package transit provides functions for fetching and processing public transit
// information.
package transit

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// Downloads all sources passed as argument.
// func DownloadAll(sources []TransitSource) {
// 	for _, s in range sources {
// 		Download(s.Uri, s.Path)
// 	}
// }

// Pulls down the content of the provided url and stores in filepath.
// If anything goes wrong it returns an error otherwise nil
func Download(url string, fileFullPath string) (err error) {
	// Create the file
	os.MkdirAll(filepath.Dir(fileFullPath), os.ModePerm)
	out, err := os.Create(fileFullPath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}
