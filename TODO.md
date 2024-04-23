# TODO

Short project checklist for the next few days.

## Done

- [x] Started the base Go HTTP service.
- [x] Added Docker setup for local development.
- [x] Added WebSocket channels.
- [x] Published WebSocket events to Redis.
- [x] Broadcast Redis PubSub events to WebSocket clients.
- [x] Added lead modal state events.
- [x] Added tests for lead modal events.
- [x] Handled graceful server shutdown.
- [x] Moved HTTP server setup out of main.

## Next

- [ ] Add a Redis readiness endpoint.
- [ ] Add Redis connection timeouts.
- [ ] Add WebSocket ping/pong handling.
- [ ] Add tests for message normalization.
- [ ] Document runtime checks and local commands.

## Later

- [ ] Add more production notes after the next few changes.
- [ ] Add GitHub Actions when the public repository is ready.
- [ ] Review WebSocket origin checks before production use.
