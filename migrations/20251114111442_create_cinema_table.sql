-- +goose Up
-- +goose StatementBegin
CREATE TABLE cinema
(
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) not null,
    details_id integer not null,

    foreign key (details_id) references cinema_details (id) on delete cascade
)


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists cinema
-- +goose StatementEnd
