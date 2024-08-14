package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Define a global secret key for JWT
var jwtKey = []byte("your_secret_key_here") // Ensure this is consistent

// CustomClaims defines the JWT claims structure
type CustomClaims struct {
	IsAdmin bool `json:"is_admin"`
	UserID  uint `json:"user_id"`
	jwt.StandardClaims
}

// Authenticate is a middleware function that checks for a valid JWT token
func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing authorization header", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(*CustomClaims)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		// Store claims in context
		ctx := context.WithValue(r.Context(), "user", claims.UserID)
		ctx = context.WithValue(ctx, "is_admin", claims.IsAdmin)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// CheckAdmin checks if the user has admin privileges
func CheckAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isAdmin, ok := r.Context().Value("is_admin").(bool)
		if !ok || !isAdmin {
			http.Error(w, "Access denied", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// GenerateToken generates a JWT token with custom claims
func GenerateToken(userID uint, isAdmin bool) (string, error) {
	claims := CustomClaims{
		UserID:  userID,
		IsAdmin: isAdmin,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			Issuer:    "myanimeapi",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}
