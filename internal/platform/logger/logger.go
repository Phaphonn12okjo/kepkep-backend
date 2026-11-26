package logger

import (
	"io"
	"log/slog"
)

// New returns a logger that writes JSON lines to w.
// Messages below level are dropped.
func New(w io.Writer, level slog.Level) *slog.Logger {
	opts := &slog.HandlerOptions{
		Level: level,
	}

	handler := slog.NewJSONHandler(w, opts)

	return slog.New(handler)
}
