package auth

import (
	"context"
	"net/http"
)

// A private key for context that only this package can access. This is important
// to prevent collisions between different context uses
var adminCodeCtxKey = &contextKey{"adminCode"}

type contextKey struct {
	name string
}

func Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c := r.Header["X-Admin-Code"]

			// Allow unauthenticated users in
			if len(c) == 0 {
				next.ServeHTTP(w, r)
				return
			}

			adminCode := c[0]

			// put it in context
			ctx := context.WithValue(r.Context(), adminCodeCtxKey, adminCode)

			// and call the next with our new context
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func ForContext(ctx context.Context) *string {
	raw, _ := ctx.Value(adminCodeCtxKey).(string)
	return &raw
}
