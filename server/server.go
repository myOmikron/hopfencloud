package server

import (
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	textTemplate "text/template"
	"time"

	"github.com/myOmikron/hopfencloud/models/conf"
	"github.com/myOmikron/hopfencloud/models/db"

	"github.com/labstack/echo/v4"
	mw "github.com/labstack/echo/v4/middleware"
	"github.com/myOmikron/echotools/color"
	"github.com/myOmikron/echotools/execution"
	"github.com/myOmikron/echotools/middleware"
	"github.com/myOmikron/echotools/utilitymodels"
	"github.com/myOmikron/echotools/worker"
	"github.com/pelletier/go-toml"
)

func StartServer(configPath string, isReloading bool) {
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

	// Initialize database
	database := initializeDatabase(config)

	// Setup directory structure if necessary
	if err := initializeDirStructure(database, config); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

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
		templates: template.Must(template.ParseFS(os.DirFS("templates/html/"), "*/*.gohtml")),
	}
	e.Renderer = renderer
	mailTemplates := textTemplate.Must(textTemplate.ParseFS(os.DirFS("templates/mail/"), "*.gotxt"))

	// Middleware
	e.Use(mw.LoggerWithConfig(mw.LoggerConfig{
		Format: "${time_rfc3339} :: ${status} ${method} ${uri} :: ${latency_human} ${error}\n",
	}))
	e.Use(mw.Recover())
	e.Use(mw.BodyLimit(config.Server.MaxFileUpload))

	duration := time.Hour * 24
	e.Use(middleware.Session(database, &middleware.SessionConfig{
		CookieName: "sessionid",
		CookieAge:  &duration,
	}))

	allowedHosts := make([]middleware.AllowedHost, 0)
	for _, host := range config.Server.AllowedHosts {
		hostPortCombination := host.Host
		if host.Port != 80 && host.Port != 443 {
			hostPortCombination = net.JoinHostPort(host.Host, strconv.Itoa(int(host.Port)))
		}
		allowedHosts = append(allowedHosts, middleware.AllowedHost{
			Host:  hostPortCombination,
			Https: host.Https,
		})
	}
	secConfig := &middleware.SecurityConfig{
		AllowedHosts:            allowedHosts,
		UseForwardedProtoHeader: config.Server.UseForwardedProtoHeader,
	}
	e.Use(middleware.Security(secConfig))

	// Allowed authenticated backends
	middleware.RegisterAuthProvider(utilitymodels.GetLocalUser(database))
	middleware.RegisterAuthProvider(utilitymodels.GetLDAPUser(database))

	// Settings
	var settings db.Settings
	settingsReloadFunc := getSettingsReloadFunc(&settings, database)

	// Router
	defineRoutes(e, database, config, wp, mailTemplates, settingsReloadFunc, &settings)

	// Start RPC listener
	var cliSock net.Listener
	go initializeRPC(&cliSock, database, config, wp, mailTemplates, settingsReloadFunc, &settings, isReloading)

	// Start database cleaner
	go cleanupDatabase(database)

	execution.SignalStart(e, config.Server.ListenAddress, &execution.Config{
		ReloadFunc: func() {
			if cliSock != nil {
				cliSock.Close()
			}
			StartServer(configPath, true)
		},
		StopFunc: func() {
			if cliSock != nil {
				cliSock.Close()
			}
		},
		TerminateFunc: func() {
			if cliSock != nil {
				cliSock.Close()
			}
		},
	})
}
