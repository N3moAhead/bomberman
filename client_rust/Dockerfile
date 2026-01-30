FROM docker.io/rust:latest AS builder

WORKDIR /usr/src/client_rust

COPY Cargo.toml Cargo.lock ./

COPY src ./src

RUN cargo build --release

FROM docker.io/debian:bookworm-slim

RUN apt-get update && apt-get install -y openssl ca-certificates && rm -rf /var/lib/apt/lists/*

COPY --from=builder /usr/src/client_rust/target/release/client_rust .

CMD ["./client_rust"]
