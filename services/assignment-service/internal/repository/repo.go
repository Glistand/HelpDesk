package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var ErrNotFound = errors.New("assignment not found")

type Agent struct {
	ID   string
	Name string
	Role string
}

type Assignment struct {
	TicketID     string
	AssigneeID   string
	AssigneeName string
	AssignedAt   time.Time
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

func (r *Repo) GetByTicket(ctx context.Context, ticketID string) (Assignment, error) {
	var a Assignment
	err := r.db.QueryRowContext(ctx, `
		SELECT ticket_id, assignee_id, assignee_name, assigned_at, source_event
		FROM assignments WHERE ticket_id = $1`, ticketID).Scan(
		&a.TicketID, &a.AssigneeID, &a.AssigneeName, &a.AssignedAt, &a.SourceEvent,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return Assignment{}, ErrNotFound
	}
	return a, err
}

// NextL1Agent picks the next L1 agent round-robin and advances the counter.
func (r *Repo) NextL1Agent(ctx context.Context) (Agent, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return Agent{}, err
	}
	defer func() { _ = tx.Rollback() }()

	var idx int
	if err := tx.QueryRowContext(ctx, `
		SELECT idx FROM rr_state WHERE queue = 'l1' FOR UPDATE`).Scan(&idx); err != nil {
		return Agent{}, err
	}

	rows, err := tx.QueryContext(ctx, `
		SELECT id, name, role FROM agents WHERE role LIKE 'L1%' ORDER BY id`)
	if err != nil {
		return Agent{}, err
	}
	defer rows.Close()

	var agents []Agent
	for rows.Next() {
		var a Agent
		if err := rows.Scan(&a.ID, &a.Name, &a.Role); err != nil {
			return Agent{}, err
		}
		agents = append(agents, a)
	}
	if err := rows.Err(); err != nil {
		return Agent{}, err
	}
	if len(agents) == 0 {
		return Agent{}, errors.New("no L1 agents")
	}

	picked := agents[idx%len(agents)]
	next := (idx + 1) % len(agents)
	if _, err := tx.ExecContext(ctx, `UPDATE rr_state SET idx = $1 WHERE queue = 'l1'`, next); err != nil {
		return Agent{}, err
	}
	if err := tx.Commit(); err != nil {
		return Agent{}, err
	}
	return picked, nil
}

// SaveAssignment stores the assignment and marks the source event processed.
// Returns false if the ticket was already assigned (idempotent).
func (r *Repo) SaveAssignment(ctx context.Context, a Assignment) (bool, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return false, err
	}
	defer func() { _ = tx.Rollback() }()

	res, err := tx.ExecContext(ctx, `
		INSERT INTO assignments (ticket_id, assignee_id, assignee_name, assigned_at, source_event)
		VALUES ($1,$2,$3,$4,$5)
		ON CONFLICT (ticket_id) DO NOTHING`,
		a.TicketID, a.AssigneeID, a.AssigneeName, a.AssignedAt, a.SourceEvent,
	)
	if err != nil {
		return false, err
	}
	n, _ := res.RowsAffected()
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO processed_events (event_id) VALUES ($1)
		ON CONFLICT (event_id) DO NOTHING`, a.SourceEvent); err != nil {
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
