package currentuser

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMiddleware_ExtractsUserIDFromBearer(t *testing.T) {
	const userID = "11111111-1111-1111-1111-111111111111"

	var captured string
	handler := Middleware(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		captured = FromContext(r.Context())
	}))

	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	req.Header.Set("Authorization", "Bearer dev-token-"+userID)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if captured != userID {
		t.Fatalf("got %q want %q", captured, userID)
	}
}

func TestMiddleware_FallsBackToDefault(t *testing.T) {
	var captured string
	handler := Middleware(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		captured = FromContext(r.Context())
	}))

	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if captured != DefaultUserID {
		t.Fatalf("got %q want %q", captured, DefaultUserID)
	}
}

func TestMiddleware_IgnoresUnknownToken(t *testing.T) {
	var captured string
	handler := Middleware(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		captured = FromContext(r.Context())
	}))

	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	req.Header.Set("Authorization", "Bearer some-other-token")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if captured != DefaultUserID {
		t.Fatalf("got %q want %q", captured, DefaultUserID)
	}
}
