#!/usr/bin/env bash
set -euo pipefail
base=${BASE_URL:-http://127.0.0.1:8080}
curl -fsS "$base/healthz" >/dev/null
curl -fsS -X POST "$base/api/v1/voyages" -H 'content-type: application/json' -d '{"id":"voy-1","vessel_id":"vessel-demo","port_id":"port-1","eta":"2026-08-22T00:00:00Z","etd":"2026-08-23T00:00:00Z"}' >/dev/null
curl -fsS -X POST "$base/api/v1/containers" -H 'content-type: application/json' -d '{"id":"box-1","voyage_id":"voy-1","iso_size":"40","weight_kg":12000,"destination":"TYO","on_deck":true,"priority":1}' >/dev/null
curl -fsS -X POST "$base/api/v1/plans/solve" -H 'content-type: application/json' -d '{"id":"plan-1","voyage_id":"voy-1"}' >/dev/null
curl -fsS -X POST "$base/api/v1/plans/plan-1/publish" >/dev/null
curl -fsS "$base/api/v1/plans/plan-1/explain" >/dev/null
echo smoke-ok
