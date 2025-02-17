package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"

	database "github.com/andreiz53/cookinator/database/handlers"
	"github.com/andreiz53/cookinator/server"
	"github.com/andreiz53/cookinator/util"
)

const ()

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("could not load env:", err)
	}
	conn, err := pgx.Connect(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("could not connect to db:", err)
	}

	store := database.NewStore(conn)
	server, err := server.NewServer(config, store)
	if err != nil {
		log.Fatal("could not start server:", err)
	}

	err = server.Run(config.ServerAddress)
	if err != nil {
		log.Fatal("server stopped:", err)
	}
}
