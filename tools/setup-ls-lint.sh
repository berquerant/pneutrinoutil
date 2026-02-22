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
set -x
readonly url="https://github.com/loeffel-io/ls-lint/releases/download/${version}/ls-lint-${__os}-${__arch}"

curl -L -s -o "$binary" "$url"
chmod +x "$binary"
