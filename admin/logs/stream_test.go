package logs

import (
	"testing"
	"time"
)

func TestStreamBroadcastDoesNotBlockOnSlowSubscriber(t *testing.T) {
	stream := NewStream()
	stuck := stream.Subscribe()
	defer stream.Unsubscribe(stuck)

	done := make(chan struct{})
	go func() {
		defer close(done)
		for i := 0; i < subscriberQueueSize+50; i++ {
			stream.Broadcast(Record{Message: "msg"})
		}
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Broadcast blocked on a subscriber that was never drained")
	}

	select {
	case <-stuck.Dropped():
	default:
		t.Fatal("expected the stuck subscriber to observe a dropped-records signal")
	}
}

func TestStreamDeliversToAnActivelyDrainedSubscriber(t *testing.T) {
	stream := NewStream()
	sub := stream.Subscribe()
	defer stream.Unsubscribe(sub)

	for i := 0; i < subscriberQueueSize+50; i++ {
		stream.Broadcast(Record{Message: "msg"})
		select {
		case <-sub.Records():
		case <-time.After(time.Second):
			t.Fatalf("expected record %d to be delivered to a subscriber drained in lockstep", i)
		}
	}
}

func TestStreamUnsubscribeStopsDelivery(t *testing.T) {
	stream := NewStream()
	sub := stream.Subscribe()
	stream.Unsubscribe(sub)

	stream.Broadcast(Record{Message: "after unsubscribe"})

	select {
	case r := <-sub.Records():
		t.Fatalf("expected no delivery after unsubscribe, got %v", r)
	default:
	}

	if got := stream.SubscriberCount(); got != 0 {
		t.Fatalf("expected 0 subscribers after unsubscribe, got %d", got)
	}
}
