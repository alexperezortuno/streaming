package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
	written    int64
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.written += int64(n)
	return n, err
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(rw, r)

		attrs := []slog.Attr{
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Int("status", rw.statusCode),
			slog.Duration("duration", time.Since(start)),
			slog.Int64("bytes", rw.written),
			slog.String("remote", r.RemoteAddr),
			slog.String("user_agent", r.UserAgent()),
		}

		if rw.statusCode >= 500 {
			slog.LogAttrs(nil, slog.LevelError, "request completed", attrs...)
		} else if rw.statusCode >= 400 {
			slog.LogAttrs(nil, slog.LevelWarn, "request completed", attrs...)
		} else {
			slog.LogAttrs(nil, slog.LevelInfo, "request completed", attrs...)
		}
	})
}
