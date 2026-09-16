package store

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

var ErrEmptyTitle = errors.New("title is required")
var ErrNotFound = errors.New("ticket not found")

type Store struct {
	nextID int
	byID   map[string]Ticket
}

type Ticket struct {
	ID        string
	Title     string
	Status    string
	CreatedAt time.Time
}

func NewStore() *Store {
	return &Store{
		nextID: 0,
		byID:   make(map[string]Ticket),
	}
}

func (s *Store) Create(title string) (Ticket, error) {
	title = strings.TrimSpace(title)

	if title == "" {
		return Ticket{}, ErrEmptyTitle
	}

	s.nextID++
	id := fmt.Sprintf("t-%d", s.nextID)

	ticket := Ticket{
		ID:        id,
		Title:     title,
		Status:    "new",
		CreatedAt: time.Now().UTC(),
	}

	s.byID[id] = ticket

	return ticket, nil
}

func (s *Store) Get(id string) (Ticket, error) {
	t, ok := s.byID[id]
	if !ok {
		return Ticket{}, ErrNotFound
	}

	return t, nil
}

func (s *Store) List() []Ticket {
	out := make([]Ticket, 0, len(s.byID))

	for _, t := range s.byID {
		out = append(out, t)
	}

	return out
}
