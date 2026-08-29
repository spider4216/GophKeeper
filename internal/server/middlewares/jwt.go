package middlewares

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
)

type claims struct {
	jwt.RegisteredClaims
	UserID int64
}

func (m Middleware) WithJwt(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ow := w

		// Извлекаем из заголовка токен
		tokenStr := r.Header.Get("Authorization")

		// Если его нет, ошибка
		if tokenStr == "" {
			m.logger.Error("Header Authorization was not provided")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		claims := &claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims,
			func(t *jwt.Token) (interface{}, error) {
				if t.Method != jwt.SigningMethodHS256 {
					m.logger.Error("JWT Alg not match")
					return nil, fmt.Errorf("JWT Alg not match: %s", t.Method.Alg())
				}

				return []byte(m.cfg.JWTKey), nil
			})
		if err != nil {
			m.logger.Errorf("Something went wrong with token: %s", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			m.logger.Error("Invalid token")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		m.logger.Debug("User id is ", claims.UserID, " set to ctx")

		// Устанавливаем идентификатор пользователя в контекст
		ctx := m.service.SetUserIdToCtx(r.Context(), int64(claims.UserID))
		r = r.WithContext(ctx)

		h.ServeHTTP(ow, r)
	}

	return http.HandlerFunc(fn)
}
