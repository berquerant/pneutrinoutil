#!/bin/bash

readonly d="$(cd "$(dirname "$0")/.." || exit 1; pwd)"

log() {
    echo >&2 "$*"
}

database() {
    "${d}/bin/mysql.sh" root "$@"
}

init_db() {
    log "init_db"
    database -e "$(cat "${d}/ddl/mysql/db.sql")"
}

init_tables() {
    local -r db="$1"
    log "init_tables ${db}"
    database "$db" -e "$(cat "${d}/ddl/mysql/tables.sql")"
}

init_users() {
    log "init_users"
    database -e "$(cat "${d}/ddl/mysql/users.sql")"
}

storage() {
    "${d}/bin/s3.sh" "$@"
}

exist_storage() {
    local -r bucket="$1"
    storage ls | awk '{print $3}' | grep -q "${bucket}"
}

init_storage() {
    log "init_storage"
    while read -r bucket ; do
        log "bucket($bucket)"
        if ! exist_storage "${bucket}" ; then
            storage mb "s3://${bucket}"
        fi
    done < "${d}/ddl/s3/bucket.txt"
}

kvs() {
    "${d}/bin/redis.sh" "$@"
}

drop_db() {
    local -r db="$1"
    log "drop_db ${db}"
    database -e "drop database if exists ${db};"
}

drop_storage() {
    local -r bucket="$1"
    log "drop_storage ${bucket}"
    if exist_storage "${bucket}" ; then
        storage rb --force "s3://${bucket}"
    fi
}

drop_kvs() {
    local -r db="$1"
    log "drop_kvs ${db}"
    kvs -n "${db}" FLUSHDB
}


cmd="$1"
shift
case "$cmd" in
    db) init_db ;;
    tables) init_tables "$1" ;;
    users) init_users ;;
    storage) init_storage ;;
    drop)
        cmd="$1"
        shift
        case "$cmd" in
            db) drop_db "$1" ;;
            storage) drop_storage "$1" ;;
            kvs) drop_kvs "$1" ;;
            *) exit 1 ;;
        esac
        ;;
    *)
        cat <<EOS >&2
$0 db
$0 tables DATABASE
$0 users
$0 storage
$0 drop (db|storage|kvs) DB_OR_BUCKET
EOS
        exit 1
        ;;
esac
