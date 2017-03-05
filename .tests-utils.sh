#!/usr/bin/env bash

set -e

USAGE="
Usage $(basename "$0") COMMAND

Avaliable commands:
   coverage-report     Run unit tests and generate coverage.txt file (useful for codecov)
"

function coverage_report() {
    echo "" > coverage.txt

    for d in $(go list ./... | grep -v vendor); do
        go test -race -coverprofile=profile.out -covermode=atomic $d
        if [ -f profile.out ]; then
            cat profile.out >> coverage.txt
            rm profile.out
        fi
    done
}

case "$1" in
    "coverage-report") coverage_report ;;
    *) echo "$USAGE" ;;
esac
