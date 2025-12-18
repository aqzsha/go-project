-- +goose Up
-- +goose StatementBegin
CREATE TABLE ticket
(
    id                BIGSERIAL PRIMARY KEY,
    seat_screening_id BIGINT       NOT NULL UNIQUE,
    qr_code           VARCHAR(255) NOT NULL,
    status            VARCHAR(50)  NOT NULL,
    scan_at           TIMESTAMP,
    created_at        TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at        TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_ticket_seat_screening
        FOREIGN KEY (seat_screening_id)
            REFERENCES seat_screening (id)
            ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS ticket;
-- +goose StatementEnd
