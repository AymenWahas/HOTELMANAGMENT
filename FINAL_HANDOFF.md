# Final Handoff Package

## 1) Project Summary
Staywise is a hotel management system with:
- customer-facing booking page
- internal hotel management dashboard
- REST API for hotel operations
- PostgreSQL data storage
- Dockerized startup for delivery and testing

## 2) Live Access
- Public website: http://localhost:8089/
- Admin panel: http://localhost:8089/admin/
- Health check: http://localhost:8089/health

## 3) Default Demo Credentials
- Token: demo-token-12345
- This token is intended for local/demo delivery only.
- For real client deployment, replace it with a strong secret in the environment.

## 4) Runtime Setup
### Option A: Quick client delivery
```bash
cd /home/ayman-wahas/HotelManagement
docker compose up --build -d
```

### Option B: Local Go development
```bash
export DATABASE_URL='postgres://hotel:hotel@localhost:5433/hotel_management?sslmode=disable'
export API_TOKEN='replace-this-with-a-long-random-secret'
export APP_PORT=8089
export HOTEL_TIMEZONE='Asia/Riyadh'
go run ./cmd/api
```

## 5) Environment Settings
File used for local setup: `.env.example`

```env
APP_PORT=8089
DATABASE_URL=postgres://hotel:hotel@localhost:5433/hotel_management?sslmode=disable
API_TOKEN=replace-this-with-a-long-random-secret
HOTEL_TIMEZONE=Asia/Riyadh
```

## 6) Docker Services
- PostgreSQL on port 5433
- API on port 8089

## 7) Common Commands
```bash
# start all
cd /home/ayman-wahas/HotelManagement
docker compose up --build -d

# stop all
docker compose down

# restart
 docker compose restart

# check health
curl http://localhost:8089/health
```

## 8) Deployment Notes
- The API is built with Go 1.26.5.
- The database is PostgreSQL 16.
- The app is configured to use `Asia/Riyadh` timezone by default.
- For production client delivery, place the app behind a reverse proxy such as Nginx or a managed ingress and replace the demo token with a real secret.

## 9) Verification Status
Validated successfully:
- tests passed via `go test ./...`
- Docker startup succeeded
- health endpoint returned `{"status":"ok"}`

## 10) Recommended Client Handoff Message
"The hotel management system is ready for demo and client delivery. The application is packaged to run locally with Docker, includes the booking portal, admin dashboard and API layer, and is ready to be reviewed on http://localhost:8089."
