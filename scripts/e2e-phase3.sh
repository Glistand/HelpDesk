#!/usr/bin/env bash
# Phase 3 e2e: create → SLA warn/breach → escalate to L2 (a-3)
set -euo pipefail

GATEWAY="${GATEWAY_URL:-http://localhost:8080}"

echo "==> login"
LOGIN=$(curl -sf -X POST "$GATEWAY/auth/login" \
  -H 'Content-Type: application/json' \
  -d '{"email":"agent@helpdesk.local","password":"password"}')
TOKEN=$(echo "$LOGIN" | python3 -c 'import json,sys; print(json.load(sys.stdin)["access_token"])')

echo "==> create ticket"
CREATE=$(curl -sf -X POST "$GATEWAY/tickets" \
  -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"title":"Phase3 SLA e2e","description":"wait for breach","priority":"high","category":"e2e","requester":"e2e"}')
TICKET_ID=$(echo "$CREATE" | python3 -c 'import json,sys; print(json.load(sys.stdin)["id"])')
echo "ticket_id=$TICKET_ID"

echo "==> wait for SLA breach + L2 escalation"
OK=0
for i in $(seq 1 60); do
  TICKET=$(curl -sf -H "Authorization: Bearer $TOKEN" "$GATEWAY/tickets/$TICKET_ID")
  ASSIGNEE=$(echo "$TICKET" | python3 -c 'import json,sys; print(json.load(sys.stdin).get("assignee_id") or "")')
  SLA=$(curl -sf -H "Authorization: Bearer $TOKEN" "$GATEWAY/tickets/$TICKET_ID/sla" || echo '{}')
  STATE=$(echo "$SLA" | python3 -c 'import json,sys; print(json.load(sys.stdin).get("state") or "")' 2>/dev/null || echo "")
  TL=$(curl -sf -H "Authorization: Bearer $TOKEN" "$GATEWAY/tickets/$TICKET_ID/timeline")
  TYPES=$(echo "$TL" | python3 -c 'import json,sys; print(",".join(e["event_type"] for e in json.load(sys.stdin).get("events") or []))')
  echo "  try=$i assignee=$ASSIGNEE sla=$STATE types=$TYPES"
  if [[ "$ASSIGNEE" == "a-3" && "$TYPES" == *"sla.breached"* && "$TYPES" == *"ticket.escalated"* ]]; then
    OK=1
    break
  fi
  sleep 0.5
done

if [[ "$OK" -ne 1 ]]; then
  echo "FAIL: expected L2 assignee a-3, sla.breached and ticket.escalated in timeline"
  exit 1
fi

echo "PASS: SLA breach escalated to L2"
