CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TYPE roles AS ENUM ('user', 'admin');

CREATE TABLE users (
    uuid UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    username TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    "role" roles NOT NULL DEFAULT 'user'
);