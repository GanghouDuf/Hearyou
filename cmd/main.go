package main

import (
	"context"
	"log"
	"project_chat/internel/auth"
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
	auth.SetSecret(config.JwtSecret)
	userRepo := storage.NewUserReporitory(conn)
	messageRepo := storage.NewMessageRepository(conn)
	defer conn.Close(ctx)
	hub := ws.NewHub()

	go hub.Run()
	srv := server.NewServer(config.Addr, hub, userRepo, messageRepo)
	log.Fatal(srv.Start())

}
