#!/bin/bash
# This script performs tests against the dcos-signal project, specifically:
#
#   * gofmt         (https://golang.org/cmd/gofmt)
#   * goimports     (https://godoc.org/cmd/goimports)
#   * golint        (https://github.com/golang/lint)
#   * go vet        (https://golang.org/cmd/vet)
#   * test coverage (https://blog.golang.org/cover)
#
# It outputs test and coverage reports in a way that Jenkins can understand,
# with test results in JUnit format and test coverage in Cobertura format.
# The reports are saved to build/$SUBDIR/{test-reports,coverage-reports}/*.xml 
#
set -e
set -o pipefail
export PATH="${GOPATH}/bin:${PATH}"

SUBDIRS="api config cmd"
SOURCE_DIR=$(git rev-parse --show-toplevel)
BUILD_DIR="${SOURCE_DIR}/build"


function logmsg {
    echo -e "\n\n*** $1 ***\n"
}

function _lint {
    logmsg "Running linter..."
    golangci-lint -v run
}

function _unittest_with_coverage {
    logmsg "Running unit tests..."
    go test -v -cover -mod=vendor -tags=unit ./...
}


# Main.
function main {
    _lint
    _unittest_with_coverage
}

main
