package subjects

// JetStream subject constants for Helpdesk Event Hub.
const (
	TicketCreated     = "helpdesk.ticket.created"
	TicketAssigned    = "helpdesk.ticket.assigned"
	TicketUpdated     = "helpdesk.ticket.updated"
	SLAWarned         = "helpdesk.sla.warned"
	SLABreached       = "helpdesk.sla.breached"
	TicketEscalated   = "helpdesk.ticket.escalated"
	DLQPrefix         = "helpdesk.dlq"
)

// Stream names created by deploy/compose init.
const (
	StreamEvents = "HELP_DESK_EVENTS"
	StreamDLQ    = "HELP_DESK_DLQ"
)
