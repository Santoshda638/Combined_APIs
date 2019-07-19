package main

import (
	"sync"
)

var (
	peopleCache = make(chan *Person, 10)
	jokeCache   = make(chan *Joke, 10)
)

//Build cache for atleast f
func buildCache() {
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			setPeopleCache()
			getPeopleCache()

		}()
	}
	wg.Wait()
}

// Get people from backend and update people cache
func setPeopleCache() {
	if p := getNameFromApi(); p != nil {
		peopleCache <- p
	}
}

// getPeopleCache get people either from cache or default.
// each time you read from cache. this method try to update cache
// from backend for the next randon name
func getPeopleCache() (p *Person) {
	select {
	case p = <-peopleCache:
	default:
		// for now taking randonm name from cache.
		// log.Printf("Incase no person detail is available, Use a cached name")
		if p == nil {
			p = new(Person)
		}
		p.Name, p.Surname = getStaticName()
	}
	// Let's update the cache with the latest one.
	go setPeopleCache()
	return
}

//Update Joke Cache
func setJokeCache(j *Joke) {
	if j != nil {
		jokeCache <- j
	}
}
