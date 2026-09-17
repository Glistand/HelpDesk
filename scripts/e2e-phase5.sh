#!/usr/bin/env bash
# Phase 5 e2e: Next.js BFF cookie auth → create ticket → card → search proxy
set -euo pipefail

WEB="${WEB_URL:-http://localhost:3000}"
MARKER="phase5-web-$(date +%s)"
COOKIE_JAR=$(mktemp)
trap 'rm -f "$COOKIE_JAR"' EXIT

echo "==> web health (login page)"
curl -sf "$WEB/login" -o /dev/null

echo "==> login via Next.js /api/auth/login"
LOGIN=$(curl -sf -c "$COOKIE_JAR" -b "$COOKIE_JAR" -X POST "$WEB/api/auth/login" \
  -H 'Content-Type: application/json' \
  -d '{"email":"agent@helpdesk.local","password":"password"}')
echo "$LOGIN" | python3 -c 'import json,sys; u=json.load(sys.stdin).get("user") or {}; assert u.get("email")=="agent@helpdesk.local", u; print("user ok:", u.get("email"))'

if ! grep -q "hd_token" "$COOKIE_JAR"; then
  echo "FAIL: hd_token cookie not set"
  cat "$COOKIE_JAR"
  exit 1
fi

echo "==> create ticket via /api/hd/tickets"
CREATE=$(curl -sf -c "$COOKIE_JAR" -b "$COOKIE_JAR" -X POST "$WEB/api/hd/tickets" \
  -H 'Content-Type: application/json' \
  -d "{\"title\":\"$MARKER VPN\",\"description\":\"Phase5 web e2e\",\"priority\":\"high\",\"category\":\"Network\"}")
TICKET_ID=$(echo "$CREATE" | python3 -c 'import json,sys; print(json.load(sys.stdin)["id"])')
echo "ticket_id=$TICKET_ID"

echo "==> ticket card via /api/hd/tickets/{id}/card"
CARD=$(curl -sf -c "$COOKIE_JAR" -b "$COOKIE_JAR" "$WEB/api/hd/tickets/$TICKET_ID/card")
echo "$CARD" | python3 -c "
import json,sys
card=json.load(sys.stdin)
assert card.get('ticket',{}).get('id')=='$TICKET_ID', card
assert isinstance(card.get('timeline'), list), card
print('card ok: timeline=', len(card['timeline']))
"

echo "==> SSR ticket page"
curl -sf -c "$COOKIE_JAR" -b "$COOKIE_JAR" "$WEB/tickets/$TICKET_ID" -o /dev/null

echo "==> wait for search via /api/hd/search"
FOUND=0
for i in $(seq 1 40); do
  SEARCH=$(curl -sf -c "$COOKIE_JAR" -b "$COOKIE_JAR" --get "$WEB/api/hd/search" --data-urlencode "q=$MARKER")
  IDS=$(echo "$SEARCH" | python3 -c 'import json,sys; print(",".join(h["id"] for h in json.load(sys.stdin).get("hits") or []))')
  echo "  try=$i ids=$IDS"
  if echo "$IDS" | grep -q "$TICKET_ID"; then
    FOUND=1
    break
  fi
  sleep 0.5
done

if [[ "$FOUND" -ne 1 ]]; then
  echo "FAIL: ticket not found via web search proxy"
  exit 1
fi

echo "==> logout"
curl -sf -c "$COOKIE_JAR" -b "$COOKIE_JAR" -X POST "$WEB/api/auth/logout" -o /dev/null

echo "PASS: Next.js cookie auth + BFF proxy + SSR card + search"
