package middleware

import (
	"cmp"
	"log/slog"
	"net/http"
	"strconv"
)

func LogRequest(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.InfoContext(r.Context(), "Start")
		writer := &responseMiddleware{
			ResponseWriter: w,
		}
		h.ServeHTTP(writer, r)
		slog.InfoContext(r.Context(), strconv.Itoa(cmp.Or(writer.status, 200)))
	})
}

type responseMiddleware struct {
	http.ResponseWriter
	status int
	done bool
}

func (w * responseMiddleware) WriteHeader(status int) {
	w.done = true
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *responseMiddleware) Write(b []byte) (int, error) {
	w.done = true
	return w.ResponseWriter.Write(b)
}
