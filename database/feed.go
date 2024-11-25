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
		return nil, fmt.Errorf("error parsing indexedAt int: %w", err)
	}
	indexedAtTime := time.Unix(indexedAtUnixTime/1000, 0)

	return &Cursor{IndexedAt: indexedAtTime, CID: cid}, nil
}

func GetFeedResults(db *sql.DB, cursorString string, limit int) (*Results, error) {
	if cursorString == CURSOR_EOF {
		return &Results{Cursor: CURSOR_EOF, Feed: []FeedItem{}}, nil
	}

	cursor := &Cursor{}
	query := ""

	if cursorString == "" {
		query = `
			SELECT * FROM post
			ORDER BY indexed_at DESC
			LIMIT :limit
		`
	} else {
		query = `
			SELECT * FROM post
			WHERE (indexed_at = :indexed_at AND cid < :cid) OR (indexed_at < :indexed_at)
			LIMIT :limit
		`
		var err error
		cursor, err = parseCursor(cursorString)
		if err != nil {
			return nil, fmt.Errorf("Cursor parse error: %w", err)
		}
	}

	rows, err := db.Query(
		query,
		sql.Named("indexed_at", cursor.IndexedAt),
		sql.Named("cid", cursor.CID),
		sql.Named("limit", limit),
	)
	if err != nil {
		return nil, fmt.Errorf("db query error: %w", err)
	}
	defer rows.Close()

	newCursor := CURSOR_EOF
	feed := make([]FeedItem, 0, limit)
	var post Post

	for rows.Next() {
		if err := rows.Scan(
			&post.URI,
			&post.CID,
			&post.ReplyParent,
			&post.ReplyRoot,
			&post.IndexedAt,
		); err != nil {
			return nil, fmt.Errorf("row scan error: %w", err)
		}

		feed = append(feed, FeedItem{URI: post.URI})
		newCursor = fmt.Sprintf("%d::%s", post.IndexedAt.UnixNano(), post.CID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error from rows.Err(): %w", err)
	}

	return &Results{Cursor: newCursor, Feed: feed}, nil
}
