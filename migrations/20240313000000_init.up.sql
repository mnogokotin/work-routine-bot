CREATE TABLE IF NOT EXISTS users
(
    id       SERIAL PRIMARY KEY,
    username TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS projects
(
    id   SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS tasks
(
    id          SERIAL PRIMARY KEY,
    user_id     INT       NOT NULL,
    project_id  INT       NOT NULL,
    description TEXT      NOT NULL,
    duration    FLOAT8    NOT NULL,
    date        TIMESTAMP NOT NULL,
    created_at  TIMESTAMP NOT NULL,
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users (id) ON UPDATE CASCADE ON DELETE RESTRICT,
    CONSTRAINT fk_project FOREIGN KEY (project_id) REFERENCES projects (id) ON UPDATE CASCADE ON DELETE RESTRICT
);
