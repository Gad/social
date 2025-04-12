CREATE TABLE IF NOT EXISTS roles (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    level int NOT NULL DEFAULT 0,
    description TEXT
);

INSERT INTO roles (name, level, description) VALUES
('admin',3, 'an admin can update and delete any post and comment'),
('user', 1, 'a User can create and manage posts and comments'),
('moderator', 2, ' a moderator can update posts and comments from other users'),
('guest', 0, 'a guest can only read posts and comments');