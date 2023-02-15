CREATE TABLE IF NOT EXISTS statuses(
    id   SERIAL NOT NULL PRIMARY KEY,
    name VARCHAR
);

INSERT INTO statuses (id, name) VALUES(1, 'new');
INSERT INTO statuses (id, name) VALUES(2, 'in process');
INSERT INTO statuses (id, name) VALUES(3, 'error');
INSERT INTO statuses (id, name) VALUES(4, 'done');