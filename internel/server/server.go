package server

import (
	"context"
	"log"
	"net/http"
	"project_chat/internel/ws"
	"time"

	"github.com/coder/websocket"
)

type Server struct {
	httpServer *http.Server
	mux        *http.ServeMux
	hub        *ws.Hub
}

func NewServer(addr string, hub *ws.Hub) *Server {
	mux := http.NewServeMux()

	s := &Server{mux: mux, hub: hub}
	s.routes() // регистрация путей вынесена отдельно от конструктора

	s.httpServer = &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second, // защита от slowloris, не мешает долгим WS-соединениям

	}

	return s
}

func (s *Server) Start() error {
	log.Printf("server started at http://localhost%s", s.httpServer.Addr)
	err := s.httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func (s *Server) routes() {

	s.mux.HandleFunc("GET /health", s.handleHealth)
	s.mux.HandleFunc("GET /ws", s.handleWS)
	fs := http.FileServer(http.Dir("web/dist"))
	s.mux.Handle("GET /", fs)
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/index.html")
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Ok"))
}

func (s *Server) handleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		log.Println("accept error:", err)
		return
	}

	client := ws.NewClient(s.hub, conn)
	s.hub.Register(client)

	go client.WritePump()
	go client.ReadPump()
}
