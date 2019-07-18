package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"runtime/debug"
	"time"
)

type jokesApi struct {
	// lets not expose this. if need we can allow clients to set it for custom person name response.
	person *Person
}

var (
	httpClient   *http.Client
	transport    *http.Transport
	randomPeople = make(chan *Person, 10)
)

// Initiate net http object.
func init() {
	transport = &http.Transport{}
	httpClient = &http.Client{
		Timeout:   time.Duration(1) * time.Minute,
		Transport: transport,
	}
	rand.Seed(42)
	// Fetch a random name in seperate go routine
	// TODO: wait or wait group is not required as i am using default offline names in case random names not available.
	// i am updating random name each time api respond with new Joke.
	// but its safe to wait until we get response
	// OR we can fill some random name at inititlization (at api start) and use those name to
	// fetch randon jokes. and again refill with anothother random name before response to client. Please have look at

	go getRandomName()

}
func NewJokesApi() (j *jokesApi) {
	j = new(jokesApi)

	// Ideally this param should come from caller when caller wat to provide custom name for search.
	// but for now we are reading from global variable, which is random and dynamic so let it be.
	j.person = new(Person)
	return
}

//GetJoke will return a random joke from api.
//TODO in case of err we should return custom err message to client instead of backend error response.
func (this *jokesApi) Getjoke() *Joke {

	select {
	case this.person = <-randomPeople:
	default:
		// for now taking randonm name from cache.
		log.Printf("Incase no person detail is provided, Use a cached name")
		this.person = new(Person)
		this.person.Name, this.person.Surname = getCachedName()
	}
	err, joke := fetchJoke(this.person)
	if err != nil {
		//For now error handling is logging err only
		//TODO use custom logging and use seperate log level.
		log.Printf("Failed to fetch joke. for %v,%v Error %v", this.person.Name, this.person.Surname, err)

	}
	// Fetch and Update Random Name in gloabl variable for next call.
	//TODO put this logic in main.go and add waitgroup 'for' loop for contineaous fetching and filling random names
	go getRandomName()
	fmt.Println(joke)
	return joke

}
func buildURL(p *Person) (string, error) {
	params := url.Values{}
	params.Add("firstName", p.Name)
	params.Add("lastName", p.Surname)
	params.Add("limitTo", "[nerdy]")
	baseURL, err := url.Parse(JOKES_URL)
	if err != nil {
		return "", err
	}
	baseURL.RawQuery = params.Encode()
	return baseURL.String(), nil
}

// fetchJoke fetches jokes from jokeApi url and returns.
func fetchJoke(person *Person) (e error, joke *Joke) {
	var jurl string
	if jurl, e = buildURL(person); e != nil {
		return
	}

	var response *http.Response
	//Joke url call.
	response, e = httpClient.Get(jurl)
	if e != nil {
		return
	}
	var jk []byte
	jk, e = ioutil.ReadAll(response.Body)
	if e != nil {
		return
	}
	e = json.Unmarshal(jk, &joke)

	return
}

// getRandomName fetch name randomly and update to cache object object.
func getRandomName() {
	defer func() {
		if r := recover(); r != nil {
			log.Fatalf("PANIC Recovery. Error %v and stack trace  %s", r, debug.Stack())
		}
	}()
	var name *Person
	var nameBytes []byte
	var err error
	response, err := httpClient.Get(USER_DETAILS_URL)
	if err != nil {
		log.Println(err)
		return
	}

	defer response.Body.Close()
	nameBytes, err = ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		return
	}

	err = json.Unmarshal(nameBytes, &name)
	if err != nil {
		log.Println(err)
		return
	}

	randomPeople <- name
	//wg.Done()
	return

}

// getCachedName returns a randon name from cache.
func getCachedName() (string, string) {
	return firstNames[rand.Intn(len(firstNames))], lastNames[rand.Intn(len(lastNames))]
}
