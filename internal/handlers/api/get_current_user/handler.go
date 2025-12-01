package get_current_user

import (
	"errors"
	"net/http"

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

// Handle GET /users/me
func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		logger.Warn("GET /users/me - Unauthorized access attempt")
		api.RespondUnauthorized(w, "Unauthorized")
		return
	}

	user, err := h.service.GetUserWithCars(r.Context(), userID)
	if err != nil {
		if errors.Is(err, userservice.ErrUserNotFound) {
			logger.Warn("GET /users/me - User not found: user_id=%d", userID)
			api.RespondUserNotFound(w)
			return
		}
		logger.Error("GET /users/me - Failed to get user: user_id=%d, error=%v", userID, err)
		api.RespondInternalError(w)
		return
	}

	logger.Info("GET /users/me - User retrieved successfully: user_id=%d", userID)
	api.RespondJSON(w, http.StatusOK, user)
}
