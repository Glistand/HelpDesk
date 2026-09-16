package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var ErrNotFound = errors.New("escalation not found")

type Escalation struct {
	TicketID     string
	Reason       string
	FromAssignee string
	ToAssignee   string
	EscalatedAt  time.Time
	SourceEvent  string
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

func (r *Repo) Get(ctx context.Context, ticketID string) (Escalation, error) {
	var e Escalation
	err := r.db.QueryRowContext(ctx, `
		SELECT ticket_id, reason, from_assignee, to_assignee, escalated_at, source_event
		FROM escalations WHERE ticket_id = $1`, ticketID).Scan(
		&e.TicketID, &e.Reason, &e.FromAssignee, &e.ToAssignee, &e.EscalatedAt, &e.SourceEvent,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return Escalation{}, ErrNotFound
	}
	return e, err
}

func (r *Repo) Save(ctx context.Context, e Escalation) (bool, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return false, err
	}
	defer func() { _ = tx.Rollback() }()

	res, err := tx.ExecContext(ctx, `
		INSERT INTO escalations (ticket_id, reason, from_assignee, to_assignee, escalated_at, source_event)
		VALUES ($1,$2,$3,$4,$5,$6)
		ON CONFLICT (ticket_id) DO NOTHING`,
		e.TicketID, e.Reason, e.FromAssignee, e.ToAssignee, e.EscalatedAt, e.SourceEvent,
	)
	if err != nil {
		return false, err
	}
	n, _ := res.RowsAffected()
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO processed_events (event_id) VALUES ($1)
		ON CONFLICT (event_id) DO NOTHING`, e.SourceEvent); err != nil {
		return false, err
	}
	if err := tx.Commit(); err != nil {
		return false, err
	}
	return n > 0, nil
}

func (r *Repo) MarkProcessed(ctx context.Context, eventID string) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO processed_events (event_id) VALUES ($1)
		ON CONFLICT (event_id) DO NOTHING`, eventID)
	return err
}
