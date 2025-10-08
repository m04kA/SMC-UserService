package get_selected_car

import (
	"errors"
	"net/http"

	"github.com/m04kA/SMK-UserService/internal/handlers/api"
	"github.com/m04kA/SMK-UserService/internal/handlers/middleware"
	userservice "github.com/m04kA/SMK-UserService/internal/service/user"
	"github.com/m04kA/SMK-UserService/pkg/logger"
)

type Handler struct {
	service *userservice.Service
}

func NewHandler(service *userservice.Service) *Handler {
	return &Handler{service: service}
}

// Handle GET /users/me/cars/selected
func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		logger.Warn("GET /users/me/cars/selected - Unauthorized access attempt")
		api.RespondUnauthorized(w, "Unauthorized")
		return
	}

	car, err := h.service.GetSelectedCar(r.Context(), userID)
	if err != nil {
		if errors.Is(err, userservice.ErrCarNotFound) {
			logger.Warn("GET /users/me/cars/selected - No selected car: user_id=%d", userID)
			api.RespondJSON(w, http.StatusNotFound, map[string]string{
				"error": "No selected car found",
			})
			return
		}
		logger.Error("GET /users/me/cars/selected - Failed to get selected car: user_id=%d, error=%v", userID, err)
		api.RespondInternalError(w)
		return
	}

	logger.Info("GET /users/me/cars/selected - Selected car retrieved: user_id=%d, car_id=%d", userID, car.ID)
	api.RespondJSON(w, http.StatusOK, car)
}
