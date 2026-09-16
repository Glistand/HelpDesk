package envelope_test

import (
	"encoding/json"
	"testing"

	"github.com/Glistand/HelpDesk/libs/eventkit/envelope"
)

func TestNewAndRoundTrip(t *testing.T) {
	payload := map[string]string{"title": "VPN down"}
	ev, err := envelope.New("ticket.created", "t-1", "corr-1", payload)
	if err != nil {
		t.Fatal(err)
	}
	if ev.EventID == "" || ev.Type != "ticket.created" || ev.AggregateID != "t-1" {
		t.Fatalf("unexpected event: %+v", ev)
	}

	raw, err := ev.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	got, err := envelope.Unmarshal(raw)
	if err != nil {
		t.Fatal(err)
	}
	if got.EventID != ev.EventID || got.CorrelationID != "corr-1" {
		t.Fatalf("round-trip mismatch: %+v", got)
	}
	var decoded map[string]string
	if err := json.Unmarshal(got.Payload, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded["title"] != "VPN down" {
		t.Fatalf("payload: %v", decoded)
	}
}
