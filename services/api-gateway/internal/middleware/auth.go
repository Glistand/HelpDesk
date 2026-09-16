package middleware

import (
	"context"
	"net/http"
	"strings"

	authv1 "github.com/Glistand/HelpDesk/api/gen/go/helpdesk/auth/v1"
)

type ctxKey string

const UserKey ctxKey = "user"

func UserFromContext(ctx context.Context) *authv1.User {
	u, _ := ctx.Value(UserKey).(*authv1.User)
	return u
}

func Auth(auth authv1.AuthServiceClient) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h := r.Header.Get("Authorization")
			if !strings.HasPrefix(h, "Bearer ") {
				http.Error(w, `{"error":"missing bearer token"}`, http.StatusUnauthorized)
				return
			}
			token := strings.TrimPrefix(h, "Bearer ")
			resp, err := auth.ValidateToken(r.Context(), &authv1.ValidateTokenRequest{AccessToken: token})
			if err != nil || !resp.GetValid() {
				http.Error(w, `{"error":"invalid token"}`, http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), UserKey, resp.GetUser())
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
