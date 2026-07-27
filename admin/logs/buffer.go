package logs

import (
	"strings"
	"sync"
	"time"
)

type Filter struct {
	Level         string
	Component     string
	Query         string
	CorrelationID string
	Limit         int
}

type Buffer struct {
	mu        sync.Mutex
	records   []Record
	size      int
	retention time.Duration
}

func NewBuffer(size int, retention time.Duration) *Buffer {
	if size <= 0 {
		size = 1
	}
	return &Buffer{size: size, retention: retention}
}

func (b *Buffer) Push(r Record) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.records = append(b.records, r)
	if len(b.records) > b.size {
		b.records = b.records[len(b.records)-b.size:]
	}
	b.evictLocked()
}

func (b *Buffer) evictLocked() {
	if b.retention <= 0 {
		return
	}
	cutoff := time.Now().Add(-b.retention)
	idx := 0
	for idx < len(b.records) && b.records[idx].Timestamp.Before(cutoff) {
		idx++
	}
	if idx > 0 {
		b.records = b.records[idx:]
	}
}

func (b *Buffer) Len() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return len(b.records)
}

func (b *Buffer) Snapshot(filter Filter) []Record {
	b.mu.Lock()
	defer b.mu.Unlock()

	matched := make([]Record, 0, len(b.records))
	for i := len(b.records) - 1; i >= 0; i-- {
		record := b.records[i]
		if matches(record, filter) {
			matched = append(matched, record)
		}
	}

	limit := filter.Limit
	if limit <= 0 || limit > len(matched) {
		limit = len(matched)
	}
	return matched[:limit]
}

func matches(record Record, filter Filter) bool {
	if filter.Level != "" && !strings.EqualFold(record.Level, filter.Level) {
		return false
	}
	if filter.Component != "" && !strings.EqualFold(record.Component, filter.Component) {
		return false
	}
	if filter.Query != "" && !strings.Contains(strings.ToLower(record.Message), strings.ToLower(filter.Query)) {
		return false
	}
	if filter.CorrelationID != "" {
		id, ok := record.Fields["correlation_id"].(string)
		if !ok || id != filter.CorrelationID {
			return false
		}
	}
	return true
}
