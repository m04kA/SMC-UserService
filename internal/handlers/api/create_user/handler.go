package create_user

import (
	"errors"
	"net/http"

	"github.com/m04kA/SMK-UserService/internal/handlers/api"
	userservice "github.com/m04kA/SMK-UserService/internal/service/user"
	"github.com/m04kA/SMK-UserService/internal/service/user/models"
)

type Handler struct {
	service *userservice.Service
}

func NewHandler(service *userservice.Service) *Handler {
	return &Handler{service: service}
}

// Handle POST /users
func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	var input models.CreateUserInputDTO
	if err := api.DecodeJSON(r, &input); err != nil {
		api.RespondBadRequest(w, "Invalid request body")
		return
	}

	user, err := h.service.CreateUser(r.Context(), input)
	if err != nil {
		if errors.Is(err, userservice.ErrUserAlreadyExists) {
			api.RespondUserAlreadyExists(w)
			return
		}
		api.RespondInternalError(w)
		return
	}

	api.RespondJSON(w, http.StatusCreated, user)
}
