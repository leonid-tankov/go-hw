package internalhttp

import (
	"net/http"
	"time"

	"github.com/leonid-tankov/go-hw/hw12_13_14_15_calendar/internal/repository"
)

type Middleware struct {
	logger repository.Logger
}

func NewMiddleware(logger repository.Logger) *Middleware {
	return &Middleware{
		logger: logger,
	}
}

func (m *Middleware) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := newLoggingResponseWriter(w)
		next.ServeHTTP(lrw, r)
		m.logger.Info("%s [%s] \"%s %s %s\" %d %0.fs %s",
			r.RemoteAddr,
			start.Format(time.RFC3339),
			r.Method,
			r.URL.Path,
			r.Proto,
			lrw.statusCode,
			time.Since(start).Seconds(),
			r.UserAgent())
	})
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}
