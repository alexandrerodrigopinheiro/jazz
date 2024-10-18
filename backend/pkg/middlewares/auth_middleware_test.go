// backend/pkg/middlewares/auth_middleware_test.go
package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"jazz/backend/pkg/auth"
)

func TestAuthMiddleware(t *testing.T) {
	validToken, _ := auth.GenerateJWT(1)
	reqWithToken, _ := http.NewRequest("GET", "/protected", nil)
	reqWithToken.Header.Set("Authorization", "Bearer "+validToken)

	rr := httptest.NewRecorder()
	handler := AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(rr, reqWithToken)
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d but got %d", http.StatusOK, rr.Code)
	}

	// Testar sem token
	reqWithoutToken, _ := http.NewRequest("GET", "/protected", nil)
	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, reqWithoutToken)
	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d but got %d", http.StatusUnauthorized, rr.Code)
	}
}
