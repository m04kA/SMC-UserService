package delete_car

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/m04kA/SMK-UserService/internal/handlers/api"
	"github.com/m04kA/SMK-UserService/internal/handlers/middleware"
	userservice "github.com/m04kA/SMK-UserService/internal/service/user"
)

type Handler struct {
	service *userservice.Service
}

func NewHandler(service *userservice.Service) *Handler {
	return &Handler{service: service}
}

// Handle DELETE /users/me/cars/{car_id}
func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		api.RespondUnauthorized(w, "Unauthorized")
		return
	}

	vars := mux.Vars(r)
	carID := vars["car_id"]
	if carID == "" {
		api.RespondBadRequest(w, "Car ID is required")
		return
	}

	err = h.service.DeleteCar(r.Context(), userID, carID)
	if err != nil {
		if errors.Is(err, userservice.ErrCarNotFound) {
			api.RespondCarNotFound(w)
			return
		}
		if errors.Is(err, userservice.ErrCarAccessDenied) {
			api.RespondCarAccessDenied(w)
			return
		}
		api.RespondInternalError(w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
