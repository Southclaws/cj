package heartbeat

import (
	"errors"
	"testing"
)

func TestJobRegistryTracksRunsAndErrors(t *testing.T) {
	r := newJobRegistry()
	r.register("aggregator", "aggregator", "@hourly")

	r.recordRun("aggregator", nil)
	snapshot := r.snapshot()
	if len(snapshot) != 1 {
		t.Fatalf("expected 1 job, got %d", len(snapshot))
	}
	if snapshot[0].RunCount != 1 || snapshot[0].LastError != "" {
		t.Fatalf("expected a clean successful run, got %+v", snapshot[0])
	}
	if snapshot[0].LastRun == nil || snapshot[0].LastRun.IsZero() {
		t.Fatal("expected LastRun to be set")
	}

	r.recordRun("aggregator", errors.New("boom"))
	snapshot = r.snapshot()
	if snapshot[0].RunCount != 2 || snapshot[0].LastError != "boom" {
		t.Fatalf("expected run count 2 with last error recorded, got %+v", snapshot[0])
	}
}

func TestJobRegistrySnapshotIsSortedByName(t *testing.T) {
	r := newJobRegistry()
	r.register("zzz", "zzz", "@hourly")
	r.register("aaa", "aaa", "@hourly")

	snapshot := r.snapshot()
	if len(snapshot) != 2 || snapshot[0].Name != "aaa" || snapshot[1].Name != "zzz" {
		t.Fatalf("expected sorted snapshot, got %+v", snapshot)
	}
}

func TestJobRegistryIgnoresRunForUnknownJob(t *testing.T) {
	r := newJobRegistry()
	r.recordRun("does-not-exist", nil)
	if len(r.snapshot()) != 0 {
		t.Fatal("expected no jobs to be tracked")
	}
}

func TestJobRegistryNeverRunJobHasNilLastRun(t *testing.T) {
	r := newJobRegistry()
	r.register("readme", "readme", "@every 6h")

	snapshot := r.snapshot()
	if len(snapshot) != 1 || snapshot[0].LastRun != nil {
		t.Fatalf("expected a never-run job to have a nil LastRun, got %+v", snapshot[0])
	}
}
