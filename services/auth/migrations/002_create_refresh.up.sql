CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE refresh (
    token_id   UUID PRIMARY KEY DEFAULT gen_random_uuid (),
	user_id    UUID NOT NULL,
    token_hash TEXT UNIQUE NOT NULL,
    expires_at timestamptz NOT NULL,
	revoked_at timestamptz,
    CONSTRAINT fk_user_identify FOREIGN KEY (user_id) REFERENCES users (uuid) ON DELETE CASCADE
);
