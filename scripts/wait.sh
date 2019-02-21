#!/bin/bash
set -e

for service in "$@"; do
    "$(dirname "$0")/wait-for-it.sh" -t 30 "$service"
done
