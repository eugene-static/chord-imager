package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

func Logging(next http.HandlerFunc, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		log.Info("new request",
			//			slog.String("host", r.Host),
			slog.String("method", r.Method))
		next(w, r)
		log.Info("sending response",
			//			slog.String("host", r.Host),
			slog.String("method", r.Method),
			//			slog.String("uri", r.RequestURI),
			slog.String("response time", time.Since(t).String()))
	}
}
