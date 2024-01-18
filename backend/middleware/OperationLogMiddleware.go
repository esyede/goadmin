package middleware

import (
	"github.com/esyede/goadmin/backend/config"
	"github.com/esyede/goadmin/backend/model"
	"github.com/esyede/goadmin/backend/repository"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Operation log channel
var OperationLogChan = make(chan *model.OperationLog, 30)

func OperationLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Starting time
		startTime := time.Now()

		// handle the request
		c.Next()

		// End Time
		endTime := time.Now()

		// Execution time
		timeCost := endTime.Sub(startTime).Milliseconds()

		// Get the current logged in user
		var username string
		ctxUser, exists := c.Get("user")
		if !exists {
			username = "Not logged in"
		}
		user, ok := ctxUser.(model.User)
		if !ok {
			username = "Not logged in"
		}
		username = user.Username

		// Get access path
		path := strings.TrimPrefix(c.FullPath(), "/"+config.Conf.System.UrlPathPrefix)
		// Request method
		method := c.Request.Method

		// Get interface description
		apiRepository := repository.NewApiRepository()
		apiDesc, _ := apiRepository.GetApiDescByPath(path, method)

		operationLog := model.OperationLog{
			Username:   username,
			Ip:         c.ClientIP(),
			IpLocation: "",
			Method:     method,
			Path:       path,
			Desc:       apiDesc,
			Status:     c.Writer.Status(),
			StartTime:  startTime,
			TimeCost:   timeCost,
			// UserAgent:  c.Request.UserAgent(),
		}

		// It is best to send the logs to rabbitmq or kafka
		// Here it is sent to the channel and three goroutine processing are enabled.
		OperationLogChan <- &operationLog
	}
}
