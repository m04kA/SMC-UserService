package create_car

import (
	"errors"
	"net/http"

	"github.com/m04kA/SMK-UserService/internal/handlers/api"
	"github.com/m04kA/SMK-UserService/internal/handlers/middleware"
	userservice "github.com/m04kA/SMK-UserService/internal/service/user"
	"github.com/m04kA/SMK-UserService/internal/service/user/models"
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
		api.RespondUnauthorized(w, "Unauthorized")
		return
	}

	var input models.CreateCarInputDTO
	if err := api.DecodeJSON(r, &input); err != nil {
		api.RespondBadRequest(w, "Invalid request body")
		return
	}

	car, err := h.service.CreateCar(r.Context(), userID, input)
	if err != nil {
		if errors.Is(err, userservice.ErrUserNotFound) {
			api.RespondUserNotFound(w)
			return
		}
		api.RespondInternalError(w)
		return
	}

	api.RespondJSON(w, http.StatusCreated, car)
}
