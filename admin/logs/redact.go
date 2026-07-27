package logs

import (
	"regexp"
	"strings"
	"sync"
)

var sensitiveKeyPattern = regexp.MustCompile(`(?i)(token|password|passwd|pass|secret|credential|authorization|apikey|api_key)`)

const redactedPlaceholder = "[REDACTED]"

type Redactor struct {
	mu           sync.RWMutex
	secretValues []string
}

func NewRedactor(secrets ...string) *Redactor {
	r := &Redactor{}
	r.SetSecrets(secrets...)
	return r
}

func (r *Redactor) SetSecrets(secrets ...string) {
	filtered := make([]string, 0, len(secrets))
	for _, s := range secrets {
		if s != "" {
			filtered = append(filtered, s)
		}
	}
	r.mu.Lock()
	r.secretValues = filtered
	r.mu.Unlock()
}

func (r *Redactor) Redact(record Record) Record {
	record.Message = r.redactString(record.Message)
	record.Fields = r.redactMap(record.Fields)
	return record
}

func (r *Redactor) redactMap(m map[string]any) map[string]any {
	if m == nil {
		return nil
	}
	out := make(map[string]any, len(m))
	for k, v := range m {
		if sensitiveKeyPattern.MatchString(k) {
			out[k] = redactedPlaceholder
			continue
		}
		out[k] = r.redactValue(v)
	}
	return out
}

func (r *Redactor) redactValue(v any) any {
	switch val := v.(type) {
	case string:
		return r.redactString(val)
	case map[string]any:
		return r.redactMap(val)
	case []any:
		out := make([]any, len(val))
		for i, item := range val {
			out[i] = r.redactValue(item)
		}
		return out
	default:
		return v
	}
}

func (r *Redactor) redactString(s string) string {
	r.mu.RLock()
	secrets := r.secretValues
	r.mu.RUnlock()

	for _, secret := range secrets {
		s = strings.ReplaceAll(s, secret, redactedPlaceholder)
	}
	return s
}
