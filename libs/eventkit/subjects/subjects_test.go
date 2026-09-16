package subjects_test

import (
	"testing"

	"github.com/Glistand/HelpDesk/libs/eventkit/subjects"
)

func TestDLQSubject(t *testing.T) {
	got := subjects.DLQSubject(subjects.TicketCreated)
	want := "helpdesk.dlq.ticket.created"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}
