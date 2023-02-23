CREATE TABLE IF NOT EXISTS roles (
    role_id uuid DEFAULT gen_random_uuid () primary key,
    name varchar(128) not null unique
);

CREATE TABLE IF NOT EXISTS users_roles (
    user_id uuid REFERENCES users (user_id) on delete cascade,
    role_id uuid REFERENCES roles (role_id) on delete cascade,
    primary key (user_id, role_id)
);

INSERT INTO roles (name) VALUES (
    unnest(array['head-doctor', 'administrator', 'doctor'])
) ON CONFLICT DO NOTHING;

INSERT INTO users_roles (user_id, role_id) VALUES (
    (SELECT user_id FROM users WHERE login = 'head'),
    (SELECT role_id FROM roles WHERE name = 'head-doctor')
) ON CONFLICT (user_id, role_id) DO NOTHING;
