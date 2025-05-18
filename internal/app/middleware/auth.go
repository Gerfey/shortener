package middleware

import (
	"net/http"

	"github.com/Gerfey/shortener/internal/app/handler"
	"github.com/google/uuid"
)

// AuthMiddleware проверяет наличие куки с идентификатором пользователя и создает её, если она отсутствует
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(handler.UserIDCookieName)
		if err != nil || cookie == nil {
			userID := uuid.New().String()

			cookie = &http.Cookie{
				Name:     handler.UserIDCookieName,
				Value:    userID,
				Path:     "/",
				MaxAge:   86400 * 30, // 30 дней
				HttpOnly: true,
			}
			http.SetCookie(w, cookie)
		}

		next.ServeHTTP(w, r)
	}
}
