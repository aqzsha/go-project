-- +goose Up
-- +goose StatementBegin
create table screening
(
    id bigserial primary key,
    cinema_id integer not null,
    hall_id integer not null,
    film_id integer not null,
    start_at timestamp not null,
    end_at timestamp not null,
    language varchar(255) not null,
    format varchar(255) not null,
    price float not null,

    foreign key (hall_id) references hall (id) on delete cascade,
    foreign key (film_id) references film (id) on delete cascade,
    foreign key (cinema_id) references cinema (id) on delete cascade
)


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists screening
-- +goose StatementEnd
