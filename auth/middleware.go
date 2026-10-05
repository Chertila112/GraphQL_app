package auth

import (
	"net/http"
	"os"
)

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := os.Getenv("TEST_USER_ID")
		if userID == "" {
			next.ServeHTTP(w, r)
			return
		}
		ctx := WithUserID(r.Context(), userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
