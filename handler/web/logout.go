package web

import (
	"errors"

	"github.com/labstack/echo/v4"
	"github.com/myOmikron/echotools/middleware"
)

func (w *Wrapper) Logout(c echo.Context) error {
	if err := middleware.Logout(w.DB, c); err != nil {
		if !errors.Is(err, middleware.ErrCookieNotFound) {
			//TODO: Display error
			return err
		}
	}

	return c.Redirect(302, "/login")
}
