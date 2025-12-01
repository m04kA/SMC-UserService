package get_selected_car

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/m04kA/SMC-UserService/internal/handlers/api"
	userservice "github.com/m04kA/SMC-UserService/internal/service/user"
	"github.com/m04kA/SMC-UserService/pkg/logger"
)

type Handler struct {
	service *userservice.Service
}

func NewHandler(service *userservice.Service) *Handler {
	return &Handler{service: service}
}

// Handle GET /internal/users/{tg_user_id}/cars/selected
func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["tg_user_id"]

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		logger.Warn("GET /internal/users/{tg_user_id}/cars/selected - Invalid user_id format: %s", userIDStr)
		api.RespondJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid user_id format",
		})
		return
	}

	car, err := h.service.GetSelectedCar(r.Context(), userID)
	if err != nil {
		if errors.Is(err, userservice.ErrCarNotFound) {
			logger.Warn("GET /internal/users/{tg_user_id}/cars/selected - No selected car: user_id=%d", userID)
			api.RespondJSON(w, http.StatusNotFound, map[string]string{
				"error": "No selected car found",
			})
			return
		}
		logger.Error("GET /internal/users/{tg_user_id}/cars/selected - Failed to get selected car: user_id=%d, error=%v", userID, err)
		api.RespondInternalError(w)
		return
	}

	logger.Info("GET /internal/users/{tg_user_id}/cars/selected - Selected car retrieved: user_id=%d, car_id=%d", userID, car.ID)
	api.RespondJSON(w, http.StatusOK, car)
}
