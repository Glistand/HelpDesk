#!/usr/bin/env bash
# Chaos: bounce a consumer mid-flight and verify assignment still completes.
set -euo pipefail

GATEWAY="${GATEWAY_URL:-http://localhost:8080}"

LOGIN=$(curl -sf -X POST "$GATEWAY/auth/login" \
  -H 'Content-Type: application/json' \
  -d '{"email":"agent@helpdesk.local","password":"password"}')
TOKEN=$(echo "$LOGIN" | python3 -c 'import json,sys; print(json.load(sys.stdin)["access_token"])')

echo "==> restart assignment-service"
docker restart helpdesk-assignment

CREATE=$(curl -sf -X POST "$GATEWAY/tickets" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"chaos-bounce","description":"chaos","priority":"high","category":"Network"}')
TICKET_ID=$(echo "$CREATE" | python3 -c 'import json,sys; print(json.load(sys.stdin)["id"])')
echo "ticket_id=$TICKET_ID"

for i in $(seq 1 50); do
  CARD=$(curl -sf -H "Authorization: Bearer $TOKEN" "$GATEWAY/tickets/$TICKET_ID/card")
  AID=$(echo "$CARD" | python3 -c 'import json,sys; c=json.load(sys.stdin); a=c.get("assignment") or {}; print(a.get("assignee_id") or "")')
  echo "  try=$i assignee=$AID"
  if [[ -n "$AID" ]]; then
    echo "PASS: assigned after chaos ($AID)"
    exit 0
  fi
  sleep 0.5
done

echo "FAIL: not assigned"
exit 1
