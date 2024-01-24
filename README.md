# gravity api

## Description

An API Written in Go using the [Go Fiber Framework](https://gofiber.io/), and using the PostgreSQL gravity database at [Database Star's SQL Code Repository](https://github.com/bbrumm/databasestar).

## Run - Local

### Prerequisites

1. [PostgreSQL](https://www.postgresql.org/download/)
2. Gravity database running, as described at [Database Star's SQL Code Repository](https://github.com/bbrumm/databasestar)
3. Clone the repository locally
4. Set `GRAVITY_API_APP_HOST` to choose where the web app is hosted 
5. Set `GRAVITY_API_DB_CONNECTION_STRING` to indicate the connection string for the PostgresQL db
5. `make local-run` OR `make build` and run the resulting `gravityapi` binary
6. Navigate to the URL you set in step 4 (`GRAVITY_API_APP_HOST`)

## Run - Docker

### Prerequisites

1. [Docker](https://docs.docker.com/engine/install/)
2. Run `make docker-build-run`
3. Navigate to the URL in `.env.docker` for `GRAVITY_API_APP_HOST`. Defaults to `http://127.0.0.1:3000`

Using Docker compose this should create the database and populate it (saving you doing steps 1-2 that you have to do when running locally), then run both the database and app containers.

## Attributions

* [Go Fiber Framework](https://gofiber.io/)
* [Database Star's SQL Code Repository](https://github.com/bbrumm/databasestar)
* [godotenv](https://github.com/joho/godotenv)
* [pgx](https://github.com/jackc/pgx)