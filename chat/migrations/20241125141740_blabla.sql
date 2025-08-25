-- +goose Up
-- +goose StatementBegin
CREATE TABLE chatv1(
    id serial primary key,
    users text[],
    created_at timestamp not null default now(),
    updated_at timestamp
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table chatv1
-- +goose StatementEnd
