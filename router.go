package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

var controller = NewController()

// Route defines a route
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes defines the list of routes of our API. we can add more handllers here.
type Routes []Route

var routes = Routes{
	Route{
		"GetNewJoke",
		"GET",
		"/GetNewJoke",
		controller.GetNewJokeHandller,
	}}

// newRouter set routing web methods http handllers.
func newRouter() *mux.Router {
	router := mux.NewRouter()
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}
