#!/usr/bin/env bash

CC="zig cc -target x86_64-linux-gnu" CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -tags sqlite_math_functions -o bin/socialmaps-api-linux-amd64 cmd/socialmaps-api/main.go

scp bin/socialmaps-api-linux-amd64 john@49.13.30.5:/home/john/
