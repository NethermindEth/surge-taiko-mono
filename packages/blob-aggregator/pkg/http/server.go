package http

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4/middleware"
	"github.com/patrickmn/go-cache"
	"github.com/taikoxyz/taiko-mono/packages/blob-aggregator/pkg/queue"

	echo "github.com/labstack/echo/v4"
)

type Server struct {
	queue queue.Queue
	echo  *echo.Echo
	cache *cache.Cache
}

type NewServerOpts struct {
	Queue       queue.Queue
	Echo        *echo.Echo
	CorsOrigins []string
}

func (opts NewServerOpts) Validate() error {
	if opts.Echo == nil {
		return ErrNoHTTPFramework
	}

	return nil
}

func NewServer(opts NewServerOpts) (*Server, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}

	cache := cache.New(5*time.Minute, 10*time.Minute)

	srv := &Server{
		queue: opts.Queue,
		echo:  opts.Echo,
		cache: cache,
	}

	corsOrigins := opts.CorsOrigins
	if corsOrigins == nil {
		corsOrigins = []string{"*"}
	}

	srv.configureMiddleware(corsOrigins)
	srv.configureRoutes()

	return srv, nil
}

// Start starts the HTTP server
func (srv *Server) Start(address string) error {
	return srv.echo.Start(address)
}

// Shutdown shuts down the HTTP server
func (srv *Server) Shutdown(ctx context.Context) error {
	return srv.echo.Shutdown(ctx)
}

// ServeHTTP implements the `http.Handler` interface which serves HTTP requests
func (srv *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	srv.echo.ServeHTTP(w, r)
}

// Health endpoints for probes
func (srv *Server) Health(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}

func (srv *Server) returnError(c echo.Context, statusCode int, err error) error {
	return c.JSON(statusCode, map[string]string{"error": err.Error()})
}

func LogSkipper(c echo.Context) bool {
	switch c.Request().URL.Path {
	case "/healthz":
		return true
	case "/metrics":
		return true
	default:
		return false
	}
}

func (srv *Server) configureMiddleware(corsOrigins []string) {
	srv.echo.Use(middleware.RequestID())

	srv.echo.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Skipper: LogSkipper,
		Format: `{"time":"${time_rfc3339_nano}","level":"INFO","message":{"id":"${id}","remote_ip":"${remote_ip}",` + //nolint:lll
			`"host":"${host}","method":"${method}","uri":"${uri}","user_agent":"${user_agent}",` + //nolint:lll
			`"response_status":${status},"error":"${error}","latency":${latency},"latency_human":"${latency_human}",` +
			`"bytes_in":${bytes_in},"bytes_out":${bytes_out}}}` + "\n",
		Output: os.Stdout,
	}))

	srv.echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: corsOrigins,
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowMethods: []string{http.MethodGet, http.MethodHead},
	}))
}
