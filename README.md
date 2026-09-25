# Hanon Systems API

Go/Chi HTTP API exposing production, schedule, scrap, heartbeat and graph data
stored in SQL Server. Interactive API documentation is served with Redoc at
`/docs`.

## Features

- [x] Chi router with sub-routers per domain
- [x] SQL Server access via `go-mssqldb` (single connection pool)
- [x] OpenAPI/Swagger documentation rendered with Redoc
- [x] Static file serving (`/static`)
- [x] Daily production/NOK graph endpoints driven by `prodFAssy_config.json`

## Requirements

- Go 1.21+
- Access to the SQL Server instance holding the production databases

## Configuration

Copy `.example.env` to `.env` and fill in the database credentials:

```
DB_HOST=
DB_PORT=
DB_USER=
DB_PASS=
DB_NAME=
```

Optional:

- `PROD_CONFIG_PATH` — absolute path to `prodFAssy_config.json`. If unset, the
  API looks next to the running binary and then in the project root.
- `LOG_REQUESTS=true` — enable per-request access logging (off by default).

## Run

```bash
go run .
```

The server listens on `:4000`. Open <http://localhost:4000/docs> for the API
documentation.

## Endpoints

| Method | Path                                     | Description                        |
| ------ | ---------------------------------------- | ---------------------------------- |
| GET    | `/production/now/{line_id}`              | Current shift production           |
| GET    | `/production/day/{line_id}`              | Production for the current day     |
| GET    | `/production/yesterday/{line_id}`        | Production for the previous day    |
| GET    | `/schedule/shift/{line_id}`              | Current shift schedule             |
| GET    | `/schedule/now/{line_id}`                | Schedule progress within the shift |
| GET    | `/schedule/day/{line_id}`                | Schedule progress for the day      |
| GET    | `/schedule/yesterday/{line_id}`          | Schedule for the previous day      |
| GET    | `/scrap/insertTicket`                    | Insert a scrap ticket              |
| GET    | `/scrap/updatePerson`                    | Update the person on a ticket      |
| GET    | `/scrap/updateRequester`                 | Update the requester on a ticket   |
| GET    | `/com/heartbeat`                         | Record an application heartbeat    |
| GET    | `/graph/api/daily/{line_id}`             | Hourly production by shift         |
| GET    | `/graph/api/dailynok/{line_id}`          | Hourly NOK by shift                |

The graph endpoints accept `startDate` and `endDate` query parameters in
`YYYY-MM-DD` format.

## Build

```bash
go build -o go-api.exe .
```

## Test

```bash
go test ./...
```

## Docker

```bash
docker build -f dockerfile -t hanon-api .
docker run --env-file .env -p 4000:4000 hanon-api
```

## License

MIT
