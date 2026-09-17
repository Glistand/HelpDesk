#!/usr/bin/env bash
# Phase 4 e2e: create ticket → indexed in Meilisearch → search hit + card aggregation
set -euo pipefail

GATEWAY="${GATEWAY_URL:-http://localhost:8080}"
MARKER="phase4-search-$(date +%s)"

echo "==> login"
LOGIN=$(curl -sf -X POST "$GATEWAY/auth/login" \
  -H 'Content-Type: application/json' \
  -d '{"email":"agent@helpdesk.local","password":"password"}')
TOKEN=$(echo "$LOGIN" | python3 -c 'import json,sys; print(json.load(sys.stdin)["access_token"])')

echo "==> create ticket"
CREATE=$(curl -sf -X POST "$GATEWAY/tickets" \
  -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d "{\"title\":\"$MARKER unique VPN\",\"description\":\"Phase4 search e2e\",\"priority\":\"high\",\"category\":\"Network\",\"requester\":\"e2e\"}")
TICKET_ID=$(echo "$CREATE" | python3 -c 'import json,sys; print(json.load(sys.stdin)["id"])')
echo "ticket_id=$TICKET_ID"

echo "==> wait for search index"
FOUND=0
for i in $(seq 1 40); do
  SEARCH=$(curl -sf -H "Authorization: Bearer $TOKEN" --get "$GATEWAY/search" --data-urlencode "q=$MARKER")
  COUNT=$(echo "$SEARCH" | python3 -c 'import json,sys; print(len(json.load(sys.stdin).get("hits") or []))')
  IDS=$(echo "$SEARCH" | python3 -c 'import json,sys; print(",".join(h["id"] for h in json.load(sys.stdin).get("hits") or []))')
  echo "  try=$i hits=$COUNT ids=$IDS"
  if echo "$IDS" | grep -q "$TICKET_ID"; then
    FOUND=1
    break
  fi
  sleep 0.5
done

if [[ "$FOUND" -ne 1 ]]; then
  echo "FAIL: ticket not found in search"
  exit 1
fi

echo "==> card aggregation"
CARD=$(curl -sf -H "Authorization: Bearer $TOKEN" "$GATEWAY/tickets/$TICKET_ID/card")
echo "$CARD" | python3 -c "
import json,sys
card=json.load(sys.stdin)
assert card.get('ticket',{}).get('id')=='$TICKET_ID', card
assert isinstance(card.get('timeline'), list), card
print('card ok: timeline=', len(card['timeline']), 'has_sla=', card.get('sla') is not None, 'has_assignment=', card.get('assignment') is not None)
"

echo "PASS: search + BFF card aggregation"
