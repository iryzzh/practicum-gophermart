CREATE TABLE IF NOT EXISTS users
(
    id    serial PRIMARY KEY,
    login varchar(255) unique not null,
    encrypted_password varchar(255) not null,
    deleted bool default false,
    deleted_at timestamp
);

create unique index idx_users_login on users (login);