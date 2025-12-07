package get_superusers

import (
	"net/http"

	"github.com/m04kA/SMC-UserService/internal/handlers/api"
	"github.com/m04kA/SMC-UserService/internal/service/user"
	"github.com/m04kA/SMC-UserService/pkg/logger"
)

const (
	ErrInternalServer = "внутренняя ошибка сервера"
)

type Handler struct {
	service *user.Service
}

func NewHandler(service *user.Service) *Handler {
	return &Handler{service: service}
}

// Response структура для ответа со списком superuser ID
type Response struct {
	SuperUserIDs []int64 `json:"super_user_ids"`
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	// Получаем список всех superuser ID
	superUserIDs, err := h.service.GetSuperUsers(r.Context())
	if err != nil {
		logger.Error("GET /internal/users/superusers - failed to get superusers: %v", err)
		api.RespondError(w, http.StatusInternalServerError, ErrInternalServer)
		return
	}

	logger.Info("GET /internal/users/superusers - success, found %d superusers", len(superUserIDs))
	api.RespondJSON(w, http.StatusOK, Response{SuperUserIDs: superUserIDs})
}
