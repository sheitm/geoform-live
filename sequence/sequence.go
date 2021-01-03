package sequence

import (
	"sync"
)

// Event is a construct to let sequence dependencies handle an event in turn.
type Event struct {
	Payload interface{}

	DoneChan chan<- struct{}
}

// Adder is a func letting dependencies make their interest in participating in ths sequential processing known.
type Adder func(chan<- *Event)

// Start starts the internal sequential processor.
func Start(trigger <-chan interface{}, doneChan chan<- struct{}) Adder {
	i := &impl{
		mux:      sync.Mutex{},
		doneChan: doneChan,
	}
	go i.start(trigger)
	return i.add
}

type impl struct {
	channels []chan<- *Event
	doneChan chan<- struct{}
	mux      sync.Mutex
}

func (s *impl) start(trigger <-chan interface{}) {
	for {
		payload := <- trigger
		for _, channel := range s.channels {
			doneChan := make(chan struct{})
			e := &Event{
				Payload:  payload,
				DoneChan: doneChan,
			}
			channel <- e
			<- doneChan
		}
		s.doneChan <- struct{}{}
	}
}

func (s *impl) add(ch chan<- *Event) {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.channels = append(s.channels, ch)
}