package resolve

type Stream interface {
	Start(input []byte, subscriber StreamSubscriber)
}

type StreamSubscriber interface {
	Results() chan<- []byte
	Done() <-chan struct{}
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
