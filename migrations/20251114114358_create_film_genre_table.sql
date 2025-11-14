-- +goose Up
-- +goose StatementBegin
CREATE TABLE film_genre
(
    id BIGSERIAL PRIMARY KEY,
    film_id INTEGER NOT NULL,
    genre_id INTEGER NOT NULL,
    FOREIGN KEY (genre_id) REFERENCES genre (id) ON DELETE CASCADE,
    FOREIGN KEY (film_id) REFERENCES film (id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS film_genre;
-- +goose StatementEnd
