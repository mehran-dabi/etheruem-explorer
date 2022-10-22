package main

import (
	"energi-challenge/console/rest/api"
	"log"
)

func main() {
	app, err := api.NewAPI()
	if err != nil {
		log.Fatalf("failed to create new REST APIs: %s", err)
	}

	err = app.Start()
	if err != nil {
		log.Fatalf("failed to start REST APIs: %s", err)
		return
	}

	app.Shutdown()
}
