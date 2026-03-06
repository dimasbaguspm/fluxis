FROM golang:1.25-alpine AS dev

WORKDIR /app

RUN go install github.com/air-verse/air@latest

COPY go.mod ./
RUN go mod download

CMD ["air"]

FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
RUN go build -o ./tmp/main ./cmd/fluxis/main.go

FROM alpine:3.21

WORKDIR /app

COPY --from=builder /app/tmp/main ./main

EXPOSE 8080

CMD ["./main"]
