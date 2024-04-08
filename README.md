# go-ws-server

Realtime event logger for a production system.

## Contents

- [English](#english)
- [Русский](#russian)

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

## Running

```bash
go run .
```

Open:

```text
http://localhost:8080
```

## Russian

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

## Запуск

```bash
go run .
```

Открой:

```text
http://localhost:8080
```
