package router

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/propolizboy/assessment/handler"
)

func Setup() *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		AuthScheme: "November",
		Validator:  handler.AuthMiddleware,
	}))
	return e
}

func SetupGracefulShutdown(e *echo.Echo) {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	errSignal := <-shutdown
	switch errSignal {
	case os.Interrupt:
		fmt.Println("Server Interrupt...")
	case syscall.SIGTERM:
		fmt.Println("Server Sigterm...")
	}
	fmt.Println("App is shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal("Error shutting down: %v\n", err)
	}
}
