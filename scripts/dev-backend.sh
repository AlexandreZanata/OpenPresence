#!/usr/bin/env bash
# Start local OpenPresence backend (Postgres + biometric gRPC + attendance HTTP).
#
# Usage:
#   ./scripts/dev-backend.sh          # start all (postgres via Docker, Go/Rust on host)
#   ./scripts/dev-backend.sh stop     # stop host processes + docker postgres
#   ./scripts/dev-backend.sh status   # show health
#
# Endpoints:
#   Postgres   localhost:5433
#   Biometric  localhost:9090 (gRPC)
#   Attendance localhost:8088 (HTTP)

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
COMPOSE_FILE="${ROOT}/infra/docker-compose.e2e.yml"
PID_DIR="${ROOT}/.local/run"
ATTENDANCE_PID="${PID_DIR}/attendance.pid"
BIOMETRIC_PID="${PID_DIR}/biometric.pid"
ATTENDANCE_LOG="${PID_DIR}/attendance.log"
BIOMETRIC_LOG="${PID_DIR}/biometric.log"

ATTENDANCE_URL="${ATTENDANCE_HTTP_URL:-http://127.0.0.1:8088}"
BIOMETRIC_GRPC="${BIOMETRIC_GRPC_ADDR:-127.0.0.1:9090}"
DATABASE_URL="${DATABASE_URL:-postgres://attendance_app:attendance_app@127.0.0.1:5433/openpresence?sslmode=disable}"

mkdir -p "$PID_DIR"

attendance_port() {
  local p="${ATTENDANCE_URL##*:}"
  echo "${p%%/*}"
}

biometric_port() {
  echo "${BIOMETRIC_GRPC##*:}"
}

is_listening() {
  local host="$1" port="$2"
  (echo >/dev/tcp/"$host"/"$port") 2>/dev/null
}

kill_port() {
  local port="$1"
  if command -v fuser >/dev/null 2>&1; then
    fuser -k "${port}/tcp" 2>/dev/null || true
  elif command -v lsof >/dev/null 2>&1; then
    lsof -ti ":${port}" 2>/dev/null | xargs -r kill 2>/dev/null || true
  fi
}

stop_pid_file() {
  local pid_file="$1"
  if [[ -f "$pid_file" ]]; then
    local pid
    pid="$(cat "$pid_file")"
    kill "$pid" 2>/dev/null || true
    rm -f "$pid_file"
  fi
}

stop_host() {
  stop_pid_file "$ATTENDANCE_PID"
  stop_pid_file "$BIOMETRIC_PID"
  kill_port "$(attendance_port)"
  kill_port "$(biometric_port)"
}

wait_http() {
  local url="$1" label="$2" tries="${3:-60}"
  for _ in $(seq 1 "$tries"); do
    if curl -sf "$url" >/dev/null 2>&1; then
      echo "OK: $label"
      return 0
    fi
    sleep 1
  done
  echo "TIMEOUT: $label ($url)" >&2
  return 1
}

wait_tcp() {
  local host="$1" port="$2" label="$3" tries="${4:-60}"
  for _ in $(seq 1 "$tries"); do
    if is_listening "$host" "$port"; then
      echo "OK: $label"
      return 0
    fi
    sleep 1
  done
  echo "TIMEOUT: $label ($host:$port)" >&2
  return 1
}

service_status() {
  local name="$1" pid_file="$2" host="$3" port="$4"
  if [[ -f "$pid_file" ]] && kill -0 "$(cat "$pid_file")" 2>/dev/null; then
    echo "RUNNING: $name (pid $(cat "$pid_file"), port $port)"
    return 0
  fi
  if is_listening "$host" "$port"; then
    echo "RUNNING: $name (port $port, pid file stale)"
    return 0
  fi
  echo "STOPPED: $name"
  return 1
}

cmd_status() {
  echo "=== OpenPresence dev backend ==="
  docker compose -f "$COMPOSE_FILE" ps postgres 2>/dev/null || true
  service_status attendance "$ATTENDANCE_PID" 127.0.0.1 "$(attendance_port)" || true
  service_status biometric "$BIOMETRIC_PID" 127.0.0.1 "$(biometric_port)" || true
  if curl -sf "${ATTENDANCE_URL}/health/live" | grep -q '"status":"ok"'; then
    echo "health/live: ok"
  else
    echo "health/live: down"
  fi
}

cmd_stop() {
  echo "Stopping host services..."
  stop_host
  echo "Stopping postgres container..."
  docker compose -f "$COMPOSE_FILE" stop postgres 2>/dev/null || true
  echo "Done."
}

cmd_start() {
  command -v docker >/dev/null || { echo "docker required" >&2; exit 1; }
  command -v go >/dev/null || { echo "go required" >&2; exit 1; }
  command -v cargo >/dev/null || { echo "cargo required" >&2; exit 1; }

  stop_host

  echo "Starting Postgres (docker)..."
  docker compose -f "$COMPOSE_FILE" up -d postgres --wait

  local grpc_port biometrics_http_port http_port
  grpc_port="$(biometric_port)"
  biometrics_http_port=$((grpc_port + 1))
  http_port="$(attendance_port)"

  echo "Starting biometric-server (host)..."
  nohup bash -c "cd '${ROOT}/services/biometric' && exec env \
    BIOMETRIC_USE_STUB=true \
    BIOMETRIC_GRPC_ADDR='0.0.0.0:${grpc_port}' \
    BIOMETRIC_HTTP_ADDR='127.0.0.1:${biometrics_http_port}' \
    RUST_LOG=warn \
    cargo run --quiet --bin biometric-server" >>"$BIOMETRIC_LOG" 2>&1 &
  echo $! >"$BIOMETRIC_PID"

  wait_tcp 127.0.0.1 "$grpc_port" "biometric gRPC"

  echo "Starting attendance-server (host)..."
  nohup bash -c "cd '${ROOT}/services/attendance' && exec env \
    DATABASE_URL='${DATABASE_URL}' \
    BIOMETRIC_GRPC_ADDR='${BIOMETRIC_GRPC}' \
    ATTENDANCE_HTTP_ADDR=':${http_port}' \
    go run ./cmd/attendance-server" >>"$ATTENDANCE_LOG" 2>&1 &
  echo $! >"$ATTENDANCE_PID"

  wait_http "${ATTENDANCE_URL}/health/live" "attendance HTTP"

  cat <<EOF

=== Backend ready ===
  Postgres:   127.0.0.1:5433
  Biometric:  ${BIOMETRIC_GRPC} (gRPC)
  Attendance: ${ATTENDANCE_URL}
  Logs:       ${PID_DIR}/*.log
  Stop:       ./scripts/dev-backend.sh stop
  Verify:     ./scripts/verify-dev-backend.sh

Note: POST /v1/auth/login is not implemented yet — admin panel will use dev auth mock until auth service exists.
EOF
}

case "${1:-start}" in
  start) cmd_start ;;
  stop) cmd_stop ;;
  status) cmd_status ;;
  *) echo "usage: $0 [start|stop|status]" >&2; exit 2 ;;
esac
