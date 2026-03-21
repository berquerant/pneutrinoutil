#!/bin/bash

set -e

readonly d="$(cd "$(dirname "$0")" || exit 1; pwd)"

readonly logfile="${d}/../tmp/k8s_worker.log"
readonly pidfile="${d}/../tmp/k8s_worker.pid"
mkdir -p "$(dirname "$pidfile")"

export MYSQLDSN="pneutrinoutil:userpass@tcp(127.0.0.1:3306)/pneutrinoutil?parseTime=true&loc=Asia%2FTokyo"
export REDISDSN="redis://127.0.0.1:6379/1"
export AWS_ENDPOINT_URL="http://127.0.0.1:9000"
export AWS_ACCESS_KEY_ID=admin
export AWS_SECRET_ACCESS_KEY=key
export STORAGEBUCKET=pneutrinoutil
export STORAGES3=true
export AWS_DEFAULT_REGION=us-east-1
export AWS_S3_DISABLE_HTTPS=true
export AWS_USE_PATH_STYLE_ENDPOINT=true

start() {
    if [[ -s "$pidfile" ]] ; then
        echo >&2 "Cannot start worker because ${pidfile} exist"
        return 1
    fi
    "${d}/../dist/pneutrinoutil-worker" \
        --shell /bin/bash \
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
