package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"strings"

	"github.com/michaelcardoza/dummy-smtp/internal/core/mail"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{db: db}
}

func (s *Storage) Save(ctx context.Context, message *mail.Message) error {
	return WithTransaction(ctx, s.db, func(tx *sql.Tx) error {
		stmp := `
		INSERT INTO messages (id, from_addr, subject, text_body, html_body, raw, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		`
		_, err := tx.ExecContext(ctx, stmp,
			message.ID,
			message.From,
			message.Subject,
			message.TextBody,
			message.HTMLBody,
			message.Raw,
			message.CreatedAt,
		)
		if err != nil {
			return err
		}

		if len(message.To) > 0 {
			placeholders := make([]string, 0, len(message.To))
			recipients := make([]any, 0, len(message.To)*3)

			for i, addr := range message.To {
				placeholders = append(placeholders, "(?, ?, ?)")
				recipients = append(recipients, message.ID, i, addr)
			}

			stmp = "INSERT INTO recipients (message_id, position, addr) VALUES" + strings.Join(placeholders, ",")
			_, err = tx.ExecContext(ctx, stmp, recipients...)
			if err != nil {
				return err
			}
		}

		if len(message.Attachments) > 0 {
			placeholders := make([]string, 0, len(message.Attachments))
			attachments := make([]any, 0, len(message.Attachments)*4)
			for i, a := range message.Attachments {
				placeholders = append(placeholders, "(?, ?, ?, ?, ?)")
				attachments = append(attachments, message.ID, i, a.Filename, a.ContentType, a.Size)
			}

			stmp = "INSERT INTO attachments (message_id, position, filename, content_type, size) VALUES" + strings.Join(placeholders, ",")
			if _, err = tx.ExecContext(ctx, stmp, attachments...); err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *Storage) List(ctx context.Context) ([]*mail.Message, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT
		    m.id,
		    m.from_addr,
		    json_group_array(s.addr ORDER BY s.position) FILTER (WHERE s.addr IS NOT NULL) AS to_recipients,
		    COALESCE((
		        SELECT json_group_array(
		            json_object('filename', a.filename, 'contentType', a.content_type, 'size', a.size)
		        )
		        FROM (SELECT * FROM attachments WHERE message_id = m.id ORDER BY position) a
		    ), '[]') AS attachments,
		    m.subject,
		    m.text_body,
		    m.html_body,
		    m.raw,
		    m.created_at
		FROM messages m
		LEFT JOIN recipients s ON s.message_id = m.id
		GROUP BY m.id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := make([]*mail.Message, 0)
	for rows.Next() {
		var message mail.Message
		var recipients string
		var attachments string
		err = rows.Scan(
			&message.ID,
			&message.From,
			&recipients,
			&attachments,
			&message.Subject,
			&message.TextBody,
			&message.HTMLBody,
			&message.Raw,
			&message.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		json.Unmarshal([]byte(recipients), &message.To)
		json.Unmarshal([]byte(attachments), &message.Attachments)

		messages = append(messages, &message)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

func (s *Storage) Get(ctx context.Context, id string) (*mail.Message, error) {
	smtp := `
		SELECT
		    m.id,
		    m.from_addr,
		    json_group_array(s.addr ORDER BY s.position) FILTER (WHERE s.addr IS NOT NULL) AS to_recipients,
		    COALESCE((
		        SELECT json_group_array(
		            json_object('filename', a.filename, 'contentType', a.content_type, 'size', a.size)
		        )
		        FROM (SELECT * FROM attachments WHERE message_id = m.id ORDER BY position) a
		    ), '[]') AS attachments,
		    m.subject,
		    m.text_body,
		    m.html_body,
		    m.raw,
		    m.created_at
		FROM messages m
		LEFT JOIN recipients s ON s.message_id = m.id
		WHERE m.id = ?
	`
	row := s.db.QueryRowContext(ctx, smtp, id)

	var message mail.Message
	var recipients string
	var attachments string
	err := row.Scan(
		&message.ID,
		&message.From,
		&recipients,
		&attachments,
		&message.Subject,
		&message.TextBody,
		&message.HTMLBody,
		&message.Raw,
		&message.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	json.Unmarshal([]byte(recipients), &message.To)
	json.Unmarshal([]byte(attachments), &message.Attachments)

	return &message, nil
}

func (s *Storage) DeleteByID(ctx context.Context, id string) error {
	stmp := "DELETE FROM messages WHERE id = ?"
	if _, err := s.db.ExecContext(ctx, stmp, id); err != nil {
		return err
	}
	return nil
}

func (s *Storage) DeleteAll(ctx context.Context) error {
	stmp := "DELETE FROM messages"
	_, err := s.db.ExecContext(ctx, stmp)
	return err
}

func (s *Storage) Close(_ context.Context) error {
	return s.db.Close()
}
