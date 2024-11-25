package database

import (
    "database/sql"
    "log"
	"os"
    "time"
    _ "github.com/mattn/go-sqlite3"
)

type Post struct {
    URI         string
    CID         string
    ReplyParent sql.NullString
    ReplyRoot   sql.NullString
    IndexedAt   time.Time
}

type SubscriptionState struct {
    Service string
    Cursor  int64
}

func getFileLocation() string {
	// In production, set DBLOCATION environment variable to a persistent volume location (e.g. /db/feed_database.db)
    if dbLocation, exists := os.LookupEnv("DBLOCATION"); exists {
        return dbLocation
    }

	// in local dev, write to a local file
    return "feed_database.db"
}

func CreateTables(db *sql.DB) {
    _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS "post" (
			"uri" VARCHAR(255) PRIMARY KEY,
			"cid" VARCHAR(255) NOT NULL,
			"reply_parent" VARCHAR(255),
			"reply_root" VARCHAR(255),
			"indexed_at" DATETIME NOT NULL
		);
	
		CREATE TABLE IF NOT EXISTS "subscriptionstate" (
			"service" VARCHAR(255) PRIMARY KEY,
			"cursor" INTEGER NOT NULL
		);
	`)

    if err != nil {
        log.Fatal(err)
    }
}

func GetConnection() (*sql.DB, error) {
    return sql.Open("sqlite3", getFileLocation())
}
