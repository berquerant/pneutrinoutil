#!/bin/bash

set -e
set -o pipefail

readonly d="$(cd "$(dirname "$0")" || exit 1; pwd)"
readonly bind="${d}/../bin/tools"
mkdir -p "$bind"

log() {
    echo >&2 "$(basename "$0"): $*"
}

readonly name="$1"
if [[ -z "$name" ]] ; then
    log "name(\$1) required"
    exit 1
fi
shift

if command -v mise >/dev/null 2>&1 ; then
    exec mise exec -- "$name" "$@"
elif [ -x "${HOME}/.local/bin/mise" ]; then
    exec "${HOME}/.local/bin/mise" exec -- "$name" "$@"
fi

exec "$name" "$@"
