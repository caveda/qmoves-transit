// Package transit provides functions for fetching and processing public transit
// information.
package transit

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/secsy/goftp"
	"log"
)

var MaxDownloadRetries int = 10

// Pulls down the content of the provided url and stores in filepath.
// If anything goes wrong it returns an error otherwise nil
func Download(address string, fileFullPath string, validateFunc func(string)bool) (err error) {
	// Create the folder
	os.MkdirAll(filepath.Dir(fileFullPath), os.ModePerm)


	// Sources may have different protocols
	u, err := url.Parse(address)
	if err != nil {
		return err
	}

	retries := 0
	validated := false
	for retries<=MaxDownloadRetries && !validated {
		retries++
		if u.Scheme == "http" || u.Scheme == "https" {
			err = fetchHTTP(u, fileFullPath)
		} else if u.Scheme == "ftp" {
			err = fetchFTP(u, fileFullPath)
		} else {
			err = fmt.Errorf("Unknown protocol for url %v", address)
		}
		validated = validateFunc(fileFullPath)
		log.Printf("%v downloaded OK", fileFullPath)
	}
	return err
}

// fetchHTTP retrieves the file url points to via HTTP.
// Output file is created in destPath.
func fetchHTTP(u *url.URL, destPath string) error {
	out, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(u.String())
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

// fetchFTP retrieves the file url points to via FTP.
// Output file is created in destPath.
func fetchFTP(u *url.URL, destPath string) (err error) {

	// Create client object with default config
	client, err := goftp.Dial(u.Hostname())
	if err != nil {
		return err
	}

	// Download a file to disk
	readme, err := os.Create(destPath)
	if err != nil {
		return err
	}

	err = client.Retrieve(u.RequestURI(), readme)
	if err != nil {
		return err
	}

	return nil
}
