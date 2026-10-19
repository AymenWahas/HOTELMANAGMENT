# Hotel Management API

**Languages:** [العربية](README.md) | [English](README.en.md)

Go REST API for the hotel’s day-to-day operations, backed by PostgreSQL. The code follows the requested dependency flow:

```text
HTTP handler -> use case -> repository interface -> PostgreSQL repository
```

## Operations included

- Room inventory and room status
- Guest profiles
- Bookings/reservations, availability conflict checks, check-in, check-out, and cancellation
- Reception arrival and departure queues, with room and guest details
- Housekeeping tasks and room turn-over
- Guest folios and payments
- Employees and departments
- Maintenance requests
- Inventory items and stock movements
- Operations dashboard, occupancy, and revenue reports
- Audit events for changes

PostgreSQL keeps each operational area in its own relational table. Foreign keys, check constraints, and a date-range exclusion constraint protect the data. Inventory movements and reservation lifecycle changes run in database transactions.

## Run locally

Requirements: Go 1.26.5+ and Docker Compose.

1. Start PostgreSQL: `docker compose up -d postgres`
2. Configure the service:

   ```sh
   export DATABASE_URL='postgres://hotel:hotel@localhost:5433/hotel_management?sslmode=disable'
   export API_TOKEN='replace-this-with-a-long-random-secret'
   export APP_PORT=8081
   export HOTEL_TIMEZONE='Asia/Riyadh'
   ```

3. Start the API and browser interface from the project root: `go run ./cmd/api`

The API applies `migrations/001_initial_schema.sql` on startup and records the applied migration. Open **http://localhost:8081** in your browser and enter the configured `API_TOKEN` to use the interface. Health check: `GET http://localhost:8081/health`.

The embedded browser interface includes overview, reception, bookings, rooms, guests, and payments screens. Its forms call the corresponding `/api/v1` endpoints on the same server.

The screen-to-API connections are:

| Browser screen | API calls |
| --- | --- |
| Overview | `GET /dashboard`, `GET /reports/occupancy`, `GET /reports/revenue`, `GET /reception/arrivals`, `GET /rooms` |
| Reception | `GET /reception/arrivals`, `GET /reception/departures`, `POST /reception/{id}/check-in`, `POST /reception/{id}/check-out`, `POST /bookings/{id}/cancel` |
| Bookings | `GET /bookings`, `POST /bookings`, `POST /bookings/{id}/check-in`, `POST /bookings/{id}/check-out`, `POST /bookings/{id}/cancel` |
| Rooms | `GET /rooms`, `POST /rooms`, `PATCH /rooms/{id}` |
| Guests | `GET /guests`, `POST /guests` |
| Payments | `GET /payments`, `GET /folios`, `POST /folios`, `POST /payments`, `PATCH /payments/{id}` |

All `/api/v1` routes require `Authorization: Bearer <API_TOKEN>`. `HOTEL_TIMEZONE` controls daily dashboard and revenue boundaries and defaults to `Asia/Riyadh`. The token is a shared service credential for this first backend version; staff login and role-based permissions can be added as a later layer.

## API

All request and response bodies use JSON. Resource records return `{ "id", "kind", "data", "created_at", "updated_at" }`.

| Resource | Routes |
| --- | --- |
| Rooms | `/api/v1/rooms` |
| Guests | `/api/v1/guests` |
| Bookings / reservations | `/api/v1/bookings` or `/api/v1/reservations` |
| Reception | `/api/v1/reception/arrivals`, `/api/v1/reception/departures` |
| Housekeeping | `/api/v1/housekeeping` |
| Folios | `/api/v1/folios` |
| Payments | `/api/v1/payments` |
| Employees | `/api/v1/employees` |
| Maintenance | `/api/v1/maintenance` |
| Inventory items | `/api/v1/inventory/items` |
| Inventory movements | `/api/v1/inventory/movements` |

Resources support `GET` list, `POST` create, `GET /{id}`, and `PATCH /{id}`. Most also support `DELETE /{id}` (soft archive); reservations, payments, and inventory movements stay immutable in the operational history. Lists accept `limit` (default 50, maximum 200) and `offset`. Reservation status changes use the lifecycle actions below.

Reservation actions:

- `POST /api/v1/reservations/{id}/check-in`
- `POST /api/v1/reservations/{id}/check-out`
- `POST /api/v1/reservations/{id}/cancel`

The same check-in, check-out, and cancellation actions are available under `/api/v1/bookings/{id}/...` or `/api/v1/reception/{id}/...`. Reception queues return today's arrivals and departures by default; pass `?date=YYYY-MM-DD` to select another day.

Reports:

- `GET /api/v1/dashboard`
- `GET /api/v1/reports/occupancy`
- `GET /api/v1/reports/revenue?from=2026-10-01&to=2026-10-31`

### Example requests

```sh
curl -H "Authorization: Bearer $API_TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"number":"101","room_type":"standard","floor":1,"rate":350}' \
  http://localhost:8081/api/v1/rooms
```

```json
{
  "room_id": "ROOM_UUID",
  "guest_id": "GUEST_UUID",
  "check_in": "2026-10-21",
  "check_out": "2026-10-24",
  "adults": 2,
  "rate": 350,
  "source": "direct"
}
```

Other create payloads use the database field names for their resource. For example, guests use `first_name` and `last_name`; a payment uses `folio_id`, `amount`, and `method`; inventory movements use `item_id`, `quantity`, and `movement_type` (`in`, `out`, `restock`, or `use`).

## Project layout

- `internal/domain`: domain records and report types
- `internal/handler`: HTTP routing, JSON, and authentication
- `internal/usecase`: input validation and operation rules
- `internal/repository`: persistence interface used by the use cases
- `internal/postgres`: PostgreSQL queries, transactions, and audit writes
- `internal/webui`: embedded browser interface for hotel staff
- `migrations`: PostgreSQL schema
- `cmd/api`: dependency wiring and process lifecycle
