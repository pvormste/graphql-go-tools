package graphql_datasource

import (
	"sync"

	"github.com/jensneuse/graphql-go-tools/pkg/engine/datasource/httpclient"
	"github.com/jensneuse/graphql-go-tools/pkg/engine/resolve"
)

type Stream struct {
	subscribers map[resolve.StreamSubscriber]struct{}
	results chan []byte
	sync.Mutex
}

func (s *Stream) removeSubscriber(subscriber resolve.StreamSubscriber) {
	s.Lock()
	delete(s.subscribers, subscriber)
	s.Unlock()
}

func (s *Stream) addSubscriber(subscriber resolve.StreamSubscriber) {
	s.Lock()
	s.subscribers[subscriber] = struct{}{}
	s.Unlock()
}

func (s *Stream) Start(input []byte, subscriber resolve.StreamSubscriber) {
	go func() {
		s.addSubscriber(subscriber)
	}()
	defer func() {
		s.removeSubscriber(subscriber)
	}()

	for {
		select {
		case <-subscriber.Done():
			return
		case res := <-s.results:
			for streamSubscriber, _ := range s.subscribers {
				streamSubscriber.Results() <- res
			}
		}
	}

	// reading from websocket and push to the next - are different threads

	// multiple trigger instances will call start

	// trigger instance is uniq per request
	// collect next and a stop

	// if I'm first websocket connection is nil

	// we should create it

	// if more users wan't to use a new stream

	// we iterate overall all next channels

	// if someone calls stop remove it from list

	// if all stops stop websocker connection
}

type StreamFactory struct {
}

func (s *StreamFactory) Stream(input []byte) resolve.Stream {
	return &Stream{}
}

func (s *StreamFactory) UniqueIdentifier(input []byte) []byte {
	_, host, _, _, _ := httpclient.GetSubscriptionInput(input)
	return host
}

type StreamManager struct {
	streams map[string]resolve.Stream
	sync.Mutex
}

func (s *StreamManager) Stream(input []byte, f resolve.StreamFactory) resolve.Stream {
	uid := string(f.UniqueIdentifier(input))

	s.Lock()
	defer s.Unlock()

	stream, ok := s.streams[uid]
	if !ok {
		stream = f.Stream(input)
		s.streams[uid] = stream
	}

	return stream
}
