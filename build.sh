#!/usr/bin/env sh
set -eu
cd "$(dirname "$0")"
version="${VERSION:-2.0.0}"
go test -mod=vendor ./...
mkdir -p build
go build -mod=vendor -trimpath -ldflags "-X main.Version=${version}" -o build/fint-graphql-cli .
./build/fint-graphql-cli --version
