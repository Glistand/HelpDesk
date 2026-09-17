#!/usr/bin/env bash
# Phase 6 e2e: hardening headers/auth + OTel traces in Jaeger + light load + chaos recovery
set -euo pipefail

GATEWAY="${GATEWAY_URL:-http://localhost:8080}"
JAEGER="${JAEGER_URL:-http://localhost:16686}"
CORR="phase6-$(date +%s)-$RANDOM"

echo "==> healthz"
curl -sf "$GATEWAY/healthz" >/dev/null

echo "==> unauthorized without token"
CODE=$(curl -s -o /dev/null -w '%{http_code}' "$GATEWAY/tickets")
[[ "$CODE" == "401" ]] || { echo "FAIL: expected 401 got $CODE"; exit 1; }

echo "==> security headers on healthz"
HEADERS=$(curl -sI "$GATEWAY/healthz")
echo "$HEADERS" | grep -qi 'X-Content-Type-Options: nosniff' || { echo "FAIL: missing nosniff"; echo "$HEADERS"; exit 1; }
echo "$HEADERS" | grep -qi 'X-Frame-Options: DENY' || { echo "FAIL: missing frame options"; echo "$HEADERS"; exit 1; }
echo "$HEADERS" | grep -qi 'X-Correlation-Id:' || { echo "FAIL: missing correlation header"; echo "$HEADERS"; exit 1; }

echo "==> login"
LOGIN=$(curl -sf -X POST "$GATEWAY/auth/login" \
  -H 'Content-Type: application/json' \
  -d '{"email":"agent@helpdesk.local","password":"password"}')
TOKEN=$(echo "$LOGIN" | python3 -c 'import json,sys; print(json.load(sys.stdin)["access_token"])')

echo "==> create ticket with correlation id"
CREATE=$(curl -sf -D /tmp/hd-phase6-headers -X POST "$GATEWAY/tickets" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -H "X-Correlation-Id: $CORR" \
  -d "{\"title\":\"phase6-trace $CORR\",\"description\":\"otel e2e\",\"priority\":\"normal\",\"category\":\"Network\"}")
TICKET_ID=$(echo "$CREATE" | python3 -c 'import json,sys; print(json.load(sys.stdin)["id"])')
echo "ticket_id=$TICKET_ID corr=$CORR"
grep -qi "X-Correlation-Id: $CORR" /tmp/hd-phase6-headers || { echo "FAIL: correlation not echoed"; cat /tmp/hd-phase6-headers; exit 1; }

echo "==> wait for Jaeger traces (api-gateway)"
FOUND=0
for i in $(seq 1 30); do
  TRACES=$(curl -sf "$JAEGER/api/traces?service=api-gateway&limit=20" || echo '{}')
  HIT=$(echo "$TRACES" | python3 -c "
import json,sys
data=json.load(sys.stdin)
traces=data.get('data') or []
print(len(traces))
")
  echo "  try=$i traces=$HIT"
  if [[ "$HIT" -gt 0 ]]; then
    FOUND=1
    break
  fi
  sleep 1
done
[[ "$FOUND" -eq 1 ]] || { echo "FAIL: no traces for api-gateway"; exit 1; }

echo "==> light load (20 creates)"
LOAD_OK=0
LOAD_FAIL=0
for i in $(seq 1 20); do
  if curl -sf -X POST "$GATEWAY/tickets" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d "{\"title\":\"phase6-load-$i-$CORR\",\"description\":\"load\",\"priority\":\"low\",\"category\":\"Network\"}" >/dev/null; then
    LOAD_OK=$((LOAD_OK+1))
  else
    LOAD_FAIL=$((LOAD_FAIL+1))
  fi
done
echo "load ok=$LOAD_OK fail=$LOAD_FAIL"
[[ "$LOAD_OK" -ge 18 ]] || { echo "FAIL: too many load failures"; exit 1; }

echo "==> chaos: restart assignment-service, expect new ticket assigned"
docker restart helpdesk-assignment >/dev/null
sleep 3
CREATE2=$(curl -sf -X POST "$GATEWAY/tickets" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"title\":\"phase6-chaos $CORR\",\"description\":\"chaos\",\"priority\":\"high\",\"category\":\"Network\"}")
TICKET2=$(echo "$CREATE2" | python3 -c 'import json,sys; print(json.load(sys.stdin)["id"])')
ASSIGNED=0
for i in $(seq 1 40); do
  CARD=$(curl -sf -H "Authorization: Bearer $TOKEN" "$GATEWAY/tickets/$TICKET2/card")
  AID=$(echo "$CARD" | python3 -c 'import json,sys; c=json.load(sys.stdin); a=c.get("assignment") or {}; print(a.get("assignee_id") or "")')
  echo "  try=$i assignee=$AID"
  if [[ -n "$AID" ]]; then
    ASSIGNED=1
    break
  fi
  sleep 0.5
done
[[ "$ASSIGNED" -eq 1 ]] || { echo "FAIL: ticket not assigned after chaos restart"; exit 1; }

echo "PASS: hardening + tracing + load + chaos recovery"
