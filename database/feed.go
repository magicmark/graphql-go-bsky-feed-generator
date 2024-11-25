package database

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const CURSOR_EOF = "eof"

type FeedItem struct {
	URI string `json:"uri"`
}

type Results struct {
	Cursor string     `json:"cursor"`
	Feed   []FeedItem `json:"feed"`
}

type Cursor struct {
	IndexedAt time.Time
	CID       string
}

func parseCursor(cursor string) (*Cursor, error) {
	cursorParts := strings.Split(cursor, "::")
	if len(cursorParts) != 2 {
		return nil, fmt.Errorf("Malformed cursor")
	}

	indexedAt, cid := cursorParts[0], cursorParts[1]
	indexedAtUnixTime, err := strconv.ParseInt(indexedAt, 10, 64)
	if err != nil {
		return nil, err
	}
	indexedAtTime := time.Unix(indexedAtUnixTime/1000, 0)

	return &Cursor{IndexedAt: indexedAtTime, CID: cid}, nil
}

func GetFeedResults(db *sql.DB, cursorString string, limit int) (*Results, error) {
	if cursorString == CURSOR_EOF {
		return &Results{Cursor: CURSOR_EOF, Feed: []FeedItem{}}, nil
	}

	cursor, err := parseCursor(cursorString)
	if err != nil {
		return nil, err
	}

	getPosts := `
		SELECT * FROM post
		WHERE (indexed_at = :1 AND cid < :2) OR (indexed_at < :1)
	`

	rows, err := db.Query(getPosts, cursor.IndexedAt, cursor.CID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	newCursor := CURSOR_EOF
	feed := make([]FeedItem, 0, limit)
	var post Post
	for rows.Next() {
		if err := rows.Scan(&post); err != nil {
			return nil, err
		}

		feed = append(feed, FeedItem{URI: post.URI})
		newCursor = fmt.Sprintf("%d::%s", post.IndexedAt.UnixNano(), post.CID)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &Results{Cursor: newCursor, Feed: feed}, nil
}
