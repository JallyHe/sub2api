package routes

import (
	"github.com/Wei-Shaw/sub2api/internal/handler"
	"github.com/Wei-Shaw/sub2api/internal/server/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterClientRoutes(
	r *gin.Engine,
	h *handler.Handlers,
	jwtAuth middleware.JWTAuthMiddleware,
) {
	r.GET("/client-auth", h.ClientAuth.ShowLoginPage)
	r.POST("/client-auth/login", h.ClientAuth.HandleLogin)
	r.GET("/client-auth/oauth/success", h.ClientAuth.OAuthSuccess)

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
