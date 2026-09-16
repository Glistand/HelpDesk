package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Glistand/HelpDesk/libs/eventkit/subjects"
	"github.com/Glistand/HelpDesk/services/ticket-service/internal/domain"
	"github.com/google/uuid"
)

var ErrNotFound = errors.New("ticket not found")
var ErrEmptyTitle = errors.New("title is required")
var ErrEmptyAssignee = errors.New("assignee_id is required")

type TicketRepo struct {
	db *sql.DB
}

func NewTicketRepo(db *sql.DB) *TicketRepo {
	return &TicketRepo{db: db}
}

type CreateInput struct {
	Title         string
	Description   string
	Priority      domain.Priority
	Category      string
	Requester     string
	CorrelationID string
}

func (r *TicketRepo) Create(ctx context.Context, in CreateInput) (domain.Ticket, error) {
	title := strings.TrimSpace(in.Title)
	if title == "" {
		return domain.Ticket{}, ErrEmptyTitle
	}
	if in.Priority == "" {
		in.Priority = domain.PriorityNormal
	}

	now := time.Now().UTC()
	t := domain.Ticket{
		ID:          "t-" + uuid.NewString()[:8],
		Title:       title,
		Description: strings.TrimSpace(in.Description),
		Status:      domain.StatusNew,
		Priority:    in.Priority,
		Category:    strings.TrimSpace(in.Category),
		Requester:   strings.TrimSpace(in.Requester),
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	payload, err := json.Marshal(map[string]any{
		"id":          t.ID,
		"title":       t.Title,
		"description": t.Description,
		"status":      t.Status,
		"priority":    t.Priority,
		"category":    t.Category,
		"requester":   t.Requester,
		"created_at":  t.CreatedAt.Format(time.RFC3339),
	})
	if err != nil {
		return domain.Ticket{}, err
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return domain.Ticket{}, err
	}
	defer func() { _ = tx.Rollback() }()

	_, err = tx.ExecContext(ctx, `
		INSERT INTO tickets (
			id, title, description, status, priority, category, requester, assignee_id, created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
		t.ID, t.Title, t.Description, t.Status, t.Priority, t.Category, t.Requester, t.AssigneeID, t.CreatedAt, t.UpdatedAt,
	)
	if err != nil {
		return domain.Ticket{}, fmt.Errorf("insert ticket: %w", err)
	}

	eventID := uuid.NewString()
	_, err = tx.ExecContext(ctx, `
		INSERT INTO outbox (event_id, aggregate_id, event_type, subject, correlation_id, payload, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7)`,
		eventID, t.ID, "ticket.created", subjects.TicketCreated, in.CorrelationID, payload, now,
	)
	if err != nil {
		return domain.Ticket{}, fmt.Errorf("insert outbox: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return domain.Ticket{}, err
	}
	return t, nil
}

func (r *TicketRepo) Get(ctx context.Context, id string) (domain.Ticket, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, title, description, status, priority, category, requester, assignee_id, created_at, updated_at
		FROM tickets WHERE id = $1`, id)
	t, err := scanTicket(row)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Ticket{}, ErrNotFound
	}
	return t, err
}

