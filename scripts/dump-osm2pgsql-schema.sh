#!/usr/bin/env sh
set -eux

pg_dump --schema-only --dbname=socialmaps --schema=osm2pgsql \
    | sed '/^--/d; /^$/d; /^\\restrict/d; /^\\unrestrict/d;' \
    | tee internal/database/osm2pgsql-schema.sql
