package middleware

import (
	"log/slog"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

// TracingMiddleware adds OpenTelemetry tracing to the request.
func TracingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// 1. Extract TraceContext from headers
		propagator := otel.GetTextMapPropagator()
		ctx = propagator.Extract(ctx, propagation.HeaderCarrier(r.Header))

		// 2. Start Span
		tracer := otel.Tracer("api-gateway")
		spanName := fmtSpanName(r)
		ctx, span := tracer.Start(ctx, spanName, trace.WithSpanKind(trace.SpanKindServer))
		defer span.End()

		// 3. Add trace_id to response header
		traceID := span.SpanContext().TraceID().String()
		w.Header().Set("X-Trace-ID", traceID)

		// 4. Log with trace_id
		slog.Info("Incoming Request",
			"method", r.Method,
			"path", r.URL.Path,
			"trace_id", traceID,
		)

		// 5. Pass context to next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func fmtSpanName(r *http.Request) string {
	return "api-gateway." + r.Method + " " + r.URL.Path
}
