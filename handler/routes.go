package handler

import (
	"net/http"

	"github.com/labstack/echo"
)

func (h *Handler) gethealthHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, "OK")
}

func (h *Handler) SetupRoute(e *echo.Echo) {
	e.GET("/healths", h.gethealthHandler)
}
