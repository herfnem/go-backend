package middleware

import (
	"context"
	"net/http"
	"sync"
	"time"

	"learn/internal/api/response"
)

type timeoutWriter struct {
	http.ResponseWriter
	mu          sync.Mutex
	timedOut    bool
	wroteHeader bool
}

func (tw *timeoutWriter) WriteHeader(statusCode int) {
	tw.mu.Lock()
	defer tw.mu.Unlock()
	if tw.timedOut {
		return
	}
	tw.wroteHeader = true
	tw.ResponseWriter.WriteHeader(statusCode)
}

func (tw *timeoutWriter) Write(p []byte) (int, error) {
	tw.mu.Lock()
	defer tw.mu.Unlock()
	if tw.timedOut {
		return 0, nil
	}
	if !tw.wroteHeader {
		tw.wroteHeader = true
		tw.ResponseWriter.WriteHeader(http.StatusOK)
	}
	return tw.ResponseWriter.Write(p)
}

func Timeout(timeout time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if timeout <= 0 {
				next.ServeHTTP(w, r)
				return
			}

			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer cancel()

			tw := &timeoutWriter{ResponseWriter: w}
			done := make(chan struct{})

			go func() {
				next.ServeHTTP(tw, r.WithContext(ctx))
				close(done)
			}()

			select {
			case <-done:
				return
			case <-ctx.Done():
				tw.mu.Lock()
				tw.timedOut = true
				if !tw.wroteHeader {
					tw.wroteHeader = true
					response.WriteError(tw.ResponseWriter, http.StatusGatewayTimeout, "request timeout")
				}
				tw.mu.Unlock()
				return
			}
		})
	}
}
