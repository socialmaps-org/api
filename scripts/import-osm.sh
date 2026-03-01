#!/usr/bin/env sh
set -eux

: "${PG_HOST:="127.0.0.1"}"
: "${PG_PORT:="5432"}"

extract="${1:-assets/monaco-260226.osm.pbf}"

osm2pgsql \
    --host "${PG_HOST}" \
    --port "${PG_PORT}" \
    --user osm2pgsql \
    --schema osm2pgsql \
    --database socialmaps \
    -O flex -S scripts/osm2pgsql.lua "${extract}"
