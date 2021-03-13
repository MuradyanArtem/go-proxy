DROP TABLE IF EXISTS requests CASCADE;
DROP TABLE IF EXISTS headers;

CREATE TABLE requests
(
    id      BIGSERIAL PRIMARY KEY NOT NULL,
    host    TEXT NOT NULL,
    request TEXT NOT NULL
);

CREATE TABLE headers
(
    req_id BIGINT REFERENCES requests (id) ON DELETE CASCADE NOT NULL,
    key    TEXT NOT NULL,
    val    TEXT
);
