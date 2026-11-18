#!/bin/sh 
  
level="$1"
msg="$2"
shift 2

printf '%s level=%s msg="%s"' \
  "$(date -u '+%Y-%m-%dT%H:%M:%SZ')" "$level" "$msg"

for field in "$@"; do
  printf ' %s' "$field"
done

printf '\n'
