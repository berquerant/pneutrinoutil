#!/bin/bash

set -e

readonly d="$(cd "$(dirname "$0")/.." || exit 1; pwd)"

if [[ "$TEST" = "true" ]] ; then
    echo >&2 "env.sh: use test environment!"
    export MYSQL_DATABASE="$TEST_MYSQL_DATABASE"
    export MYSQL_USER="$TEST_MYSQL_USER"
    export MYSQL_PASSWORD="$TEST_MYSQL_PASSWORD"
    export REDIS_DB="$TEST_REDIS_DB"
    export STORAGE_BUCKET="$TEST_STORAGE_BUCKET"
    export MYSQL_DSN="$TEST_MYSQL_DSN"
    export REDIS_DSN="$TEST_REDIS_DSN"
    export PNEUTRINOUTIL="${d}/dist/pneutrinoutil-mockcli"
    worker_work_dir="$(mktemp -d)"
    export WORKDIR="$worker_work_dir"
fi

"$@"
