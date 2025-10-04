package update_current_user

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

// Handle PUT /users/me
func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		api.RespondUnauthorized(w, "Unauthorized")
		return
	}

	var input models.UpdateUserInputDTO
	if err := api.DecodeJSON(r, &input); err != nil {
		api.RespondBadRequest(w, "Invalid request body")
		return
	}

	user, err := h.service.UpdateUser(r.Context(), userID, input)
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
