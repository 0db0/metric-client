package http

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"metric-client/config"
	"metric-client/internal/models"
	"net/http"
	"strconv"
	"time"
)

const updateEndpoint = "/updates"

type Client struct {
	interval       time.Duration
	httpClient     *http.Client
	serviceAddress string
	userAgentName  string
}

func NewClient(cfg config.Config) *Client {
	return &Client{
		interval: cfg.Client.Interval,
		httpClient: &http.Client{
			Timeout: cfg.Client.Timeout,
		},
		serviceAddress: cfg.Client.Address,
		userAgentName:  cfg.Client.UserAgentName + "/" + cfg.App.Version,
	}
}

func (c Client) SendMetrics(ctx context.Context, ch <-chan []models.Metric) {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case metric, ok := <-ch:
			if !ok {
				return
			}

			go func(metric []models.Metric) {
				<-ticker.C

				req, err := c.createRequest(metric)
				if err != nil {
					return
				}

				resp, err := c.httpClient.Do(req)
				if err != nil {
					log.Printf("error requesting metrics: %v", err)
					return
				}

				err = resp.Body.Close()
				if err != nil {
					return
				}

				if resp.StatusCode != http.StatusOK {
					return
				}
			}(metric)
		}
	}
}

func (c Client) createRequest(metrics []models.Metric) (*http.Request, error) {
	body, err := json.Marshal(metrics)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost, c.serviceAddress+updateEndpoint, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	request.Header.Set("User-Agent", c.userAgentName)
	request.Header.Set("Content-Type", "application/json; charset=utf-8")
	request.Header.Add("Content-Length", strconv.Itoa(len(body)))
	request.Header.Set("Accept", "application/json; charset=utf-8")

	return request, nil
}
