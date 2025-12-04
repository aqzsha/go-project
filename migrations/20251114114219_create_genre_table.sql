-- +goose Up
-- +goose StatementBegin
create table genre
(
    id bigserial primary key,
    name varchar(255) not null
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists genre
-- +goose StatementEnd
