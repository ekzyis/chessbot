package db

import (
	"database/sql"
	"log"

	sn "github.com/ekzyis/snappy"
	_ "github.com/mattn/go-sqlite3"
)

var (
	db = getDb()
)

func getDb() *sql.DB {
	var (
		db  *sql.DB
		err error
	)
	if db, err = sql.Open("sqlite3", "chessbot.sqlite3?_foreign_keys=on"); err != nil {
		log.Fatal(err)
	} else {
		if err = migrate(db); err != nil {
			log.Fatal(err)
		}
	}
	return db
}

func migrate(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS items (
			id INTEGER PRIMARY KEY,
			user_id INTEGER NOT NULL,
			text TEXT NOT NULL,
			parent_id INTEGER REFERENCES items(id),
			created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`)

	return err
}

func ItemHasReply(parentId int, userId int) (bool, error) {
	var (
		count int
		err   error
	)

	if err = db.QueryRow(`SELECT COUNT(1) FROM items WHERE parent_id = ? AND user_id = ?`, parentId, userId).Scan(&count); err != nil {
		return true, err
	}

	return count > 0, nil
}

func InsertItem(item *sn.Item) error {
	if _, err := db.Exec(
		`INSERT INTO items(id, user_id, text, parent_id) VALUES (?, ?, ?, NULLIF(?, 0))`,
		item.Id, item.User.Id, item.Text, item.ParentId); err != nil {
		return err
	}

	return nil
}
