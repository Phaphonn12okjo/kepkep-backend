package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"testing"
)

func TestNewWritesJSON(t *testing.T) {
	var buf bytes.Buffer
	logger := New(&buf, slog.LevelInfo)

	logger.Info("hello", "user_id", 42)

	var got map[string]any
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("output is not valid JSON: %v\noutput: %s", err, buf.String())
	}

	if got["msg"] != "hello" {
		t.Errorf("msg = %v, expected %q", got["msg"], "hello")
	}
	if got["level"] != "INFO" {
		t.Errorf("level = %v, expected %q", got["level"], "INFO")
	}
	if got["user_id"] != float64(42) {
		t.Errorf("user_id = %v, expected 42", got["user_id"])
	}
	if _, ok := got["time"]; !ok {
		t.Errorf("expected a time field, got none")
	}
}

func TestNewLevelFiltering(t *testing.T) {
	tests := []struct {
		name           string
		level          slog.Level // the logger's minimum level
		logAt          slog.Level // the level we log a message at
		expectedOutput bool       // should anything be written
	}{
		{name: "debug hidden at info", level: slog.LevelInfo, logAt: slog.LevelDebug, expectedOutput: false},
		{name: "info shown at info", level: slog.LevelInfo, logAt: slog.LevelInfo, expectedOutput: true},
		{name: "debug shown at debug", level: slog.LevelDebug, logAt: slog.LevelDebug, expectedOutput: true},
		{name: "info hidden at warn", level: slog.LevelWarn, logAt: slog.LevelInfo, expectedOutput: false},
		{name: "error shown at warn", level: slog.LevelWarn, logAt: slog.LevelError, expectedOutput: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := New(&buf, tt.level)

			logger.Log(context.Background(), tt.logAt, "test message")

			gotOutput := buf.Len() > 0
			if gotOutput != tt.expectedOutput {
				t.Errorf("logger at %s, logged at %s: output = %t, expected %t",
					tt.level, tt.logAt, gotOutput, tt.expectedOutput)
			}
		})
	}
}
