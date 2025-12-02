package handler

import (
	"github.com/gorilla/mux"
	"github.com/ucups/go-user-service/internal/usecase"
)

// SetupRoutes configures all HTTP routes
func SetupRoutes(userUseCase *usecase.UserUseCase) *mux.Router {
	handler := NewUserHandler(userUseCase)

	router := mux.NewRouter()

	// User routes
	router.HandleFunc("/users/ping", handler.Ping).Methods("GET")
	router.HandleFunc("/users/{id}", handler.GetUser).Methods("GET")
	router.HandleFunc("/users", handler.GetAllUsers).Methods("GET")
	router.HandleFunc("/users", handler.CreateUser).Methods("POST")

	return router
}
