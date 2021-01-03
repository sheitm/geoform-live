package competitions

import (
	"github.com/3lvia/telemetry-go"
	"github.com/sheitm/ofever/athletes"
	"github.com/sheitm/ofever/persist"
	"github.com/sheitm/ofever/sequence"
	"sync"
)

func Start(
	add sequence.Adder,
	athleteID athletes.AthleteIDFunc,
	persistFunc persist.Persist,
	reader persist.ReadFunc,
	readContainersFunc persist.ReadContainersFunc,
	logChannels telemetry.LogChans) telemetry.RequestHandler {
	seqChan := make(chan *sequence.Event)
	add(seqChan)

	i := &impl{
		comps:              map[string]*comp{},
		athleteID:          athleteID,
		persistFunc:        persistFunc,
		mux:                &sync.Mutex{},
		logChannels:        logChannels,
	}
	go i.start(seqChan, reader, readContainersFunc)

	h := &handler {
		get: i.get,
		getAll: i.getAll,
	}
	return h
}
