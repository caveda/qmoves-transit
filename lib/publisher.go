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
const EnvDatabaseUri string = "DATABASE_URI"
const EnvDatabaseCredentials string = "DATABASE_CREDENTIALS_PATH"

// Publish deploys lines in the correct format in path.
// The format is determined by the presenter.
func Publish(td TransitData, destPath string, p Presenter) error {

	if err := publishLocally(td.lines, destPath, p); err != nil {
		log.Printf("Error publishing lines locally: %v", err)
		return err
	}

/*	if err := publishRemote(td); err != nil {
		log.Printf("Error publishing lines remotely: %v", err)
		return err
	}*/

	return nil
}

func publishLocally(lines []Line, destPath string, p Presenter) error {
	log.Printf("Publishing %d lines locally", len(lines))
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

// publishRemote reads the json documents generated in the given paths
// and publishes them in remote storage for the clients to consume.
func publishRemote(td TransitData) error {
	ctx := context.Background()
	client, err := getFirebaseClient(ctx)
	if err != nil {
		return err
	}

	basePath := "Bilbobus/" + string(td.metadata[0].PathData)

	if err = postFullLines(client, ctx, td.lines, basePath+"/AllLines"); err != nil {
		return err
	}
	log.Printf("Published remotely %v lines", len(td.lines))

	if err = postStopList(client, ctx, td.stops, basePath+"/Stops"); err != nil {
		return err
	}
	log.Printf("Published remotely %v stop list", len(td.stops))

	if err = postMetadata(client, ctx, td.metadata, "Bilbobus/Metadata"); err != nil {
		return err
	}
	log.Printf("Published remotely version %v", td.metadata)

	return nil
}


func postMetadata(c *db.Client, ctx context.Context, version []MetadataItem, path string) error {
	c.NewRef(path).Delete(ctx)
	if err := c.NewRef(path).Set(ctx, version); err != nil {
		log.Printf("Error publishing version %v : %v", version, err)
		return err
	}
	return nil
}

func postFullLines(c *db.Client, ctx context.Context, lines []Line, path string) error {
	c.NewRef(path).Delete(ctx)
	for _, l := range lines {
		if err := c.NewRef(path+"/"+l.Id).Set(ctx, l); err != nil {
			log.Printf("Error publishing line %v : %v", l.Id, err)
			return err
		}
	}
	return nil
}

func postStopList(c *db.Client, ctx context.Context, stops []Stop, path string) error {
	c.NewRef(path).Delete(ctx)
	if err := c.NewRef(path).Set(ctx, stops); err != nil {
		log.Printf("Error publishing stopList at %v : %v", path, err)
		return err
	}
	return nil
}

func getFirebaseClient(ctx context.Context) (*db.Client, error) {
	config := &firebase.Config{
		DatabaseURL: os.Getenv(EnvDatabaseUri),
	}

	credentials := os.Getenv(EnvDatabaseCredentials)
	log.Printf("Credentials file: %v", credentials)
	opt := option.WithCredentialsFile(credentials)
	app, err := firebase.NewApp(ctx, config, opt)
	if err != nil {
		log.Printf("Error creating Firebase App: %v", err)
		return nil, err
	}

	return app.Database(ctx)
}
