#!/usr/bin/env bash
# Concurrent create load against api-gateway.
set -euo pipefail

GATEWAY="${GATEWAY_URL:-http://localhost:8080}"
N="${LOAD_N:-50}"
CONCURRENCY="${LOAD_CONCURRENCY:-10}"

LOGIN=$(curl -sf -X POST "$GATEWAY/auth/login" \
  -H 'Content-Type: application/json' \
  -d '{"email":"agent@helpdesk.local","password":"password"}')
TOKEN=$(echo "$LOGIN" | python3 -c 'import json,sys; print(json.load(sys.stdin)["access_token"])')

TMP=$(mktemp -d)
trap 'rm -rf "$TMP"' EXIT

echo "load: n=$N concurrency=$CONCURRENCY"
seq 1 "$N" | xargs -P "$CONCURRENCY" -I{} bash -c '
  start=$(date +%s%3N)
  code=$(curl -s -o /dev/null -w "%{http_code}" -X POST "'"$GATEWAY"'/tickets" \
    -H "Authorization: Bearer '"$TOKEN"'" \
    -H "Content-Type: application/json" \
    -d "{\"title\":\"load-{}-$(date +%s)\",\"description\":\"load\",\"priority\":\"low\",\"category\":\"Network\"}")
  end=$(date +%s%3N)
  echo "$code $((end-start))" >> "'"$TMP"'/results"
'

python3 - <<PY
from pathlib import Path
rows=[l.split() for l in Path("$TMP/results").read_text().splitlines() if l.strip()]
codes=[int(c) for c,_ in rows]
ms=sorted(int(m) for _,m in rows)
ok=sum(1 for c in codes if 200<=c<300)
print(f"total={len(codes)} ok={ok} errors={len(codes)-ok}")
if ms:
  def pct(p):
    i=min(len(ms)-1, max(0, int(round((p/100)*(len(ms)-1)))))
    return ms[i]
  print(f"latency_ms p50={pct(50)} p95={pct(95)} max={ms[-1]}")
raise SystemExit(0 if ok==len(codes) else 1)
PY
