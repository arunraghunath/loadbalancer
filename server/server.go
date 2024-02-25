package server

import (
	"fmt"
	"net/http"
	"sync"
)

type Server struct {
	Address string `json:"port"`
	Name    string `json:"name"`
	URL     string `json:"url"`
	healthy bool
	mu      sync.RWMutex
}

func (s *Server) DisconnectServer() {
	s.mu.Lock()
	s.healthy = false
	s.mu.Unlock()
}
func (s *Server) SetHealthy(status bool) {
	s.mu.Lock()
	s.healthy = true
	s.mu.Unlock()
}

func (s *Server) IsHealthy() bool {
	s.mu.RLock()
	isHealthy := s.healthy
	s.mu.RUnlock()
	return isHealthy
}

func (s *Server) StartServer() {
	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "This request is handled by server : %v\n", s.Address)
	}))
	s.SetHealthy(true)
	fmt.Printf("Starting server with name-->%s, and url --> %s\n", s.Name, s.URL)
	http.ListenAndServe(s.Address, mux)

}
