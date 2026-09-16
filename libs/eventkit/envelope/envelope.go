package envelope

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Event is the JSON envelope published to NATS JetStream.
type Event struct {
	EventID       string          `json:"event_id"`
	CorrelationID string          `json:"correlation_id"`
	OccurredAt    time.Time       `json:"occurred_at"`
	Type          string          `json:"type"`
	AggregateID   string          `json:"aggregate_id"`
	Payload       json.RawMessage `json:"payload"`
}

// New builds an envelope with a fresh event_id and UTC occurred_at.
func New(eventType, aggregateID, correlationID string, payload any) (Event, error) {
	raw, err := json.Marshal(payload)
	if err != nil {
		return Event{}, err
	}
	if correlationID == "" {
		correlationID = uuid.NewString()
	}
	return Event{
		EventID:       uuid.NewString(),
		CorrelationID: correlationID,
		OccurredAt:    time.Now().UTC(),
		Type:          eventType,
		AggregateID:   aggregateID,
		Payload:       raw,
	}, nil
}

// Marshal serializes the envelope as JSON.
func (e Event) Marshal() ([]byte, error) {
	return json.Marshal(e)
}

// Unmarshal parses an envelope from JSON bytes.
func Unmarshal(data []byte) (Event, error) {
	var e Event
	err := json.Unmarshal(data, &e)
	return e, err
}
