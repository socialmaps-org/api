#!/usr/bin/env sh
set -eux

: "${PG_HOST:="127.0.0.1"}"
: "${PG_PORT:="5432"}"

: "${PG_USER_DEFAULT:="postgres"}"
: "${PG_USER_SM_API:="socialmaps_api"}"

: "${PG_DB_DEFAULT:="postgres"}"
: "${PG_DB_SOCIALMAPS:="socialmaps"}"


psql -h "${PG_HOST}" -p "${PG_PORT}" -v "ON_ERROR_STOP=1" \
    -U "${PG_USER_DEFAULT}" -d "${PG_DB_DEFAULT}" \
    -c "DROP DATABASE IF EXISTS ${PG_DB_SOCIALMAPS};" \
    -c "DROP USER IF EXISTS ${PG_USER_SM_API};" \
    -c "CREATE USER ${PG_USER_SM_API};" \
    -c "CREATE DATABASE ${PG_DB_SOCIALMAPS} OWNER = ${PG_USER_SM_API};"

psql -h "${PG_HOST}" -p "${PG_PORT}" -v "ON_ERROR_STOP=1" \
    -U "${PG_USER_SM_API}" -d "${PG_DB_SOCIALMAPS}" \
    --file=internal/database/schema.sql
