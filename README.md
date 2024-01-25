# WORK IN PROGRESS - gravity api

## Warning

This is an in progress project, so functionality will be in constant flux!

## Description

A work in progress API written in Go using the [Go Fiber Framework](https://gofiber.io/), and using the PostgreSQL gravity database at [Database Star's SQL Code Repository](https://github.com/bbrumm/databasestar).

My goal is to put my learning to work and write Go more efficiently than my previous API projects.

## Features

* Both database and app are dockerised, with a working docker compose file
* Able to switch between docker and local runs using only a command line flag
* Makes use of interfaces to provide a more generic search function, reducing repetition of code for searching different models - see `search.go` for how this is implemented
* Makes use of test sets to avoid declaring dozens of repeated test functions, one for each route
* Uses a larger data set than before, requiring more thought on response size and handling.

## Incoming Features

* Paging and offset support, as some return payloads are > 10,000 lines!
* POST support to add new data
* Expanded use of related tables to provide fuller responses to queries, instead of just IDs for some fields
* Once the endpoints have stabilised somewhat, documentation in the form of OpenAPI specs and HTML docs likely generated automatically
* Improved error response - currently just text, but will be a JSON response with an error code and message for consistency with other JSON responses
* Potentially implement useful Fiber middleware including caching, monitoring, auth for sensitive endpoints and rate limiting.


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