package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var (
	server          *httptest.Server
	reader          io.Reader
	addWordUrl      string
	getWordsUrl     string
	getWordsInclUrl string
	getWordInfoUrl  string
)

// Testing constants
const (
	SERVER_URL       = "http://localhost:8083"
	HTTP_REQ_SUCCESS = 200
)

// Testing function called once at the beginning
func init() {
	// Empy DB
	emptyDB()

	// Get the API URL
	addWordUrl = SERVER_URL + "/words/add"
	getWordsUrl = SERVER_URL + "/words"
	getWordsInclUrl = SERVER_URL + "/words/incl"
	getWordInfoUrl = SERVER_URL + "/word/"
}

// emptyDB empties the MongoDB dictionary database
func emptyDB() {
	// Get DB session
	dbs, err := mgo.Dial(HOST)
	if err != nil {
		panic(err)
	}
	defer dbs.Close()

	// Switch the session to a monotonic behavior.
	dbs.SetMode(mgo.Monotonic, true)
	dbs.DB(DICTIONARY_DB).C(WORD_COLLECTION).RemoveAll(nil)
}

// addWord is an API to add a new word to the dictionary
func TestAddWord(t *testing.T) {
	fmt.Println("-- addWord --")

	// Valid case 1
	wordJSON := `{
		"word": "test",
		"definition": "Definition of test",
		"creation_time": 1456530010
	}`

	reader = strings.NewReader(wordJSON)                        //Convert string to reader
	request, err := http.Post(addWordUrl, CONTENT_TYPE, reader) //Create request with JSON body
	if err != nil {
		t.Error(err.Error())
	}
	if request.StatusCode != HTTP_REQ_SUCCESS {
		t.Errorf("Success expected: %d", request.StatusCode) // HTTP request failed: Test failed
	}
}

// getWords is an API to get all the dictionary words
func TestGetWords(t *testing.T) {
	fmt.Println("-- getWords --")

	// Valid case 1
	request, err := http.Get(getWordsUrl) //Create request with JSON body
	if err != nil {
		t.Error(err.Error())
	}
	if request.StatusCode != HTTP_REQ_SUCCESS {
		t.Errorf("Success expected: %d", request.StatusCode) // HTTP request failed: Test failed
	}
}

// getWordsIncl is an API to get all the dictionary words that include a
// user specified string
func TestGetWordsIncl(t *testing.T) {
	fmt.Println("-- getWordsIncl --")

	// Valid case 1
	filterJSON := `{
		"filter_str": "est"
	}`
	reader = strings.NewReader(filterJSON)                           //Convert string to reader
	request, err := http.Post(getWordsInclUrl, CONTENT_TYPE, reader) //Create request with JSON body
	if err != nil {
		t.Error(err.Error())
	}
	if request.StatusCode != HTTP_REQ_SUCCESS {
		t.Errorf("Success expected: %d", request.StatusCode) // HTTP request failed: Test failed
	}
}

// getWordInfo is an API to get a user specified word information (creation date,
// definition, comments, rating)
func TestGetWordInfo(t *testing.T) {
	fmt.Println("-- getWordInfo --")

	// Valid case 1
	searchedWord := "test"
	t.Log(getWordInfoUrl)
	request, err := http.Get(getWordInfoUrl + searchedWord) //Create request with JSON body
	if err != nil {
		t.Error(err.Error())
	}
	if request.StatusCode != HTTP_REQ_SUCCESS {
		t.Errorf("Success expected: %d", request.StatusCode) // HTTP request failed: Test failed
	}
}
