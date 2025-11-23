#!/bin/bash

client() {
    aws s3 "$@"
}

ping() {
    echo >&2 "ping s3"
    client ls
}

wait_ping() {
    for c in $(seq 30) ; do
        if (( c > 1 )) ; then
            sleep 1
        fi
        if ping ; then
            echo >&2 "ping s3 success!"
            return
        fi
    done
    return 1
}

case "$1" in
    "wait") wait_ping ;;
    "ping") ping ;;
    *) client "$@" ;;
esac
