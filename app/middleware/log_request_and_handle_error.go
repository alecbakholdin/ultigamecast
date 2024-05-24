package middleware

import (
	"cmp"
	"log/slog"
	"net/http"
	"strconv"
)

func LogRequestAndHandleError(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.InfoContext(r.Context(), "Start")
		writer := &responseMiddleware{
			ResponseWriter: w,
		}
		h.ServeHTTP(writer, r)
		if status := cmp.Or(writer.status, 200); status < 400 {
			slog.InfoContext(r.Context(), strconv.Itoa(status))
		} else {
			slog.ErrorContext(r.Context(), strconv.Itoa(status))
		}
	})
}

type responseMiddleware struct {
	http.ResponseWriter
	status int
	err    string
	done   bool
}

func (w *responseMiddleware) WriteHeader(status int) {
	w.done = true
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *responseMiddleware) Write(b []byte) (int, error) {
	w.done = true
	if w.status >= 400 {
		w.err = string(b[:50])
	}
	return w.ResponseWriter.Write(b)
}
