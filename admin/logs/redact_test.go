package logs

import "testing"

func TestRedactRedactsSensitiveKeys(t *testing.T) {
	r := NewRedactor()
	record := Record{
		Message: "starting up",
		Fields: map[string]any{
			"DiscordToken": "abc123",
			"MongoPass":    "hunter2",
			"username":     "southclaws",
			"nested": map[string]any{
				"api_key": "sekrit",
				"safe":    "value",
			},
		},
	}

	got := r.Redact(record)

	if got.Fields["DiscordToken"] != redactedPlaceholder {
		t.Fatalf("expected DiscordToken to be redacted, got %v", got.Fields["DiscordToken"])
	}
	if got.Fields["MongoPass"] != redactedPlaceholder {
		t.Fatalf("expected MongoPass to be redacted, got %v", got.Fields["MongoPass"])
	}
	if got.Fields["username"] != "southclaws" {
		t.Fatalf("expected username to survive redaction, got %v", got.Fields["username"])
	}
	nested, ok := got.Fields["nested"].(map[string]any)
	if !ok {
		t.Fatalf("expected nested map to survive, got %T", got.Fields["nested"])
	}
	if nested["api_key"] != redactedPlaceholder {
		t.Fatalf("expected nested api_key to be redacted, got %v", nested["api_key"])
	}
	if nested["safe"] != "value" {
		t.Fatalf("expected nested safe value to survive, got %v", nested["safe"])
	}
}

func TestRedactScrubsKnownSecretValuesWhereverTheyAppear(t *testing.T) {
	r := NewRedactor("supersecrettoken")
	record := Record{
		Message: "failed request with token supersecrettoken",
		Fields: map[string]any{
			"detail": "auth header was Bearer supersecrettoken",
		},
	}

	got := r.Redact(record)

	if got.Message != "failed request with token [REDACTED]" {
		t.Fatalf("expected secret value scrubbed from message, got %q", got.Message)
	}
	if got.Fields["detail"] != "auth header was Bearer [REDACTED]" {
		t.Fatalf("expected secret value scrubbed from field, got %v", got.Fields["detail"])
	}
}

func TestRedactHandlesNilFields(t *testing.T) {
	r := NewRedactor()
	got := r.Redact(Record{Message: "no fields here"})
	if got.Fields != nil {
		t.Fatalf("expected nil fields to stay nil, got %v", got.Fields)
	}
}
