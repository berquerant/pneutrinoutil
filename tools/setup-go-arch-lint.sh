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
        os_name="$(echo "$__os" | tr '[:upper:]' '[:lower:]')"
        ;;
    *)
        log "${__os} not supported"
        exit 1
        ;;
esac

__arch="$(arch | sed -e 's/x86_64/amd64/' -e 's/aarch64/arm64/')"
readonly url="https://github.com/fe3dback/go-arch-lint/releases/download/${version}/go-arch-lint_${version#v}_${__os}_${__arch}.tar.gz"

tmpd="$(mktemp -d)"
cd "$tmpd"
readonly dest="pkg.tar.gz"
curl -L -s -o "$dest" "$url"
tar xzf "$dest"
mv "./go-arch-lint" "$binary"
