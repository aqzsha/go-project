-- +goose Up
-- +goose StatementBegin
create table person
(
    id bigserial primary key,
    full_name varchar(255) not null
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists person
-- +goose StatementEnd
