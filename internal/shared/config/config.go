// Package config reads every setting the service needs, once, at startup.
package config

import (
	"fmt"
	"os"
	"strings"
)

// Config is flat and strict on purpose: a service that starts without knowing
// where its database is will answer 500 to everything, which is a worse failure
// than not starting at all.
type Config struct {
	Port               string
	PostgresDSN        string
	KafkaBrokers       []string
	GameServersBaseURL string
	AccountsBaseURL    string
	AccountsAPIKey     string
}

func Load() (Config, error) {
	settings := Config{Port: valueOr("PORT", "8082")}

	var err error

	if settings.PostgresDSN, err = required("POSTGRES_DSN"); err != nil {
		return Config{}, err
	}

	if settings.KafkaBrokers, err = list("KAFKA_BROKERS"); err != nil {
		return Config{}, err
	}

	if settings.GameServersBaseURL, err = required("GAMESERVERS_BASE_URL"); err != nil {
		return Config{}, err
	}

	if settings.AccountsBaseURL, err = required("ACCOUNTS_BASE_URL"); err != nil {
		return Config{}, err
	}

	if settings.AccountsAPIKey, err = required("ACCOUNTS_API_KEY"); err != nil {
		return Config{}, err
	}

	return settings, nil
}

func required(key string) (string, error) {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return "", fmt.Errorf("%s was not configured", key)
	}
	return value, nil
}

func valueOr(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func list(key string) ([]string, error) {
	var out []string
	for _, part := range strings.Split(os.Getenv(key), ",") {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			out = append(out, trimmed)
		}
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("%s was not configured", key)
	}
	return out, nil
}
