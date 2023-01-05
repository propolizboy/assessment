package main

import (
	"net/http"
	"os"

	"github.com/propolizboy/assessment/handler"
	"github.com/propolizboy/assessment/router"
)

func main() {
	e := router.Setup()
	handler.SetupRoute(e)

	addr := ":" + os.Getenv("Port")
	go func() {
		if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	router.SetupGracefulShutdown(e)
}
