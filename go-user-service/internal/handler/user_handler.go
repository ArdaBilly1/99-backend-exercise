package handler

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ucups/go-user-service/internal/usecase"
)

// UserHandler handles user HTTP requests
type UserHandler struct {
	userUseCase *usecase.UserUseCase
}

// NewUserHandler creates a new user handler
func NewUserHandler(userUseCase *usecase.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
	}
}

// CreateUser handles POST /users
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	if err := r.ParseForm(); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid form data")
		return
	}

	name := r.FormValue("name")

	// Create user via use case
	user, err := h.userUseCase.CreateUser(name)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Return success response
	WriteSuccess(w, map[string]interface{}{
		"user": user,
	})
}

// GetUser handles GET /users/{id}
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	// Get user via use case
	user, err := h.userUseCase.GetUserByID(id)
	if err != nil {
		WriteError(w, http.StatusNotFound, err.Error())
		return
	}

	// Return success response
	WriteSuccess(w, map[string]interface{}{
		"user": user,
	})
}

// GetAllUsers handles GET /users
func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	// Parse pagination params
	pageNumStr := r.URL.Query().Get("page_num")
	pageSizeStr := r.URL.Query().Get("page_size")

	pageNum := 1
	pageSize := 10

	if pageNumStr != "" {
		if val, err := strconv.Atoi(pageNumStr); err == nil {
			pageNum = val
		} else {
			WriteError(w, http.StatusBadRequest, "invalid page_num")
			return
		}
	}

	if pageSizeStr != "" {
		if val, err := strconv.Atoi(pageSizeStr); err == nil {
			pageSize = val
		} else {
			WriteError(w, http.StatusBadRequest, "invalid page_size")
			return
		}
	}

	// Get users via use case
	users, err := h.userUseCase.GetAllUsers(pageNum, pageSize)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Return success response
	WriteSuccess(w, map[string]interface{}{
		"users": users,
	})
}

// Ping handles GET /users/ping
func (h *UserHandler) Ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong!"))
}
