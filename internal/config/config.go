package config

import (
	"fmt"
	"os"
	"time"
)

// Environment is the environment the app runs in, set by APP_ENV.
type Environment string

// The environments KepKep can run in.
const (
	EnvDevelopment Environment = "development"
	EnvStaging     Environment = "staging"
	EnvProduction  Environment = "production"
)

// Config holds all settings the app reads from environment variables.
type Config struct {
	AppEnv          Environment   // APP_ENV, default "development"
	HTTPAddr        string        // HTTP_ADDR, default ":8080"
	LogLevel        string        // LOG_LEVEL, default "info"
	ShutdownTimeout time.Duration // SHUTDOWN_TIMEOUT, default 10s
}

// Validate returns an error if e is not development, staging or production.
func (e Environment) Validate() error {
	switch e {
	case EnvDevelopment, EnvStaging, EnvProduction:
		return nil
	default:
		return fmt.Errorf("invalid APP_ENV %q: must be development, staging or production", e)
	}
}

// Load reads the config from environment variables, applies defaults for
// unset values and validates the result. It returns an error if any value is
// invalid, so the app can stop at startup instead of failing later.
func Load() (*Config, error) {
	shutdownTimeout, err := getEnvDuration("SHUTDOWN_TIMEOUT", time.Second*10)
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	cfg := &Config{
		AppEnv:          Environment(getEnvString("APP_ENV", "development")),
		HTTPAddr:        getEnvString("HTTP_ADDR", ":8080"),
		LogLevel:        getEnvString("LOG_LEVEL", "info"),
		ShutdownTimeout: shutdownTimeout,
	}

	if err := cfg.AppEnv.Validate(); err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	return cfg, nil
}

// getEnvString returns the value of key, or fallback if it is unset or empty.
func getEnvString(key, fallback string) string {
	val, ok := os.LookupEnv(key)
	if !ok || val == "" {
		return fallback
	}

	return val
}

// getEnvDuration parses key as a duration such as "10s" or "1m30s". It returns
// fallback if key is unset or empty, and an error if the value is invalid or not positive.
func getEnvDuration(key string, fallback time.Duration) (time.Duration, error) {
	val, ok := os.LookupEnv(key)
	if !ok || val == "" {
		return fallback, nil
	}

	d, err := time.ParseDuration(val)
	if err != nil {
		return 0, fmt.Errorf("invalid %s %q: %w", key, val, err)
	}

	if d <= 0 {
		return 0, fmt.Errorf("invalid %s %q: must be greater than 0", key, val)
	}

	return d, nil
}
