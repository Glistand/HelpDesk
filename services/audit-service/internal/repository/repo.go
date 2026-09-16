package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type TimelineEvent struct {
	ID         string
	TicketID   string
	EventID    string
	EventType  string
	Title      string
	Detail     string
	Actor      string
	OccurredAt time.Time
}

type Repo struct {
	db *sql.DB
}

func New(db *sql.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) Insert(ctx context.Context, e TimelineEvent) (bool, error) {
	res, err := r.db.ExecContext(ctx, `
		INSERT INTO timeline_events (id, ticket_id, event_id, event_type, title, detail, actor, occurred_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		ON CONFLICT (event_id) DO NOTHING`,
		e.ID, e.TicketID, e.EventID, e.EventType, e.Title, e.Detail, e.Actor, e.OccurredAt,
	)
	if err != nil {
		return false, err
	}
	n, _ := res.RowsAffected()
	return n > 0, nil
}

func (r *Repo) ListByTicket(ctx context.Context, ticketID string) ([]TimelineEvent, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, ticket_id, event_id, event_type, title, detail, actor, occurred_at
		FROM timeline_events
		WHERE ticket_id = $1
		ORDER BY occurred_at ASC, created_at ASC`, ticketID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]TimelineEvent, 0)
	for rows.Next() {
		var e TimelineEvent
		if err := rows.Scan(
			&e.ID, &e.TicketID, &e.EventID, &e.EventType, &e.Title, &e.Detail, &e.Actor, &e.OccurredAt,
		); err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}

func (r *Repo) Exists(ctx context.Context, eventID string) (bool, error) {
	var n int
	err := r.db.QueryRowContext(ctx, `SELECT 1 FROM timeline_events WHERE event_id = $1`, eventID).Scan(&n)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return err == nil, err
}
