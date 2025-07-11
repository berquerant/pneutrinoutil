#!/bin/bash

readonly d="$(cd "$(dirname "$0")" || exit; pwd)"

client() {
    "${d}/docker.sh" exec -it mysql mysql -h "$MYSQL_HOST" "$@"
}

root() {
    client -uroot -p"$MYSQL_ROOT_PASSWORD" "$@"
}

user() {
    client -u"$MYSQL_USER" -p"$MYSQL_PASSWORD" "$@"
}

ping() {
    echo >&2 "ping mysql"
    "${d}/docker.sh" exec -it mysql mysqladmin ping -uroot -p"$MYSQL_ROOT_PASSWORD"
}

wait_ping() {
    for c in $(seq 30) ; do
        if (( c > 1 )) ; then
            sleep 1
        fi
        if ping ; then
            echo >&2 "ping mysql success!"
            return
        fi
    done
    return 1
}

readonly user="$1"
shift
case "$user" in
    "root") root "$@" ;;
    "user") user "$@" ;;
    "ping") ping ;;
    "wait") wait_ping ;;
    *)
        echo >&2 "USER requried: root or user"
        exit 1
        ;;
esac
