CREATE TABLE IF NOT EXISTS roles (
    id   uuid DEFAULT gen_random_uuid () primary key,
    name varchar(128) not null unique
);

CREATE TABLE IF NOT EXISTS users_roles (
    user_id uuid,
    role_id uuid,
    primary key (user_id, role_id),
    foreign key (user_id) REFERENCES users (id) on delete cascade,
    foreign key (role_id) REFERENCES roles (id) on delete cascade
);

INSERT INTO roles (name) VALUES ('head-doctor') ON CONFLICT (name) DO NOTHING;
INSERT INTO roles (name) VALUES ('administrator') ON CONFLICT (name) DO NOTHING;
INSERT INTO roles (name) VALUES ('doctor') ON CONFLICT (name) DO NOTHING;

INSERT INTO users_roles (user_id, role_id) VALUES (
    (SELECT id FROM users WHERE login = 'head'),
    (SELECT id FROM roles WHERE name = 'head-doctor')
) ON CONFLICT (user_id, role_id) DO NOTHING;
