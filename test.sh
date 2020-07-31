#!/usr/bin/env bash

cd generated/go || exit 1
go mod tidy
go get -d -v
go fmt . -s
go test ./...
