package routes

import (
	"backend/common"
	"backend/config"
	"backend/middleware"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func InitRoutes() *gin.Engine {
	gin.SetMode(config.Conf.System.Mode)

	// Create a route with default middleware:
	// Log and recovery middleware
	r := gin.Default()
	// Create a route without middleware:
	// r := gin.New()
	// r.Use(gin.Recovery())

	// Enable rate limiting middleware
	// Default is to fill one token every 50 milliseconds, up to 200
	fillInterval := time.Duration(config.Conf.RateLimit.FillInterval)
	capacity := config.Conf.RateLimit.Capacity
	r.Use(middleware.RateLimitMiddleware(time.Millisecond*fillInterval, capacity))

	// Enable cors middleware
	r.Use(middleware.CORSMiddleware())

	// Enable operation log middleware
	r.Use(middleware.OperationLogMiddleware())

	// Initialize JWT auth middleware
	authMiddleware, err := middleware.InitAuth()

	if err != nil {
		common.Log.Panicf("Failed to initialize JWT middleware: %v", err)
		panic(fmt.Sprintf("Failed to initialize JWT middleware: %v", err))
	}

	// 路由分组
	apiGroup := r.Group("/" + config.Conf.System.UrlPathPrefix)

	// 注册路由
	InitBaseRoutes(apiGroup, authMiddleware)         // Register basic routes, no jwt auth middleware, no casbin middleware required
	InitUserRoutes(apiGroup, authMiddleware)         // Register user routes, jwt auth middleware, casbin auth middleware
	InitRoleRoutes(apiGroup, authMiddleware)         // Register role routes, jwt auth middleware, casbin auth middleware
	InitMenuRoutes(apiGroup, authMiddleware)         // Registration menu routes, jwt auth middleware, casbin auth middleware
	InitApiRoutes(apiGroup, authMiddleware)          // Register interface routes, jwt auth middleware, casbin auth middleware
	InitOperationLogRoutes(apiGroup, authMiddleware) // Register operation log routes, jwt auth middleware, casbin auth middleware

	common.Log.Info("Initial routing is completed!")
	return r
}
