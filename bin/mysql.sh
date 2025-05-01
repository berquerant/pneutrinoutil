#!/bin/bash

client() {
    docker compose exec -it mysql mysql "$@"
}

root() {
    client -uroot -p"$MYSQL_ROOT_PASSWORD" "$@"
}

user() {
    client -u"$MYSQL_USER" -p"$MYSQL_PASSWORD" "$@"
}

readonly user="$1"
shift
case "$user" in
    "root") root "$@" ;;
    "user") user "$@" ;;
    *)
        echo >&2 "USER requried: root or user"
        exit 1
        ;;
esac
