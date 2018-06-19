package transit

import (
	"archive/zip"
	"io"
	"log"
	"os"
	"path/filepath"
)

// UnzipFromArchive is a function that unzips a file from the supplied zip archive.
// The unzipped file is save to dest path.
func UnzipFromArchive(archive, file, dest string) error {
	r, err := zip.OpenReader(archive)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	os.MkdirAll(dest, os.ModePerm)

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(dest, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), f.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		if f.Name == file {
			log.Printf("%v found inside %v: unzipping", file, archive)
			err = extractAndWriteFile(f)
			if err != nil {
				return err
			}
			break
		}
	}

	return nil
}
