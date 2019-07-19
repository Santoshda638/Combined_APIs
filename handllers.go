package main

import (
	"encoding/json"
	"net/http"
)

type handller struct {
	japi *jokesApi
}

func NewHandller() *handller {
	return &handller{japi: NewJokesApi()}
}

// GetNewJoke a http handller returns joke response
func (this *handller) GetNewJokeHandller(w http.ResponseWriter, r *http.Request) {
	//TODO better handle errors and logging is required.
	joke, e := this.japi.Getjoke()
	if e != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data, _ := json.Marshal(joke)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	return
}
