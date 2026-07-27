package api

import (
	"net/http"

	"github.com/Southclaws/cj/admin/apierrors"
	"github.com/Southclaws/cj/admin/middleware"
	"github.com/Southclaws/cj/admin/readmodel"
)

func handleProfile(provider *readmodel.ProfileProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestID := middleware.RequestIDFromContext(r.Context())

		profile, found, err := provider.Get(r.PathValue("id"))
		if err != nil {
			writeDataError(w, err, requestID)
			return
		}
		if !found {
			apierrors.WriteError(w, http.StatusNotFound, apierrors.CodeNotFound,
				"No such user is known to CJ or the guild.", requestID)
			return
		}
		apierrors.WriteJSON(w, http.StatusOK, profile)
	}
}
