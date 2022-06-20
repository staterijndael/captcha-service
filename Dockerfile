FROM golang:latest

COPY ./ /app

WORKDIR /app

RUN go mod download

ENTRYPOINT go run cmd/main.go