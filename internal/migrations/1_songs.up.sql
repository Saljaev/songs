CREATE TABLE IF NOT EXISTS groups (
    id BIGSERIAL PRIMARY KEY NOT NULL,
    name TEXT UNIQUE NOT NULL
);

CREATE INDEX name_idx ON groups(name);

CREATE TABLE IF NOT EXISTS songs (
    id BIGSERIAL PRIMARY KEY NOT NULL,
    title TEXT NOT NULL,
    group_id BIGINT REFERENCES groups(id),
    release_date DATE,
    link TEXT,
    lyric TEXT
);

ALTER TABLE songs ADD CONSTRAINT song_unique UNIQUE (title, group_id);
CREATE INDEX release_date_idx ON songs(release_date);
