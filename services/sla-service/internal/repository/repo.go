package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var ErrNotFound = errors.New("sla not found")

type TicketSLA struct {
	TicketID         string
	State            string
	Policy           string
	FirstResponseDue time.Time
	ResolveDue       time.Time
	WarnedAt         *time.Time
	BreachedAt       *time.Time
	CancelledAt      *time.Time
	CorrelationID    string
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

func (r *Repo) MarkProcessed(ctx context.Context, eventID string) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO processed_events (event_id) VALUES ($1)
		ON CONFLICT (event_id) DO NOTHING`, eventID)
	return err
}

func (r *Repo) Get(ctx context.Context, ticketID string) (TicketSLA, error) {
	var s TicketSLA
	var warned, breached, cancelled sql.NullTime
	err := r.db.QueryRowContext(ctx, `
		SELECT ticket_id, state, policy, first_response_due, resolve_due,
		       warned_at, breached_at, cancelled_at, correlation_id
		FROM ticket_sla WHERE ticket_id = $1`, ticketID).Scan(
		&s.TicketID, &s.State, &s.Policy, &s.FirstResponseDue, &s.ResolveDue,
		&warned, &breached, &cancelled, &s.CorrelationID,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return TicketSLA{}, ErrNotFound
	}
	if err != nil {
		return TicketSLA{}, err
	}
	if warned.Valid {
		t := warned.Time
		s.WarnedAt = &t
	}
	if breached.Valid {
		t := breached.Time
		s.BreachedAt = &t
	}
	if cancelled.Valid {
		t := cancelled.Time
		s.CancelledAt = &t
	}
	return s, nil
}

func (r *Repo) Upsert(ctx context.Context, s TicketSLA) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO ticket_sla (
			ticket_id, state, policy, first_response_due, resolve_due, correlation_id, created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,NOW(),NOW())
		ON CONFLICT (ticket_id) DO UPDATE SET
			state = EXCLUDED.state,
			policy = EXCLUDED.policy,
			first_response_due = EXCLUDED.first_response_due,
			resolve_due = EXCLUDED.resolve_due,
			correlation_id = EXCLUDED.correlation_id,
			warned_at = NULL,
			breached_at = NULL,
			cancelled_at = NULL,
			updated_at = NOW()`,
		s.TicketID, s.State, s.Policy, s.FirstResponseDue, s.ResolveDue, s.CorrelationID,
	)
	return err
}

func (r *Repo) MarkWarned(ctx context.Context, ticketID string, at time.Time) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE ticket_sla
		SET state = CASE WHEN state = 'ok' THEN 'warning' ELSE state END,
		    warned_at = COALESCE(warned_at, $2),
		    updated_at = NOW()
		WHERE ticket_id = $1 AND cancelled_at IS NULL`, ticketID, at)
	return err
}

func (r *Repo) MarkBreached(ctx context.Context, ticketID string, at time.Time) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE ticket_sla
		SET state = 'breached',
		    breached_at = COALESCE(breached_at, $2),
		    updated_at = NOW()
		WHERE ticket_id = $1 AND cancelled_at IS NULL`, ticketID, at)
	return err
}

func (r *Repo) Cancel(ctx context.Context, ticketID string, at time.Time) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE ticket_sla
		SET state = 'cancelled', cancelled_at = $2, updated_at = NOW()
		WHERE ticket_id = $1 AND cancelled_at IS NULL`, ticketID, at)
	return err
}

func (r *Repo) MarkFired(ctx context.Context, ticketID, kind string) (bool, error) {
	res, err := r.db.ExecContext(ctx, `
		INSERT INTO fired_timers (ticket_id, kind) VALUES ($1,$2)
		ON CONFLICT (ticket_id, kind) DO NOTHING`, ticketID, kind)
	if err != nil {
		return false, err
	}
	n, _ := res.RowsAffected()
	return n > 0, nil
}

func (r *Repo) ClearFired(ctx context.Context, ticketID string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM fired_timers WHERE ticket_id = $1`, ticketID)
	return err
}
