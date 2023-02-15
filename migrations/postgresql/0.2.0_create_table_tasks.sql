CREATE TABLE IF NOT EXISTS tasks(
    id               SERIAL NOT NULL PRIMARY KEY,
    status_id        INT REFERENCES statuses(id),
    method           VARCHAR,
    content_length   INT,
    http_status_code INT,
    url              VARCHAR,
    header           JSONB,
    created_at       BIGINT,
    updated_at       BIGINT,
    deleted_at       BIGINT,
);