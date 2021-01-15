package resolve

// assumption: the resolver (stateful) keeps track of open host connections


type StopSubscription struct {
	ch <-chan struct{}
}

type Subscription struct {
	id      uint64
	Results <-chan []byte
}

type StreamSubscriber interface {
	Results() chan<- []byte
	Stop() <-chan struct{}
}

type StreamFactory interface {
	Stream(input []byte) Stream
	// UniqueIdentifier helpes the StreamManager to distinguish between uniq stream hosts
	// Good values are: hostname, url, ip address
	UniqueIdentifier(input []byte) []byte
}

// StreamManager - is responsible for storing and sharing streams between planers
// We could have a multple planners creating subscriptions to the same host and have to manage streams separately
type StreamManager interface {
	Stream(input []byte, f StreamFactory) Stream
}

// TODO: we could have a logic of handling multiple connected client inside a stream implementation