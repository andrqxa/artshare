package currentuser

import (
	"context"
	"net/http"
	"strings"
)

const DefaultUserID = "00000000-0000-0000-0000-000000000001"

const tokenPrefix = "dev-token-"

type contextKey struct{}

func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, contextKey{}, userID)
}

func FromContext(ctx context.Context) string {
	if value, ok := ctx.Value(contextKey{}).(string); ok && value != "" {
		return value
	}
	return DefaultUserID
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := userIDFromAuthorization(r.Header.Get("Authorization"))
		if userID == "" {
			userID = DefaultUserID
		}
		next.ServeHTTP(w, r.WithContext(WithUserID(r.Context(), userID)))
	})
}

func userIDFromAuthorization(header string) string {
	if header == "" {
		return ""
	}
	const bearer = "Bearer "
	if !strings.HasPrefix(header, bearer) {
		return ""
	}
	token := strings.TrimPrefix(header, bearer)
	if !strings.HasPrefix(token, tokenPrefix) {
		return ""
	}
	return strings.TrimPrefix(token, tokenPrefix)
}
