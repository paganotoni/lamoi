package messages

import (
	"github.com/leapkit/leapkit/core/db"
	"github.com/oklog/ulid/v2"
)

func NewService(fn db.ConnFn) *service {
	return &service{
		db: fn,
	}
}

type service struct {
	db db.ConnFn
}

func (s *service) ContentOf(id string) (string, error) {
	conn, err := s.db()
	if err != nil {
		return "", err
	}

	var message string
	err = conn.QueryRow("SELECT message FROM messages WHERE id = ?", id).Scan(&message)
	if err != nil {
		return "", err
	}

	return message, nil
}

func (s *service) IsComplete(id string) (bool, error) {
	conn, err := s.db()
	if err != nil {
		return false, err
	}

	var complete bool
	err = conn.QueryRow("SELECT complete FROM messages WHERE id = ?", id).Scan(&complete)
	if err != nil {
		return false, err
	}

	return complete, nil
}

func (s *service) AppendTo(conversationID, message, kind, context string) (string, error) {
	conn, err := s.db()
	if err != nil {
		return "", err
	}

	id := ulid.Make()
	_, err = conn.Exec("INSERT INTO messages (id, conversation_id, message, kind) VALUES (?, ?, ?, ?)", id.String(), conversationID, message, kind)
	if err != nil {
		return "", err
	}

	if context != "" {
		_, err = conn.Exec("UPDATE conversations SET context = ? WHERE id = ?", context, conversationID)
		if err != nil {
			return "", err
		}
	}

	return id.String(), nil
}

func (s *service) AppendPendingTo(conversationID string) (string, error) {
	conn, err := s.db()
	if err != nil {
		return "", err
	}

	id := ulid.Make()
	_, err = conn.Exec("INSERT INTO messages (id, conversation_id, message, kind, complete) VALUES (?, ?, ?, ?, ?)", id.String(), conversationID, "", "ollama", false)
	if err != nil {
		return "", err
	}

	return id.String(), nil
}
