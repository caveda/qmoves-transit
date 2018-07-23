package transit

import (
	"log"
	"os"
	"path"

	firebase "firebase.google.com/go"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
)

// Constants
const EnvDatabaseUri string = "DATABASE_URI"
const EnvDatabaseCredentials string = "DATABASE_CREDENTIALS_PATH"

// Publish deploys lines in the correct format in path.
// The format is determined by the presenter.
func Publish(lines []Line, destPath string, p Presenter) error {
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

	log.Printf("Publishing %d lines remotely", len(lines))
	publishToFirebase(lines)

	return nil
}

// PublishFirebase reads the json documents generated in the given paths
// and publishes them in Firebase Cloud storage for the clients to consume.
func publishToFirebase(lines []Line) error {
	ctx := context.Background()
	config := &firebase.Config{
		DatabaseURL: os.Getenv(EnvDatabaseUri),
	}

	credentials := os.Getenv(EnvDatabaseCredentials)
	log.Printf("Credentials file: %v", credentials)
	opt := option.WithCredentialsFile(credentials)
	app, err := firebase.NewApp(ctx, config, opt)
	if err != nil {
		log.Fatal(err)
	}

	client, err := app.Database(ctx)
	if err != nil {
		log.Fatal(err)
	}

	for _, l := range lines {
		// Write formatted line as a file in destination
		if err := client.NewRef("bilbobus/"+l.Id).Set(ctx, l); err != nil {
			log.Fatal(err)
		}
	}

	return nil
}
