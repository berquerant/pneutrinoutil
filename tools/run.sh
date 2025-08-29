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
readonly binary="${bind}/${name}"

readonly mod="${d}/go.mod"
package_name() {
    grep -F "$name" "$mod" | grep -v indirect | xargs
}
case "$name" in
    go-arch-lint)
        "${d}/setup-go-arch-lint.sh" "$GO_ARCH_LINT_VERSION" "$binary"
        ;;
    *)
        if [[ "$(package_name | wc -l | xargs)" != "1" ]] ; then
            log "name(${name}) is invalid"
            exit 1
        fi
        readonly pkg="$(package_name)"
        if [[ ! -x "$binary" ]] ; then
            log "build ${binary}"
            go -C "$d" build -o "$binary" "$pkg"
        fi
        ;;
esac
"$binary" "$@"
