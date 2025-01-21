package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"time"
)

type (
	Config struct {
		App      App      `yaml:"app"`
		Reporter Reporter `yaml:"reporter"`
		Client   Client   `yaml:"client"`
	}

	App struct {
		Name     string        `env-required:"true" yaml:"name" env:"APP_NAME"`
		Version  string        `env-required:"true" yaml:"version" env:"APP_VERSION"`
		Lifetime time.Duration `yaml:"lifetime" env:"APP_LIFETIME"`
	}

	Reporter struct {
		PollInterval time.Duration `yaml:"poll_interval" env:"REPORT_POLL_INTERVAL"`
	}

	Client struct {
		Interval      time.Duration `yaml:"interval" env:"CLIENT_INTERVAL"`
		Timeout       time.Duration `yaml:"timeout" env:"CLIENT_TIMEOUT"`
		Address       string        `yaml:"metric_server_address" env:"METRIC_SERVER_ADDRESS"`
		UserAgentName string        `yaml:"user_agent_name"`
	}
)

func MustLoad() Config {
	var cfg Config

	if err := cleanenv.ReadConfig("./config/config.yaml", &cfg); err != nil {
		panic(fmt.Errorf("yaml config error: %w", err))
	}

	if err := cleanenv.ReadConfig(".env", &cfg); err != nil {
		panic(fmt.Errorf("env config error: %w", err))
	}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic(err)
	}

	return cfg
}
