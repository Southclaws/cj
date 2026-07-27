package logs

import (
	"encoding/json"
	"fmt"
	"time"

	"go.uber.org/zap/zapcore"
)

type Record struct {
	Timestamp time.Time      `json:"timestamp"`
	Level     string         `json:"level"`
	Component string         `json:"component,omitempty"`
	Message   string         `json:"message"`
	Fields    map[string]any `json:"fields,omitempty"`
}

func fieldsToMap(fields []zapcore.Field) map[string]any {
	enc := zapcore.NewMapObjectEncoder()
	for _, f := range fields {
		f.AddTo(enc)
	}

	out := make(map[string]any, len(enc.Fields))
	for k, v := range enc.Fields {
		out[k] = normalize(v)
	}
	return out
}

func normalize(v any) any {
	b, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprintf("%v", v)
	}
	var normalized any
	if err := json.Unmarshal(b, &normalized); err != nil {
		return fmt.Sprintf("%v", v)
	}
	return normalized
}
