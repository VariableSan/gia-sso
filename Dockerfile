FROM golang:1.24 AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=1 go build -o sso ./cmd/sso
RUN CGO_ENABLED=1 go build -o migrator ./cmd/migrator

FROM debian:bookworm-slim
WORKDIR /app
COPY --from=builder /app/sso ./sso
COPY --from=builder /app/migrator ./migrator
COPY ./config/dev.yaml ./config/dev.yaml

ENTRYPOINT ["./sso", "--config=./config/dev.yaml"]

