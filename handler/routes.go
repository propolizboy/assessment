package handler

import (
	"github.com/labstack/echo"
)

func (h *Handler) SetupRoute(e *echo.Echo) {
	e.GET("/healths", h.gethealthHandler)
	e.POST("/expenses", h.CreateExpenseHandler)
	e.GET("/expenses/:id", h.GetExpenseByIDHandler)
	e.PUT("/expenses/:id", h.UpdateExpenseByIDHandler)
	e.GET("/expenses", h.GETAllExpenseHandler)
}
