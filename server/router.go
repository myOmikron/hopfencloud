package server

import (
	"github.com/myOmikron/hopfencloud/handler/web"
	"github.com/myOmikron/hopfencloud/models/conf"

	"github.com/labstack/echo/v4"
	"github.com/myOmikron/echotools/worker"
	"gorm.io/gorm"
)

func defineRoutes(e *echo.Echo, db *gorm.DB, config *conf.Config, wp worker.Pool) {
	e.GET("/", handler.FileHandler)

	e.Static("/static/", "static/")
}
