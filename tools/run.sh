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
    golangci-lint)
        "${d}/setup-golangci-lint.sh" "$GOLANGCI_LINT_VERSION" "$binary"
        ;;
    go-arch-lint)
        "${d}/setup-go-arch-lint.sh" "$GO_ARCH_LINT_VERSION" "$binary"
        ;;
    shellcheck)
        "${d}/shellcheck.sh" "$SHELLCHECK_VERSION" "$binary" "$@"
        exit
        ;;
    ls-lint)
        "${d}/setup-ls-lint.sh" "$LS_LINT_VERSION" "$binary"
        ;;
    kind)
        "${d}/setup-kind.sh" "$KIND_VERSION" "$binary"
        ;;
    helm)
        "${d}/setup-helm.sh" "$HELM_VERSION" "$binary"
        ;;
    kubectl)
        "${d}/setup-kubectl.sh" "$KUBECTL_VERSION" "$binary"
        ;;
    stern)
        "${d}/setup-stern.sh" "$STERN_VERSION" "$binary"
        ;;
    yq)
        readonly yq_version=4.52.5
        "${d}/setup-yq.sh" "$yq_version" "$binary"
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
