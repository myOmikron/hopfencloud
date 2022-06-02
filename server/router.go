package server

import (
	"github.com/myOmikron/hopfencloud/handler/web"
	"github.com/myOmikron/hopfencloud/models/conf"

	"github.com/labstack/echo/v4"
	"github.com/myOmikron/echotools/middleware"
	"github.com/myOmikron/echotools/worker"
	"gorm.io/gorm"
)

func loginRequired(f func(c echo.Context) error) echo.HandlerFunc {
	return func(c echo.Context) error {
		sessionContext, err := middleware.GetSessionContext(c)
		if err != nil {
			return err
		}

		if !sessionContext.IsAuthenticated() {
			return c.Redirect(302, "/login?redirect_to="+c.Path())
		}

		return f(c)
	}
}

func unauthenticatedOnly(f func(c echo.Context) error) echo.HandlerFunc {
	return func(c echo.Context) error {
		sessionContext, err := middleware.GetSessionContext(c)
		if err != nil {
			return err
		}

		if sessionContext.IsAuthenticated() {
			return c.Redirect(302, "/")
		}

		return f(c)
	}
}

func defineRoutes(e *echo.Echo, db *gorm.DB, config *conf.Config, wp worker.Pool) {
	webWrapper := web.Wrapper{
		DB:         db,
		Config:     config,
		WorkerPool: wp,
	}
	e.GET("/login", unauthenticatedOnly(webWrapper.LoginGet))
	e.POST("/login", unauthenticatedOnly(webWrapper.LoginPost))

	e.GET("/logout", webWrapper.Logout)

	e.Static("/static/", "static/")
}
