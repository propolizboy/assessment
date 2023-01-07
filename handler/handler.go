package handler

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo"
)

type Handler struct {
	DB *sql.DB
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{db}
}

func (h *Handler) gethealthHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, "Hello Expenses")
}
