package api

import (
	"context"
	"fmt"
	"log/slog"
	nethttp "net/http"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/taikoxyz/taiko-mono/packages/blob-aggregator/pkg/http"
	"github.com/urfave/cli/v2"
)

type API struct {
	srv      *http.Server
	httpPort uint64
	wg       sync.WaitGroup
}

func (api *API) InitFromCli(ctx context.Context, c *cli.Context) error {
	cfg, err := NewConfigFromCliContext(c)
	if err != nil {
		return err
	}

	return InitFromConfig(ctx, api, cfg)
}

func InitFromConfig(ctx context.Context, api *API, cfg *Config) (err error) {
	q, err := cfg.OpenQueueFunc()
	if err != nil {
		return err
	}

	srv, err := http.NewServer(http.NewServerOpts{
		Queue:       q,
		Echo:        echo.New(),
		CorsOrigins: cfg.CORSOrigins,
	})
	if err != nil {
		return err
	}

	api.srv = srv
	api.httpPort = cfg.HTTPPort

	return nil
}

func (api *API) Name() string {
	return "api"
}

func (api *API) Close(ctx context.Context) {
	if err := api.srv.Shutdown(ctx); err != nil {
		slog.Error("srv shutdown", "error", err)
	}

	api.wg.Wait()
}

// nolint: funlen
func (api *API) Start() error {
	go func() {
		if err := api.srv.Start(fmt.Sprintf(":%v", api.httpPort)); err != nethttp.ErrServerClosed {
			slog.Error("http srv start", "error", err.Error())
		}
	}()

	return nil
}
