package middleware

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/Southclaws/cj/admin/apierrors"
)

func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				id := RequestIDFromContext(r.Context())
				if id == "" {
					id = r.Header.Get("X-Request-Id")
				}
				if id == "" {
					id = newRequestID()
				}
				zap.L().Error("dashboard handler panic",
					zap.Any("panic", rec),
					zap.String("request_id", id),
					zap.String("path", r.URL.Path))
				apierrors.WriteError(w, http.StatusInternalServerError, apierrors.CodeInternal,
					"An internal error occurred.", id)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
