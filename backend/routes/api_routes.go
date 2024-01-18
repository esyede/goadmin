package routes

import (
	"backend/controller"
	"backend/middleware"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

func InitApiRoutes(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) gin.IRoutes {
	apiController := controller.NewApiController()
	router := r.Group("/api")
	// Enable jwt auth middleware
	router.Use(authMiddleware.MiddlewareFunc())
	// Enable casbin auth middleware
	router.Use(middleware.CasbinMiddleware())
	{
		router.GET("/list", apiController.GetApis)
		router.GET("/tree", apiController.GetApiTree)
		router.POST("/create", apiController.CreateApi)
		router.PATCH("/update/:apiId", apiController.UpdateApiById)
		router.DELETE("/delete/batch", apiController.BatchDeleteApiByIds)
	}

	return r
}
