package main

import (
	"github.com/esyede/goadmin/backend/common"
	"github.com/esyede/goadmin/backend/config"
	"github.com/esyede/goadmin/backend/middleware"
	"github.com/esyede/goadmin/backend/repository"
	"github.com/esyede/goadmin/backend/routes"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	config.InitConfig()
	common.InitLogger()
	common.InitMysql()
	common.InitCasbinEnforcer()
	common.InitValidate()
	common.InitData()

	logRepository := repository.NewOperationLogRepository()

	for i := 0; i < 3; i++ {
		go logRepository.SaveOperationLogChannel(middleware.OperationLogChan)
	}

	r := routes.InitRoutes()
	host := "localhost"
	port := config.Conf.System.Port
	srv := &http.Server{Addr: fmt.Sprintf("%s:%d", host, port), Handler: r}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			common.Log.Fatalf("Listen: %s\n", err)
		}
	}()

	common.Log.Info(fmt.Sprintf("Server is running at %s:%d/%s", host, port, config.Conf.System.UrlPathPrefix))

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	common.Log.Info("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		common.Log.Fatal("Server forced to shutdown:", err)
	}

	common.Log.Info("Server exiting!")
}
