package series

import (
	"github.com/sheitm/ofever/sequence"
	"github.com/sheitm/ofever/types"
)

func Start(add sequence.Adder) {
	seqChan := make(chan *sequence.Event)
	r := &router{seriesMap: map[string]chan<- *types.SeasonFetch{}}
	r.start(seqChan)
	add(seqChan)
}
