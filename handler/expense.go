package handler

import (
	"github.com/labstack/echo"
	"github.com/propolizboy/assessment/expense"
)

func (h *Handler) CreateExpenseHandler(c echo.Context) error {
	return expense.Create(c, h.DB)
}

func (h *Handler) GetExpenseByIDHandler(c echo.Context) error {
	return expense.GetById(c, h.DB)
}
