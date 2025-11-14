-- +goose Up
-- +goose StatementBegin
create table cinema_details
(
    id bigserial primary key,
    name varchar(255) not null,
    description text,
    address varchar(255) not null,
    latitude  numeric(9,6) not null,
    longitude numeric(9,6) not null,
    created_at timestamp default current_date
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists cinema_details
-- +goose StatementEnd
