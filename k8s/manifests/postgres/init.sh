#!/bin/sh
set -eu

psql -v ON_ERROR_STOP=1 -U "$POSTGRES_USER" -d "$POSTGRES_DB" <<EOF
CREATE SCHEMA IF NOT EXISTS internal;
EOF
