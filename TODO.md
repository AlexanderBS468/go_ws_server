# TODO

Short project checklist for the next few days.

## Done

Some runtime checks were added after a production smoke test.

- [x] Started the base Go HTTP service.
- [x] Added Docker setup for local development.
- [x] Added WebSocket channels.
- [x] Published WebSocket events to Redis.
- [x] Broadcast Redis PubSub events to WebSocket clients.
- [x] Added lead modal state events.
- [x] Added tests for lead modal events.
- [x] Handled graceful server shutdown.
- [x] Moved HTTP server setup out of main.
- [x] Added a Redis readiness endpoint.
- [x] Added Redis connection timeouts.
- [x] Added WebSocket ping/pong handling after stale connections showed up.
- [x] Added tests for message normalization.
- [x] Documented runtime checks and local commands.

## Later

- [ ] Add more production notes after the next few changes.
- [ ] Add GitHub Actions when the public repository is ready.
- [ ] Review WebSocket origin checks before production use.
