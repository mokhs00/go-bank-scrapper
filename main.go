package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go-job-scrapper/scrapper"
	"os"
	"strings"
)

const fileName string = "jobs.csv"

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", handleHome)
	e.POST("/scrape", handleScrape)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}

// Handler
func handleHome(c echo.Context) error {
	return c.File("statics/home.html")
}

func handleScrape(c echo.Context) error {
	defer os.Remove(fileName)
	term := scrapper.CleanString(strings.ToLower(c.FormValue("term")))
	scrapper.Scrape(term)

	return c.Attachment(fileName, fileName)
}
