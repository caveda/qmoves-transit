package transit

import (
	"log"
	"os"
	"path"
)

// Publish deploys lines in the correct format in path.
// The format is determined by the presenter.
func Publish(lines []Line, destPath string, p Presenter) error {
	os.MkdirAll(destPath, os.ModePerm)

	for _, l := range lines {
		fl, err := p.Format(l)
		if err != nil {
			log.Printf("Error formatting line %v. Error:%v", l, err)
			return err
		}

		// Write formatted line as a file in destination
		err = CreateFile(path.Join(destPath, l.Id), fl)
		if err != nil {
			log.Printf("Error creating file for line %v. Error:%v", l.Id, err)
			return err
		}
	}

	return nil
}
