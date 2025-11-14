-- +goose Up
-- +goose StatementBegin
CREATE TABLE film_details
(
    id BIGSERIAL PRIMARY KEY,
    duration text,
    premier timestamp not null,
    production varchar(255),
    director varchar(255),
    rate float,
    age_limit integer
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP table if exists film_details
-- +goose StatementEnd

