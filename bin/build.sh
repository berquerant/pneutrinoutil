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

build() {
    go build -v -trimpath -ldflags "$(ldflags)" "$@"
}

build "$@"
