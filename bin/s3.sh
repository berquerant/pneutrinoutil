#!/bin/bash

client() {
    aws s3 "$@"
}

case "$1" in
    *) client "$@" ;;
esac
