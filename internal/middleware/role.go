package middleware

import "net/http"

func RequireRole(role string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			user, err := GetUser(r)
			if err != nil {
				http.Error(w, "Not Authorized", http.StatusForbidden)
				return
			}
			if user.Role != role {
				http.Error(w, "Unauthorized", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func RequireRoles(roles ...string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, err := GetUser(r)
			if err != nil {
				http.Error(w, "Not Authorized", http.StatusForbidden)
				return
			}

			allowed := false
			for _, role := range roles {
				if user.Role == role {
					allowed = true
					break
				}
			}

			if !allowed {
				http.Error(w, "Unauthorized", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
