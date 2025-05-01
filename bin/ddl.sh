#!/bin/bash

readonly d="$(cd "$(dirname "$0")/.." || exit; pwd)"

client() {
    "${d}/bin/mysql.sh" root "$@"
}

init_db() {
    client -e "$(cat ${d}/ddl/mysql/db.sql)"
}

init_tables() {
    local -r db="$1"
    client "$db" -e "$(cat ${d}/ddl/mysql/tables.sql)"
}

init_users() {
    client -e "$(cat ${d}/ddl/mysql/users.sql)"
}

readonly cmd="$1"
shift
case "$cmd" in
    db) init_db ;;
    tables) init_tables "$1" ;;
    users) init_users ;;
    *)
        cat <<EOS >&2
$0 db
$0 tables DATABASE
$0 users
EOS
        exit 1
        ;;
esac
