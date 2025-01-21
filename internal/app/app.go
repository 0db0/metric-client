package app

import (
	"context"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"metric-client/config"
	"metric-client/internal/adapters/http"
	"metric-client/internal/pkg/logger"
	"metric-client/internal/services/reporter"
	"os"
	"os/signal"
	"syscall"
)

func Run(cfg config.Config) {
	log := logger.New()

	ctx, cancel := context.WithTimeout(context.Background(), cfg.App.Lifetime)
	defer cancel()

	log.Info("start metrics crawling")

	r := reporter.New(cfg)
	c := http.NewClient(cfg)

	metrics := r.GetMetrics(ctx)
	c.SendMetrics(ctx, metrics)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case <-ctx.Done():
		log.Info("stop metrics crawling", ctx.Err())
	case s := <-interrupt:
		log.Info(fmt.Sprintf("client interrupt by signal %s", s.String()))
		cancel()
	}

	exitCode := 0
	defer func() {
		if err := recover(); err != nil {
			log.Error("client shutdown due to panic", err)

			exitCode = 1
		}

		os.Exit(exitCode)
	}()

	log.Info("client shutdown")
}
