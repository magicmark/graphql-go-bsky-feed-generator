package firehose

import (
	"log"
	"encoding/json"
	"github.com/gorilla/websocket"
	"go-bsky-feed/database"
	"go-bsky-feed/algorithm"
    "database/sql"
	"time"
)

func createHandler(db *sql.DB, repoMessage *RepoMessage) {
	var record Record
	if err := json.Unmarshal([]byte(repoMessage.Commit.Record), &record); err != nil {
		log.Fatal(err)
	}

	if algorithm.PostFilterMatch(&record.Text) == true {
		post := &database.Post{
			URI: repoMessage.URI(),
			CID: repoMessage.Commit.CID,
			IndexedAt: time.Now().UTC(),
		}

		if record.Reply.Parent.URI != "" {
			post.ReplyParent = sql.NullString{String: record.Reply.Parent.URI, Valid: true}
			post.ReplyRoot = sql.NullString{String: record.Reply.Root.URI, Valid: true}
		}

		log.Println("Found post! Inserting: %s", record.Text)
		database.InsertPost(db, post)
	}
}

func SubscribeToFirehose() {
	uri := "wss://jetstream1.us-west.bsky.network/subscribe?wantedCollections=app.bsky.feed.post"
	log.Printf("connecting to %s", uri)

	c, _, err := websocket.DefaultDialer.Dial(uri, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	db, err := database.GetConnection()
	if err != nil {
		log.Fatal("couldn't get database connection:", err)
	}

	repoMessage := new(RepoMessage)

	for {
		_, message, err := c.ReadMessage()

		if err != nil {
			log.Fatal(err)
		}

		if err := json.Unmarshal(message, &repoMessage); err != nil {
			log.Fatal(err)
		}

		if (repoMessage.Commit.Operation == "create") {
			// we only filter for app.bsky.feed.post in the wss url, so we know it's always a post record type.
			createHandler(db, repoMessage)
		} else if (repoMessage.Commit.Operation == "delete") {
			database.DeletePost(db, repoMessage.URI())
		}
	}
}
