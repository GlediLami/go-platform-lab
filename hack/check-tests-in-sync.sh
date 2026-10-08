#!/usr/bin/env bash
# Every exercise has the same tests as its solution. This catches drift.
set -euo pipefail
cd "$(dirname "$0")/.."
status=0
for sol in solutions/*/; do
  ex="exercises/$(basename "$sol")"
  for t in "$sol"*_test.go; do
    if ! cmp -s "$t" "$ex/$(basename "$t")"; then
      echo "out of sync: $t vs $ex/$(basename "$t")"
      status=1
    fi
  done
done
[ $status -eq 0 ] && echo "all exercise tests match their solutions"
exit $status
