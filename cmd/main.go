package main

import (
	"context"
	"log"
	"project_chat/internel/config"
	"project_chat/internel/server"
	"project_chat/internel/storage"
	"project_chat/internel/ws"
)

func main() {
	ctx := context.Background()
	config := config.Load()
	conn, err := storage.Connect(config.DatabaseUrl, ctx)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer conn.Close(ctx)
	hub := ws.NewHub()
	go hub.Run()
	srv := server.NewServer(config.Addr, hub)
	log.Fatal(srv.Start())

}
