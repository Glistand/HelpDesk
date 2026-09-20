#!/usr/bin/env bash
# Smoke: widget session → visitor message (+ bot) → handoff → agent reply
set -euo pipefail

GATEWAY="${GATEWAY_URL:-http://localhost:8080}"
SITE_KEY="${WIDGET_SITE_KEY:-demo-site}"

echo "== health =="
curl -sf "$GATEWAY/healthz" | grep -q ok

echo "== widget session =="
SESSION=$(curl -sf -X POST "$GATEWAY/widget/session" \
  -H "Content-Type: application/json" \
  -H "X-Site-Key: $SITE_KEY" \
  -d "{\"site_key\":\"$SITE_KEY\"}")
VISITOR=$(echo "$SESSION" | jq -r .visitor_id)
CONV=$(echo "$SESSION" | jq -r .conversation.id)
echo "visitor=$VISITOR conversation=$CONV"
test -n "$VISITOR" && test "$VISITOR" != null
test -n "$CONV" && test "$CONV" != null

echo "== visitor message (may invoke OpenRouter) =="
MSG=$(curl -sf -X POST "$GATEWAY/widget/conversations/$CONV/messages" \
  -H "Content-Type: application/json" \
  -H "X-Visitor-Id: $VISITOR" \
  -H "X-Site-Key: $SITE_KEY" \
  -d '{"body":"Привет, нужна помощь с заказом"}')
echo "$MSG" | jq -r '.message.role,.bot_reply.role // "no-bot"'
test "$(echo "$MSG" | jq -r .message.role)" = "visitor"

echo "== handoff =="
HAND=$(curl -sf -X POST "$GATEWAY/widget/conversations/$CONV/handoff" \
  -H "Content-Type: application/json" \
  -H "X-Visitor-Id: $VISITOR" \
  -d '{}')
STATUS=$(echo "$HAND" | jq -r .conversation.status)
echo "status=$STATUS"
test "$STATUS" = "waiting_agent"

echo "== agent login + reply =="
TOKEN=$(curl -sf -X POST "$GATEWAY/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"agent@helpdesk.local","password":"password"}' | jq -r .access_token)
test -n "$TOKEN" && test "$TOKEN" != null

AGENT=$(curl -sf -X POST "$GATEWAY/conversations/$CONV/messages" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"body":"Здравствуйте! Я агент, чем помочь?"}')
echo "$AGENT" | jq -r .message.role
test "$(echo "$AGENT" | jq -r .message.role)" = "agent"

echo "== list conversations =="
curl -sf "$GATEWAY/conversations" -H "Authorization: Bearer $TOKEN" | jq -e --arg id "$CONV" \
  '.conversations | map(.id) | index($id) != null' >/dev/null

echo "OK smoke widget→bot→handoff→agent ($CONV)"
