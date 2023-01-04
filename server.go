package main

import (
	"log"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	port := os.Getenv("Port")
	addr := ":" + port
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	log.Println("Server started at:", port)
	log.Fatal(e.Start(addr))
}
