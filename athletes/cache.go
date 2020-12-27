package athletes

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
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
	mux               sync.Mutex
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