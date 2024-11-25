package database

import (
    "database/sql"
    "log"
)

func InsertPost(db *sql.DB, post *Post) {
	insertPost := `
	INSERT INTO post(uri, cid, reply_parent, reply_root, indexed_at)
	VALUES (?, ?, ?, ?, ?)
	`

	_, err := db.Exec(insertPost,
		post.URI,
		post.CID,
		post.ReplyParent,
		post.ReplyRoot,
		post.IndexedAt.Format("2006-01-02 15:04:05.999999"),
	)

	if err != nil {
		log.Fatal(err)
	}
}

func DeletePost(db *sql.DB, uri string) {
	_, err := db.Exec(`DELETE FROM post WHERE uri = ?`, uri)

	if err != nil {
		log.Fatal(err)
	}
}