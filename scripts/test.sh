#!/usr/bin/env sh
set -eux

go test -tags sqlite_math_functions -count=1 ./...
