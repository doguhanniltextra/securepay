// Package logger provides a context-aware slog.Handler that enriches every
// log record with OpenTelemetry trace correlation fields (trace_id, span_id,
// trace_flags). All log output is JSON-formatted for easy ingestion by
// centralised logging systems such as ELK or Grafana Loki.
package logger

import (
	"context"
	"log/slog"
	"os"

	"go.opentelemetry.io/otel/trace"
)

// otelHandler wraps an inner slog.Handler and, on every Handle call,
// extracts the active span from the context and injects the W3C
// trace-correlation fields into the log record before forwarding it.
type otelHandler struct {
	inner slog.Handler
}

// Enabled delegates level filtering to the inner handler.
func (h *otelHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.inner.Enabled(ctx, level)
}

// Handle enriches the record with trace fields then forwards it.
func (h *otelHandler) Handle(ctx context.Context, r slog.Record) error {
	if span := trace.SpanFromContext(ctx); span.SpanContext().IsValid() {
		sc := span.SpanContext()
		r.AddAttrs(
			slog.String("trace_id", sc.TraceID().String()),
			slog.String("span_id", sc.SpanID().String()),
			slog.String("trace_flags", sc.TraceFlags().String()),
		)
	}
	return h.inner.Handle(ctx, r)
}

// WithAttrs returns a new handler whose inner handler has the given attrs.
func (h *otelHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &otelHandler{inner: h.inner.WithAttrs(attrs)}
}

// WithGroup returns a new handler whose inner handler opens the given group.
func (h *otelHandler) WithGroup(name string) slog.Handler {
	return &otelHandler{inner: h.inner.WithGroup(name)}
}

// New creates a new JSON slog.Logger that automatically injects
// OpenTelemetry trace_id / span_id fields from the request context.
//
// Usage:
//
//	logger := logger.New()
//	slog.SetDefault(logger)
//
//	// Inside a handler, use the context-aware variants:
//	slog.InfoContext(ctx, "payment processed", "payment_id", id)
func New() *slog.Logger {
	jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	return slog.New(&otelHandler{inner: jsonHandler})
}
