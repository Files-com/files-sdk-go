#!/usr/bin/env bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
echo "${DIR}"
cd "${DIR}" || exit 1
go mod tidy
go get -d -v
go get golang.org/x/tools/cmd/goimports
goimports -w .
gofmt -s -w .
