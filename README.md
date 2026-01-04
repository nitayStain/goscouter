# goscouter

## Run with Docker
```
docker compose build
docker compose up
```

Services:
- Backend at http://localhost:8080 returns a JSON hello message.
- Frontend at http://localhost:3000 shows a basic welcome page.
- PostgreSQL exposed on 5432 with user/password/db `postgres/postgres/goscouter`.
