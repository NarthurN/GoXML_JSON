// Пакет для middleware
package middleware

import (
	"net/http"
	"strings"

	"github.com/NarthurN/GoXML_JSON/pkg/logger"
	"github.com/NarthurN/GoXML_JSON/settings"
)

// Auth - middleware для авторизации
func Auth(logger *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")

			const prefix = "Bearer "
			if !strings.HasPrefix(authHeader, prefix) {
				logger.Logf("❌ неверный формат токена: remote=%s, path=%s, got_header=%q", r.RemoteAddr, r.URL.Path, authHeader)
				w.Header().Set("WWW-Authenticate", `Bearer realm="Access to the API"`)
				http.Error(w, "❌ не авторизованный доступ", http.StatusUnauthorized)
				return
			}

			// Получение токена без "Bearer "
			token := strings.TrimSpace(authHeader[len(prefix):])
			if token != settings.AuthKey {
				logger.Logf("❌ неверный токен: remote=%s, path=%s, got_token=%q", r.RemoteAddr, r.URL.Path, token)
				w.Header().Set("WWW-Authenticate", `Bearer realm="Access to the API"`)
				http.Error(w, "❌ не авторизованный доступ", http.StatusUnauthorized)
				return
			}

			logger.Logf("✅ авторизованный доступ: remote=%s, path=%s", r.RemoteAddr, r.URL.Path)
			next.ServeHTTP(w, r)
		})
	}
}
