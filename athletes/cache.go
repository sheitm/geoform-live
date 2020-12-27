package athletes

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/3lvia/telemetry-go"
	"github.com/sheitm/ofever/persist"
	"strings"
	"sync"
)

type athleteWithID struct {
	ID   string `json:"id"`
	SHA  string `json:"sha"`
	Name string `json:"name"`
	Club string `json:"club"`
}

type cache struct {
	competitorsBySHA  map[string]*athleteWithID
	competitorsByGuid map[string]*athleteWithID
	logChannels       telemetry.LogChans
	mux               sync.Mutex
}

func (c *cache) all() []*athleteWithID {
	var result []*athleteWithID
	for _, athlete := range c.competitorsBySHA {
		result = append(result, athlete)
	}
	return result
}

func (c *cache) competitor(name, club string) (*athleteWithID, bool) {
	sha := sha(name, club)
	if a, ok := c.competitorsBySHA[sha]; ok {
		return a, true
	}

	c.mux.Lock()
	defer c.mux.Unlock()

	if a, ok := c.competitorsBySHA[sha]; ok {
		return a, true
	}

	a := &athleteWithID{
		ID:   guid(),
		SHA:  sha,
		Name: name,
		Club: club,
	}
	c.competitorsBySHA[sha] = a
	c.competitorsByGuid[a.ID] = a

	return a, false
}

func (c *cache) init(reader persist.ReadFunc) {
	c.mux = sync.Mutex{}
	send := make(chan []byte)
	done := make(chan struct{})
	sha := map[string]*athleteWithID{}
	guid := map[string]*athleteWithID{}
	r := persist.Read{
		Series: container,
		Path:   "",
		Send:   send,
		Done:   done,
	}

	go func(s <-chan []byte) {
		for {
			b := <- s
			var athlete athleteWithID
			err := json.Unmarshal(b, &athlete)
			if err != nil {
				c.logChannels.ErrorChan <- err
				continue
			}
			sha[athlete.SHA] = &athlete
			guid[athlete.ID] = sha[athlete.SHA]
		}
	}(send)

	reader(r)

	<- done
	c.competitorsBySHA = sha
	c.competitorsByGuid = guid
}

func sha(name, club string) string {
	n := strings.TrimSpace(strings.ToLower(name))
	c := strings.TrimSpace(strings.ToLower(club))
	s := fmt.Sprintf("%s%s", n, c)
	hash := sha1.New()
	hash.Write([]byte(s))
	sha := base64.URLEncoding.EncodeToString(hash.Sum(nil))
	return sha
}

func guid() string {
	b := make([]byte, 16)
	rand.Read(b)
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid
}