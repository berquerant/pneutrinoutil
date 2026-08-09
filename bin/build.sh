#!/bin/bash

short_sha() {
    git rev-parse --short HEAD 2>/dev/null || echo "unknown"
}

current_tag() {
    if [ -f "VERSION" ]; then
        echo "v$(cat VERSION | tr -d ' \n\r')"
    else
        echo "unknown"
    fi
}

build_date() {
    date -u +'%Y-%m-%dT%H:%M:%SZ'
}

readonly version_package="github.com/berquerant/pneutrinoutil/pkg/version"

ldflags() {
    echo "-X ${version_package}.Version=$(current_tag) -X ${version_package}.Revision=$(short_sha) -X ${version_package}.BuildDate=$(build_date)"
}

go_bin() {
    if command -v go >/dev/null 2>&1; then
        echo "go"
    elif command -v mise >/dev/null 2>&1; then
        echo "mise exec -- go"
    elif [ -x "${HOME}/.local/bin/mise" ]; then
        echo "${HOME}/.local/bin/mise exec -- go"
    else
        echo "go"
    fi
}

build() {
    $(go_bin) build -v -trimpath -ldflags "$(ldflags)" "$@"
}

build "$@"
