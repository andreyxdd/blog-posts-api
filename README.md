# Local dev

To start the app with hot reload on code updates for local dev, run:
```bash
docker-compose -f docker-compose.dev.yaml up --build
```

Check if the server has started with:
```bash
curl localhost:8080/api/health
```
The result should be `{"status":"ok"}`.

Access swagger API specifications at:
```bash
curl localhost:8080/api/docs/index.html
```

# Possible improvements

- Implement pagination for `GET /api/v1/posts` endpoint
- Add customized logger (such as [zaplog](https://github.com/uber-go/zap))
- Implement various middlewares for rate limiting, auth, etc.
- Optimize docker image
- Extract config from environmental variables (such as port, credentials, etc.) or config files
- Handle errors on server startup
