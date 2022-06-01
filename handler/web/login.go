package web

import (
	"github.com/labstack/echo/v4"
)

type LoginData struct {
	PageTitle string
}

func (w *Wrapper) Login(c echo.Context) error {
	return c.Render(200, "login", &LoginData{
		PageTitle: "Login - " + w.Config.General.SiteName,
	})
}
