#!/bin/bash

set -e

readonly d="$(cd "$(dirname "$0")" || exit 1; pwd)"

readonly logfile="${d}/../tmp/server.log"
readonly pidfile="${d}/../tmp/server.pid"
mkdir -p "$(dirname "$pidfile")"

start() {
    if [[ -s "$pidfile" ]] ; then
        echo >&2 "Cannot start server because ${pidfile} exist"
        return 1
    fi
    "${d}/../dist/pneutrinoutil-server" \
        --port "$SERVER_PORT" \
        --mysqlDSN "$MYSQL_DSN" \
        --redisDSN "$REDIS_DSN" \
        --storageBucket "$STORAGE_BUCKET" \
        --storageS3 \
        --debug >> "$logfile" 2>&1 &
    echo $! > "$pidfile"
    echo >&2 "Server started with pid $(cat "$pidfile")"
}

stop() {
    if [[ ! -s "$pidfile" ]] ; then
        echo >&2 "Cannot stop server because ${pidfile} is not found"
        return
    fi
    kill "$(cat "$pidfile")" || true
    rm -f "$pidfile"
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
