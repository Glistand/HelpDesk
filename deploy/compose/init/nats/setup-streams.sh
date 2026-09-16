#!/bin/sh
set -eu

NATS_URL="${NATS_URL:-nats://nats:4222}"

echo "Waiting for NATS JetStream at ${NATS_URL}..."
until nats --server "${NATS_URL}" server check connection >/dev/null 2>&1; do
  sleep 1
done

echo "Creating JetStream streams..."

if ! nats --server "${NATS_URL}" stream info HELP_DESK_EVENTS >/dev/null 2>&1; then
  nats --server "${NATS_URL}" stream add HELP_DESK_EVENTS \
    --subjects "helpdesk.ticket.>" \
    --subjects "helpdesk.sla.>" \
    --subjects "helpdesk.assignment.>" \
    --subjects "helpdesk.notification.>" \
    --subjects "helpdesk.audit.>" \
    --storage file \
    --retention limits \
    --max-msgs=-1 \
    --max-age=720h \
    --dupe-window=2m \
    --defaults
else
  echo "Stream HELP_DESK_EVENTS already exists"
fi

if ! nats --server "${NATS_URL}" stream info HELP_DESK_DLQ >/dev/null 2>&1; then
  nats --server "${NATS_URL}" stream add HELP_DESK_DLQ \
    --subjects "helpdesk.dlq.>" \
    --storage file \
    --retention limits \
    --max-msgs=-1 \
    --max-age=168h \
    --defaults
else
  echo "Stream HELP_DESK_DLQ already exists"
fi

echo "JetStream streams ready:"
nats --server "${NATS_URL}" stream ls
