package middleware

import (
	"context"
	"go-authentication/src/pkg"
	"net/http"
)

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("at")

		if err != nil {
			http.Error(w, "Missing authorization token", http.StatusUnauthorized)
			return
		}

		claims, err := pkg.ValidateJWT(cookie.Value)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Set claims in context or request

		ctx := context.WithValue(r.Context(), "user", claims.Uid)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
