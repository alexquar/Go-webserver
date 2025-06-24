package middleware

import (
	"context"
	"github.com/google/uuid"
	"net/http"
)

const userIDKey = "userID"

func UUIDCookieMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("user_id")
		var userID string
		if err != nil || cookie.Value == "" {
			userID = uuid.NewString()

			http.SetCookie(w, &http.Cookie{
				Name:     "user_id",
				Value:    userID,
				Path:     "/",
				HttpOnly: true,
				MaxAge:   60 * 60 * 24 * 365, // 1 year
				SameSite: http.SameSiteLaxMode,
				Secure:   false, // true if using HTTPS
			})
		} else {
			userID = cookie.Value
		}

		// Store in request context
		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserID(r *http.Request) string {
	if val, ok := r.Context().Value(userIDKey).(string); ok {
		return val
	}
	return ""
}
