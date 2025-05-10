
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install git and other tools required by go mod
RUN apk add --no-cache git

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o score-engine ./cmd/server

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/score-engine .

COPY database.db .

EXPOSE 50051

CMD ["./score-engine"]
