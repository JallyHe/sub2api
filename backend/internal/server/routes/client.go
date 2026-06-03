package routes

import (
	"github.com/Wei-Shaw/sub2api/internal/handler"
	"github.com/Wei-Shaw/sub2api/internal/server/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterClientRoutes registers the StoryClaw desktop client API.
// Public: POST /api/client/auth/login (no auth required)
// Protected: GET  /api/client/models, GET /api/client/credits (JWT required)
func RegisterClientRoutes(
	r *gin.Engine,
	h *handler.Handlers,
	jwtAuth middleware.JWTAuthMiddleware,
) {
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
