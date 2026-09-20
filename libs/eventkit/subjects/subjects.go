package subjects

// JetStream subject constants for Helpdesk Event Hub.
const (
	TicketCreated    = "helpdesk.ticket.created"
	TicketAssigned   = "helpdesk.ticket.assigned"
	TicketUpdated    = "helpdesk.ticket.updated"
	SLAWarned        = "helpdesk.sla.warned"
	SLABreached      = "helpdesk.sla.breached"
	TicketEscalated  = "helpdesk.ticket.escalated"
	NotificationSent = "helpdesk.notification.sent"
	ConversationCreated = "helpdesk.conversation.created"
	ConversationMessage = "helpdesk.conversation.message"
	ConversationHandoff = "helpdesk.conversation.handoff"
	DLQPrefix        = "helpdesk.dlq"
)

// Stream names created by deploy/compose init.
const (
	StreamEvents = "HELP_DESK_EVENTS"
	StreamDLQ    = "HELP_DESK_DLQ"
)

// DLQSubject maps an original subject onto the DLQ stream.
// helpdesk.ticket.created → helpdesk.dlq.ticket.created
func DLQSubject(original string) string {
	const prefix = "helpdesk."
	if len(original) > len(prefix) && original[:len(prefix)] == prefix {
		return DLQPrefix + "." + original[len(prefix):]
	}
	return DLQPrefix + "." + original
}
