package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var e = createMux()

func main() {
	e.GET("/", hello)
	e.GET("/public", public)
	e.GET("/private", private)

	e.Logger.Fatal(e.Start(":8082"))
}

func createMux() *echo.Echo {
	e := echo.New()

	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(middleware.Gzip())

	return e
}

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!\n")
}

func public(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, public!\n")
}

func private(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, private!\n")
}
