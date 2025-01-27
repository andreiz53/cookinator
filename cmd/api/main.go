package main

import (
	"context"
	"log"

	database "github.com/andreiz53/cookinator/database/handlers"
	"github.com/andreiz53/cookinator/server"
	"github.com/jackc/pgx/v5"
)

const (
	dbSource      = "postgres://root:secret@localhost:5432/cookinator?sslmode=disable"
	serverAddress = ":8080"
)

func main() {
	conn, err := pgx.Connect(context.Background(), dbSource)
	if err != nil {
		log.Fatal("could not connect to db:", err)
	}

	store := database.NewStore(conn)
	server := server.NewServer(store)

	err = server.Run(serverAddress)
	if err != nil {
		log.Fatal("server stopped:", err)
	}
}
