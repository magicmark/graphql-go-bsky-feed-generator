package main

import (
	"go-bsky-feed/database"
	"go-bsky-feed/firehose"
	"go-bsky-feed/webserver"
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
