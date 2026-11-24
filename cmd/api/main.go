// Command api runs the KepKep HTTP server (and, from Phase 6, the River workers).
package main

// Phase 0 — write this yourself. Suggested order:
//
//  1. Load config from env (internal/config).
//  2. Create the slog JSON logger (internal/platform/logger).
//  3. Create the Echo server, add request-ID middleware.
//  4. Add GET /healthz that returns 200 {"status":"ok"}.
//  5. Start the server in a goroutine.
//  6. Wait for SIGINT/SIGTERM (signal.NotifyContext), then call e.Shutdown
//     with a timeout so in-flight requests can finish.
func main() {
}
