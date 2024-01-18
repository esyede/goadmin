package middleware

import (
	"backend/common"
	"backend/config"
	"backend/repository"
	"backend/response"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

var checkLock sync.Mutex

// Casbin middleware, RBAC-based permission access control model
func CasbinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ur := repository.NewUserRepository()
		user, err := ur.GetCurrentUser(c)
		if err != nil {
			response.Response(c, 401, 401, nil, "User is not logged in")
			c.Abort()
			return
		}
		if user.Status != 1 {
			response.Response(c, 401, 401, nil, "The current user has been disabled")
			c.Abort()
			return
		}
		// Get all roles of the user
		roles := user.Roles
		// Get the keywords of all the user's roles that are not disabled
		var subs []string
		for _, role := range roles {
			if role.Status == 1 {
				subs = append(subs, role.Keyword)
			}
		}
		// Get the request path URL
		// obj := strings.Replace(c.Request.URL.Path, "/"+config.Conf.System.UrlPathPrefix, "", 1)
		obj := strings.TrimPrefix(c.FullPath(), "/"+config.Conf.System.UrlPathPrefix)
		// 获取请求方式
		act := c.Request.Method

		isPass := check(subs, obj, act)
		if !isPass {
			response.Response(c, 401, 401, nil, "Permission denied")
			c.Abort()
			return
		}

		c.Next()
	}
}

func check(subs []string, obj string, act string) bool {
	// Only one request is allowed to perform verification at the same time, otherwise the verification may fail.
	checkLock.Lock()
	defer checkLock.Unlock()
	isPass := false
	for _, sub := range subs {
		pass, _ := common.CasbinEnforcer.Enforce(sub, obj, act)
		if pass {
			isPass = true
			break
		}
	}
	return isPass
}
