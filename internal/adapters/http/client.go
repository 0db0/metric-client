package http

import (
	"bytes"
	"context"
	"encoding/json"
	"metric-client/config"
	"metric-client/internal/models"
	"metric-client/internal/pkg/logger"
	"net/http"
	"strconv"
	"time"
)

const updateEndpoint = "/v0.1/updates"

type Client struct {
	interval       time.Duration
	httpClient     *http.Client
	serviceAddress string
	userAgentName  string
	log            logger.Interface
}

func NewClient(cfg config.Config, log logger.Interface) *Client {
	return &Client{
		interval: cfg.Client.Interval,
		httpClient: &http.Client{
			Timeout: cfg.Client.Timeout,
		},
		serviceAddress: cfg.Client.Address,
		userAgentName:  cfg.Client.UserAgentName + "/" + cfg.App.Version,
		log:            log,
	}
}

func (c Client) SendMetrics(ctx context.Context, ch <-chan []models.Metric) {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case metrics, ok := <-ch:
			if !ok {
				return
			}

			go func(metrics []models.Metric) {
				<-ticker.C

				req, err := c.createRequest(metrics)
				if err != nil {
					c.log.Error("create request failed", err)
					return
				}

				resp, err := c.httpClient.Do(req)
				if err != nil {
					c.log.Error("http request failed", err)
					return
				}

				err = resp.Body.Close()
				if err != nil {
					c.log.Error("close response body failed", err)
					return
				}

				if resp.StatusCode != http.StatusOK {
					c.log.Error("http response status code invalid. Got code", resp.StatusCode)
					return
				}
			}(metrics)
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
