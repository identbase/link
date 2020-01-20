package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	identity "github.com/identbase/link/pkg/matrix/federation"
	"github.com/identbase/link/pkg/store"
	"github.com/identbase/serv/pkg/server"
)

/*
Config provides the standard arguments that can be provided at execution either
through ENV variables or command line arguments. */
type Config struct {
	Host     string
	Port     string
	NeedHelp bool
	Debug    bool
}

var config Config
var s *server.Server

func help() {
	fmt.Println("Identbase")
	fmt.Println()
	fmt.Println("Flags:")
	flag.PrintDefaults()
}

func init() {
	var defaults Config
	var ok bool

	if defaults.Host, ok = os.LookupEnv("HOST"); !ok {
		defaults.Host = "localhost"
	}

	flag.StringVar(&config.Host, "host", defaults.Host, "IP to serve application traffic on")

	if defaults.Port, ok = os.LookupEnv("PORT"); !ok {
		defaults.Port = "8000"
	}

	flag.StringVar(&config.Port, "port", defaults.Port, "Port to serve traffic on")

	// No need to include these in the defaults.
	flag.BoolVar(&config.NeedHelp, "h", false, "Show help text")

	flag.BoolVar(&config.Debug, "v", false, "Verbose output")
}

func main() {
	flag.Parse()

	if config.NeedHelp {
		help()

		return
	}

	s = server.New(server.Config{
		HideBanner: true,
		HidePort:   true,
		Debug:      config.Debug,
	})

	s.Use(middleware.Logger())
	s.Use(middleware.Recover())
	s.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderXRequestedWith,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
		},
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodOptions,
		},
	}))

	db := store.New(s.Logger())
	s.Load(s.NewModule("/_matrix/federation", &identity.Matrix{Database: db}))

	s.Logger().Infof("Launching identbase service on %s:%s", config.Host, config.Port)
	s.Logger().Fatal(s.Start(&http.Server{
		Addr:         fmt.Sprintf("%s:%s", config.Host, config.Port),
		ReadTimeout:  2 * time.Minute,
		WriteTimeout: 10 * time.Minute,
	}))
}
