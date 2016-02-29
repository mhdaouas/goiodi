package main

import (
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"net/http"
)

// Route structure
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Route slice
type Routes []Route

// createRouter creates router with several routes and sets up the logger
func createRouter(dbs *mgo.Session) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	routes := generateRoutes(dbs)

	for _, route := range routes {
		handler := Logger(route.HandlerFunc, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))
	return router
}

// generateRoutes defines route names, methods, paths and handlers
func generateRoutes(dbs *mgo.Session) Routes {
	return Routes{
		// Index route
		Route{
			"Index page",
			"GET",
			"/",
			serveIndex,
		},
		// Route to get all dictionary words
		Route{
			"Get all dictionary words",
			"GET",
			"/words",
			getWords(dbs),
		},
		// Route to get all dictionary words including a specific string
		Route{
			"Get all dictionary words including a specific string",
			"POST",
			"/words/incl",
			getWordsIncl(dbs),
		},
		// Route to get a specific word (with relative info) from the dictionary
		Route{
			"Get a specific word (with relative info) from the dictionary",
			"GET",
			"/word/{word:[a-z]+}",
			getWordInfo(dbs),
		},
		// Route to add a new word to the dictionary
		Route{
			"Add a new word to the dictionary",
			"POST",
			"/words/add",
			addWord(dbs),
		},
		// Route to add a new comment for a specific word of the dictionary
		Route{
			"Add a new comment for a specific word of the dictionary",
			"POST",
			"/comments/add",
			addComment(dbs),
		},
	}
}

// serveIndex defines the handler of the index route
func serveIndex(w http.ResponseWriter, r *http.Request) {
	Log.Notice("Serve Index")
	http.ServeFile(w, r, "static/index.html")
}
