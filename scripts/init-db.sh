#!/usr/bin/env sh
set -eux

: "${PG_HOST:="127.0.0.1"}"
: "${PG_PORT:="5432"}"

: "${PG_USER_DEFAULT:="postgres"}"

: "${PG_DB_DEFAULT:="postgres"}"

psql -h "${PG_HOST}" -p "${PG_PORT}" -v "ON_ERROR_STOP=1" \
    -U "${PG_USER_DEFAULT}" -d "${PG_DB_DEFAULT}" \
    -c "DROP DATABASE IF EXISTS socialmaps;" \
    -c "DROP USER IF EXISTS socialmaps_api;" \
    -c "DROP USER IF EXISTS osm2pgsql;" \
    -c "CREATE USER socialmaps_api;" \
    -c "CREATE USER osm2pgsql;" \
    -c "CREATE DATABASE socialmaps OWNER = socialmaps_api;" \
    -c "GRANT CONNECT, CREATE ON DATABASE socialmaps TO osm2pgsql;"

psql -h "${PG_HOST}" -p "${PG_PORT}" -v "ON_ERROR_STOP=1" \
    -U "${PG_USER_DEFAULT}" -d "socialmaps" \
    -c "CREATE EXTENSION postgis;"

psql -h "${PG_HOST}" -p "${PG_PORT}" -v "ON_ERROR_STOP=1" \
    -U "osm2pgsql" -d "socialmaps" \
    -c "CREATE SCHEMA osm2pgsql AUTHORIZATION osm2pgsql;" \
    -c "GRANT USAGE ON SCHEMA osm2pgsql TO socialmaps_api;" \
    -c "ALTER DEFAULT PRIVILEGES IN SCHEMA osm2pgsql GRANT SELECT ON TABLES TO socialmaps_api;"

psql -h "${PG_HOST}" -p "${PG_PORT}" -v "ON_ERROR_STOP=1" \
    -U "socialmaps_api" -d "socialmaps" \
    --file=internal/database/schema.sql
