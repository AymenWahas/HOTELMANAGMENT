# Quick Start Guide

## Run the project
```bash
cd /home/ayman-wahas/HotelManagement
docker compose up --build -d
```

## Open the app
- Public site: http://localhost:8089/
- Admin site: http://localhost:8089/admin/
- Health check: http://localhost:8089/health

## Environment variables
```env
APP_PORT=8089
DATABASE_URL=postgres://hotel:hotel@localhost:5433/hotel_management?sslmode=disable
API_TOKEN=replace-this-with-a-long-random-secret
HOTEL_TIMEZONE=Asia/Riyadh
```

## Default demo token
```text
demo-token-12345
```

## Useful commands
```bash
docker compose down
docker compose restart
docker compose logs -f api
```

## If port 8089 is busy
Stop the local process or old container using the port, then restart:
```bash
lsof -i :8089
kill <PID>
docker compose up -d
```

## Health verification
```bash
curl http://localhost:8089/health
```

Expected output:
```json
{"status":"ok"}
```
