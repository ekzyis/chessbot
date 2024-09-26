package db

import (
	"database/sql"
	"errors"
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

	if count > 0 {
		return true, nil
	}

	// check if parent already exists, this means we ignored it
	if err = db.QueryRow(`SELECT COUNT(1) FROM items WHERE id = ?`, parentId).Scan(&count); err != nil {
		return true, err
	}

	if count > 0 {
		log.Printf("ignoring known item %d", parentId)
	}

	return count > 0, nil
}

func InsertItem(item *sn.Item) error {
	if _, err := db.Exec(``+
		`INSERT INTO items(id, user_id, text, parent_id) VALUES (?, ?, ?, NULLIF(?, 0)) `+
		`ON CONFLICT DO UPDATE SET text = EXCLUDED.text, updated_at = CURRENT_TIMESTAMP`,
		item.Id, item.User.Id, item.Text, item.ParentId); err != nil {
		return err
	}

	return nil
}

func GetThread(id int) ([]sn.Item, error) {
	var (
		items []sn.Item
		err   error
	)

	var item sn.Item
	item.ParentId = id

	for item.ParentId > 0 {
		// TODO: can't select created_at because sqlite3 doesn't support timestamps natively
		// see https://github.com/mattn/go-sqlite3/issues/142
		if err = db.QueryRow(
			`SELECT id, user_id, text, COALESCE(parent_id, 0) FROM items WHERE id = ?`, item.ParentId).
			Scan(&item.Id, &item.User.Id, &item.Text, &item.ParentId); err != nil {
			if err == sql.ErrNoRows {
				return nil, errors.New("item not found in db")
			}
			return nil, err
		}

		items = append([]sn.Item{item}, items...)
	}

	return items, nil
}