func (r *TicketRepo) List(ctx context.Context, status domain.Status, assigneeID string, limit int) ([]domain.Ticket, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	q := `
		SELECT id, title, description, status, priority, category, requester, assignee_id, created_at, updated_at
		FROM tickets WHERE 1=1`
	args := []any{}
	n := 1
	if status != "" {
		q += fmt.Sprintf(" AND status = $%d", n)
		args = append(args, status)
		n++
	}
	if assigneeID != "" {
		q += fmt.Sprintf(" AND assignee_id = $%d", n)
		args = append(args, assigneeID)
		n++
	}
	q += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d", n)
	args = append(args, limit)

	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]domain.Ticket, 0)
	for rows.Next() {
		t, err := scanTicket(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

func (r *TicketRepo) UpdateStatus(ctx context.Context, id string, status domain.Status, correlationID string) (domain.Ticket, error) {
	now := time.Now().UTC()
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return domain.Ticket{}, err
	}
	defer func() { _ = tx.Rollback() }()

	res, err := tx.ExecContext(ctx, `
		UPDATE tickets SET status = $1, updated_at = $2 WHERE id = $3`,
		status, now, id,
	)
	if err != nil {
		return domain.Ticket{}, err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return domain.Ticket{}, ErrNotFound
	}

	row := tx.QueryRowContext(ctx, `
		SELECT id, title, description, status, priority, category, requester, assignee_id, created_at, updated_at
		FROM tickets WHERE id = $1`, id)
	t, err := scanTicket(row)
	if err != nil {
		return domain.Ticket{}, err
	}

	payload, err := json.Marshal(map[string]any{
		"id":         t.ID,
		"status":     t.Status,
		"updated_at": t.UpdatedAt.Format(time.RFC3339),
	})
	if err != nil {
		return domain.Ticket{}, err
	}
	_, err = tx.ExecContext(ctx, `
		INSERT INTO outbox (event_id, aggregate_id, event_type, subject, correlation_id, payload, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7)`,
		uuid.NewString(), t.ID, "ticket.updated", subjects.TicketUpdated, correlationID, payload, now,
	)
	if err != nil {
		return domain.Ticket{}, err
	}
	if err := tx.Commit(); err != nil {
		return domain.Ticket{}, err
	}
	return t, nil
}

func (r *TicketRepo) Assign(ctx context.Context, id, assigneeID, correlationID string) (domain.Ticket, error) {
	assigneeID = strings.TrimSpace(assigneeID)
	if assigneeID == "" {
		return domain.Ticket{}, ErrEmptyAssignee
	}
	now := time.Now().UTC()
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return domain.Ticket{}, err
	}
	defer func() { _ = tx.Rollback() }()

	res, err := tx.ExecContext(ctx, `
		UPDATE tickets
		SET assignee_id = $1, status = CASE WHEN status = 'new' THEN 'open' ELSE status END, updated_at = $2
		WHERE id = $3`,
		assigneeID, now, id,
	)
	if err != nil {
		return domain.Ticket{}, err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return domain.Ticket{}, ErrNotFound
	}

	row := tx.QueryRowContext(ctx, `
		SELECT id, title, description, status, priority, category, requester, assignee_id, created_at, updated_at
		FROM tickets WHERE id = $1`, id)
	t, err := scanTicket(row)
	if err != nil {
		return domain.Ticket{}, err
	}

	payload, err := json.Marshal(map[string]any{
		"id":          t.ID,
		"assignee_id": t.AssigneeID,
		"status":      t.Status,
		"updated_at":  t.UpdatedAt.Format(time.RFC3339),
	})
	if err != nil {
		return domain.Ticket{}, err
	}
	_, err = tx.ExecContext(ctx, `
		INSERT INTO outbox (event_id, aggregate_id, event_type, subject, correlation_id, payload, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7)`,
		uuid.NewString(), t.ID, "ticket.assigned", subjects.TicketAssigned, correlationID, payload, now,
	)
	if err != nil {
		return domain.Ticket{}, err
	}
	if err := tx.Commit(); err != nil {
		return domain.Ticket{}, err
	}
	return t, nil
}

func (r *TicketRepo) ClaimUnpublished(ctx context.Context, limit int) ([]domain.OutboxEvent, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, event_id, aggregate_id, event_type, subject, correlation_id, payload, created_at
		FROM outbox
		WHERE published_at IS NULL
		ORDER BY id
		LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]domain.OutboxEvent, 0)
	for rows.Next() {
		var e domain.OutboxEvent
		if err := rows.Scan(&e.ID, &e.EventID, &e.AggregateID, &e.EventType, &e.Subject, &e.CorrelationID, &e.Payload, &e.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}

func (r *TicketRepo) MarkPublished(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `UPDATE outbox SET published_at = NOW() WHERE id = $1`, id)
	return err
}

type scanner interface {
	Scan(dest ...any) error
}

func scanTicket(s scanner) (domain.Ticket, error) {
	var t domain.Ticket
	var status, priority string
	err := s.Scan(
		&t.ID, &t.Title, &t.Description, &status, &priority,
		&t.Category, &t.Requester, &t.AssigneeID, &t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		return domain.Ticket{}, err
	}
	t.Status = domain.Status(status)
	t.Priority = domain.Priority(priority)
	return t, nil
}
