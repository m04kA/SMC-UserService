package get_current_user

import (
	"errors"
	"net/http"

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

// Handle GET /users/me
func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		api.RespondUnauthorized(w, "Unauthorized")
		return
	}

	user, err := h.service.GetUserWithCars(r.Context(), userID)
	if err != nil {
		if errors.Is(err, userservice.ErrUserNotFound) {
			api.RespondUserNotFound(w)
			return
		}
		api.RespondInternalError(w)
		return
	}

	api.RespondJSON(w, http.StatusOK, user)
}
