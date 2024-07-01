package conversations

import (
	"github.com/leapkit/leapkit/core/db"
	"github.com/oklog/ulid/v2"
)

type conversation struct {
	ID       string
	Name     string
	Messages []message
}

type message struct {
	Message string
	Kind    string
}

func NewService(fn db.ConnFn) *service {
	return &service{
		db: fn,
	}
}

type service struct {
	db db.ConnFn
}

// Create creates a new conversation and returns the ID
// of the conversation
func (s *service) Create(message string) (string, error) {
	conn, err := s.db()
	if err != nil {
		return "", err
	}

	id := ulid.Make()
	_, err = conn.Exec("INSERT INTO conversations (id, name) VALUES (?, ?)", id.String(), message)
	if err != nil {
		return "", err
	}

	_, err = conn.Exec("INSERT INTO messages (conversation_id, message) VALUES (?, ?)", id.String(), message)
	if err != nil {
		return "", err
	}

	return id.String(), nil
}

func (s *service) Find(id string) (*conversation, error) {
	conn, err := s.db()
	if err != nil {
		return nil, err
	}

	rows, err := conn.Query("SELECT message, kind FROM messages WHERE conversation_id = ?", id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	conv := &conversation{
		Name:     "",
		Messages: []message{},
	}

	for rows.Next() {
		var m message
		err = rows.Scan(&m.Message, &m.Kind)
		if err != nil {
			return nil, err
		}

		conv.Messages = append(conv.Messages, m)
	}

	err = conn.QueryRow("SELECT name FROM conversations WHERE id = ?", id).Scan(&conv.Name)
	if err != nil {
		return nil, err
	}

	return conv, nil
}

func (s *service) ContextFor(id string) (string, error) {
	conn, err := s.db()
	if err != nil {
		return "", err
	}

	var context string
	err = conn.QueryRow("SELECT context FROM conversations WHERE id = ?", id).Scan(&context)
	if err != nil {
		return "", err
	}

	return context, nil
}

func (s *service) List() ([]conversation, error) {
	conn, err := s.db()
	if err != nil {
		return nil, err
	}

	rows, err := conn.Query("SELECT id, name FROM conversations ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	conversations := []conversation{}
	for rows.Next() {
		var c conversation
		err = rows.Scan(&c.ID, &c.Name)
		if err != nil {
			return nil, err
		}

		conversations = append(conversations, c)
	}

	return conversations, nil
}

func (s *service) Update(messageID, message, context string) error {
	conn, err := s.db()
	if err != nil {
		return err
	}

	_, err = conn.Exec("UPDATE messages SET message = ? WHERE id = ?", message, messageID)
	if err != nil {
		return err
	}

	if context != "" {
		_, err = conn.Exec("UPDATE conversations SET context = ? WHERE id = (SELECT conversation_id FROM messages WHERE id = ?)", context, messageID)
		if err != nil {
			return err
		}

		_, err = conn.Exec("UPDATE messages SET complete = ? WHERE id = ?", true, messageID)
		if err != nil {
			return err
		}
	}

	return nil
}
