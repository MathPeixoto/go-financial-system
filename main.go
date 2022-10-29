package main

import (
	"database/sql"
	"github.com/MathPeixoto/go-financial-system/api"
	db "github.com/MathPeixoto/go-financial-system/db/sqlc"
	"github.com/MathPeixoto/go-financial-system/util"
	_ "github.com/golang/mock/mockgen/model"
	_ "github.com/lib/pq"
	"log"
)

func main() {

	config, err := util.LoadConfig("app.env")
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}

	var conn *sql.DB
	conn, err = sql.Open(config.DatabaseDriver, config.DatabaseSource)
	if err != nil {
		log.Fatalln("cannot connect to database:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatalln("cannot start server:", err)
		return
	}
}
