-- +goose Up
-- +goose StatementBegin
create table review
(
    id bigserial primary key,
    film_id integer not null,
    user_id integer not null,
    title varchar(255) not null,
    body text not null,
    rating float not null,
    created_at timestamp default current_date,

    foreign key (film_id) references film (id) on delete cascade,
    foreign key (user_id) references users (id) on delete cascade
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists review
-- +goose StatementEnd
