package handlers

import (
	"context"
	"net/http"
)

type MiddlewareProfile struct {
	UserID int
	Roles  []string
}

type Middleware struct {
}

func NewMiddleware() *Middleware {
	return &Middleware{}
}

func (mw *Middleware) Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		r = r.WithContext(context.WithValue(ctx, "profile", MiddlewareProfile{
			UserID: 1,
			Roles:  []string{"task_creator", "authenticated"},
		}))
		next.ServeHTTP(w, r)
	})
}
