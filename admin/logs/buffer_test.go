package logs

import (
	"testing"
	"time"
)

func TestBufferEvictsOldestWhenFull(t *testing.T) {
	b := NewBuffer(3, 0)
	for i := 0; i < 5; i++ {
		b.Push(Record{Message: string(rune('a' + i)), Timestamp: time.Now()})
	}

	got := b.Snapshot(Filter{})
	if len(got) != 3 {
		t.Fatalf("expected 3 records, got %d", len(got))
	}
	if got[0].Message != "e" || got[2].Message != "c" {
		t.Fatalf("expected newest-first [e d c], got %v", got)
	}
}

func TestBufferEvictsByRetention(t *testing.T) {
	b := NewBuffer(10, 10*time.Millisecond)
	b.Push(Record{Message: "old", Timestamp: time.Now().Add(-time.Hour)})
	b.Push(Record{Message: "new", Timestamp: time.Now()})

	got := b.Snapshot(Filter{})
	if len(got) != 1 || got[0].Message != "new" {
		t.Fatalf("expected only the recent record, got %v", got)
	}
}

func TestBufferSnapshotFiltersByLevelComponentAndQuery(t *testing.T) {
	b := NewBuffer(10, 0)
	b.Push(Record{Level: "info", Component: "bot", Message: "hello world", Timestamp: time.Now()})
	b.Push(Record{Level: "error", Component: "storage", Message: "failed to connect", Timestamp: time.Now()})

	if got := b.Snapshot(Filter{Level: "error"}); len(got) != 1 || got[0].Component != "storage" {
		t.Fatalf("level filter failed: %v", got)
	}
	if got := b.Snapshot(Filter{Component: "bot"}); len(got) != 1 || got[0].Message != "hello world" {
		t.Fatalf("component filter failed: %v", got)
	}
	if got := b.Snapshot(Filter{Query: "connect"}); len(got) != 1 || got[0].Level != "error" {
		t.Fatalf("query filter failed: %v", got)
	}
	if got := b.Snapshot(Filter{}); len(got) != 2 {
		t.Fatalf("expected both records with no filter, got %d", len(got))
	}
}

func TestBufferSnapshotFiltersByCorrelationID(t *testing.T) {
	b := NewBuffer(10, 0)
	b.Push(Record{Message: "a", Fields: map[string]any{"correlation_id": "abc"}, Timestamp: time.Now()})
	b.Push(Record{Message: "b", Fields: map[string]any{"correlation_id": "xyz"}, Timestamp: time.Now()})

	got := b.Snapshot(Filter{CorrelationID: "abc"})
	if len(got) != 1 || got[0].Message != "a" {
		t.Fatalf("correlation filter failed: %v", got)
	}
}

func TestBufferSnapshotRespectsLimit(t *testing.T) {
	b := NewBuffer(10, 0)
	for i := 0; i < 5; i++ {
		b.Push(Record{Message: "x", Timestamp: time.Now()})
	}
	if got := b.Snapshot(Filter{Limit: 2}); len(got) != 2 {
		t.Fatalf("expected 2 records, got %d", len(got))
	}
}
