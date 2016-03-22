package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"net/http"
	"time"
)

type APIResponse struct {
	Error    string      `json:"error"`
	Meta     interface{} `json:"meta"`
	Response interface{} `json:"response"`
}

func (r APIResponse) String() (s string) {
	b, err := json.MarshalIndent(r, "", "   ")
	if err != nil {
		s = ""
		return
	}
	s = string(b)
	return s
}

func printAPIResponse(w http.ResponseWriter, data interface{}) {
	fmt.Fprint(w, APIResponse{Response: data})
}

// getWords is an API to get all the dictionary words
func getWords(dbs *mgo.Session) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var words []Word

		// Check HTTP request content type as JSON
		w.Header().Set("Content-Type", CONTENT_TYPE)

		// Check HTTP request method is the correct one
		if r.Method != "GET" {
			http.Error(w, "HTTP method not allowed", http.StatusMethodNotAllowed)
			return
		}

		err := dbs.DB(GOIODI_DB).C(WORD_COLLECTION).Find(nil).All(&words)
		if err != nil {
			Log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		printAPIResponse(w, words)
	}
}

// getWordsIncl is an API to get all the dictionary words that include a
// user specified string
func getWordsIncl(dbs *mgo.Session) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var words []Word

		// Check HTTP request content type as JSON
		w.Header().Set("Content-Type", CONTENT_TYPE)

		// Check HTTP request method is the correct one
		if r.Method != "POST" {
			http.Error(w, "HTTP Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		byteData, err := ioutil.ReadAll(r.Body)
		if err != nil {
			Log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		type Data struct {
			FilterStr string `json:"filter_str" bson:"filter_str" db:"filter_str"`
		}
		var d Data
		err = json.Unmarshal(byteData, &d)
		if err != nil {
			Log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		Log.Debug("Filter string: ", d.FilterStr)

		err = dbs.DB(GOIODI_DB).C(WORD_COLLECTION).
			Find(bson.M{"word": bson.RegEx{d.FilterStr + ".*", ""}}).
			All(&words)
		if err != nil {
			Log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		printAPIResponse(w, words)
	}
}

// addComment is an API to add a comment for a specific word in the dictionary
func addComment(dbs *mgo.Session) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var comment Comment

		// Check HTTP request content type as JSON
		w.Header().Set("Content-Type", CONTENT_TYPE)

		// Check HTTP request method is the correct one
		if r.Method != "POST" {
			http.Error(w, "HTTP method not allowed", http.StatusMethodNotAllowed)
			return
		}

		byteData, err := ioutil.ReadAll(r.Body)
		if err != nil {
			Log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = json.Unmarshal(byteData, &comment)
		if err != nil {
			Log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Add the creation time of the word based on Epoch time
		now := time.Now()
		secs := now.Unix()
		comment.CreationTime = secs

		if comment.Word != "" && comment.Content != "" {
			c := dbs.DB(GOIODI_DB).C(COMMENT_COLLECTION)
			err = c.Insert(comment)
			if err != nil {
				Log.Error(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			err := errors.New("The comment has no related word and/or content")
			Log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// addWord is an API to add a new word to the dictionary
func addWord(dbs *mgo.Session) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var word Word

		// Check HTTP request content type as JSON
		w.Header().Set("Content-Type", CONTENT_TYPE)

		// Check HTTP request method is the correct one
		if r.Method != "POST" {
			http.Error(w, "HTTP method not allowed", http.StatusMethodNotAllowed)
			return
		}

		byteData, err := ioutil.ReadAll(r.Body)
		if err != nil {
			Log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = json.Unmarshal(byteData, &word)
		if err != nil {
			Log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Add the creation time of the word based on Epoch time
		now := time.Now()
		secs := now.Unix()
		word.CreationTime = secs

		if word.Word != "" && word.Definition != "" {
			c := dbs.DB(GOIODI_DB).C(WORD_COLLECTION)
			err = c.Insert(word)
			if err != nil {
				Log.Error(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			err := errors.New("The new word and/or its definition is an empty string")
			Log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// getWordInfo is an API to get a user specified word information (creation date,
// definition, comments, rating)
func getWordInfo(dbs *mgo.Session) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var word Word

		// Check HTTP request content type as JSON
		w.Header().Set("Content-Type", CONTENT_TYPE)

		// Check HTTP request method is the correct one
		if r.Method != "GET" {
			http.Error(w, "HTTP Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		vars := mux.Vars(r)
		if searchedWord, ok := vars["word"]; ok {
			Log.Debug("Searched word: ", searchedWord)

			err := dbs.DB(GOIODI_DB).C(WORD_COLLECTION).
				Find(bson.M{"word": searchedWord}).
				One(&word)
			if err != nil {
				Log.Error(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			printAPIResponse(w, word)
		} else {
			err := errors.New("No word in request")
			Log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
