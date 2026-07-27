package api

import (
	"io/fs"
	"net/http"

	"github.com/Southclaws/cj/admin/apierrors"
	"github.com/Southclaws/cj/admin/middleware"
	"github.com/Southclaws/cj/web"
)

func staticHandler() http.Handler {
	if !web.Embedded {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			apierrors.WriteError(w, http.StatusNotFound, apierrors.CodeNotFound,
				"The dashboard frontend is not embedded in this build.",
				middleware.RequestIDFromContext(r.Context()))
		})
	}

	sub, err := fs.Sub(web.Dist, "dist")
	if err != nil {
		panic(err)
	}
	fileServer := http.FileServer(http.FS(sub))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if len(path) > 0 && path[0] == '/' {
			path = path[1:]
		}
		if path != "" {
			if _, statErr := fs.Stat(sub, path); statErr != nil {
				r = r.Clone(r.Context())
				r.URL.Path = "/"
			}
		}
		fileServer.ServeHTTP(w, r)
	})
}

func handleAPINotFound(w http.ResponseWriter, r *http.Request) {
	apierrors.WriteError(w, http.StatusNotFound, apierrors.CodeNotFound,
		"The requested resource does not exist.", middleware.RequestIDFromContext(r.Context()))
}
