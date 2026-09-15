FROM golang:1.26.8

COPY go.mod go.sum /app/
WORKDIR /app

ENV CGO_ENABLED=0
ENV GOPATH=/go

RUN go mod download


