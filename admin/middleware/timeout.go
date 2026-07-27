package middleware

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/Southclaws/cj/admin/apierrors"
)

type timeoutWriter struct {
	mu       sync.Mutex
	w        http.ResponseWriter
	timedOut bool
}

func (tw *timeoutWriter) Header() http.Header {
	return tw.w.Header()
}

func (tw *timeoutWriter) WriteHeader(status int) {
	tw.mu.Lock()
	defer tw.mu.Unlock()
	if tw.timedOut {
		return
	}
	tw.w.WriteHeader(status)
}

func (tw *timeoutWriter) Write(b []byte) (int, error) {
	tw.mu.Lock()
	defer tw.mu.Unlock()
	if tw.timedOut {
		return len(b), nil
	}
	return tw.w.Write(b)
}

func (tw *timeoutWriter) markTimedOut() {
	tw.mu.Lock()
	defer tw.mu.Unlock()
	tw.timedOut = true
}

func Timeout(d time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), d)
			defer cancel()
			r = r.WithContext(ctx)

			tw := &timeoutWriter{w: w}
			done := make(chan struct{})
			panicChan := make(chan any, 1)

			go func() {
				defer func() {
					if rec := recover(); rec != nil {
						panicChan <- rec
					}
				}()
				next.ServeHTTP(tw, r)
				close(done)
			}()

			select {
			case rec := <-panicChan:
				panic(rec)
			case <-done:
			case <-ctx.Done():
				tw.markTimedOut()
				id := RequestIDFromContext(ctx)
				apierrors.WriteError(w, http.StatusServiceUnavailable, apierrors.CodeRequestTimeout,
					"The request took too long to process.", id)
			}
		})
	}
}
