# NATS JetStream subjects for Helpdesk Event Hub.
# Partition / ordering key: ticket_id (aggregate_id in EventEnvelope).

helpdesk.ticket.created
helpdesk.ticket.assigned
helpdesk.ticket.updated
helpdesk.sla.warned
helpdesk.sla.breached
helpdesk.ticket.escalated
helpdesk.dlq.>
