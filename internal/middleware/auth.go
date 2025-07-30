package middleware

import (
	"context"
	"net/http"
	"testi/internal/session"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Получить session_id из cookie
		cookie, err := r.Cookie("session_id")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Проверить сессию в Redis
		sessionData, err := session.GetSession(cookie.Value)
		if err != nil {
			// Удалить невалидную cookie
			http.SetCookie(w, &http.Cookie{
				Name:   "session_id",
				Value:  "",
				Path:   "/",
				MaxAge: -1,
			})
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Добавить данные пользователя в контекст запроса
		ctx := r.Context()
		ctx = context.WithValue(ctx, "user", sessionData)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Функция для получения пользователя из контекста
func GetUserFromContext(r *http.Request) *session.Session {
	user := r.Context().Value("user")
	if user == nil {
		return nil
	}
	return user.(*session.Session)
}
