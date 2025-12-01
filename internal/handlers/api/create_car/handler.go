package create_car

import (
	"errors"
	"net/http"

	"github.com/m04kA/SMC-UserService/internal/handlers/api"
	"github.com/m04kA/SMC-UserService/internal/handlers/middleware"
	userservice "github.com/m04kA/SMC-UserService/internal/service/user"
	"github.com/m04kA/SMC-UserService/internal/service/user/models"
	"github.com/m04kA/SMC-UserService/pkg/logger"
)

type Handler struct {
	service *userservice.Service
}

func NewHandler(service *userservice.Service) *Handler {
	return &Handler{service: service}
}

// Handle POST /users/me/cars
func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		logger.Warn("POST /users/me/cars - Unauthorized access attempt")
		api.RespondUnauthorized(w, "Unauthorized")
		return
	}

	var input models.CreateCarInputDTO
	if err := api.DecodeJSON(r, &input); err != nil {
		logger.Warn("POST /users/me/cars - Invalid request body: user_id=%d, error=%v", userID, err)
		api.RespondBadRequest(w, "Invalid request body")
		return
	}

	car, err := h.service.CreateCar(r.Context(), userID, input)
	if err != nil {
		if errors.Is(err, userservice.ErrUserNotFound) {
			logger.Warn("POST /users/me/cars - User not found: user_id=%d", userID)
			api.RespondUserNotFound(w)
			return
		}
		logger.Error("POST /users/me/cars - Failed to create car: user_id=%d, error=%v", userID, err)
		api.RespondInternalError(w)
		return
	}

	logger.Info("POST /users/me/cars - Car created successfully: user_id=%d, car_id=%d", userID, car.ID)
	api.RespondJSON(w, http.StatusCreated, car)
}
