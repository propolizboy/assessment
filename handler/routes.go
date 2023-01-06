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
	e.POST("/expenses", h.CreateExpenseHandler)
	e.GET("/expenses/:id", h.GetExpenseByIDHandler)
	e.PUT("/expenses/:id", h.UpdateExpenseByIDHandler)
}
