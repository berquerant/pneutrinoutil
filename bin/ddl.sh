#!/bin/bash

readonly d="$(cd "$(dirname "$0")/.." || exit 1; pwd)"

log() {
    echo >&2 "$*"
}

database() {
    "${d}/bin/mysql.sh" root "$@"
}

s3() {
    "${d}/bin/s3.sh" "$@"
}

exist_bucket() {
    local -r bucket="$1"
    s3 ls | awk '{print $3}' | grep -q "${bucket}"
}

redis() {
    "${d}/bin/redis.sh" "$@"
}

drop_mysql() {
    local -r db="$1"
    log "drop_mysql ${db}"
    database -e "drop database if exists ${db};"
}

drop_s3() {
    local -r bucket="$1"
    log "drop_s3 ${bucket}"
    if exist_bucket "${bucket}" ; then
        s3 rb --force "s3://${bucket}"
    fi
    create_s3 "$bucket"
}

create_s3() {
    local -r bucket="$1"
    log "create_s3 ${bucket}"
    if exist_bucket "${bucket}" ; then
        s3 mb "s3://${bucket}"
    fi
}

drop_redis() {
    local -r db="$1"
    log "drop_redis ${db}"
    redis -n "${db}" FLUSHDB
}


cmd="$1"
shift
case "$cmd" in
    s3)
        create_s3 "$1"
        ;;
    drop)
        cmd="$1"
        shift
        case "$cmd" in
            mysql) drop_mysql "$1" ;;
            s3) drop_s3 "$1" ;;
            redis) drop_redis "$1" ;;
            *) exit 1 ;;
        esac
        ;;
    *)
        cat <<EOS >&2
$0 drop (mysql|s3|redis) DB_OR_BUCKET
EOS
        exit 1
        ;;
esac
