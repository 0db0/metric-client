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
			Value: toFloat64Pointer(stats.Alloc),
		},
		models.Metric{
			ID:    "BuckHashSys",
			MType: models.TypeGauge,
			Value: toFloat64Pointer(stats.BuckHashSys),
		},
		models.Metric{
			ID:    "Frees",
			MType: models.TypeGauge,
			Value: toFloat64Pointer(stats.Frees),
		},
		models.Metric{
			ID:    "GCCPUFraction",
			MType: models.TypeGauge,
			Value: toFloat64Pointer(stats.GCCPUFraction),
		},
		models.Metric{
			ID:    "GCSys",
			MType: models.TypeGauge,
			Value: toFloat64Pointer(stats.GCSys),
		},
		models.Metric{
			ID:    "HeapAlloc",
			MType: models.TypeGauge,
			Value: toFloat64Pointer(stats.HeapAlloc),
		},
		models.Metric{
			ID:    "HeapIdle",
			MType: models.TypeGauge,
			Value: toFloat64Pointer(stats.HeapIdle),
		},
		models.Metric{
			ID:    "HeapInuse",
			MType: models.TypeGauge,
			Value: toFloat64Pointer(stats.HeapInuse),
		},
		models.Metric{
			ID:    "HeapObjects",
			MType: models.TypeGauge,
			Value: toFloat64Pointer(stats.HeapObjects),
		},
		models.Metric{
			ID:    "HeapReleased",
			MType: models.TypeGauge,
			Value: toFloat64Pointer(stats.HeapReleased),
		},
		models.Metric{
			ID:    "HeapSys",
			MType: models.TypeGauge,
			Value: toFloat64Pointer(stats.HeapSys),
		},
		models.Metric{
			ID:    "LastGC",
			MType: models.TypeGauge,
			Value: toFloat64Pointer(stats.LastGC),
		},
		models.Metric{
			ID:    "Lookups",
			MType: models.TypeGauge,
			Value: toFloat64Pointer(stats.Lookups),
		}, models.Metric{
			ID:    "MCacheInuse",
			MType: models.TypeGauge,
			Value: toFloat64Pointer(stats.MCacheInuse),
		}, models.Metric{
			ID:    "MCacheSys",
			MType: models.TypeGauge,
			Value: toFloat64Pointer(stats.MCacheSys),
		}, models.Metric{
			ID:    "MSpanInuse",
			MType: models.TypeGauge,
			Value: toFloat64Pointer(stats.MSpanInuse),
		},
		models.Metric{
			ID:    "MSpanSys",
			MType: models.TypeGauge,
			Value: toFloat64Pointer(stats.MSpanSys),
		},
		models.Metric{
			ID:    "Mallocs",
			MType: models.TypeGauge,
			Value: toFloat64Pointer(stats.Mallocs),
		},
		models.Metric{
			ID:    "NextGC",
			MType: models.TypeGauge,
			Value: toFloat64Pointer(stats.NextGC),
		},
		models.Metric{
			ID:    "NumForcedGC",
			MType: models.TypeGauge,
			Value: toFloat64Pointer(stats.NumForcedGC),
		},
		models.Metric{
			ID:    "NumGC",
			MType: models.TypeGauge,
			Value: toFloat64Pointer(stats.NumGC),
		},
		models.Metric{
			ID:    "OtherSys",
			MType: models.TypeGauge,
			Value: toFloat64Pointer(stats.OtherSys),
		},
		models.Metric{
			ID:    "PauseTotalNs",
			MType: models.TypeGauge,
			Value: toFloat64Pointer(stats.PauseTotalNs),
		},
		models.Metric{
			ID:    "StackInuse",
			MType: models.TypeGauge,
			Value: toFloat64Pointer(stats.StackInuse),
		},
		models.Metric{
			ID:    "StackSys",
			MType: models.TypeGauge,
			Value: toFloat64Pointer(stats.StackSys),
		},
		models.Metric{
			ID:    "Sys",
			MType: models.TypeGauge,
			Value: toFloat64Pointer(stats.Sys),
		},
		models.Metric{
			ID:    "TotalAlloc",
			MType: models.TypeGauge,
			Value: toFloat64Pointer(stats.TotalAlloc),
		},
		models.Metric{
			ID:    "RandomValue",
			MType: models.TypeGauge,
			Value: toFloat64Pointer(rand.Float64()),
		},
	)

	atomic.AddInt64(counter, 1)
	pollCountDelta := atomic.LoadInt64(counter)

	metrics = append(metrics, models.Metric{
		ID:    "PollCount",
		MType: models.TypeCounter,
		Delta: &pollCountDelta,
	})

	return metrics
}

func toFloat64Pointer(stat any) *float64 {
	switch value := stat.(type) {
	case float64:
		return &value
	case uint64:
		f := float64(value)
		return &f
	case uint32:
		f := float64(value)
		return &f
	}

	return nil
}
