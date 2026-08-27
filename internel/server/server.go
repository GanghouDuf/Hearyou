package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"project_chat/internel/auth"
	"project_chat/internel/dto"
	"project_chat/internel/storage"
	"project_chat/internel/validation"
	"project_chat/internel/ws"
	"time"

	"github.com/coder/websocket"
)

type Server struct {
	httpServer  *http.Server
	mux         *http.ServeMux
	roomManager *ws.RoomManager
	userRepo    *storage.UserRepository
	messageRepo *storage.MessageRepository
	roomRepo    *storage.RoomRepository
}

func NewServer(addr string, roomManager *ws.RoomManager, userRepo *storage.UserRepository, messageRepo *storage.MessageRepository, roomRepo *storage.RoomRepository) *Server {
	mux := http.NewServeMux()

	s := &Server{mux: mux, roomManager: roomManager, userRepo: userRepo, messageRepo: messageRepo, roomRepo: roomRepo}
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

//func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
//	http.ServeFile(w, r, "web/index.html")
//}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Ok"))
}

func (s *Server) handleWS(w http.ResponseWriter, r *http.Request) {
	tokenString := r.URL.Query().Get("token")
	if tokenString == "" {
		http.Error(w, "missing token", http.StatusUnauthorized)
		return
	}

	roomName := r.URL.Query().Get("room")
	if roomName == "" {
		http.Error(w, "missing room", http.StatusBadRequest)
		return
	}

	claims, err := auth.ValidateToken(tokenString)
	if err != nil {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	userid, ok := claims["user_id"].(float64)
	if !ok {
		http.Error(w, "invalid token claims", http.StatusUnauthorized)
		return
	}

	username, ok := claims["username"].(string)
	if !ok {
		http.Error(w, "invalid token claims", http.StatusUnauthorized)
		return
	}
	userID := int(userid)
	// находим комнату, или создаём, если её ещё нет
	room, err := s.roomRepo.GetByName(r.Context(), roomName)
	if err != nil {
		newID, createErr := s.roomRepo.Create(r.Context(), roomName, userID)
		if createErr != nil {
			http.Error(w, "failed to create room", http.StatusInternalServerError)
			return
		}
		room = &storage.Room{ID: newID, Name: roomName, CreatedBy: userID}
	}

	hub := s.roomManager.GetOrCreate(roomName)

	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		log.Println("accept error:", err)
		return
	}

	client := ws.NewClient(hub, conn, username, userID, room.ID, s.messageRepo)
	hub.Register(client)

	history, err := s.messageRepo.GetHistory(r.Context(), room.ID, 50)
	if err != nil {
		log.Println("failed to load history:", err)
	} else {
		for i, j := 0, len(history)-1; i < j; i, j = i+1, j-1 {
			history[i], history[j] = history[j], history[i]
		}
		for _, m := range history {
			out, _ := json.Marshal(dto.Message{
				Type:    "chat",
				Author:  m.Author,
				Payload: m.Payload,
			})
			client.Send(out)
		}
	}

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
		http.Error(w, "Имя уже занято", http.StatusConflict)
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
		http.Error(w, "Неверные учётные данные", http.StatusUnauthorized)
		return
	}
	log.Printf("login: found user, hash=%s", user.PasswordHash)
	log.Printf("login: comparing with password=%s", req.Password)

	if !auth.CheckPassword(user.PasswordHash, req.Password) {
		log.Println("login: password mismatch")
		http.Error(w, "Неверные учётные данные", http.StatusUnauthorized)
		return
	}

	token, err := auth.GenerateToken(user.ID, user.Username)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(auth.TokenResponse{Token: token})
}
