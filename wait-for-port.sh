#!/bin/sh

set -e

host="$1"
port="$2"
shift 2
cmd="$@"

until nc -z "$host" "$port"; do
  echo "Waiting for port $port to be ready..."
  sleep 1
done

>&2 echo "Port $port is ready - executing command"
exec "$cmd"
