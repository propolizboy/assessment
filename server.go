package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func AuthMiddleware(key string, c echo.Context) (bool, error) {
	log.Println("[", key, "]")
	if key == "10, 2009" {
		return true, nil
	}
	return false, nil
}

func gethealthHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, "OK")
}

func main() {
	port := os.Getenv("Port")
	addr := ":" + port
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	//set auth
	e.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		AuthScheme: "November",
		Validator:  AuthMiddleware,
	}))
	e.GET("/healths", gethealthHandler)

	log.Println("Server started at:", port)
	go func() {
		if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// graceful shutdown
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
