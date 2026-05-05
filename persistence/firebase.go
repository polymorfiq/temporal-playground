package persistence

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/db"
	"google.golang.org/api/option"
)

const db_url = "https://jobtrackr-9b53c-default-rtdb.firebaseio.com"

func SaveJob(namespace string, key string, job interface{}) error {
	client, err := firebaseDb()
	if err != nil {
		return err
	}

	ref := client.NewRef(fmt.Sprintf("jobs/%s/%s", namespace, key))
	err = ref.Set(context.Background(), job)
	if err != nil {
		return err
	}

	return nil
}

func firebaseDb() (*db.Client, error) {
	app, err := firebaseApp()
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	return app.Database(ctx)
}

func firebaseApp() (*firebase.App, error) {
	opt := option.WithCredentialsFile("/Users/mylaptop/firebase/jobtrackr-9b53c-firebase-adminsdk-fbsvc-f2765f0a8e.json")

	conf := &firebase.Config{
		DatabaseURL: db_url,
	}

	return firebase.NewApp(context.Background(), conf, opt)
}
