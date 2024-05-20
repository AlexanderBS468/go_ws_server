# go-ws-server

Realtime event logger for a production system.

## Contents

- [English](#english)
- [Русский](#russian)

<a id="english"></a>

## English

`go-ws-server` is a small Go service for realtime logging and delivery of
production system events over WebSocket.

The service accepts browser WebSocket connections, receives event messages, and
broadcasts them to connected clients. It is intended for operational events that
need to appear in the interface immediately, without polling the main backend.

Planned runtime flow:

1. Clients connect to a WebSocket endpoint.
2. The service receives JSON event messages.
3. Events are published to a channel.
4. Connected clients subscribed to that channel receive the update.

Current WebSocket endpoint:

```text
ws://localhost:8080/ws/{channel}
```

Message format:

```json
{
  "event": "message",
  "data": "Hello WebSocket"
}
```

Supported state events:

```text
login
lead-details-open
lead-details-close
lead-page-out
```

Incoming WebSocket events are published to Redis using the WebSocket channel
name as the Redis channel. The app also subscribes to Redis and broadcasts
received events to browser clients in the matching channel.

## Running

```bash
go run .
```

Configuration defaults:

```text
SERVER_PORT=8080
REDIS_ADDR=localhost:6379
WS_ALLOWED_ORIGINS=http://localhost:8080,http://127.0.0.1:8080
```

Open:

```text
http://localhost:8080
```

The page includes a small browser client for sending messages to the WebSocket
endpoint and reading the echoed response.

Open the page in two browser tabs to see messages broadcast between connected
clients in the same channel.

## Runtime checks

```text
GET /health
```

Checks that the HTTP process is running.

```text
GET /ready
```

Checks that the service can reach Redis with `PING`. It returns
`503 Service Unavailable` when Redis is unavailable.

The server also handles `SIGINT` and `SIGTERM` for graceful shutdown, and keeps
WebSocket connections alive with ping/pong checks.

## Docker

```bash
docker compose up --build
```

Open:

```text
http://localhost:8080
```

The compose setup also starts a local Redis container on `localhost:6379`.

Useful local commands:

```bash
docker compose build
docker compose up
docker compose down
docker run --rm -v "$PWD":/app -w /app golang:1.20-alpine go test ./...
```

## Production notes

- Redis is a required runtime dependency because messages are delivered through
  Pub/Sub.
- Use `/health` for process liveness and `/ready` for Redis readiness checks.
- Stop the container with `SIGTERM` so the HTTP server can shut down
  gracefully.
- Set `WS_ALLOWED_ORIGINS` before exposing the service outside a trusted
  network.

<a id="russian"></a>

## Русский

`go-ws-server` - небольшой Go-сервис для realtime-логирования и доставки
событий боевой системы через WebSocket.

Сервис принимает WebSocket-подключения из браузера, получает сообщения о
событиях и рассылает их подключенным клиентам. Он нужен для операционных
событий, которые должны появляться в интерфейсе сразу, без постоянного опроса
основного backend.

Планируемая схема работы:

1. Клиенты подключаются к WebSocket endpoint.
2. Сервис принимает JSON-сообщения с событиями.
3. События публикуются в канал.
4. Подключенные клиенты, подписанные на этот канал, получают обновление.

Текущий WebSocket endpoint:

```text
ws://localhost:8080/ws/{channel}
```

Формат сообщения:

```json
{
  "event": "message",
  "data": "Hello WebSocket"
}
```

Поддерживаемые события состояния:

```text
login
lead-details-open
lead-details-close
lead-page-out
```

Входящие WebSocket-события публикуются в Redis, где имя WebSocket-канала
используется как Redis-канал. Приложение также подписывается на Redis и
рассылает полученные события браузерным клиентам в соответствующем канале.

## Запуск

```bash
go run .
```

Настройки по умолчанию:

```text
SERVER_PORT=8080
REDIS_ADDR=localhost:6379
WS_ALLOWED_ORIGINS=http://localhost:8080,http://127.0.0.1:8080
```

Открой:

```text
http://localhost:8080
```

Страница содержит небольшой браузерный клиент для отправки сообщений в
WebSocket endpoint и просмотра echo-ответа.

Открой страницу в двух вкладках браузера, чтобы увидеть рассылку сообщений
между подключенными клиентами в одном канале.

## Runtime-проверки

```text
GET /health
```

Проверяет, что HTTP-процесс запущен.

```text
GET /ready
```

Проверяет доступность Redis через `PING`. Если Redis недоступен, возвращает
`503 Service Unavailable`.

Сервер также обрабатывает `SIGINT` и `SIGTERM` для graceful shutdown и
поддерживает WebSocket-соединения через ping/pong-проверки.

## Docker

```bash
docker compose up --build
```

Открой:

```text
http://localhost:8080
```

Compose-конфигурация также запускает локальный Redis на `localhost:6379`.

Полезные локальные команды:

```bash
docker compose build
docker compose up
docker compose down
docker run --rm -v "$PWD":/app -w /app golang:1.20-alpine go test ./...
```

## Production-заметки

- Redis является обязательной runtime-зависимостью, потому что сообщения
  доставляются через Pub/Sub.
- Используй `/health` для проверки процесса и `/ready` для проверки готовности
  Redis.
- Останавливай контейнер через `SIGTERM`, чтобы HTTP-сервер завершался
  аккуратно.
- Перед доступом вне доверенной сети настрой `WS_ALLOWED_ORIGINS`.
