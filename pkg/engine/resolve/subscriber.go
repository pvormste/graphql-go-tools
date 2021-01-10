package resolve

import "context"

type subscriber struct {
	next chan []byte
	stop chan struct{}
}

func newSubscriber() *subscriber {
	return &subscriber{
		next: make(chan []byte),
		stop: make(chan struct{}),
	}
}

func (s *subscriber) Results() chan<- []byte {
	return s.next
}

func (s *subscriber) Done() <-chan struct{} {
	return s.stop
}

func (s *subscriber) Stop() {
	close(s.stop)
}

func (s *subscriber) Next(ctx context.Context) (data []byte, ok bool) {
	done := ctx.Done()
	select {
	case <-done:
		return nil, false
	case result, ok := <-s.next:
		return result, ok
	}
}
