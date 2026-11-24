package config

import (
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name        string
		env         map[string]string
		expected    Config
		expectedErr bool
	}{
		{
			name: "defaults when nothing is set",
			env:  map[string]string{},
			expected: Config{
				AppEnv:          Environment("development"),
				HTTPAddr:        ":8080",
				LogLevel:        "info",
				ShutdownTimeout: time.Second * 10,
			},
			expectedErr: false,
		},
		{
			name: "all env vars are set",
			env: map[string]string{
				"APP_ENV":          string(EnvProduction),
				"HTTP_ADDR":        ":9090",
				"LOG_LEVEL":        "debug",
				"SHUTDOWN_TIMEOUT": "30s",
			},
			expected: Config{
				AppEnv:          Environment("production"),
				HTTPAddr:        ":9090",
				LogLevel:        "debug",
				ShutdownTimeout: time.Second * 30,
			},
			expectedErr: false,
		},
		{
			name: "invalid APP_ENV returns an error",
			env: map[string]string{
				"APP_ENV": "producktion",
			},
			expectedErr: true,
		},
		{
			name: "bad duration returns an error",
			env: map[string]string{
				"SHUTDOWN_TIMEOUT": "thirteen secondssss",
			},
			expectedErr: true,
		},
		{
			name: "negative duration returns an error",
			env: map[string]string{
				"SHUTDOWN_TIMEOUT": "-5s",
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, key := range []string{"APP_ENV", "HTTP_ADDR", "LOG_LEVEL", "SHUTDOWN_TIMEOUT"} {
				t.Setenv(key, "")
			}

			for key, value := range tt.env {
				t.Setenv(key, value)
			}

			result, err := Load()

			if tt.expectedErr {
				if err == nil {
					t.Fatalf("expected an error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if *result != tt.expected {
				t.Errorf("got %+v, expected %+v", *result, tt.expected)
			}
		})
	}
}
