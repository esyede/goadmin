package routes

import (
	"github.com/esyede/goadmin/backend/controller"
	"github.com/esyede/goadmin/backend/middleware"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

func InitMenuRoutes(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) gin.IRoutes {
	menuController := controller.NewMenuController()
	router := r.Group("/menu")
	// Enable jwt auth middleware
	router.Use(authMiddleware.MiddlewareFunc())
	// Enable casbin auth middleware
	router.Use(middleware.CasbinMiddleware())
	{
		router.GET("/tree", menuController.GetMenuTree)
		router.GET("/list", menuController.GetMenus)
		router.POST("/create", menuController.CreateMenu)
		router.PATCH("/update/:menuId", menuController.UpdateMenuById)
		router.DELETE("/delete/batch", menuController.BatchDeleteMenuByIds)
		router.GET("/access/list/:userId", menuController.GetUserMenusByUserId)
		router.GET("/access/tree/:userId", menuController.GetUserMenuTreeByUserId)
	}

	return r
}
