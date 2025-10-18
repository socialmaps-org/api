#!/usr/bin/env bash

out=$(deadcode -test ./...)

if [[ $out ]]; then
    # The output is not empty -> there is dead code
    echo "$out"
    exit 1
else
    exit 0
fi
