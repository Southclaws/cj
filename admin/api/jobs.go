package api

import (
	"net/http"

	"github.com/Southclaws/cj/admin/apierrors"
	"github.com/Southclaws/cj/admin/readmodel"
)

func handleJobs(provider *readmodel.JobsProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apierrors.WriteJSON(w, http.StatusOK, provider.List())
	}
}
