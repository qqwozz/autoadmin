FROM golang:1.22-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /build
COPY app/go.mod app/go.sum ./
RUN go mod download
COPY app/ .

RUN CGO_ENABLED=1 go build -o /autoadmin-api ./cmd/api
RUN CGO_ENABLED=1 go build -o /autoadmin-bot ./cmd/bot

FROM alpine:3.20

RUN apk add --no-cache sqlite bash

WORKDIR /app
COPY --from=builder /autoadmin-api /app/autoadmin-api
COPY --from=builder /autoadmin-bot /app/autoadmin-bot
COPY app/data.sql /app/data.sql

EXPOSE 8080

CMD ["sh", "-c", "if [ ! -s /data/database.sqlite ]; then sqlite3 /data/database.sqlite < /app/data.sql && echo 'SQLite DB initialized'; else echo 'SQLite DB already exists'; fi && /app/autoadmin-api"]
