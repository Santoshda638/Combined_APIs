package main

type Person struct {
	Gender  string `json:"gender" bson:"gender"`
	Name    string `json:"name" bson:"name"`
	Surname string `json:"surname" bson:"surname"`
	Region  string `json:"region" bson:"region"`
}

type Joke struct {
	Type  string `json:"type"`
	Value struct {
		ID         int      `json:"id"`
		Joke       string   `json:"joke"`
		Categories []string `json:"categories"`
	} `json:"value"`
}

const USER_DETAILS_URL = "http://uinames.com/api/"
const JOKES_URL = "http://api.icndb.com/jokes/random?firstName=%v&lastName=%v&limitTo=[nerdy]"

var firstNames = []string{"sam", "denial", "Michel"}
var lastNames = []string{"mac", "piotr", "Sandy"}
