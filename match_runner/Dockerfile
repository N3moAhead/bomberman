FROM docker.io/golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o /app/match-runner ./cmd/match_runner

FROM alpine:latest

RUN apk update && \
  apk add --no-cache podman

COPY --from=builder /app/match-runner /usr/local/bin/match-runner

ENTRYPOINT ["/usr/local/bin/match-runner"]
