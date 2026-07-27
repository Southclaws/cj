package middleware

import (
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/Southclaws/cj/admin/apierrors"
)

func hostOnly(hostport string) string {
	host, _, err := net.SplitHostPort(hostport)
	if err != nil {
		return strings.ToLower(hostport)
	}
	return strings.ToLower(host)
}

func isSafeMethod(method string) bool {
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodOptions:
		return true
	default:
		return false
	}
}

func NetworkGuard(allowedHosts []string) func(http.Handler) http.Handler {
	allowed := make(map[string]struct{}, len(allowedHosts))
	for _, h := range allowedHosts {
		allowed[hostOnly(h)] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := RequestIDFromContext(r.Context())

			if _, ok := allowed[hostOnly(r.Host)]; !ok {
				apierrors.WriteError(w, http.StatusForbidden, apierrors.CodeHostNotAllowed,
					"This host is not permitted to access the dashboard.", id)
				return
			}

			if !isSafeMethod(r.Method) {
				origin := r.Header.Get("Origin")
				if origin == "" {
					apierrors.WriteError(w, http.StatusForbidden, apierrors.CodeOriginNotAllowed,
						"A matching Origin header is required for this request.", id)
					return
				}
				parsed, err := url.Parse(origin)
				if err != nil {
					apierrors.WriteError(w, http.StatusForbidden, apierrors.CodeOriginNotAllowed,
						"The Origin header could not be parsed.", id)
					return
				}
				if _, ok := allowed[hostOnly(parsed.Host)]; !ok {
					apierrors.WriteError(w, http.StatusForbidden, apierrors.CodeOriginNotAllowed,
						"This origin is not permitted to make this request.", id)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
