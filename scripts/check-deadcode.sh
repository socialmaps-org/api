#!/usr/bin/env sh
set -eux

out=$(deadcode -test ./...)

if [ -n "${out}" ]; then
    # The output is not empty -> there is dead code
    echo "${out}"
    exit 1
else
    exit 0
fi
