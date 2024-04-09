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
ws://localhost:8080/ws
```

Message format:

```json
{
  "event": "message",
  "data": "Hello WebSocket"
}
```

## Running

```bash
go run .
```

Open:

```text
http://localhost:8080
```

The page includes a small browser client for sending messages to the WebSocket
endpoint and reading the echoed response.

Open the page in two browser tabs to see messages broadcast between connected
clients.

## Docker

```bash
docker compose up --build
```

Open:

```text
http://localhost:8080
```

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
ws://localhost:8080/ws
```

Формат сообщения:

```json
{
  "event": "message",
  "data": "Hello WebSocket"
}
```

## Запуск

```bash
go run .
```

Открой:

```text
http://localhost:8080
```

Страница содержит небольшой браузерный клиент для отправки сообщений в
WebSocket endpoint и просмотра echo-ответа.

Открой страницу в двух вкладках браузера, чтобы увидеть рассылку сообщений
между подключенными клиентами.

## Docker

```bash
docker compose up --build
```

Открой:

```text
http://localhost:8080
```
