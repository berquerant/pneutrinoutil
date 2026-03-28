#!/bin/bash

readonly d="$(cd "$(dirname "$0")/.." || exit 1; pwd)"


client() {
    "${d}/tools/run.sh" kubectl exec -it deploy/pneutrinoutil-redis -- redis-cli "$@"
}

case "$1" in
    *) client "$@" ;;
esac
