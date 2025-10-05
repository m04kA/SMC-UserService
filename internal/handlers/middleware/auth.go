package middleware

import (
	"context"
	"errors"
	"net/http"
	"strconv"
)

type contextKey string

const (
	UserIDKey contextKey = "userID"
)

var (
	ErrMissingUserID = errors.New("missing X-User-ID header")
	ErrInvalidUserID = errors.New("invalid user ID format")
)

// UserIDAuth извлекает user ID из заголовка X-User-ID
func UserIDAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userIDStr := r.Header.Get("X-User-ID")
		if userIDStr == "" {
			http.Error(w, ErrMissingUserID.Error(), http.StatusUnauthorized)
			return
		}

		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			http.Error(w, ErrInvalidUserID.Error(), http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserIDFromContext извлекает user ID из контекста
func GetUserIDFromContext(ctx context.Context) (int64, error) {
	userID, ok := ctx.Value(UserIDKey).(int64)
	if !ok {
		return 0, ErrInvalidUserID
	}
	return userID, nil
}
