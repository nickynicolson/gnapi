package web

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const withLogs = true

//go:embed static
var static embed.FS

// Run starts HTTP/1 service for scientific names verification.
func Run(domain string, port int) {
	var err error
	log.Printf("Starting the HTTP API server on port %d.", port)
	e := echo.New()
	e.Use(middleware.Gzip())
	e.Use(middleware.CORS())
	if withLogs {
		e.Use(middleware.Logger())
	}

	e.Renderer, err = NewTemplate()
	if err != nil {
		e.Logger.Fatal(err)
	}

	e.GET("/", home(domain))
	e.GET("/gnparser", gnparser(domain))
	e.GET("/gnames", gnames(domain))
	e.GET("/gnmatcher", gnmatcher(domain))
	e.GET("/gnfinder", gnfinder(domain))

	fs := http.FileServer(http.FS(static))
	e.GET("/static/*", echo.WrapHandler(fs))

	addr := fmt.Sprintf(":%d", port)
	s := &http.Server{
		Addr:         addr,
		ReadTimeout:  5 * time.Minute,
		WriteTimeout: 5 * time.Minute,
	}
	e.Logger.Fatal(e.StartServer(s))
}

type Data struct {
	Domain, DocJSON string
}

func home(domain string) func(echo.Context) error {
	return func(c echo.Context) error {
		data := Data{Domain: domain}
		return c.Render(http.StatusOK, "home", data)
	}
}

func gnfinder(domain string) func(echo.Context) error {
	return func(c echo.Context) error {
		data := Data{Domain: domain, DocJSON: "static/gnfinder/openapi.json"}
		return c.Render(http.StatusOK, "api", data)
	}
}

func gnparser(domain string) func(echo.Context) error {
	return func(c echo.Context) error {
		data := Data{Domain: domain, DocJSON: "static/gnparser/openapi.json"}
		return c.Render(http.StatusOK, "api", data)
	}
}

func gnmatcher(domain string) func(echo.Context) error {
	return func(c echo.Context) error {
		data := Data{Domain: domain, DocJSON: "static/gnmatcher/openapi.json"}
		return c.Render(http.StatusOK, "api", data)
	}
}

func gnames(domain string) func(echo.Context) error {
	return func(c echo.Context) error {
		data := Data{Domain: domain, DocJSON: "static/gnames/openapi.json"}
		return c.Render(http.StatusOK, "api", data)
	}
}
