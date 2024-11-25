package main

import (
	"go-bsky-feed/app/database"
	"go-bsky-feed/app/firehose"
	"go-bsky-feed/app/webserver"
	"log"
)

func main() {
	db, err := database.GetConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	database.CreateTables(db)
	go firehose.SubscribeToFirehose()
	webserver.ServeHTTP(db)
}
