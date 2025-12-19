## Run

Dockerize

> docker compose up --build -d

Make migrations

> go run cmd/main.go -cmd up

Reset migrations

> go run cmd/main.go -cmd reset
