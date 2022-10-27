postgres:
	docker run --name bank -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=postgres -d postgres:14.3

createdb:
	docker exec -it bank createdb --username=root --owner=root bank

dropdb:
	docker exec -it bank dropdb bank

migrate:
	curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz \
 	&&  sudo mv migrate /usr/bin/migrate \
	&& which migrate

migrateup:
	migrate -path db/migration -database "postgresql://root:postgres@localhost:5432/bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:postgres@localhost:5432/bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

run:
	go run .

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test run