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

	if err := publishRemote(td); err != nil {
		log.Printf("Error publishing lines remotely: %v", err)
		return err
	}

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

	basePath := "bilbobus/" + string(td.version.ver)

	if err = postFullLines(client, ctx, td.lines, basePath+"/allLines"); err != nil {
		return err
	}
	log.Printf("Published remotely %v full lines", len(td.lines))

	if err = postLinesList(client, ctx, td.dayLines, basePath+"/dayLines"); err != nil {
		return err
	}
	log.Printf("Published remotely %v day lines", len(td.dayLines))

	if err = postLinesList(client, ctx, td.nightLines, basePath+"/nightLines"); err != nil {
		return err
	}
	log.Printf("Published remotely %v night lines", len(td.nightLines))

	if err = postVersion(client, ctx, td.version, basePath+"/ver"); err != nil {
		return err
	}
	log.Printf("Published remotely version %v", td.version)

	return nil
}

func postLinesList(c *db.Client, ctx context.Context, lines []Line, path string) error {
	if err := c.NewRef(path).Set(ctx, lines); err != nil {
		log.Printf("Error publishing lineList at %v : %v", path, err)
		return err
	}
	return nil
}

func postVersion(c *db.Client, ctx context.Context, version Version, path string) error {
	var vs []Version
	vs = append(vs, version)
	if err := c.NewRef(path).Set(ctx, vs); err != nil {
		log.Printf("Error publishing version %v : %v", vs, err)
		return err
	}
	return nil
}

func postFullLines(c *db.Client, ctx context.Context, lines []Line, path string) error {
	for _, l := range lines {
		if err := c.NewRef(path+"/"+l.Id).Set(ctx, l); err != nil {
			log.Printf("Error publishing line %v : %v", l.Id, err)
			return err
		}
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
