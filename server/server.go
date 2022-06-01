package server

import (
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"io/ioutil"
	"os"
	"time"

	"github.com/myOmikron/hopfencloud/models/conf"

	"github.com/labstack/echo/v4"
	mw "github.com/labstack/echo/v4/middleware"
	"github.com/myOmikron/echotools/color"
	"github.com/myOmikron/echotools/execution"
	"github.com/myOmikron/echotools/middleware"
	"github.com/myOmikron/echotools/utilitymodels"
	"github.com/myOmikron/echotools/worker"
	"github.com/pelletier/go-toml"
)

func StartServer(configPath string) {
	config := &conf.Config{}

	if configBytes, err := ioutil.ReadFile(configPath); errors.Is(err, fs.ErrNotExist) {
		color.Printf(color.RED, "Config was not found at %s\n", configPath)
		b, _ := toml.Marshal(config)
		fmt.Print(string(b))
		os.Exit(1)
	} else {
		if err := toml.Unmarshal(configBytes, config); err != nil {
			panic(err)
		}
	}

	// Check for valid config values
	if err := config.CheckConfig(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	db := initializeDatabase(config)

	// Web server
	e := echo.New()
	e.HideBanner = true

	// Worker pool
	wp := worker.NewPool(&worker.PoolConfig{
		NumWorker: 10,
		QueueSize: 100,
	})
	wp.Start()

	// Template rendering
	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("templates/*.gohtml")),
	}
	e.Renderer = renderer

	// Middleware
	e.Use(mw.Logger())
	e.Use(mw.Recover())

	duration := time.Hour * 24
	e.Use(middleware.Session(db, &middleware.SessionConfig{
		CookieName: "sessionid",
		CookieAge:  &duration,
	}))

	allowedHosts := make([]middleware.AllowedHost, 0)
	for _, host := range config.Server.AllowedHosts {
		allowedHosts = append(allowedHosts, middleware.AllowedHost{
			Host:  host.Host,
			Https: host.Https,
		})
	}
	secConfig := &middleware.SecurityConfig{
		AllowedHosts:            allowedHosts,
		UseForwardedProtoHeader: config.Server.UseForwardedProtoHeader,
	}
	e.Use(middleware.Security(secConfig))

	// Router
	defineRoutes(e, db, config, wp)

	execution.SignalStart(e, config.Server.ListenAddress, &execution.Config{
		ReloadFunc: func() {
			StartServer(configPath)
		},
		StopFunc: func() {

		},
		TerminateFunc: func() {

		},
	})
}
