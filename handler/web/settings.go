package web

import (
	"fmt"
	"github.com/labstack/echo/v4"
)

func (w *Wrapper) SettingsMail(c echo.Context) error {
	return c.Render(200, "admin/settings_mail", nil)
}

func (w *Wrapper) SettingsGeneral(c echo.Context) error {
	fmt.Sprintf("")

	return c.Render(200, "admin/settings_general", nil)
}
