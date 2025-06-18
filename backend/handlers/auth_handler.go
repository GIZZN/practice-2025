package handlers

import (
	"delivery-service/middleware"
	"delivery-service/models"
	"delivery-service/services"
	"encoding/json"
	"net/http"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "неверный формат данных", http.StatusBadRequest)
		return
	}

	response, err := h.authService.Register(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	middleware.SendJSON(w, http.StatusCreated, response)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "неверный формат данных", http.StatusBadRequest)
		return
	}

	response, err := h.authService.Login(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	middleware.SendJSON(w, http.StatusOK, response)
}

func (h *AuthHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*models.User)
	if user == nil {
		http.Error(w, "пользователь не найден", http.StatusNotFound)
		return
	}

	response := struct {
		*models.User
		BirthDate *string `json:"birthDate,omitempty"`
	}{
		User: user,
	}

	if user.BirthDate != nil {
		birthDate := user.BirthDate.Format("2006-01-02")
		response.BirthDate = &birthDate
	}

	middleware.SendJSON(w, http.StatusOK, response)
}

func (h *AuthHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*models.User)
	if user == nil {
		http.Error(w, "пользователь не найден", http.StatusNotFound)
		return
	}

	var req models.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "неверный формат данных", http.StatusBadRequest)
		return
	}

	if err := h.authService.UpdateUser(user.ID, &req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	updatedUser, err := h.authService.GetUserByID(user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Форматируем даты для фронтенда
	response := struct {
		*models.User
		BirthDate *string `json:"birthDate,omitempty"`
	}{
		User: updatedUser,
	}

	if updatedUser.BirthDate != nil {
		birthDate := updatedUser.BirthDate.Format("2006-01-02")
		response.BirthDate = &birthDate
	}

	middleware.SendJSON(w, http.StatusOK, response)
}
