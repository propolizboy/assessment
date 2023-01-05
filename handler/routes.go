package handler

import (
	"net/http"

	"github.com/labstack/echo"
)

func gethealthHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, "OK")
}

func SetupRoute(e *echo.Echo) {
	e.GET("/healths", gethealthHandler)
}
