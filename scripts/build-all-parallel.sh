#!/bin/sh
set -eu

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

TEMPORARY="${TMPDIR:-/tmp}"
LOG_DIR="${TEMPORARY%/}/node-builds"
BINARY="${LOG_DIR}/scripts"
mkdir -p "$LOG_DIR"

if [ "$#" -gt 0 ]; then
  TARGETS="$*"
else
  TARGETS=""
  for directory in [0-9]*/; do
    if [ -f "${directory}Dockerfile" ]; then
      TARGETS="${TARGETS} ${directory%/}"
    fi
  done
fi

go build -o "$BINARY" ./scripts

echo "logs: ${LOG_DIR}/build-<line>.log"

PIDS=""
for line in $TARGETS; do
  (
    if "$BINARY" build "$line" >"${LOG_DIR}/build-${line}.log" 2>&1; then
      echo "OK ${line}"
    else
      echo "FAILURE ${line} (see ${LOG_DIR}/build-${line}.log)"
      exit 1
    fi
  ) &
  PIDS="${PIDS} $!"
done

FAIL=0
for pid in $PIDS; do
  if ! wait "$pid"; then
    FAIL=1
  fi
done

if [ "$FAIL" -ne 0 ]; then
  echo "One or more builds failed."
  exit 1
fi

echo "OK"
