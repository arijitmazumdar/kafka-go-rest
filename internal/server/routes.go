package server

import (
	"encoding/json"
	"log"
	"net/http"

	_ "kafka-go-rest/docs"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/", s.HelloWorldHandler)
	r.HandleFunc("/health", s.HealthCheckHandler)

	http.Handle("/swagger/", httpSwagger.WrapHandler)

	r.HandleFunc("/users", s.UsersHandler).Methods("GET")
	r.HandleFunc("/users/{id}", s.UserByKeyHandler).Methods("GET")
	/**
	r.HandleFunc("/users/{id}", s.GetUser).Methods("GET")
	r.HandleFunc("/users", s.CreateUser).Methods("POST")
	r.HandleFunc("/users/{id}", u.UpdateUser).Methods("PUT")
	r.HandleFunc("/users", s.UsersHandler).Methods("GET")
	*/
	return r
}

func (s *Server) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["status"] = "healthy"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(jsonResp)
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(jsonResp)
}
