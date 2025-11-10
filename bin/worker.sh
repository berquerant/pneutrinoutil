#!/bin/bash

set -e

readonly d="$(cd "$(dirname "$0")" || exit 1; pwd)"

readonly logfile="${d}/../tmp/worker.log"
readonly pidfile="${d}/../tmp/worker.pid"
mkdir -p "$(dirname "$pidfile")"

start() {
    if [[ -s "$pidfile" ]] ; then
        echo >&2 "Cannot start worker because ${pidfile} exist"
        return 1
    fi
    "${d}/../dist/pneutrinoutil-worker" \
        --mysqlDSN "$MYSQL_DSN" \
        --redisDSN "$REDIS_DSN" \
        --shell /bin/bash \
        --storageBucket "$STORAGE_BUCKET" \
        --storageS3 \
        --debug >> "$logfile" 2>&1 &
    echo $! > "$pidfile"
    echo >&2 "Worker started with pid $(cat "$pidfile")"
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
