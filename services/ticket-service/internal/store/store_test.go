package store

import (
	"errors"
	"testing"
)

func TestCreate_Success(t *testing.T) {
	s := NewStore()
	ticket, err := s.Create("VPN down")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ticket.ID == "" {
		t.Fatal("expected non-empty id")
	}
	if ticket.Status != "new" {
		t.Fatalf("expected status new, got %s", ticket.Status)
	}
}

func TestCreate_EmptyTitle(t *testing.T) {
	s := NewStore()
	_, err := s.Create("   ")
	if !errors.Is(err, ErrEmptyTitle) {
		t.Fatalf("expected ErrEmptyTitle, got %v", err)
	}
}

func TestGet_NotFound(t *testing.T) {
	s := NewStore()
	_, err := s.Get("t-404")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
