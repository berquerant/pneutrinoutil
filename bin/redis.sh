#!/bin/bash


client() {
    kubectl exec -it deploy/pneutrinoutil-redis -- redis-cli "$@"
}

case "$1" in
    *) client "$@" ;;
esac
