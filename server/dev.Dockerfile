FROM golang:1.24.5-alpine as builder

WORKDIR /app

COPY src/ .

RUN go mod download
RUN go install github.com/air-verse/air@latest;

EXPOSE 8080/udp
CMD ["air", "-c", ".air.toml"]
