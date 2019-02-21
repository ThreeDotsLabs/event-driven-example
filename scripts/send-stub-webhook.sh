#!/bin/bash
readonly payload="$(dirname "$0")/payload.json"
readonly time=$(date --rfc-3339=seconds | tr ' ' 'T')
readonly id=$(tr -dc 'a-f0-9' < /dev/urandom | head -c40)

cat >"$payload" <<EOF
{
  "commits": [{
    "id": "$id",
    "message": "This is an example commit message",
    "timestamp": "$time",
    "author": {
      "name": "John Doe"
    }
  }]
}
EOF

curl -X POST "http://127.0.0.1:8080" \
     -d "@$payload"
