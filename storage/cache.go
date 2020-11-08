package storage

var currentCache *cache

type jsonGetter func() []string

type cache struct {
	getter jsonGetter
	m map[int]*computedSeason
}

func (c *cache) init() {
	//m := map[int]*computedSeason{}
	//for _, j := range c.getter() {
	//	cs, err := computeSeason(j)
	//	if err == nil {
	//		m[cs.year] = cs
	//	}
	//}
	//c.m = m
}

func getJSONsFromDirectory() []string {
	return nil
}