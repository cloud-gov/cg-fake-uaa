#! /bin/sh

set -e

export GOOS=linux
export GOARCH=amd64

go build

docker build -t cg-fake-uaa:latest .

rm cg-fake-uaa
