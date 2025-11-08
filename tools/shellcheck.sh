#!/bin/bash

set -e
set -o pipefail

readonly version="$1"
readonly binary="$2"
shift 2

log() {
    echo >&2 "$(basename "$0"): $*"
}

install() {
    if [[ -x "$binary" ]] ; then
        return
    fi
    __os="$(uname -s)"
    case "$__os" in
        "Darwin" | "Linux")
            __os="$(echo "$__os" | tr '[:upper:]' '[:lower:]')"
            ;;
        *)
            log "${__os} not supported"
            exit 1
            ;;
    esac

    __arch="$(arch | sed -e 's/amd64/x86_64/' -e 's/arm64/aarch64/')"
    readonly url="https://github.com/koalaman/shellcheck/releases/download/${version}/shellcheck-${version}.${__os}.${__arch}.tar.xz"

    tmpd="$(mktemp -d)"
    cd "$tmpd"
    readonly dest="pkg.tar.xz"
    curl -L -s -o "$dest" "$url"
    tar xzf "$dest"
    mv "./shellcheck-${version}/shellcheck" "$binary"
}


if [[ -n "$CI" ]] ; then
    shellcheck "$@"
else
    install
    "$binary" "$@"
fi
