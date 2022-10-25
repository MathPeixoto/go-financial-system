postgres:
	docker run --name bank -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=postgres -d postgres:14.3

createdb:
	docker exec -it bank createdb --username=root --owner=root bank

dropdb:
	docker exec -it bank dropdb bank

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