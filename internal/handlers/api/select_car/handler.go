package select_car

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/m04kA/SMC-UserService/internal/handlers/api"
	"github.com/m04kA/SMC-UserService/internal/handlers/middleware"
	userservice "github.com/m04kA/SMC-UserService/internal/service/user"
	"github.com/m04kA/SMC-UserService/pkg/logger"
)

type Handler struct {
	service *userservice.Service
}

func NewHandler(service *userservice.Service) *Handler {
	return &Handler{service: service}
}

// Handle PUT /users/me/cars/{car_id}/select
func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		logger.Warn("PUT /users/me/cars/{car_id}/select - Unauthorized access attempt")
		api.RespondUnauthorized(w, "Unauthorized")
		return
	}

	role, err := middleware.GetRoleFromContext(r.Context())
	if err != nil {
		logger.Warn("PUT /users/me/cars/{car_id}/select - Cannot extract role: user_id=%d", userID)
		api.RespondUnauthorized(w, "Unauthorized")
		return
	}

	vars := mux.Vars(r)
	carIDStr := vars["car_id"]
	carID, err := strconv.ParseInt(carIDStr, 10, 64)
	if err != nil {
		logger.Warn("PUT /users/me/cars/{car_id}/select - Invalid car_id: user_id=%d, car_id=%s", userID, carIDStr)
		api.RespondJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid car_id",
		})
		return
	}

	car, err := h.service.SetSelectedCar(r.Context(), userID, carID, role)
	if err != nil {
		if errors.Is(err, userservice.ErrCarNotFound) {
			logger.Warn("PUT /users/me/cars/{car_id}/select - Car not found: user_id=%d, car_id=%d", userID, carID)
			api.RespondJSON(w, http.StatusNotFound, map[string]string{
				"error": "Car not found",
			})
			return
		}
		if errors.Is(err, userservice.ErrCarAccessDenied) {
			logger.Warn("PUT /users/me/cars/{car_id}/select - Access denied: user_id=%d, car_id=%d", userID, carID)
			api.RespondCarAccessDenied(w)
			return
		}
		logger.Error("PUT /users/me/cars/{car_id}/select - Failed to select car: user_id=%d, car_id=%d, error=%v", userID, carID, err)
		api.RespondInternalError(w)
		return
	}

	logger.Info("PUT /users/me/cars/{car_id}/select - Car selected: user_id=%d, car_id=%d", userID, carID)
	api.RespondJSON(w, http.StatusOK, car)
}
