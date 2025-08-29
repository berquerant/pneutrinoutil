#!/bin/bash

readonly d="$(cd "$(dirname "$0")" || exit 1; pwd)"

client() {
    "${d}/docker.sh" exec redis redis-cli -h "$REDIS_HOST" "$@"
}

ping() {
    echo >&2 "ping redis"
    client ping
}

wait_ping() {
    for c in $(seq 30) ; do
        if (( c > 1 )) ; then
            sleep 1
        fi
        if ping ; then
            echo >&2 "ping redis success!"
            return
        fi
    done
    return 1
}

case "$1" in
    "wait") wait_ping ;;
    *) client "$@" ;;
esac
