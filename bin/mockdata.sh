#!/bin/bash

set -e
set -o pipefail

readonly d="$(cd "$(dirname "$0")" || exit 1; pwd)"
readonly worker="${d}/worker.sh"
readonly docker="${d}/docker.sh"
readonly mockcli="${d}/../dist/pneutrinoutil-mockcli"
readonly gendata="${d}/../dist/pneutrinoutil-gendata"

task() {
    "${d}/../task" --dir "${d}/.." "$@"
}

tmpd="$(mktemp -d)"
mkdir -p "$tmpd"
readonly dummycli="${tmpd}/pneutrinoutil-dummycli"

usage() {
    cat <<EOS
Generate dummy artifacts using mockcli.
mockdata.sh COUNT [BASENAME] [SCORE_CONTENT]

The worker will be started using mockcli.
Please execute task restart-worker once the worker's process finishes,
as we will restart it using the normal cli.

You can create a failed artifact by setting FAIL.
You can delay artifact generation by setting DURATION.
EOS
}

prepare_dummycli() {
    local -r fail="$1"
    local -r duration="$2"
    local cmd="$mockcli"
    if [[ -n "$fail" ]] ; then
        cmd="${cmd} --fail"
    fi
    if [[ -n "$duration" ]] ; then
        cmd="${cmd} --duration ${duration}"
    fi
    cat <<EOS > "$dummycli"
#!/bin/bash
${cmd} "\$@"
EOS
    chmod +x "$dummycli"
}

main() {
    local -r count="${1}"
    local -r basename="${2:-mockdata_basename}"
    local -r content="${3:-mockdata_content}"
    prepare_dummycli "$FAIL" "$DURATION"
    task ping-infra
    task build-mockcli
    task build-gendata
    task build-worker
    "$docker" up -d server
    "$worker" stop
    sleep 3
    PNEUTRINOUTIL="$dummycli" "$worker" start
    seq "$count" | awk -v x="$basename" '{print x"_"$0}' | "$gendata" -c "$content"
    "$worker" stop
    "$worker" start
}

case "$1" in
    "" | "-h" | "--help")
        usage
        ;;
    *) main "$@" ;;
esac
