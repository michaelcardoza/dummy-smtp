package sqlite

import "database/sql"

func Migrate(db *sql.DB) error {
	_, err := db.Exec(`
		PRAGMA foreign_keys = ON;
		CREATE TABLE IF NOT EXISTS messages
		(
			id TEXT PRIMARY KEY,
			from_addr TEXT NOT NULL,
			subject TEXT,
			text_body TEXT,
			html_body TEXT,
			raw TEXT,
			created_at TIMESTAMP NOT NULL
		);
		CREATE TABLE IF NOT EXISTS recipients
		(
			message_id TEXT NOT NULL REFERENCES messages(id) ON DELETE CASCADE,
		    position INTEGER NOT NULL,
		    addr NOT NULL,
		    PRIMARY KEY(message_id, position)
		);
		CREATE INDEX IF NOT EXISTS idx_recipients_addr ON recipients(addr);
	`)
	if err != nil {
		return err
	}

	return nil
}
