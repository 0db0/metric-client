package reporter

import (
	"context"
	"math/rand"
	"metric-client/config"
	"metric-client/internal/models"
	"metric-client/internal/pkg/logger"
	"runtime"
	"sync/atomic"
	"time"
)

type Reporter struct {
	pollInterval time.Duration
	log          logger.Interface
}

func New(cfg config.Config, log logger.Interface) *Reporter {
	return &Reporter{
		pollInterval: cfg.Reporter.PollInterval,
		log:          log,
	}
}

func (r Reporter) GetMetrics(ctx context.Context) <-chan []models.Metric {
	var metricChan = make(chan []models.Metric, 1)

	go func() {
		var counter int64
		ticker := time.NewTicker(r.pollInterval)

		defer close(metricChan)
		defer ticker.Stop()

		for range ticker.C {
			var stats runtime.MemStats
			runtime.ReadMemStats(&stats)

			metrics := collectMetrics(stats, &counter)

			select {
			case metricChan <- metrics:
			case <-ctx.Done():
				return
			}
		}

	}()

	return metricChan
}

func collectMetrics(stats runtime.MemStats, counter *int64) []models.Metric {
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
			Value: float64(stats.GCSys),
		},
		models.Metric{
			ID:    "HeapAlloc",
			MType: models.TypeGauge,
			Value: float64(stats.HeapAlloc),
		},
		models.Metric{
			ID:    "HeapIdle",
			MType: models.TypeGauge,
			Value: float64(stats.HeapIdle),
		},
		models.Metric{
			ID:    "HeapInuse",
			MType: models.TypeGauge,
			Value: float64(stats.HeapInuse),
		},
		models.Metric{
			ID:    "HeapObjects",
			MType: models.TypeGauge,
			Value: float64(stats.HeapObjects),
		},
		models.Metric{
			ID:    "HeapReleased",
			MType: models.TypeGauge,
			Value: float64(stats.HeapReleased),
		},
		models.Metric{
			ID:    "HeapSys",
			MType: models.TypeGauge,
			Value: float64(stats.HeapSys),
		},
		models.Metric{
			ID:    "LastGC",
			MType: models.TypeGauge,
			Value: float64(stats.LastGC),
		},
		models.Metric{
			ID:    "Lookups",
			MType: models.TypeGauge,
			Value: float64(stats.Lookups),
		}, models.Metric{
			ID:    "MCacheInuse",
			MType: models.TypeGauge,
			Value: float64(stats.MCacheInuse),
		}, models.Metric{
			ID:    "MCacheSys",
			MType: models.TypeGauge,
			Value: float64(stats.MCacheSys),
		}, models.Metric{
			ID:    "MSpanInuse",
			MType: models.TypeGauge,
			Value: float64(stats.MSpanInuse),
		},
		models.Metric{
			ID:    "MSpanSys",
			MType: models.TypeGauge,
			Value: float64(stats.MSpanSys),
		},
		models.Metric{
			ID:    "Mallocs",
			MType: models.TypeGauge,
			Value: float64(stats.Mallocs),
		},
		models.Metric{
			ID:    "NextGC",
			MType: models.TypeGauge,
			Value: float64(stats.NextGC),
		},
		models.Metric{
			ID:    "NumForcedGC",
			MType: models.TypeGauge,
			Value: float64(stats.NumForcedGC),
		},
		models.Metric{
			ID:    "NumGC",
			MType: models.TypeGauge,
			Value: float64(stats.NumGC),
		},
		models.Metric{
			ID:    "OtherSys",
			MType: models.TypeGauge,
			Value: float64(stats.OtherSys),
		},
		models.Metric{
			ID:    "PauseTotalNs",
			MType: models.TypeGauge,
			Value: float64(stats.PauseTotalNs),
		},
		models.Metric{
			ID:    "StackInuse",
			MType: models.TypeGauge,
			Value: float64(stats.StackInuse),
		},
		models.Metric{
			ID:    "StackSys",
			MType: models.TypeGauge,
			Value: float64(stats.StackSys),
		},
		models.Metric{
			ID:    "Sys",
			MType: models.TypeGauge,
			Value: float64(stats.Sys),
		},
		models.Metric{
			ID:    "TotalAlloc",
			MType: models.TypeGauge,
			Value: float64(stats.TotalAlloc),
		},
		models.Metric{
			ID:    "RandomValue",
			MType: models.TypeGauge,
			Value: rand.Float64(),
		},
	)

	atomic.AddInt64(counter, 1)
	metrics = append(metrics, models.Metric{
		ID:    "PollCount",
		MType: models.TypeCounter,
		Delta: atomic.LoadInt64(counter),
	})

	return metrics
}
