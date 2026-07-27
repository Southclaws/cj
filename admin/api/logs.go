package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Southclaws/cj/admin/apierrors"
	"github.com/Southclaws/cj/admin/logs"
	"github.com/Southclaws/cj/admin/middleware"
)

func handleLogsList(buffer *logs.Buffer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		limit := 200
		if raw := query.Get("limit"); raw != "" {
			if parsed, err := strconv.Atoi(raw); err == nil {
				limit = parsed
			}
		}

		records := buffer.Snapshot(logs.Filter{
			Level:         query.Get("level"),
			Component:     query.Get("component"),
			Query:         query.Get("q"),
			CorrelationID: query.Get("correlationId"),
			Limit:         limit,
		})

		apierrors.WriteJSON(w, http.StatusOK, records)
	}
}

func handleLogsStream(buffer *logs.Buffer, stream *logs.Stream) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestID := middleware.RequestIDFromContext(r.Context())

		flusher, ok := w.(http.Flusher)
		if !ok {
			apierrors.WriteError(w, http.StatusInternalServerError, apierrors.CodeInternal,
				"Streaming is not supported by this server.", requestID)
			return
		}

		sub := stream.Subscribe()
		defer stream.Unsubscribe(sub)

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.WriteHeader(http.StatusOK)

		backlog := buffer.Snapshot(logs.Filter{Limit: 50})
		for i := len(backlog) - 1; i >= 0; i-- {
			writeSSERecord(w, backlog[i])
		}
		flusher.Flush()

		ctx := r.Context()
		for {
			select {
			case <-ctx.Done():
				return
			case <-sub.Dropped():
				writeSSEComment(w, "records dropped, client too slow")
				flusher.Flush()
			case record, ok := <-sub.Records():
				if !ok {
					return
				}
				writeSSERecord(w, record)
				flusher.Flush()
			}
		}
	}
}

func writeSSERecord(w http.ResponseWriter, record logs.Record) {
	body, err := json.Marshal(record)
	if err != nil {
		return
	}
	w.Write([]byte("event: log\ndata: "))
	w.Write(body)
	w.Write([]byte("\n\n"))
}

func writeSSEComment(w http.ResponseWriter, comment string) {
	w.Write([]byte(": " + comment + "\n\n"))
}
