#!/usr/bin/env sh
set -eux

shellcheck --enable=all scripts/*.sh
