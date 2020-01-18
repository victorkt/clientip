package clientip

import (
	"net/http"
)

// Middleware returns an HTTP handler that will add the client IP to the
// request context.
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := toContext(r.Context(), FromRequest(r))
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
