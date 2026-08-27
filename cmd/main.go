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
	defer conn.Close(ctx)

	userRepo := storage.NewUserReporitory(conn)
	messageRepo := storage.NewMessageRepository(conn)
	roomRepo := storage.NewRoomRepository(conn)

	roomManager := ws.NewRoomManager()

	srv := server.NewServer(config.Addr, roomManager, userRepo, messageRepo, roomRepo)

	log.Fatal(srv.Start())
}
