-- +goose Up
-- +goose StatementBegin
Create table users(
    id serial primary key ,
    username varchar(50) NOT NULL,
    password varchar not null,
    role varchar NOT NULL DEFAULT 'user'
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table users;
-- +goose StatementEnd
