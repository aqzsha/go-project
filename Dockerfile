FROM golang:1.25

WORKDIR /app

# Копируем go.mod и go.sum для кэширования зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Устанавливаем Air для hot-reload
RUN go install github.com/air-verse/air@latest

# Копируем проект
COPY . .

# Запуск через Air
CMD ["air", "-c", ".air.toml"]
