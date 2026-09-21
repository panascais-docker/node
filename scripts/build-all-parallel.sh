#!/bin/sh
set -eu

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

VERSIONS="12 14 15 16 17 18 19 20 21 22 23 24 25 26 latest lts"
LOG_DIR="${TMPDIR:-/tmp}node-builds"
mkdir -p "$LOG_DIR"

echo "logs: $LOG_DIR/build-<version>.log"

for v in $VERSIONS; do
  (
    if bun run build.ts "$v" >"$LOG_DIR/build-$v.log" 2>&1; then
      echo "OK $v"
    else
      echo "FAILURE $v (see $LOG_DIR/build-$v.log)"
      exit 1
    fi
  ) &
done

FAIL=0
for v in $VERSIONS; do
  if ! wait; then
    FAIL=1
  fi
done

if [ "$FAIL" -ne 0 ]; then
  echo "One or more builds failed."
  exit 1
fi

echo "OK"
