#!/usr/bin/env sh
set -eux

out=$(gofmt -e -l -s .)

if [ -n "${out}" ]; then
    # The output is not empty -> there are formatting issues
    echo "${out}"
    exit 1
else
    exit 0
fi
