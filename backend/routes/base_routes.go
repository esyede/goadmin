package routes

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// Register basic routing
func InitBaseRoutes(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) gin.IRoutes {
	router := r.Group("/base")
	{
		// Login, logout, refresh token (no authentication required)
		router.POST("/login", authMiddleware.LoginHandler)
		router.POST("/logout", authMiddleware.LogoutHandler)
		router.POST("/refreshToken", authMiddleware.RefreshHandler)
	}

	return r
}
