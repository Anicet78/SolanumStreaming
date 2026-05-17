CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TYPE roles AS ENUM ('user', 'admin');
CREATE TYPE languages AS ENUM ('en', 'fr');

CREATE TABLE users (
    uuid     UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    username TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    role     roles NOT NULL DEFAULT 'user',
    language languages NOT NULL DEFAULT 'fr'
);
