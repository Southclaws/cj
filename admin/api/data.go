package api

import (
	"errors"
	"net/http"
	"strconv"

	"go.uber.org/zap"

	"github.com/Southclaws/cj/admin/apierrors"
	"github.com/Southclaws/cj/admin/middleware"
	"github.com/Southclaws/cj/admin/readmodel"
)

func handleDataUsers(provider *readmodel.UsersProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestID := middleware.RequestIDFromContext(r.Context())
		query := r.URL.Query()

		limit, _ := strconv.Atoi(query.Get("limit"))
		offset, _ := strconv.Atoi(query.Get("offset"))

		page, err := provider.List(query.Get("q"), limit, offset)
		if err != nil {
			writeDataError(w, err, requestID)
			return
		}
		apierrors.WriteJSON(w, http.StatusOK, page)
	}
}

func handleDataUser(provider *readmodel.UsersProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestID := middleware.RequestIDFromContext(r.Context())

		user, found, err := provider.Get(r.PathValue("id"))
		if err != nil {
			writeDataError(w, err, requestID)
			return
		}
		if !found {
			apierrors.WriteError(w, http.StatusNotFound, apierrors.CodeNotFound,
				"No such user is on record.", requestID)
			return
		}
		apierrors.WriteJSON(w, http.StatusOK, user)
	}
}

func handleDataChat(provider *readmodel.ChatProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestID := middleware.RequestIDFromContext(r.Context())
		query := r.URL.Query()

		messages, err := provider.Messages(query.Get("userId"), query.Get("q"))
		if err != nil {
			if errors.Is(err, readmodel.ErrChatFilterRequired) {
				apierrors.WriteError(w, http.StatusBadRequest, apierrors.CodeBadRequest,
					"A userId or a search query is required.", requestID)
				return
			}
			writeDataError(w, err, requestID)
			return
		}
		apierrors.WriteJSON(w, http.StatusOK, messages)
	}
}

func handleDataMessage(provider *readmodel.ChatProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestID := middleware.RequestIDFromContext(r.Context())

		message, found, err := provider.GetByID(r.PathValue("id"))
		if err != nil {
			writeDataError(w, err, requestID)
			return
		}
		if !found {
			apierrors.WriteError(w, http.StatusNotFound, apierrors.CodeNotFound,
				"No such message is on record.", requestID)
			return
		}
		apierrors.WriteJSON(w, http.StatusOK, message)
	}
}

func writeDataError(w http.ResponseWriter, err error, requestID string) {
	zap.L().Error("data readmodel request failed", zap.Error(err), zap.String("request_id", requestID))
	apierrors.WriteError(w, http.StatusInternalServerError, apierrors.CodeInternal,
		"Failed to load data.", requestID)
}
