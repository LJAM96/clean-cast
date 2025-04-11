package app

import (
	"bytes"
	"html/template"
	"ikoyhn/podcast-sponsorblock/internal/models"
	"ikoyhn/podcast-sponsorblock/internal/services"
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
		podcasts, err := services.GetAllPodcasts()
		if err != nil {
			return c.HTML(http.StatusOK, "<h1>No podcasts found</h1>")
		}

		tmpl, err := template.ParseFiles("../../public/pages/podcasts/podcasts.html")
		if err != nil {
			return err
		}

		data := struct {
			Podcasts []models.Podcast
		}{
			Podcasts: podcasts,
		}

		buf := new(bytes.Buffer)
		err = tmpl.Execute(buf, data)
		if err != nil {
			return err
		}

		return c.HTML(http.StatusOK, buf.String())
	})

	e.GET("/popup/:id", func(c echo.Context) error {
		id := c.Param("id")
		episodes, err := services.GetPodcastEpisodesByPodcastId(id)
		if err != nil {
			return c.HTML(http.StatusOK, "<h1>No episodes found</h1>")
		}

		tmpl, err := template.ParseFiles("../../public/pages/podcasts/popup.html")
		if err != nil {
			return err
		}

		data := struct {
			Episodes []models.PodcastEpisode
		}{
			Episodes: episodes,
		}

		buf := new(bytes.Buffer)
		err = tmpl.Execute(buf, data)
		if err != nil {
			return err
		}

		return c.HTML(http.StatusOK, buf.String())
	})

	e.GET("/downloads", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "<h1>Downloads</h1>")
	})
}
