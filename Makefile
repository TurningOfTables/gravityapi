# Runs standard go tests
test:
	set -o pipefail && go test -json | tparse -all

test-verbose:
	go test -v

# Calculates test coverage and displays breakdown by file/function
test-coverage:
	go test -coverprofile=c.out
	go tool cover -func=c.out

# Standard local run
local-run:
	go run .

# Build binary
build:
	go build -o

# Standard local run with air to allow hot reloading
local-air-run:
	air

# Builds and runs app and db in docker
docker-build-run:
	docker compose up --build

# Builds only the db
docker-build-db:
	docker build --tag gravityapi_db ./db

# Builds only the app
docker-build-app:
	docker build --tag gravityapi .