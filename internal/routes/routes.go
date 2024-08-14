package routes

import (
	"myanimeapi/pkg/handlers"

	"github.com/gorilla/mux"
)

func RegisterRoutes(router *mux.Router) {
	// Register anime routes
	handlers.RegisterAnimeRoutes(router)

	// Register user routes
	handlers.RegisterUserRoutes(router)

	// Register review routes
	handlers.RegisterReviewRoutes(router)

	// Register auth routes
	handlers.RegisterAuthRoutes(router)
}
