package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

func TestAuthenticateMiddleware(t *testing.T) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  1,
		"is_admin": true,
		"exp":      time.Now().Add(time.Hour * 1).Unix(),
	})
	tokenString, _ := token.SignedString(jwtKey)

	req, err := http.NewRequest("GET", "/anime/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+tokenString)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("user")
		isAdmin := r.Context().Value("is_admin")
		assert.Equal(t, uint(1), userID)
		assert.Equal(t, true, isAdmin)
	})

	middleware := Authenticate(handler)
	middleware.ServeHTTP(rr, req)
}
