package loggingmdw

import (
	"bufio"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	requestMdw "github.com/guluzadehh/go_chat/internal/http/middlewares/requestmdw"
)

type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode   int
	bytesWritten int
}

func (w *responseWriterWrapper) WriteHeader(statusCode int) {
	if w.statusCode == http.StatusOK {
		w.ResponseWriter.WriteHeader(statusCode)
		w.statusCode = statusCode
	}
}

func (w *responseWriterWrapper) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := w.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, http.ErrNotSupported
	}
	return hijacker.Hijack()
}

func (w *responseWriterWrapper) Write(data []byte) (int, error) {
	written, err := w.ResponseWriter.Write(data)
	w.bytesWritten = written
	return written, err
}

func newResponseWriterWrapper(w http.ResponseWriter) *responseWriterWrapper {
	return &responseWriterWrapper{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
		bytesWritten:   0,
	}
}

func LogRequests(log *slog.Logger) mux.MiddlewareFunc {
	log = log.With("component", "middleware/logging")

	return func(next http.Handler) http.Handler {
		log.Info("logger middleware enabled")

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := log.With(
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("remote_addr", r.RemoteAddr),
				slog.String("user_agent", r.UserAgent()),
				slog.String("request_id", requestMdw.GetReqId(r)),
			)

			ww := newResponseWriterWrapper(w)

			t1 := time.Now()
			defer func() {
				log.Info("request comleted",
					slog.Int("status", ww.statusCode),
					slog.Int("bytes", ww.bytesWritten),
					slog.String("duration", time.Since(t1).String()),
				)
			}()

			next.ServeHTTP(ww, r)
		})
	}
}
