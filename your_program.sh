#!/bin/sh
#
# Use this script to run your program locally.
#

set -e # Exit early if any commands fail

(
  cd "$(dirname "$0")"
  go build -o /tmp/redis-go-build app/*.go
)

exec /tmp/redis-go-build "$@"
