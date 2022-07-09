package handlers

import (
	"context"
	"net/http"
	"task/internal/app/grpc/client"
	"task/internal/app/interfaces"
)

type MiddlewareProfile struct {
	UserID string
	Roles  []string
}

type Middleware struct {
	ts interfaces.TokenStorage
	vc *client.GRPCValidatorClient
}

func NewMiddleware(ts interfaces.TokenStorage, vc *client.GRPCValidatorClient) *Middleware {
	return &Middleware{
		ts: ts,
		vc: vc,
	}
}

func (mw *Middleware) Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var mwProfile MiddlewareProfile
		mw.ts.SetParams(map[string]interface{}{
			"request": r,
		})
		accessToken, refreshToken, err := mw.ts.GetTokens()
		if err == nil {
			userData, err := mw.vc.Validate(accessToken, refreshToken)
			if err == nil {
				mwProfile.UserID = userData.UserId
				mwProfile.Roles = append(mwProfile.Roles, userData.Roles...)
				_ = mw.ts.SetTokens(accessToken, refreshToken)
			}
		}
		ctx := r.Context()
		r = r.WithContext(context.WithValue(ctx, "profile", mwProfile))
		next.ServeHTTP(w, r)
	})
}
