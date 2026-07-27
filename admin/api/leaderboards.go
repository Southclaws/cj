package api

import (
	"net/http"
	"strconv"

	"github.com/Southclaws/cj/admin/apierrors"
	"github.com/Southclaws/cj/admin/middleware"
	"github.com/Southclaws/cj/admin/readmodel"
)

func handleLeaderboardsTopMessages(provider *readmodel.LeaderboardsProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestID := middleware.RequestIDFromContext(r.Context())
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

		top, err := provider.TopMessages(limit)
		if err != nil {
			writeDataError(w, err, requestID)
			return
		}
		apierrors.WriteJSON(w, http.StatusOK, top)
	}
}

func handleLeaderboardsTopReactions(provider *readmodel.LeaderboardsProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestID := middleware.RequestIDFromContext(r.Context())
		query := r.URL.Query()
		limit, _ := strconv.Atoi(query.Get("limit"))

		top, err := provider.TopReactions(limit, query.Get("reaction"))
		if err != nil {
			writeDataError(w, err, requestID)
			return
		}
		apierrors.WriteJSON(w, http.StatusOK, top)
	}
}
