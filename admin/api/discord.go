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

func writeDiscordError(w http.ResponseWriter, r *http.Request, err error) {
	requestID := middleware.RequestIDFromContext(r.Context())

	switch {
	case errors.Is(err, readmodel.ErrDiscordNotReady):
		apierrors.WriteError(w, http.StatusServiceUnavailable, apierrors.CodeInternal,
			"CJ is not currently connected to Discord.", requestID)
		return
	case errors.Is(err, readmodel.ErrMemberCacheNotReady):
		apierrors.WriteError(w, http.StatusServiceUnavailable, apierrors.CodeInternal,
			"The member list is still loading. Try again shortly.", requestID)
		return
	}

	zap.L().Error("discord readmodel request failed", zap.Error(err), zap.String("request_id", requestID))
	apierrors.WriteError(w, http.StatusBadGateway, apierrors.CodeInternal,
		"Failed to fetch data from Discord.", requestID)
}

func handleDiscordGuild(provider *readmodel.DiscordProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		guild, err := provider.Guild()
		if err != nil {
			writeDiscordError(w, r, err)
			return
		}
		apierrors.WriteJSON(w, http.StatusOK, guild)
	}
}

func handleDiscordChannels(provider *readmodel.DiscordProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		channels, err := provider.Channels()
		if err != nil {
			writeDiscordError(w, r, err)
			return
		}
		apierrors.WriteJSON(w, http.StatusOK, channels)
	}
}

func handleDiscordRoles(provider *readmodel.DiscordProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		roles, err := provider.Roles()
		if err != nil {
			writeDiscordError(w, r, err)
			return
		}
		apierrors.WriteJSON(w, http.StatusOK, roles)
	}
}

func handleDiscordMembers(provider *readmodel.DiscordProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		limit, _ := strconv.Atoi(query.Get("limit"))
		offset, _ := strconv.Atoi(query.Get("offset"))

		page, err := provider.Members(query.Get("q"), limit, offset)
		if err != nil {
			writeDiscordError(w, r, err)
			return
		}
		apierrors.WriteJSON(w, http.StatusOK, page)
	}
}

func handleDiscordMember(provider *readmodel.DiscordProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		member, err := provider.Member(r.PathValue("id"))
		if err != nil {
			writeDiscordError(w, r, err)
			return
		}
		apierrors.WriteJSON(w, http.StatusOK, member)
	}
}

func handleDiscordMemberPermissions(provider *readmodel.DiscordProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		permissions, err := provider.MemberPermissions(r.PathValue("id"))
		if err != nil {
			writeDiscordError(w, r, err)
			return
		}
		apierrors.WriteJSON(w, http.StatusOK, permissions)
	}
}

func handleDiscordAuditLog(provider *readmodel.DiscordProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		entries, err := provider.AuditLog(0)
		if err != nil {
			writeDiscordError(w, r, err)
			return
		}
		apierrors.WriteJSON(w, http.StatusOK, entries)
	}
}
