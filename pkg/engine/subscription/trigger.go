package subscription

import (
	"context"
)

func NewTrigger(subscriptionID uint64) Trigger {
	return Trigger{
		subscriptionID: subscriptionID,
		results:        make(chan []byte), // unbuffered channel
	}
}

// Trigger - is an istance of an active subscriptions
// trigger has a stream
// trigger forward results to listener
// one same input from different clients will result in a single subscription
type Trigger struct {
	subscriptionID uint64
	results        chan []byte
}

func (h *Trigger) SubscriptionID() uint64 {
	return h.subscriptionID
}

func (h *Trigger) Next(ctx context.Context) (data []byte, ok bool) {
	done := ctx.Done()
	select {
	case <-done:
		return nil, false
	case result, ok := <-h.results:
		return result, ok
	}
}
