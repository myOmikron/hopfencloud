package handler

import "github.com/labstack/echo/v4"

func FileHandler(c echo.Context) error {
	return c.Render(200, "files", nil)
}
