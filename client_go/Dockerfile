FROM docker.io/golang:1.25-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /app/client_go ./cmd/client_go/main.go

FROM docker.io/alpine:latest
WORKDIR /root/

COPY --from=builder /app/client_go .

CMD ["./client_go"]
