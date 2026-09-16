#!/usr/bin/env bash
# Phase 2 e2e: login → create ticket → wait for assign + notify → timeline ≥ 2 events
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
  -d '{"title":"Phase2 e2e","description":"auto-assign + notify + audit","priority":"high","category":"e2e","requester":"e2e"}')
TICKET_ID=$(echo "$CREATE" | python3 -c 'import json,sys; print(json.load(sys.stdin)["id"])')
echo "ticket_id=$TICKET_ID"

echo "==> wait for assignee + timeline"
OK=0
for i in $(seq 1 40); do
  TICKET=$(curl -sf -H "Authorization: Bearer $TOKEN" "$GATEWAY/tickets/$TICKET_ID")
  ASSIGNEE=$(echo "$TICKET" | python3 -c 'import json,sys; print(json.load(sys.stdin).get("assignee_id") or "")')
  TL=$(curl -sf -H "Authorization: Bearer $TOKEN" "$GATEWAY/tickets/$TICKET_ID/timeline")
  COUNT=$(echo "$TL" | python3 -c 'import json,sys; print(len(json.load(sys.stdin).get("events") or []))')
  TYPES=$(echo "$TL" | python3 -c 'import json,sys; print(",".join(e["event_type"] for e in json.load(sys.stdin).get("events") or []))')
  echo "  try=$i assignee=$ASSIGNEE events=$COUNT types=$TYPES"
  if [[ -n "$ASSIGNEE" && "$COUNT" -ge 2 ]]; then
    OK=1
    break
  fi
  sleep 0.5
done

if [[ "$OK" -ne 1 ]]; then
  echo "FAIL: expected assignee and timeline with ≥2 events"
  exit 1
fi

echo "PASS: ticket assigned and timeline populated"
