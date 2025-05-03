#!/bin/bash

client() {
    aws s3 "$@"
}

ping() {
    client ls >/dev/null 2>&1
}

wait_ping() {
    for c in $(seq 30) ; do
        if (( c > 1 )) ; then
            sleep 1
        fi
        if ping ; then
            return
        fi
    done
    return 1
}

case "$1" in
    "wait") wait_ping ;;
    *) client "$@" ;;
esac
