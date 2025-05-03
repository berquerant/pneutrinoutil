#!/bin/bash

limit=20
interval=1

while [[ $# -gt 0 ]] ; do
    case "$1" in
        "--limit")
            if [[ -z "$2" ]] ; then
                echo >&2 "$0: --limit requires an argument"
                exit 1
            fi
            limit="$2"
            shift 2
            ;;
        "--interval")
            if [[ -z "$2" ]] ; then
                echo >&2 "$0: --interval requires an argument"
                exit 1
            fi
            interval="$2"
            shift 2
            ;;
        "--")
            shift
            break
            ;;
        *)
            echo >&2 "$0: unknown arguments!"
            exit 1
            ;;
    esac
done

for c in in $(seq "$limit") ; do
    echo >&2 "retry[$1][${c}]..."
    if ((c > 1)) ; then
        sleep "$interval"
    fi
    if "$@" ; then
        exit
    fi
done

echo >&2 "retry[$1] failed!"
exit 1
