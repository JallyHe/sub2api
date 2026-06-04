package routes

import (
	"github.com/Wei-Shaw/sub2api/internal/handler"
	"github.com/Wei-Shaw/sub2api/internal/server/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterClientRoutes registers the StoryClaw desktop client API.
//
// Browser-based OAuth flow (no auth required):
//   GET  /client-auth             → login page HTML
//   POST /client-auth/login       → process login, redirect to storyclaw://auth?token=...
//
// REST API (JWT required):
//   POST /api/client/auth/login   → exchange email+password for JWT (legacy / testing)
//   GET  /api/client/models       → endpoint + api_key + model list + credits
//   GET  /api/client/credits      → credit balance
func RegisterClientRoutes(
	r *gin.Engine,
	h *handler.Handlers,
	jwtAuth middleware.JWTAuthMiddleware,
) {
	// Browser-based deep-link auth flow
	r.GET("/client-auth", h.ClientAuth.ShowLoginPage)
	r.POST("/client-auth/login", h.ClientAuth.HandleLogin)

	// REST API for the desktop client
	client := r.Group("/api/client")
	{
		auth := client.Group("/auth")
		{
			auth.POST("/login", h.Client.Login)
		}

		authed := client.Group("")
		authed.Use(gin.HandlerFunc(jwtAuth))
		{
			authed.GET("/models", h.Client.GetModels)
			authed.GET("/credits", h.Client.GetCredits)
		}
	}
}
