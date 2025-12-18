-- +goose Up
-- +goose StatementBegin
CREATE TABLE seat_screening
(
    id            BIGSERIAL PRIMARY KEY,
    hall_id       BIGINT NOT NULL,
    screening_id  BIGINT NOT NULL,
    row           BIGINT NOT NULL,
    number        BIGINT NOT NULL,
    status        VARCHAR(50) NOT NULL,
    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    UNIQUE (screening_id, row, number)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS seat_screening;
-- +goose StatementEnd
