package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"go.uber.org/zap"

	"github.com/Southclaws/cj/admin/actions"
	"github.com/Southclaws/cj/admin/apierrors"
	"github.com/Southclaws/cj/admin/middleware"
	"github.com/Southclaws/cj/admin/readmodel"
)

type actionInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Risk        string `json:"risk"`
}

func toActionInfo(a actions.Action) actionInfo {
	return actionInfo{Name: a.Name(), Description: a.Description(), Risk: string(a.Risk())}
}

func handleActionsList(runner *actions.Runner) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		list := runner.Registry().List()
		out := make([]actionInfo, 0, len(list))
		for _, a := range list {
			out = append(out, toActionInfo(a))
		}
		apierrors.WriteJSON(w, http.StatusOK, out)
	}
}

func handleActionGet(runner *actions.Runner) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestID := middleware.RequestIDFromContext(r.Context())

		action, ok := runner.Registry().Get(r.PathValue("name"))
		if !ok {
			apierrors.WriteError(w, http.StatusNotFound, apierrors.CodeActionNotFound,
				"No such action is registered.", requestID)
			return
		}
		apierrors.WriteJSON(w, http.StatusOK, toActionInfo(action))
	}
}

func readActionInput(r *http.Request) (json.RawMessage, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	if len(body) == 0 {
		return json.RawMessage("{}"), nil
	}
	return json.RawMessage(body), nil
}

func handleActionPreview(runner *actions.Runner) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestID := middleware.RequestIDFromContext(r.Context())

		input, err := readActionInput(r)
		if err != nil {
			apierrors.WriteError(w, http.StatusBadRequest, apierrors.CodeBadRequest,
				"The request body could not be read.", requestID)
			return
		}

		preview, err := runner.Preview(r.Context(), r.PathValue("name"), input)
		if err != nil {
			writeActionError(w, err, requestID)
			return
		}
		apierrors.WriteJSON(w, http.StatusOK, preview)
	}
}

func handleActionExecute(runner *actions.Runner) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestID := middleware.RequestIDFromContext(r.Context())

		input, err := readActionInput(r)
		if err != nil {
			apierrors.WriteError(w, http.StatusBadRequest, apierrors.CodeBadRequest,
				"The request body could not be read.", requestID)
			return
		}

		result, err := runner.Execute(r.Context(), r.PathValue("name"), input, requestID)
		if err != nil {
			writeActionError(w, err, requestID)
			return
		}
		apierrors.WriteJSON(w, http.StatusOK, result)
	}
}

func writeActionError(w http.ResponseWriter, err error, requestID string) {
	switch {
	case errors.Is(err, actions.ErrActionNotFound):
		apierrors.WriteError(w, http.StatusNotFound, apierrors.CodeActionNotFound,
			"No such action is registered.", requestID)
	case errors.Is(err, actions.ErrConfirmationRequired):
		apierrors.WriteError(w, http.StatusBadRequest, apierrors.CodeConfirmationRequired,
			"This action requires a matching confirmation phrase.", requestID)
	case errors.Is(err, actions.ErrInvalidInput):
		apierrors.WriteError(w, http.StatusBadRequest, apierrors.CodeInvalidActionInput,
			err.Error(), requestID)
	default:
		zap.L().Error("action execution failed", zap.Error(err), zap.String("request_id", requestID))
		apierrors.WriteError(w, http.StatusBadGateway, apierrors.CodeInternal,
			"The action could not be completed successfully.", requestID)
	}
}

func handleActionRunsList(provider *readmodel.ActionRunsProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestID := middleware.RequestIDFromContext(r.Context())

		runs, err := provider.List(0)
		if err != nil {
			zap.L().Error("failed to list action runs", zap.Error(err), zap.String("request_id", requestID))
			apierrors.WriteError(w, http.StatusInternalServerError, apierrors.CodeInternal,
				"Failed to load action history.", requestID)
			return
		}
		apierrors.WriteJSON(w, http.StatusOK, runs)
	}
}

func handleActionRunGet(provider *readmodel.ActionRunsProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestID := middleware.RequestIDFromContext(r.Context())

		run, found, err := provider.Get(r.PathValue("id"))
		if err != nil {
			zap.L().Error("failed to load action run", zap.Error(err), zap.String("request_id", requestID))
			apierrors.WriteError(w, http.StatusInternalServerError, apierrors.CodeInternal,
				"Failed to load action history.", requestID)
			return
		}
		if !found {
			apierrors.WriteError(w, http.StatusNotFound, apierrors.CodeNotFound,
				"No such action run is on record.", requestID)
			return
		}
		apierrors.WriteJSON(w, http.StatusOK, run)
	}
}
