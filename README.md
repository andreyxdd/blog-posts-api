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
