## Description
A simple golang api to fetch a random joke based on a random person name. this api calls two backend api's.
one for random names and another is for getting jokes based on these random names.

Back end api calls is devidided into two parts, on for getting random name. which is Aync call and 
run through a seperate goroutine to update the global names cache. 
and another api(when user called 'localhost:8080/GetNewJoke') fetch names from this cache and
make backend api call to get a joke by sending this name as parameter.


TODO: Performance benchmarking 

Dependencies
-------
go get github.com/gorilla/mux


## Build & Run
go build
./jokeapi

## Usage
curl http://localhost:8080/GetNewJoke


Should return a random joke

