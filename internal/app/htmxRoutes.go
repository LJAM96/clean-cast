package app

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func registerHtmxRoutes(e *echo.Echo) {
	e.Static("/", "../../public")
	e.GET("/", func(c echo.Context) error {
		return c.File("../../public/index.html")
	})
	e.GET("/dashboard", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "<h1>Dashboard</h1>")
	})

	e.GET("/podcasts", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "<h1>Podcasts</h1>")
	})

	e.GET("/downloads", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "<h1>Downloads</h1>")
	})
}
