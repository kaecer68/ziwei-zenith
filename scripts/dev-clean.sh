#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
PORTS_FILE="$REPO_ROOT/.env.ports"

ZIWEI_GRPC_PORT="${ZIWEI_GRPC_PORT:-50053}"
ZIWEI_REST_PORT="${ZIWEI_REST_PORT:-8083}"

if [[ -f "$PORTS_FILE" ]]; then
  # shellcheck disable=SC1090
  source "$PORTS_FILE"
fi

ports=("$ZIWEI_GRPC_PORT" "$ZIWEI_REST_PORT")
for port in "${ports[@]}"; do
  pids="$(lsof -tiTCP:"$port" -sTCP:LISTEN || true)"
  if [[ -n "$pids" ]]; then
    echo "[dev-clean] 清理 port $port: $pids"
    kill $pids 2>/dev/null || true
  fi
done
