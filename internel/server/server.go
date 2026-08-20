package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"project_chat/internel/auth"
	"project_chat/internel/storage"
	"project_chat/internel/validation"
	"project_chat/internel/ws"
	"time"

	"github.com/coder/websocket"
)

type Server struct {
	httpServer *http.Server
	mux        *http.ServeMux
	hub        *ws.Hub
	userRepo   *storage.UserRepository
}

func NewServer(addr string, hub *ws.Hub, userRepo *storage.UserRepository) *Server {
	mux := http.NewServeMux()

	s := &Server{mux: mux, hub: hub, userRepo: userRepo}
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
	s.mux.HandleFunc("POST /register", s.handleRegister)
	s.mux.HandleFunc("POST /login", s.handleLogin)
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

func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	var req auth.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if err := validation.V.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	if err := s.userRepo.Create(r.Context(), req.Username, hash); err != nil {
		http.Error(w, "username already taken", http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	var req auth.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if err := validation.V.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := s.userRepo.GetByUsername(r.Context(), req.Username)
	if err != nil {
		log.Println("login: GetByUsername error:", err)
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	log.Printf("login: found user, hash=%s", user.PasswordHash) // ←
	log.Printf("login: comparing with password=%s", req.Password)

	if !auth.CheckPassword(user.PasswordHash, req.Password) {
		log.Println("login: password mismatch")
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := auth.GenerateToken(user.ID, user.Username)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(auth.TokenResponse{Token: token})
}
