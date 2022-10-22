package main

import (
	"energi-challenge/console/indexer/api"
	"log"
)

func main() {
	app, err := api.NewIndexer()
	if err != nil {
		log.Fatalf("error creating new api: %s", err)
		return
	}

	err = app.Start()

	if err != nil {
		log.Fatalf("error starting the api: %s", err)
		return
	}

	app.Shutdown()

}
