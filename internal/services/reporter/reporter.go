package reporter

import (
	"context"
	"metric-client/config"
	"metric-client/internal/models"
	"runtime"
	"time"
)

type Reporter struct {
	pollInterval time.Duration
}

func New(cfg config.Config) *Reporter {
	return &Reporter{
		pollInterval: cfg.Reporter.PollInterval,
	}
}

func (r Reporter) GetMetrics(ctx context.Context) <-chan []models.Metric {
	var metricChan = make(chan []models.Metric, 1)

	go func() {
		ticker := time.NewTicker(r.pollInterval)

		defer close(metricChan)
		defer ticker.Stop()

		for range ticker.C {
			var stats runtime.MemStats
			runtime.ReadMemStats(&stats)

			metrics := collectMetrics(stats)

			select {
			case metricChan <- metrics:
			case <-ctx.Done():
				return
			}
		}
	}()

	return metricChan
}

func collectMetrics(stats runtime.MemStats) []models.Metric {
	var metrics []models.Metric

	metrics = append(
		metrics,
		models.Metric{
			ID:    "Alloc",
			MType: models.TypeGauge,
			Value: float64(stats.Alloc),
		},
		models.Metric{
			ID:    "BuckHashSys",
			MType: models.TypeGauge,
			Value: float64(stats.BuckHashSys),
		},
		models.Metric{
			ID:    "Frees",
			MType: models.TypeGauge,
			Value: float64(stats.Frees),
		},
		models.Metric{
			ID:    "GCCPUFraction",
			MType: models.TypeGauge,
			Value: stats.GCCPUFraction,
		},
		models.Metric{
			ID:    "GCSys",
			MType: models.TypeGauge,
			Value: float64(stats.Sys),
		},
		models.Metric{
			ID:    "HeapAlloc",
			MType: models.TypeGauge,
			Value: float64(stats.HeapAlloc),
		},
	)

	return metrics
}
