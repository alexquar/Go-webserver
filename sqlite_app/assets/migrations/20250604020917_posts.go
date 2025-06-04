package migrations

//create w goose -dir=assets/migrations create posts
import (
	"context"
	"database/sql"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upPosts, downPosts)
}

func upPosts(ctx context.Context, tx *sql.Tx) error {
	if _, err := tx.Exec(`CREATE TABLE posts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
		)`); err != nil {
		return err
	}

	// This code is executed when the migration is applied.
	return nil
}

func downPosts(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}
