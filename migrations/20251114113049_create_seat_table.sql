-- +goose Up
-- +goose StatementBegin
create table seat
(
    id bigserial primary key,
    hall_id integer not null,
    row varchar(255) not null,
    number integer not null,

    foreign key (hall_id) references hall (id) on delete cascade
)


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists seat
-- +goose StatementEnd
