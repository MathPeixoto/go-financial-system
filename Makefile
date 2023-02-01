DB_URL=postgresql://root:postgres@localhost:5432/bank?sslmode=disable


network:
	docker network create bank-network

postgres:
	docker run --name bankdatabase --network bank-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=postgres -d postgres:14.3

createdb:
	docker exec -it bankdatabase createdb --username=root --owner=root bank

dropdb:
	docker exec -it bankdatabase dropdb bank

migrate:
	curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz \
 	&&  sudo mv migrate /usr/bin/migrate \
	&& which migrate

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migrateup1:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

migratedown1:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1

db_docs:
	dbdocs build doc/db.dbml

db_schema:
	dbml2sql --postgres -o doc/schema.sql doc/db.dbml

installSqlc:
	go install github.com/kyleconroy/sqlc/cmd/sqlc@latest

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

run:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/MathPeixoto/go-financial-system/db/sqlc Store

proto:
	rm -f pb/*.go
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
        --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
        proto/*.proto

evans:
	evans --host localhost --port 9090 -r repl

.PHONY: network postgres createdb dropdb migrateup migrateup1 migratedown migratedown1 db_docs db_schema sqlc test run gin mock proto evans