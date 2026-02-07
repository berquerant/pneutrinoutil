#!/bin/bash

set -e
set -o pipefail

readonly version="$1"
readonly binary="$2"

if [[ -x "$binary" ]] ; then
    exit
fi

log() {
    echo >&2 "$(basename "$0"): $*"
}

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

__arch="$(arch | sed -e 's/x86_64/amd64/' -e 's/aarch64/arm64/')"
readonly url="https://github.com/golangci/golangci-lint/releases/download/${version}/golangci-lint-${version#v}-${__os}-${__arch}.tar.gz"

tmpd="$(mktemp -d)"
cd "$tmpd"
readonly dest="pkg.tar.gz"
curl -L -s -o "$dest" "$url"
tar xzf "$dest"
mv "./golangci-lint-${version#v}-${__os}-${__arch}/golangci-lint" "$binary"
