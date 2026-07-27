package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/Southclaws/cj/admin/apierrors"
	"github.com/Southclaws/cj/admin/middleware"
	"github.com/Southclaws/cj/admin/readmodel"
	"github.com/Southclaws/cj/storage"
)

var snowflakePattern = regexp.MustCompile(`^[0-9]{1,20}$`)

func validateGuildSettings(s storage.GuildSettings) error {
	fields := map[string]string{
		"searchMessageChannelId": s.SearchMessageChannelID,
		"errorReportChannelId":   s.ErrorReportChannelID,
		"leaderboardChannelId":   s.LeaderboardChannelID,
		"adsChannelId":           s.AdsChannelID,
		"ltfChannelId":           s.LTFChannelID,
	}
	for name, value := range fields {
		if value != "" && !snowflakePattern.MatchString(value) {
			return fmt.Errorf("%s must be a numeric Discord ID", name)
		}
	}
	for _, id := range s.LTFUserIDs {
		if !snowflakePattern.MatchString(id) {
			return fmt.Errorf("ltfUserIds must all be numeric Discord IDs")
		}
	}
	return nil
}

func handleGetSettings(provider *readmodel.SettingsProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestID := middleware.RequestIDFromContext(r.Context())

		settings, err := provider.Get()
		if err != nil {
			apierrors.WriteError(w, http.StatusInternalServerError, apierrors.CodeInternal,
				"Failed to load settings.", requestID)
			return
		}
		apierrors.WriteJSON(w, http.StatusOK, settings)
	}
}

func handlePutSettings(provider *readmodel.SettingsProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestID := middleware.RequestIDFromContext(r.Context())

		var input storage.GuildSettings
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			apierrors.WriteError(w, http.StatusBadRequest, apierrors.CodeBadRequest,
				"The request body is not valid JSON.", requestID)
			return
		}

		if err := validateGuildSettings(input); err != nil {
			apierrors.WriteError(w, http.StatusBadRequest, apierrors.CodeBadRequest, err.Error(), requestID)
			return
		}

		saved, err := provider.Set(input, requestID)
		if err != nil {
			apierrors.WriteError(w, http.StatusInternalServerError, apierrors.CodeInternal,
				"Failed to save settings.", requestID)
			return
		}

		apierrors.WriteJSON(w, http.StatusOK, saved)
	}
}
