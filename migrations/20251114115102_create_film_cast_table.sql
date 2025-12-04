-- +goose Up
-- +goose StatementBegin
create table film_cast
(
    id bigserial primary key,
    film_id integer not null,
    person_id integer not null,
    role_name varchar(255) not null,

    foreign key (person_id) references person (id) on delete cascade,
    foreign key (film_id) references film (id) on delete cascade
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists film_cast
-- +goose StatementEnd
