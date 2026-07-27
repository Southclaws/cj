package apierrors

import (
	"encoding/json"
	"net/http/httptest"
	"testing"
)

func TestWriteErrorEnvelopeShape(t *testing.T) {
	rec := httptest.NewRecorder()

	WriteError(rec, 403, CodeHostNotAllowed, "nope", "req-1")

	var body struct {
		Error struct {
			Code      string `json:"code"`
			Message   string `json:"message"`
			RequestID string `json:"requestId"`
		} `json:"error"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("response body is not valid JSON: %v", err)
	}

	if body.Error.Code != CodeHostNotAllowed {
		t.Fatalf("got code %q, want %q", body.Error.Code, CodeHostNotAllowed)
	}
	if body.Error.RequestID != "req-1" {
		t.Fatalf("got requestId %q, want req-1", body.Error.RequestID)
	}
	if rec.Code != 403 {
		t.Fatalf("got status %d, want 403", rec.Code)
	}
}
