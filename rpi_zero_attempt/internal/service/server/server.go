package server

import (
	"context"
	_ "embed"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	"github.com/chack93/dmx_light_ctrl/internal/domain"
	"github.com/chack93/dmx_light_ctrl/internal/service/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
)

// @Version 1.0.0
// @Title dmx_light_ctrl rest api
// @Description send dmx messages to control lights
// @LicenseName GPLv3
// @Server http://raspberrypi.local local_network
// @Security AuthorizationHeader read write
// @SecurityScheme AuthorizationHeader http bearer Input your token

var log = logger.Get()

type Server struct {
	echo *echo.Echo
}

var server *Server

func Get() *Server {
	return server
}

func New() *Server {
	server = &Server{}
	return server
}

//go:embed swagger/swagger_gen.yaml
var swaggerYaml []byte

//go:embed swagger/index.html
var swaggerHtml []byte

func (srv *Server) Init(wg *sync.WaitGroup) error {
	srv.echo = echo.New()
	srv.echo.HideBanner = true
	srv.echo.HidePort = true
	srv.echo.Use(middleware.Logger())
	srv.echo.Use(middleware.Recover())

	baseURL := "/api"
	apiGroup := srv.echo.Group(baseURL)
	apiGroup.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, struct {
			Status string
		}{"ok"})
	})
	apiGroup.GET("/doc", func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "doc/")
	})
	apiGroup.GET("/doc", func(c echo.Context) error {
		return c.HTMLBlob(http.StatusOK, swaggerHtml)
	})
	apiGroup.GET("/doc/swagger.yaml", func(c echo.Context) error {
		return c.HTMLBlob(http.StatusOK, swaggerYaml)
	})
	domain.RegisterHandlers(srv.echo, baseURL)

	address := fmt.Sprintf("%s:%s", viper.GetString("server.host"), viper.GetString("server.port"))
	go func() {
		if err := srv.echo.Start(address); err != nil && err != http.ErrServerClosed {
			log.Errorf("server start failed, err: %v", err)
			wg.Done()
		}
	}()

	defer func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt)
		<-quit
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.echo.Shutdown(ctx); err != nil {
			log.Errorf("server shutdown failed, err: %v", err)
		}
		log.Info("server shutdown")
		wg.Done()
	}()

	for _, el := range srv.echo.Routes() {
		lastSlash := strings.LastIndex(el.Name, "/")
		domainHandler := el.Name[lastSlash:]
		log.Infof("%6s %s -> %s", el.Method, el.Path, domainHandler)
	}
	log.Infof("http server started on %s", address)
	return nil
}
