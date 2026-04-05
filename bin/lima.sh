#!/bin/bash

readonly d="$(cd "$(dirname "$0")" || exit 1; pwd)"
readonly name="pneutrinoutil"

readonly uv_version="0.11.3"
readonly pnpm_version="10.33.0"

go_version() {
    grep -E "^go \d+\.\d+\.\d+" "${d}/../go.mod" | awk '{print $2}'
}

start() {
    local __go_version
    __go_version="$(go_version)"

    limactl start \
            --name "$name" \
            --yes \
            --mount-none \
            --cpus=4 \
            --memory=4 \
            --disk=50 \
            template:docker
    limactl copy "${d}/lima-setup.sh" "${name}:/tmp/"
    limactl shell "$name" /tmp/lima-setup.sh "${__go_version}" "${uv_version}" "${pnpm_version}"
    ssh
}

stop() {
    limactl stop "$name"
}

reload() {
    stop || true
    limactl delete "$name" || true
    start
}

ssh() {
    exec limactl shell "$name"
}

set -ex
readonly cmd="$1"
case "$cmd" in
    "ssh") ssh ;;
    "start") start ;;
    "stop") stop ;;
    "reload") reload ;;
    *)
        echo>&2 "Please select start, stop or reload"
        exit 1
        ;;
esac
