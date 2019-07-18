package main

import (
	"encoding/json"
	"net/http"
)

//Controller ...
type Controller struct {
	japi *jokesApi
}

func NewController() *Controller {
	return &Controller{japi: NewJokesApi()}
}

// GetNewJoke a http handller returns joke response
func (c *Controller) GetNewJokeHandller(w http.ResponseWriter, r *http.Request) {

	joke := c.japi.Getjoke()
	data, _ := json.Marshal(joke)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	return
}
