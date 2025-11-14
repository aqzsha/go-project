-- +goose Up
-- +goose StatementBegin
create table hall
(
    id bigserial primary key,
    cinema_id integer,
    name varchar(255),
    seats integer,

    foreign key (cinema_id) references cinema (id) on delete cascade
)


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists hall
-- +goose StatementEnd
