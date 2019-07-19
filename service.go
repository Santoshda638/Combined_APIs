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
}

var (
	httpClient *http.Client
	transport  *http.Transport
)

// Initiate net http object.
func init() {
	transport = &http.Transport{}
	httpClient = &http.Client{
		Timeout:   time.Duration(1) * time.Minute,
		Transport: transport,
	}
}

func NewJokesApi() *jokesApi {
	return new(jokesApi)
}

//GetJoke will return a random joke from api or cache.
func (this *jokesApi) Getjoke() (joke *Joke, err error) {
	var ok bool
	//First check in chache
	select {
	case joke, ok = <-jokeCache:
		if ok {
			// for Updating cache
			go fetchJokeFromApi(getPeopleCache())
			return
		}
	default:
		p := getPeopleCache()
		// Get it from API
		joke, err = fetchJokeFromApi(p)
		if err != nil {
			//For now error handling is logging err only
			//TODO use custom logging and use seperate log level.
			//TODO in case of error we should return custom err message to client instead of backend error response.
			log.Printf("Failed to fetch joke. for %v,%v Error %v", p.Name, p.Surname, err)

		}
	}

	//fmt.Println(joke)
	return

}

// fetchJokeFromApi fetches jokes from jokeApi url and returns.
func fetchJokeFromApi(person *Person) (joke *Joke, e error) {
	jokeUrl := fmt.Sprintf(JOKES_URL, url.QueryEscape(person.Name), url.QueryEscape(person.Surname))
	var response *http.Response
	//Joke url call.
	response, e = httpClient.Get(jokeUrl)
	if e != nil {
		return
	}

	if code := response.StatusCode; code != http.StatusOK {
		fmt.Println("Received status ", code)
	}

	var jbody []byte
	jbody, e = ioutil.ReadAll(response.Body)
	if e != nil {
		return
	}
	e = json.Unmarshal(jbody, &joke)

	//Lets update the cache.
	go setJokeCache(joke)

	return
}

// getRandomName fetch name randomly and update to cache object object.
func getNameFromApi() (p *Person) {
	defer func() {
		if r := recover(); r != nil {
			log.Fatalf("PANIC Recovery. Error %v and stack trace  %s", r, debug.Stack())
		}
	}()

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

	err = json.Unmarshal(nameBytes, &p)
	if err != nil {
		log.Println(err)
		return
	}

	return

}

// getStaticName returns a randon name from const variable.
func getStaticName() (string, string) {
	return firstNames[rand.Intn(len(firstNames))], lastNames[rand.Intn(len(lastNames))]
}
