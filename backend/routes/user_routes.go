package routes

import (
	"github.com/esyede/goadmin/backend/controller"
	"github.com/esyede/goadmin/backend/middleware"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

func InitUserRoutes(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) gin.IRoutes {
	userController := controller.NewUserController()
	router := r.Group("/user")
	// Enable jwt auth middleware
	router.Use(authMiddleware.MiddlewareFunc())
	// Enable casbin auth middleware
	router.Use(middleware.CasbinMiddleware())
	{
		router.POST("/info", userController.GetUserInfo)
		router.GET("/list", userController.GetUsers)
		router.PUT("/changePwd", userController.ChangePwd)
		router.POST("/create", userController.CreateUser)
		router.PATCH("/update/:userId", userController.UpdateUserById)
		router.DELETE("/delete/batch", userController.BatchDeleteUserByIds)
	}

	return r
}
