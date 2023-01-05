package handler

import "github.com/labstack/echo"

func AuthMiddleware(key string, c echo.Context) (bool, error) {
	if key == "10, 2009" {
		return true, nil
	}
	return false, nil
}
