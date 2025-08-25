-- +goose Up
-- +goose StatementBegin
CREATE TABLE chat_log(
    id serial primary key,
    time timestamp not null default now(),
    chat_id integer not null,
    action text not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table chat_log
-- +goose StatementEnd
