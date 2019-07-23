package transit

import (
	"log"
	"os"
	"path"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/db"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
)

// Constants
const envDatabaseUri string = "DATABASE_URI"
const envDatabaseCredentials string = "DATABASE_CREDENTIALS_PATH"
const formmatedLinesOutputName string = "alllines.json"
const envDoNotPublish string = "DO_NOT_PUBLISH_REMOTE"

// Publish deploys lines in the correct format in path.
// The format is determined by the presenter.
func Publish(td TransitData, destPath string, p Presenter) error {

	if err := publishLocally(td.lines, destPath, p); err != nil {
		log.Printf("Error publishing lines locally: %v", err)
		return err
	}

	if GetEnvVariableValueBool(envDoNotPublish) {
		log.Printf("Publishing data to remote source")
		if err := publishRemote(td); err != nil {
			log.Printf("Error publishing lines remotely: %v", err)
			return err
		}
	} else {
		log.Printf("Skipped remote publish")
	}

	return nil
}

func publishLocally(lines []Line, destPath string, p Presenter) error {
	log.Printf("Publishing %d lines locally", len(lines))
	os.MkdirAll(destPath, os.ModePerm)

	json, err := p.FormatList(lines)
	if err != nil {
		log.Printf("Error formatting list of lines. Error:%v", err)
		return err
	}

	// Write formatted line as a file in destination
	err = CreateFile(path.Join(destPath, formmatedLinesOutputName), json)
	if err != nil {
		log.Printf("Error creating file for lines. Error:%v", err)
		return err
	}

	return nil
}

// publishRemote reads the json documents generated in the given paths
// and publishes them in remote storage for the clients to consume.
func publishRemote(td TransitData) error {
	ctx := context.Background()
	client, err := getFirebaseClient(ctx)
	if err != nil {
		return err
	}

	if err = postMetadata(ctx, client, td.metadata, "Bilbobus/Metadata"); err != nil {
		return err
	}
	log.Printf("Published remotely version %v", td.metadata)

	return nil
}

func postMetadata(ctx context.Context, c *db.Client, version []MetadataItem, path string) error {
	c.NewRef(path).Delete(ctx)
	if err := c.NewRef(path).Set(ctx, version); err != nil {
		log.Printf("Error publishing version %v : %v", version, err)
		return err
	}
	return nil
}

func postFullLines(ctx context.Context, c *db.Client, lines []Line, path string) error {
	c.NewRef(path).Delete(ctx)
	for _, l := range lines {
		if err := c.NewRef(path+"/"+l.Id).Set(ctx, l); err != nil {
			log.Printf("Error publishing line %v : %v", l.Id, err)
			return err
		}
	}
	return nil
}

func postStopList(ctx context.Context, c *db.Client, stops []Stop, path string) error {
	c.NewRef(path).Delete(ctx)
	if err := c.NewRef(path).Set(ctx, stops); err != nil {
		log.Printf("Error publishing stopList at %v : %v", path, err)
		return err
	}
	return nil
}

func getFirebaseClient(ctx context.Context) (*db.Client, error) {
	config := &firebase.Config{
		DatabaseURL: os.Getenv(envDatabaseUri),
	}

	credentials := os.Getenv(envDatabaseCredentials)
	log.Printf("Credentials file: %v", credentials)
	opt := option.WithCredentialsFile(credentials)
	app, err := firebase.NewApp(ctx, config, opt)
	if err != nil {
		log.Printf("Error creating Firebase App: %v", err)
		return nil, err
	}

	return app.Database(ctx)
}
