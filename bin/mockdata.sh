#!/bin/bash

set -e
set -o pipefail

readonly d="$(cd "$(dirname "$0")" || exit 1; pwd)"
readonly mockcli="./dist/pneutrinoutil-mockcli"
readonly worker="${d}/worker.sh"
readonly docker="${d}/docker.sh"

task() {
    "${d}/../task" --dir "${d}/.." "$@"
}

tmpd="$(mktemp -d)"
mkdir -p "$tmpd"

post() {
    local -r basename="$1"
    local -r content="$2"
    local -r url="${SERVER_URI}/proc"
    local -r score="${tmpd}/${basename}.musicxml"
    echo "$content" > "$score"
    echo >&2 "Create proc from ${score}"
    curl -s -D- -X POST "$url" \
         -H 'accept: application/json' \
         -H 'Content-Type: multipart/form-data' \
         -F "score=@${score}"
}

usage() {
    cat <<EOS
Generate dummy artifacts using mockcli.
mockdata.sh COUNT [BASENAME] [SCORE_CONTENT]

The worker will be started using mockcli.
Please execute task restart-worker once the worker's process finishes,
as we will restart it using the normal cli.
EOS
}

main() {
    local -r count="${1}"
    local -r basename="${2:-mockdata_basename}"
    local -r basecontent="${3:-mockdata_content}"
    task ping-infra
    task build-mockcli
    task build-worker
    "$docker" up -d server
    "$worker" stop
    sleep 3
    for i in $(seq "$count") ; do
        post "${basename}_${i}" \
             "${basecontent}_${i}"
    done
    PNEUTRINOUTIL="$mockcli" "$worker" start
}

case "$1" in
    "" | "-h" | "--help")
        usage
        ;;
    *) main "$@" ;;
esac
