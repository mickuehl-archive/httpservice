package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/txsvc/httpservice/pkg/api"
	"github.com/txsvc/httpservice/pkg/httpserver"
	"github.com/txsvc/stdlib/pkg/env"
)

func setup() *echo.Echo {
	// create a new router instance
	e := echo.New()

	// add and configure the middlewares
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.DefaultCORSConfig))

	// TODO add your own endpoints here
	e.GET("/", api.DefaultEndpoint)

	return e
}

func shutdown(*echo.Echo) {
	// TODO: implement your own stuff here
}

func init() {
	// only needed if deployed to GCP
	if !env.Exists("PROJECT_ID") {
		log.Fatal("missing environment variable 'PROJECT_ID'")
	}
}

func main() {
	service, err := httpserver.New(setup, shutdown, nil)
	if err != nil {
		log.Fatal(err)
	}
	service.StartBlocking()
}
