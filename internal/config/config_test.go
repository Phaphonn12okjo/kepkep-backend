package config

import (
	"log/slog"
	"strings"
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name        string
		env         map[string]string
		expected    Config
		expectedErr bool
		errContains string
	}{
		{
			name: "defaults when nothing is set",
			env:  map[string]string{},
			expected: Config{
				AppEnv:          EnvDevelopment,
				HTTPAddr:        ":8080",
				LogLevel:        slog.LevelInfo,
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
				AppEnv:          EnvProduction,
				HTTPAddr:        ":9090",
				LogLevel:        slog.LevelDebug,
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
			errContains: "invalid APP_ENV",
		},
		{
			name: "bad duration returns an error",
			env: map[string]string{
				"SHUTDOWN_TIMEOUT": "thirteen secondssss",
			},
			expectedErr: true,
			errContains: "time: invalid duration",
		},
		{
			name: "negative duration returns an error",
			env: map[string]string{
				"SHUTDOWN_TIMEOUT": "-5s",
			},
			expectedErr: true,
			errContains: "must be greater than 0",
		},
		{
			name: "invalid LOG_LEVEL returns an error",
			env: map[string]string{
				"LOG_LEVEL": "debucking",
			},
			expectedErr: true,
			errContains: "unknown name",
		},
		{
			name: "accept WARN capital letters",
			env: map[string]string{
				"LOG_LEVEL": "WARN",
			},
			expected: Config{
				AppEnv:          EnvDevelopment,
				HTTPAddr:        ":8080",
				LogLevel:        slog.LevelWarn,
				ShutdownTimeout: time.Second * 10,
			},
			expectedErr: false,
		},
		{
			name: "invalid LOG_LEVEL = INFO+2 returns an error",
			env: map[string]string{
				"LOG_LEVEL": "INFO+2",
			},
			expectedErr: true,
			errContains: "must be debug, info, warn or error",
		},
		{
			name: "blank values use defaults",
			env: map[string]string{
				"HTTP_ADDR":        "  ",
				"LOG_LEVEL":        "  ",
				"SHUTDOWN_TIMEOUT": "  ",
			},
			expected: Config{
				AppEnv:          EnvDevelopment,
				HTTPAddr:        ":8080",
				LogLevel:        slog.LevelInfo,
				ShutdownTimeout: time.Second * 10,
			},
			expectedErr: false,
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
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("error %q should contain %q", err, tt.errContains)
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
