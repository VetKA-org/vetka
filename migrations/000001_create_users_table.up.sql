CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS users (
    id    uuid DEFAULT gen_random_uuid () primary key,
    login varchar(128) not null unique,
    password text not null
);

CREATE INDEX uk_users_login ON users (login);

INSERT INTO users (login, password) VALUES (
  'admin',
  crypt('1q2w3e', gen_salt('bf', 8))
) ON CONFLICT (login) DO NOTHING;
