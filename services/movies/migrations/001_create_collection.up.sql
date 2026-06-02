CREATE TABLE collection (
    user_id      UUID NOT NULL,
    movie_id     INTEGER NOT NULL,
    torrent_id   TEXT NOT NULL,
    movie_lenght INTEGER NOT NULL,
    progression  INTERVAL DEFAULT '0',

    CONSTRAINT fk_user_identify
        FOREIGN KEY (user_id)
        REFERENCES users(uuid)
        ON DELETE CASCADE
);
