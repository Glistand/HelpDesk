package domain

import "time"

type Status string

const (
	StatusNew      Status = "new"
	StatusOpen     Status = "open"
	StatusPending  Status = "pending"
	StatusResolved Status = "resolved"
	StatusClosed   Status = "closed"
)

type Priority string

const (
	PriorityLow    Priority = "low"
	PriorityNormal Priority = "normal"
	PriorityHigh   Priority = "high"
	PriorityUrgent Priority = "urgent"
)

type Source string

const (
	SourceManager Source = "manager"
	SourceBot     Source = "bot"
)

type Ticket struct {
	ID             string
	Title          string
	Description    string
	Status         Status
	Priority       Priority
	Category       string
	Requester      string
	AssigneeID     string
	Source         Source
	CreatedByID    string
	ConversationID string
	CreationReason string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type OutboxEvent struct {
	ID            int64
	EventID       string
	AggregateID   string
	EventType     string
	Subject       string
	CorrelationID string
	Payload       []byte
	CreatedAt     time.Time
}
