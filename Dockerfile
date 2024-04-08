FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY . .

RUN go build -o /app/bin/go-ws-server .

FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/bin/go-ws-server ./go-ws-server
COPY --from=builder /app/public ./public

EXPOSE 8080

CMD ["./go-ws-server"]
