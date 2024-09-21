-- +goose Up
CREATE TABLE white_list (
    id      serial PRIMARY KEY,
    cidr   text not null,
    created_at  timestamptz not null
);
CREATE TABLE black_list (
    id      serial PRIMARY KEY,
    cidr   text not null,
    created_at  timestamptz not null
);
CREATE TABLE settings (
    id      serial PRIMARY KEY,
    login_count int not null,
    password_count int not null,
    ip_count int not null
);

INSERT INTO white_list(cidr, created_at)
VALUES ('192.1.10.0/30', NOW()),
       ('192.2.10.0/30', NOW());

INSERT INTO black_list(cidr, created_at)
VALUES ('192.100.10.0/30', NOW()),
       ('192.200.10.0/30', NOW());

INSERT INTO settings (login_count, password_count, ip_count)
VALUES (10, 100, 1000);

-- +goose Down
drop table white_list;
drop table black_list;
drop table settings;
