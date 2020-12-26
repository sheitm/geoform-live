package sequence

import (
	"testing"
)

func TestStart(t *testing.T) {
	// Arrange
	trigger := make(chan interface{})
	done := make(chan struct{})

	ch1 := make(chan *Event)
	ch2 := make(chan *Event)
	ch3 := make(chan *Event)

	var ints []int

	// Act
	adder := Start(trigger, done)

	adder(ch1)
	adder(ch2)
	adder(ch3)

	go func() {
		for {
			select {
			case e := <-ch1:
				ints = append(ints, e.Payload.(int))
				e.DoneChan <- struct{}{}
			case e := <-ch2:
				ints = append(ints, e.Payload.(int))
				e.DoneChan <- struct{}{}
			case e := <-ch3:
				ints = append(ints, e.Payload.(int))
				e.DoneChan <- struct{}{}
			}
		}
	}()

	trigger <- 1

	<- done

	// Assert
	if len(ints) != 3 {
		t.Errorf("expected 3 ints, got %d", len(ints))
	}

}
