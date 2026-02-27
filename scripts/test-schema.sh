#!/usr/bin/env bash
set -eux

: "${PG_HOST:="127.0.0.1"}"
: "${PG_PORT:="5432"}"
: "${PG_USER_DEFAULT:="postgres"}"

: "${PG_DB_DEFAULT:="postgres"}"
: "${PG_DB_TEST:="socialmaps_schema_test"}"

exit_code=0

# the cleanup function will be the exit point
cleanup () {
  psql -h "${PG_HOST}" -p "${PG_PORT}" -U "${PG_USER_DEFAULT}" --variable "ON_ERROR_STOP=1" -d "${PG_DB_DEFAULT}" -c "DROP DATABASE ${PG_DB_TEST};"

  exit "${exit_code}"
}

trap cleanup ERR INT TERM

psql -h "${PG_HOST}" -p "${PG_PORT}" -U "${PG_USER_DEFAULT}" --variable "ON_ERROR_STOP=1" -d "${PG_DB_DEFAULT}" --command="CREATE DATABASE ${PG_DB_TEST};"
psql -h "${PG_HOST}" -p "${PG_PORT}" -U "${PG_USER_DEFAULT}" --variable "ON_ERROR_STOP=1" -d "${PG_DB_TEST}" --file=internal/database/schema.sql

exit_code=$?
cleanup
