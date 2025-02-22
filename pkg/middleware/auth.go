package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hashicorp/vault/api"
)

// CustomClaims defines the JWT claims structure
type CustomClaims struct {
	IsAdmin bool `json:"is_admin"`
	UserID  uint `json:"user_id"`
	jwt.RegisteredClaims
}

var jwtKey []byte

func init() {
	// Read the Vault address from the environment, or use a default value
	vaultAddr := os.Getenv("VAULT_ADDR")
	if vaultAddr == "" {
		vaultAddr = "http://vault:8200" // Default value
	}

	// Initialize Vault client
	client, err := api.NewClient(&api.Config{
		Address: vaultAddr, // Use the VAULT_ADDR environment variable or default value
	})
	if err != nil {
		panic(fmt.Errorf("failed to create Vault client: %v", err))
	}

	// Set the Vault token (use the root token for development)
	client.SetToken("root")

	// Read the JWT secret from Vault
	secret, err := client.Logical().Read("secret/data/myapp")
	if err != nil {
		panic(fmt.Errorf("failed to read secret from Vault: %v", err))
	}

	// Extract the secret value
	if secret != nil && secret.Data != nil {
		data, ok := secret.Data["data"].(map[string]interface{})
		if !ok {
			panic("invalid secret data format")
		}
		key, ok := data["JWT_SECRET_KEY"].(string)
		if !ok {
			panic("invalid JWT secret key format")
		}
		jwtKey = []byte(key)
	} else {
		panic("secret not found")
	}
}

// Authenticate is a middleware function that checks for a valid JWT token
func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Authenticate middleware triggered for:", r.URL.Path)
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

		// Store claims in context using the existing contextKey type
		ctx := context.WithValue(r.Context(), userContextKey, claims.UserID)
		ctx = context.WithValue(ctx, isAdminContextKey, claims.IsAdmin)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// CheckAdmin checks if the user has admin privileges
func CheckAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isAdmin, ok := r.Context().Value(isAdminContextKey).(bool)
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
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			Issuer:    "myanimeapi",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}
