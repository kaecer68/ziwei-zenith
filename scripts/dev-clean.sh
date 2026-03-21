#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
PORTS_FILE="$REPO_ROOT/.env.ports"

# Ensure .env.ports exists (run sync-contracts if needed)
if [[ ! -f "$PORTS_FILE" ]]; then
  echo "[dev-clean] .env.ports not found, running sync-contracts..."
  bash "$SCRIPT_DIR/sync-contracts.sh"
fi

# Source the ports file
# shellcheck disable=SC1090
source "$PORTS_FILE"

# Validate required ports are set
if [[ -z "${ZIWEI_GRPC_PORT:-}" ]] || [[ -z "${ZIWEI_REST_PORT:-}" ]]; then
  echo "[dev-clean] Error: ZIWEI_GRPC_PORT and ZIWEI_REST_PORT must be defined in .env.ports" >&2
  exit 1
fi

ports=("$ZIWEI_GRPC_PORT" "$ZIWEI_REST_PORT")
for port in "${ports[@]}"; do
  pids="$(lsof -tiTCP:"$port" -sTCP:LISTEN || true)"
  if [[ -n "$pids" ]]; then
    echo "[dev-clean] 清理 port $port: $pids"
    kill $pids 2>/dev/null || true
  fi
done
