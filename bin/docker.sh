#!/bin/bash

readonly d="$(cd "$(dirname "$0")/.." || exit; pwd)"

client() {
    docker compose "$@"
}


case "$1" in
    db | database)
        shift
        client -f "${d}/docker-database.yml" "$@"
        ;;
    s3 | storage)
        shift
        client -f "${d}/docker-storage.yml" "$@"
        ;;
    *)
        client "$@"
        ;;
esac
