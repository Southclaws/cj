package logs

import "sync"

const subscriberQueueSize = 256

type Subscriber struct {
	records chan Record
	dropped chan struct{}
}

func (s *Subscriber) Records() <-chan Record {
	return s.records
}

func (s *Subscriber) Dropped() <-chan struct{} {
	return s.dropped
}

type Stream struct {
	mu          sync.Mutex
	subscribers map[*Subscriber]struct{}
}

func NewStream() *Stream {
	return &Stream{subscribers: make(map[*Subscriber]struct{})}
}

func (s *Stream) Subscribe() *Subscriber {
	sub := &Subscriber{
		records: make(chan Record, subscriberQueueSize),
		dropped: make(chan struct{}, 1),
	}
	s.mu.Lock()
	s.subscribers[sub] = struct{}{}
	s.mu.Unlock()
	return sub
}

func (s *Stream) Unsubscribe(sub *Subscriber) {
	s.mu.Lock()
	delete(s.subscribers, sub)
	s.mu.Unlock()
}

func (s *Stream) Broadcast(r Record) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for sub := range s.subscribers {
		select {
		case sub.records <- r:
		default:
			select {
			case sub.dropped <- struct{}{}:
			default:
			}
		}
	}
}

func (s *Stream) SubscriberCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.subscribers)
}
