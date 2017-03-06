#!/usr/bin/env bash

set -e

USAGE="
Usage $(basename "$0") COMMAND

Avaliable commands:
   coverage-report     Run unit tests and generate coverage.txt file (useful for codecov)
   e2e                 Run e2e tests
"
RESET_COLOR="\033[0m"

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

function echo_green(){
    echo -e "\033[32m$1${RESET_COLOR}"
}

function echo_red() {
    echo -e "\033[1;31m$1${RESET_COLOR}"
}

function echo_yellow() {
    echo -e "\033[1;33m$1${RESET_COLOR}"
}

function check_result() {
    local got="$1"
    local expected="$2"

    if [ "$got" != "$expected" ]; then
        echo_red "expected $expected, got $got"
        exit 1
    fi
}

function _e2e_single_file() {
    printf "Testing with a single file... "
    local cmd="./sif Cgo _tests/golang.txt"

    check_result "$($cmd | wc -l | tr -d '[:space:]')" "3"
    check_result "$($cmd | head -n 1)" "_tests/golang.txt"
    check_result "$($cmd | tail -n 1)" "12: Cgo is not Go."

    echo_green OK
}

function _e2e_dir(){
    printf "Testing with a dir... "
    local cmd="./sif better _tests"

    check_result "$($cmd | grep -v '[0-9]' | xargs)" "_tests/golang.txt _tests/python.txt"
    check_result "$($cmd | wc -l | tr -d '[:space:]')" "13"
    check_result "$($cmd | tail -n 1)" "18: Although never is often better than *right* now."

    echo_green "OK"
}

function _e2e_build() {
    printf "build app... "
    make build
    echo_green "OK"
}

function _e2e_without_target() {
    printf "Testing without a target... "
    local cmd="./sif Hello"

    check_result "$($cmd | wc -l | tr -d '[:space:]')" "5"
    check_result "$($cmd | grep -v '[0-9]' | xargs)" "$(basename $0) _tests/subdir/hello.go"

    echo_green "OK"
}

function _e2e_with_two_targets(){
    printf "Testing with two targets... "
    local cmd="./sif better _tests/python.txt _tests/golang.txt"

    check_result "$($cmd | grep -v '[0-9]' | xargs)" "_tests/python.txt _tests/golang.txt"
    check_result "$($cmd | wc -l | tr -d '[:space:]')" "13"
    check_result "$($cmd | tail -n 1)" "14: Clear is better than clever."

    echo_green "OK"
}

function _e2e_case_insensitive_flag() {
    printf "Testing with -i flag (case insensitive)... "
    local cmd="./sif -i cgo _tests/golang.txt"

    check_result "$($cmd | wc -l | tr -d '[:space:]')" "3"
    check_result "$($cmd | head -n 1)" "_tests/golang.txt"
    check_result "$($cmd | tail -n 1)" "12: Cgo is not Go."

    echo_green "OK"
}

function _e2e_only_filenames_flag() {
    printf "Testing with -l flag (print only file names)... "
    local cmd="./sif -l Cgo"

    check_result "$($cmd | wc -l | tr -d '[:space:]')" "2"
    check_result "$($cmd | head -n 1)" "$(basename $0)"
    check_result "$($cmd | tail -n 1)" "_tests/golang.txt"

    echo_green "OK"
}

function e2e() {
    echo_yellow "START E2E TESTS"
    _e2e_build
    _e2e_without_target
    _e2e_single_file
    _e2e_dir
    _e2e_with_two_targets

    echo_yellow " - FLAGS TESTS"
    _e2e_case_insensitive_flag
    _e2e_only_filenames_flag

    rm sif
    echo_yellow "FINISH"
}

case "$1" in
    "coverage-report") coverage_report ;;
    "e2e") e2e ;;
    *) echo "$USAGE" ;;
esac
