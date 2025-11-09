#!/bin/bash

set -e

readonly d="$(cd "$(dirname "$0")" || exit 1; pwd)"

readonly logfile="${d}/../tmp/worker.log"
readonly pidfile="${d}/../tmp/worker.pid"
mkdir -p "$(dirname "$pidfile")"

start() {
    "${d}/../dist/pneutrinoutil-worker" \
        --mysqlDSN "$MYSQL_DSN" \
        --redisDSN "$REDIS_DSN" \
        --shell /bin/bash \
        --storageBucket "$STORAGE_BUCKET" \
        --storageS3 \
        --debug >> "$logfile" &
    echo $! > "$pidfile"
}

stop() {
    kill "$(cat "$pidfile")"
}


readonly cmd="$1"
case "$cmd" in
    start) start ;;
    stop) stop ;;
    *)
        echo >&2 'Please select start or stop'
        exit 1
        ;;
esac
