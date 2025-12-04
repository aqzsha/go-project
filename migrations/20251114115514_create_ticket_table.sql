-- +goose Up
-- +goose StatementBegin
create table ticket
(
    id bigserial primary key,
    user_id integer not null,
    screening_id integer not null,
    seat_id integer not null,
    price float not null,
    status varchar(255) not null,
    created_at timestamp default current_date,

    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (screening_id) REFERENCES screening (id) ON DELETE CASCADE,
    FOREIGN KEY (seat_id) REFERENCES seat (id) ON DELETE CASCADE
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists ticket
-- +goose StatementEnd
