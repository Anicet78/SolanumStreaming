#!/bin/sh
set -eu

psql -v ON_ERROR_STOP=1 -U "$POSTGRES_USER" -d postgres <<EOF
SELECT format('CREATE DATABASE %I', '$POSTGRES_DB')
WHERE NOT EXISTS (SELECT 1 FROM pg_database WHERE datname = '$POSTGRES_DB')\gexec
EOF