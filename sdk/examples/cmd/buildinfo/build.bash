#!/usr/bin/env bash

set -eof pipefail

gitCommit=$(git rev-parse HEAD)
gitTagRef=$(git name-rev --name-only --tags "${gitCommit}")
branch=$(git rev-parse --abbrev-ref HEAD)
echo $gitTagRef
go build -v -ldflags="-X go.clly.me/minic/sdk.branch=${branch} -X go.clly.me/minic/sdk.tag=${gitTagRef}" .
./buildinfo
