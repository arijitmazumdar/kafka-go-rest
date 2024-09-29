package server

import (
	"encoding/json"
	"log"
	"net/http"
)

type User struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Email string  `json:"email"`
	Age   float64 `json:"age"`
}

// UserRouterHandler handles routes related to users
func (s *Server) UserRouterHandler() http.Handler {
	r := http.NewServeMux()
	r.HandleFunc("/users", s.UsersHandler)
	return r
}

// UsersHandler handles the /users route
// @Summary Get users
// @Description Get a list of users
// @Tags users
// @Accept  json
// @Produce  json
// @Success 200 {array} User
// @Router /users [get]
// UsersHandler handles the /users route
func (s *Server) UsersHandler(w http.ResponseWriter, r *http.Request) {
	users := []User{
		{ID: 1, Name: "John Doe", Email: "john@example.com", Age: 30},
		{ID: 2, Name: "Jane Smith", Email: "jane@example.com", Age: 25},
	}

	jsonResp, err := json.Marshal(users)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(jsonResp)
}

// UserByKeyHandler handles the /users/{id} route
func (s *Server) UserByKeyHandler(w http.ResponseWriter, r *http.Request) {
	user := User{
		ID:    1,
		Name:  "John Doe",
		Email: "john@example.com",
		Age:   30,
	}

	jsonResp, err := json.Marshal(user)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(jsonResp)
}
