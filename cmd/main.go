package main

import (
	"metric-client/config"
	"metric-client/internal/app"
)

func main() {
	cfg := config.MustLoad()
	app.Run(cfg)
}
