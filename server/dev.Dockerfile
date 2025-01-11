FROM golang:1.23-alpine as builder

WORKDIR /app

COPY src/ .

RUN go install github.com/air-verse/air@latest;

CMD ["air", "-c", ".air.toml"]
