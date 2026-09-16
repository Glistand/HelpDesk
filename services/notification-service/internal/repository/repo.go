package repository

import (
	"context"
	"database/sql"
	"errors"
)

type Notification struct {
	EventID   string
	TicketID  string
	Channel   string
	Recipient string
	Subject   string
	Body      string
	Status    string
	Attempts  int
}

type Repo struct {
	db *sql.DB
}

func New(db *sql.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) AlreadyProcessed(ctx context.Context, eventID string) (bool, error) {
	var n int
	err := r.db.QueryRowContext(ctx, `SELECT 1 FROM processed_events WHERE event_id = $1`, eventID).Scan(&n)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return err == nil, err
}

func (r *Repo) Save(ctx context.Context, n Notification) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	_, err = tx.ExecContext(ctx, `
		INSERT INTO notifications (event_id, ticket_id, channel, recipient, subject, body, status, attempts)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		ON CONFLICT (event_id) DO NOTHING`,
		n.EventID, n.TicketID, n.Channel, n.Recipient, n.Subject, n.Body, n.Status, n.Attempts,
	)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, `
		INSERT INTO processed_events (event_id) VALUES ($1)
		ON CONFLICT (event_id) DO NOTHING`, n.EventID)
	if err != nil {
		return err
	}
	return tx.Commit()
}
