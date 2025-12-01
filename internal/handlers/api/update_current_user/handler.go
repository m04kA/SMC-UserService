package update_current_user

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

// Handle PUT /users/me
func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		logger.Warn("PUT /users/me - Unauthorized access attempt")
		api.RespondUnauthorized(w, "Unauthorized")
		return
	}

	var input models.UpdateUserInputDTO
	if err := api.DecodeJSON(r, &input); err != nil {
		logger.Warn("PUT /users/me - Invalid request body: user_id=%d, error=%v", userID, err)
		api.RespondBadRequest(w, "Invalid request body")
		return
	}

	user, err := h.service.UpdateUser(r.Context(), userID, input)
	if err != nil {
		if errors.Is(err, userservice.ErrUserNotFound) {
			logger.Warn("PUT /users/me - User not found: user_id=%d", userID)
			api.RespondUserNotFound(w)
			return
		}
		logger.Error("PUT /users/me - Failed to update user: user_id=%d, error=%v", userID, err)
		api.RespondInternalError(w)
		return
	}

	logger.Info("PUT /users/me - User updated successfully: user_id=%d", userID)
	api.RespondJSON(w, http.StatusOK, user)
}
