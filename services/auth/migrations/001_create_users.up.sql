CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE users (
	uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	username TEXT UNIQUE NOT NULL,
	password_hash TEXT NOT NULL
);
