test:
	go test

test-coverage:
	go test -coverprofile=c.out
	go tool cover -func=c.out

run:
	go run .

run-air:
	air

docker-build-db:
	docker build --tag gravityapi_db ./db

docker-build-app:
	docker build --tag gravityapi .

