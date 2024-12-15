package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	t.Run("No cookie", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rr := httptest.NewRecorder()

		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		handler := AuthMiddleware(testHandler)

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		cookies := rr.Result().Cookies()
		assert.Len(t, cookies, 1, "Expected one cookie to be set")
		assert.Equal(t, "user_id", cookies[0].Name)
		assert.NotEmpty(t, cookies[0].Value)
		assert.Equal(t, "/", cookies[0].Path)
		assert.True(t, cookies[0].HttpOnly)
		assert.Equal(t, 86400*30, cookies[0].MaxAge)
	})

	t.Run("With existing cookie", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		existingCookieValue := "existing-user-id"
		req.AddCookie(&http.Cookie{
			Name:  "user_id",
			Value: existingCookieValue,
		})
		rr := httptest.NewRecorder()

		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("user_id")
			assert.NoError(t, err)
			assert.Equal(t, existingCookieValue, cookie.Value)
			w.WriteHeader(http.StatusOK)
		})

		handler := AuthMiddleware(testHandler)

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		cookies := rr.Result().Cookies()
		assert.Empty(t, cookies, "Expected no new cookies to be set")
	})
}
