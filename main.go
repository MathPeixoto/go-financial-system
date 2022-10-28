package main

import (
	"database/sql"
	"github.com/MathPeixoto/go-financial-system/api"
	db "github.com/MathPeixoto/go-financial-system/db/sqlc"
	_ "github.com/lib/pq"
	"log"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:postgres@localhost:5432/bank?sslmode=disable"
	address  = "localhost:8080"
)

func main() {

	var conn *sql.DB
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatalln("cannot connect to database:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)
	err = server.Start(address)
	if err != nil {
		log.Fatalln("cannot start server:", err)
		return
	}
}
