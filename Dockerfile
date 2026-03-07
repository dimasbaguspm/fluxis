FROM golang:1.25-alpine AS dev

WORKDIR /app

RUN apk add --no-cache make git

# Install dev tools — cached independently of application code
RUN go install github.com/air-verse/air@latest \
    && go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest \
    && go install github.com/swaggo/swag/cmd/swag@latest

# Download dependencies — only re-runs when go.mod/go.sum change
COPY go.mod go.sum ./
RUN go mod download

COPY Makefile ./

CMD ["air"]

FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
RUN go build -o ./tmp/main ./cmd/fluxis/main.go

FROM alpine:3.21

WORKDIR /app

COPY --from=builder /app/tmp/main ./bin/main
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

CMD ["./bin/main"]
