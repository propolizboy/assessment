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

func (h *Handler) UpdateExpenseByIDHandler(c echo.Context) error {
	return expense.UpdateByID(c, h.DB)
}

func (h *Handler) GETAllExpenseHandler(c echo.Context) error {
	return expense.GetAll(c, h.DB)
}
