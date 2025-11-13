FROM golang:1.24-alpine AS build
WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 go build -o main cmd/http/main.go

FROM alpine:latest AS publish
WORKDIR /app

RUN apk add --no-cache ghostscript tzdata

ENV TZ=Asia/Almaty
RUN cp /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

COPY --from=build /app/main main
COPY --from=build /app/config.dev.yaml config.yaml

ENTRYPOINT ["/app/main"]

