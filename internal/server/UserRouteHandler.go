package server

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// User represents a user in the system
type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// UserRouteHandler handles routes related to users
type UserRouteHandler struct {
	users map[string]User
}

// NewUserRouteHandler creates a new UserRouteHandler
func NewUserRouteHandler() *UserRouteHandler {
	return &UserRouteHandler{
		users: make(map[string]User),
	}
}

// RegisterRoutes registers the user routes with the router
func (h *UserRouteHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/users", h.GetUsers).Methods("GET")
	router.HandleFunc("/users/{id}", h.GetUser).Methods("GET")
	router.HandleFunc("/users", h.CreateUser).Methods("POST")
	router.HandleFunc("/users/{id}", h.UpdateUser).Methods("PUT")
	router.HandleFunc("/users/{id}", h.DeleteUser).Methods("DELETE")
}

// GetUsers handles GET /users
func (h *UserRouteHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users := make([]User, 0, len(h.users))
	for _, user := range h.users {
		users = append(users, user)
	}
	json.NewEncoder(w).Encode(users)
}

// GetUser handles GET /users/{id}
func (h *UserRouteHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	user, exists := h.users[id]
	if !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(user)
}

// CreateUser handles POST /users
func (h *UserRouteHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	h.users[user.ID] = user
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// UpdateUser handles PUT /users/{id}
func (h *UserRouteHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if _, exists := h.users[id]; !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	h.users[id] = user
	json.NewEncoder(w).Encode(user)
}

// DeleteUser handles DELETE /users/{id}
func (h *UserRouteHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if _, exists := h.users[id]; !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	delete(h.users, id)
	w.WriteHeader(http.StatusNoContent)
}
