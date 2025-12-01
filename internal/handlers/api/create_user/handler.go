package create_user

import (
	"errors"
	"net/http"

	"github.com/m04kA/SMC-UserService/internal/handlers/api"
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

// Handle POST /users
func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	var input models.CreateUserInputDTO
	if err := api.DecodeJSON(r, &input); err != nil {
		logger.Warn("POST /users - Invalid request body: %v", err)
		api.RespondBadRequest(w, "Invalid request body")
		return
	}

	user, err := h.service.CreateUser(r.Context(), input)
	if err != nil {
		if errors.Is(err, userservice.ErrUserAlreadyExists) {
			logger.Warn("POST /users - User already exists: tg_user_id=%d", input.TGUserID)
			api.RespondUserAlreadyExists(w)
			return
		}
		logger.Error("POST /users - Failed to create user: tg_user_id=%d, error=%v", input.TGUserID, err)
		api.RespondInternalError(w)
		return
	}

	logger.Info("POST /users - User created successfully: tg_user_id=%d", user.TGUserID)
	api.RespondJSON(w, http.StatusCreated, user)
}
