# Runs standard go tests
test:
	set -o pipefail && go test -json | tparse -all

test-verbose:
	go test -v

# Calculates test coverage and displays breakdown by file/function
test-coverage:
	go test . -coverprofile=c.out
	go tool cover -func=c.out

# Calculates test coverage and displays breakdown in the browser
test-coverage-browser:
	go test . -coverprofile=c.out
	go tool cover -html=c.out

# Standard local run
run-local:
	go run .

# Build binary
build:
	go build -o

# Standard local run with air to allow hot reloading
run-air:
	air

# Builds and runs app and db in docker
build-run-docker:
	docker-compose build
	docker-compose up

test-docker:
	docker-compose build
	docker-compose up -d
	docker exec -it gravityapi-web-1 go test -v .


# Builds only the db
build-docker-db:
	docker build --tag gravityapi_db ./db

# Builds only the app
build-docker-app:
	docker build --tag gravityapi .