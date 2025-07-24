# Local dev

To start the app, run:
```bash
docker-compose up --build
```

Check if the server has started with:
```bash
curl localhost:8080/api/health
```
The result should be `{"status":"ok"}`.

# Possible improvements

- Implement pagination for `GET /api/v1/posts` endpoint
- Add customized logger (such as [zaplog](https://github.com/uber-go/zap))
- Implement various middlewares for rate , auth
- Optimize docker image
- Extract config from environmental variables (such as port, credentials, etc.) or config files
- Handle errors on server startup
- Add swagger docs auto-generated from the comments
