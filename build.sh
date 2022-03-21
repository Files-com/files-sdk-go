#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
echo "${DIR}" || exit 1
cd "${DIR}" || exit 1
go mod tidy
go get -d -v
go install golang.org/x/tools/cmd/goimports@latest
goimports -w .
gofmt -s -w .
