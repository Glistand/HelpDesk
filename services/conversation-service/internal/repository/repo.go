package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Conversation struct {
	ID         string
	SiteKey    string
	VisitorID  string
	Status     string
	AssigneeID string
	Preview    string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type Message struct {
	ID             string
	ConversationID string
	Role           string
	Body           string
	AuthorID       string
	CreatedAt      time.Time
}

type Repo struct {
	db *sql.DB
}

func New(db *sql.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) CreateConversation(ctx context.Context, siteKey, visitorID string) (Conversation, Message, error) {
	now := time.Now().UTC()
	c := Conversation{
		ID:        "c-" + uuid.NewString()[:8],
		SiteKey:   siteKey,
		VisitorID: visitorID,
		Status:    "bot",
		Preview:   "Здравствуйте! Чем могу помочь?",
		CreatedAt: now,
		UpdatedAt: now,
	}
	welcome := Message{
		ID:             "m-" + uuid.NewString()[:8],
		ConversationID: c.ID,
		Role:           "bot",
		Body:           "Здравствуйте! Я виртуальный помощник. Опишите проблему — или нажмите «Нужен человек», чтобы связаться с агентом.",
		AuthorID:       "bot",
		CreatedAt:      now,
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return Conversation{}, Message{}, err
	}
	defer func() { _ = tx.Rollback() }()

	_, err = tx.ExecContext(ctx, `
		INSERT INTO conversations (id, site_key, visitor_id, status, assignee_id, preview, created_at, updated_at)
		VALUES ($1,$2,$3,$4,'',$5,$6,$7)`,
		c.ID, c.SiteKey, c.VisitorID, c.Status, c.Preview, c.CreatedAt, c.UpdatedAt)
	if err != nil {
		return Conversation{}, Message{}, err
	}
	_, err = tx.ExecContext(ctx, `
		INSERT INTO messages (id, conversation_id, role, body, author_id, created_at)
		VALUES ($1,$2,$3,$4,$5,$6)`,
		welcome.ID, welcome.ConversationID, welcome.Role, welcome.Body, welcome.AuthorID, welcome.CreatedAt)
	if err != nil {
		return Conversation{}, Message{}, err
	}
	if err := tx.Commit(); err != nil {
		return Conversation{}, Message{}, err
	}
	return c, welcome, nil
}

func (r *Repo) GetConversation(ctx context.Context, id string) (Conversation, error) {
	var c Conversation
	err := r.db.QueryRowContext(ctx, `
		SELECT id, site_key, visitor_id, status, assignee_id, preview, created_at, updated_at
		FROM conversations WHERE id=$1`, id).Scan(
		&c.ID, &c.SiteKey, &c.VisitorID, &c.Status, &c.AssigneeID, &c.Preview, &c.CreatedAt, &c.UpdatedAt)
	if err == sql.ErrNoRows {
		return Conversation{}, fmt.Errorf("not found")
	}
	return c, err
}

func (r *Repo) ListConversations(ctx context.Context, status string, limit int) ([]Conversation, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	var rows *sql.Rows
	var err error
	if status == "" || status == "unspecified" {
		rows, err = r.db.QueryContext(ctx, `
			SELECT id, site_key, visitor_id, status, assignee_id, preview, created_at, updated_at
			FROM conversations ORDER BY updated_at DESC LIMIT $1`, limit)
	} else {
		rows, err = r.db.QueryContext(ctx, `
			SELECT id, site_key, visitor_id, status, assignee_id, preview, created_at, updated_at
			FROM conversations WHERE status=$1 ORDER BY updated_at DESC LIMIT $2`, status, limit)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Conversation
	for rows.Next() {
		var c Conversation
		if err := rows.Scan(&c.ID, &c.SiteKey, &c.VisitorID, &c.Status, &c.AssigneeID, &c.Preview, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

func (r *Repo) AddMessage(ctx context.Context, conversationID, role, body, authorID string) (Message, error) {
	now := time.Now().UTC()
	m := Message{
		ID:             "m-" + uuid.NewString()[:8],
		ConversationID: conversationID,
		Role:           role,
		Body:           body,
		AuthorID:       authorID,
		CreatedAt:      now,
	}
	preview := body
	if len([]rune(preview)) > 120 {
		preview = string([]rune(preview)[:120]) + "…"
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return Message{}, err
	}
	defer func() { _ = tx.Rollback() }()

	_, err = tx.ExecContext(ctx, `
		INSERT INTO messages (id, conversation_id, role, body, author_id, created_at)
		VALUES ($1,$2,$3,$4,$5,$6)`,
		m.ID, m.ConversationID, m.Role, m.Body, m.AuthorID, m.CreatedAt)
	if err != nil {
		return Message{}, err
	}
	_, err = tx.ExecContext(ctx, `
		UPDATE conversations SET preview=$2, updated_at=$3 WHERE id=$1`,
		conversationID, preview, now)
	if err != nil {
		return Message{}, err
	}
	if err := tx.Commit(); err != nil {
		return Message{}, err
	}
	return m, nil
}

func (r *Repo) ListMessages(ctx context.Context, conversationID, afterID string) ([]Message, error) {
	var rows *sql.Rows
	var err error
	if afterID == "" {
		rows, err = r.db.QueryContext(ctx, `
			SELECT id, conversation_id, role, body, author_id, created_at
			FROM messages WHERE conversation_id=$1 ORDER BY created_at ASC`, conversationID)
	} else {
		rows, err = r.db.QueryContext(ctx, `
			SELECT id, conversation_id, role, body, author_id, created_at
			FROM messages
			WHERE conversation_id=$1 AND created_at > (SELECT created_at FROM messages WHERE id=$2)
			ORDER BY created_at ASC`, conversationID, afterID)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Message
	for rows.Next() {
		var m Message
		if err := rows.Scan(&m.ID, &m.ConversationID, &m.Role, &m.Body, &m.AuthorID, &m.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, rows.Err()
}

func (r *Repo) SetStatus(ctx context.Context, id, status, assigneeID string) (Conversation, error) {
	now := time.Now().UTC()
	_, err := r.db.ExecContext(ctx, `
		UPDATE conversations SET status=$2, assignee_id=COALESCE(NULLIF($3,''), assignee_id), updated_at=$4
		WHERE id=$1`, id, status, assigneeID, now)
	if err != nil {
		return Conversation{}, err
	}
	return r.GetConversation(ctx, id)
}

func (r *Repo) RecentForBot(ctx context.Context, conversationID string, limit int) ([]Message, error) {
	if limit <= 0 {
		limit = 20
	}
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, conversation_id, role, body, author_id, created_at FROM (
			SELECT id, conversation_id, role, body, author_id, created_at
			FROM messages WHERE conversation_id=$1
			ORDER BY created_at DESC LIMIT $2
		) t ORDER BY created_at ASC`, conversationID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Message
	for rows.Next() {
		var m Message
		if err := rows.Scan(&m.ID, &m.ConversationID, &m.Role, &m.Body, &m.AuthorID, &m.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, rows.Err()
}

func Truncate(s string, n int) string {
	s = strings.TrimSpace(s)
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n]) + "…"
}
