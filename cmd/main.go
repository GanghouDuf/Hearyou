package main

import (
	"log"
	"project_chat/internel/server"
	"project_chat/internel/ws"
)

func main() {
	hub := ws.NewHub()
	go hub.Run()
	srv := server.NewServer(":8080", hub)
	log.Fatal(srv.Start())

}
