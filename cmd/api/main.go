package main

import (
	"context"
	"log"

	database "github.com/andreiz53/cookinator/database/handlers"
	"github.com/andreiz53/cookinator/server"
	"github.com/andreiz53/cookinator/util"
	"github.com/jackc/pgx/v5"
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
	server := server.NewServer(store)

	err = server.Run(config.ServerAddress)
	if err != nil {
		log.Fatal("server stopped:", err)
	}
}
