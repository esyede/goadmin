package routes

import (
	"github.com/esyede/goadmin/backend/controller"
	"github.com/esyede/goadmin/backend/middleware"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

func InitOperationLogRoutes(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) gin.IRoutes {
	operationLogController := controller.NewOperationLogController()
	router := r.Group("/log")
	// Enable jwt auth middleware
	router.Use(authMiddleware.MiddlewareFunc())
	// Enable casbin auth middleware
	router.Use(middleware.CasbinMiddleware())
	{
		router.GET("/operation/list", operationLogController.GetOperationLogs)
		router.DELETE("/operation/delete/batch", operationLogController.BatchDeleteOperationLogByIds)
	}

	return r
}
