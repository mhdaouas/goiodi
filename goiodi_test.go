package main

import (
	"crypto/tls"
	"fmt"
	"gopkg.in/mgo.v2"
	"io"
	// "io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"os/exec"
	"strings"
	"testing"
)

var (
	client *http.Client
	reader io.Reader

	userPwd            string
	userPwdSalt        string
	addUserUrl         string
	loginUserUrl       string
	checkUserLoggedUrl string

	addWordUrl      string
	getWordsUrl     string
	getWordsInclUrl string
	getWordInfoUrl  string

	cookieJar, _ = cookiejar.New(nil)
)

// Testing parameters
const (
	SERVER_URL = "https://localhost:8083"

	USER_USERNAME = "test_username"
	USER_PWD      = "test_password"
	USER_EMAIL    = "test_email@goiodi.com"

	INVALID_USER_USERNAME = "test_invalid_username"
	INVALID_USER_PWD      = "test_invalid_password"
)

// Testing function called once at the beginning
func init() {
	// Initialize DB (prepare it for tests)
	err := emptyDB()
	if err != nil {
		panic(err)
	}

	// Run DB initialization script
	// result := bson.M{}
	// if err = dbs.DB(GOIODI_DB).Run(bson.M{"eval": "initDB();"}, &result); err != nil {
	// 	panic(err)
	// } else {
	// 	fmt.Println(result)
	// }

	cmd := exec.Command("mongo", HOST+"/"+GOIODI_DB, DB_INIT_SCRIPT)
	err = cmd.Start()
	if err != nil {
		panic(err)
	}
	Log.Notice("Waiting for command to finish...")
	err = cmd.Wait()
	Log.Notice("Command finished with error: %v", err)

	// Get the API URL
	addUserUrl = SERVER_URL + "/users/add"
	loginUserUrl = SERVER_URL + "/user/login"
	checkUserLoggedUrl = SERVER_URL + "/user/login/check"
	addWordUrl = SERVER_URL + "/words/add"
	getWordsUrl = SERVER_URL + "/words"
	getWordsInclUrl = SERVER_URL + "/words/incl"
	getWordInfoUrl = SERVER_URL + "/word"

	// Do not check TLS certificate and set default cookie Jar
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client = &http.Client{
		Transport: tr,
		Jar:       cookieJar,
	}
}

// emptyDB empties the MongoDB dictionary database
func emptyDB() error {
	// Get DB session
	dbs, err := mgo.Dial(HOST)
	if err != nil {
		return err
	}
	defer dbs.Close()

	// Switch the session to a monotonic behavior.
	dbs.SetMode(mgo.Monotonic, true)
	dbs.DB(GOIODI_DB).C(WORD_COLLECTION).RemoveAll(nil)
	return err
}

// addWord is an API to add a new word to the dictionary
func TestAddWord(t *testing.T) {
	fmt.Println("-- addWord --")

	// Valid case 1

	wordJSON := `{
		"word": "test",
		"definition": "Definition of test"
	}`

	reader = strings.NewReader(wordJSON)                          // Convert string to reader
	request, err := client.Post(addWordUrl, CONTENT_TYPE, reader) // Create request with JSON body
	if err != nil {
		t.Error(err.Error())
	}
	if request.StatusCode != http.StatusOK {
		t.Errorf("Success expected: %d", request.StatusCode) // HTTP request failed: Test failed
	}

	// Invalid case 1: Missing word

	invalidWordJSON := `{
		"definition": "Definition of test"
	}`

	reader = strings.NewReader(invalidWordJSON)                  // Convert string to reader
	request, err = client.Post(addWordUrl, CONTENT_TYPE, reader) // Create request with JSON body
	if err != nil {
		t.Error(err.Error())
	}
	if request.StatusCode != http.StatusInternalServerError {
		t.Errorf("Internal server error expected: %d", request.StatusCode) // HTTP request failed: Test failed
	}

	// Invalid case 2: Missing definition

	invalidWordJSON = `{
		"word": "test"
	}`

	reader = strings.NewReader(invalidWordJSON)                  // Convert string to reader
	request, err = client.Post(addWordUrl, CONTENT_TYPE, reader) // Create request with JSON body
	if err != nil {
		t.Error(err.Error())
	}
	if request.StatusCode != http.StatusInternalServerError {
		t.Errorf("Internal server error expected: %d", request.StatusCode) // HTTP request failed: Test failed
	}
}

// getWords is an API to get all the dictionary words
func TestGetWords(t *testing.T) {
	fmt.Println("-- getWords --")

	// Valid case 1

	request, err := client.Get(getWordsUrl) // Create request with JSON body
	if err != nil {
		t.Error(err.Error())
	}
	defer request.Body.Close()

	if request.StatusCode != http.StatusOK {
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
	reader = strings.NewReader(filterJSON)                             //Convert string to reader
	request, err := client.Post(getWordsInclUrl, CONTENT_TYPE, reader) //Create request with JSON body
	if err != nil {
		t.Error(err.Error())
	}
	if request.StatusCode != http.StatusOK {
		t.Errorf("Success expected: %d", request.StatusCode) // HTTP request failed: Test failed
	}

	// Valid case 2: Invalid word filter

	invalidFilterJSON := `{
		"filter_str": "fff"
	}`
	reader = strings.NewReader(invalidFilterJSON)                     //Convert string to reader
	request, err = client.Post(getWordsInclUrl, CONTENT_TYPE, reader) //Create request with JSON body
	if err != nil {
		t.Error(err.Error())
	}
	if request.StatusCode != http.StatusOK {
		t.Errorf("Internal server error expected: %d", request.StatusCode) // HTTP request failed: Test failed
	}
}

// getWordInfo is an API to get a user specified word information (creation date,
// definition, comments, rating)
func TestGetWordInfo(t *testing.T) {
	fmt.Println("-- getWordInfo --")

	// Valid case 1

	searchedWord := "test"
	t.Log(getWordInfoUrl)
	request, err := client.Get(getWordInfoUrl + "/" + searchedWord) // Create request with JSON body
	if err != nil {
		t.Error(err.Error())
	}
	defer request.Body.Close()

	if request.StatusCode != http.StatusOK {
		t.Errorf("Success expected: %d", request.StatusCode) // HTTP request failed: Test failed
	}

	// Invalid case 1: Invalid searched word

	invalidSearchedWord := "test_invalid"
	t.Log(getWordInfoUrl)
	request, err = client.Get(getWordInfoUrl + "/" + invalidSearchedWord) // Create request with JSON body
	if err != nil {
		t.Error(err.Error())
	}
	defer request.Body.Close()

	if request.StatusCode != http.StatusNotFound {
		t.Errorf("Internal server error expected: %d", request.StatusCode) // HTTP request failed: Test failed
	}
}

// addUser is an API to add a new user
func TestAddUser(t *testing.T) {
	fmt.Println("-- addUser --")

	// Valid case 1

	userJSON := `{
		"email": "` + USER_EMAIL + `",
		"password": "` + USER_PWD + `",
		"username": "` + USER_USERNAME + `"
	}`

	reader = strings.NewReader(userJSON)                          // Convert string to reader
	request, err := client.Post(addUserUrl, CONTENT_TYPE, reader) // Create request with JSON body
	if err != nil {
		t.Error(err.Error())
	}
	if request.StatusCode != http.StatusOK {
		t.Errorf("Success expected: %d", request.StatusCode) // HTTP request failed: Test failed
	}

	// Invalid case 1: Missing user email

	invalidUserJSON := `{
		"password": "` + USER_PWD + `",
		"username": "` + USER_USERNAME + "_invalid_1" + `"
	}`

	reader = strings.NewReader(invalidUserJSON)                  // Convert string to reader
	request, err = client.Post(addUserUrl, CONTENT_TYPE, reader) // Create request with JSON body
	if err != nil {
		t.Error(err.Error())
	}
	if request.StatusCode != http.StatusInternalServerError {
		t.Errorf("Internal server error expected: %d", request.StatusCode) // HTTP request failed: Test failed
	}

	// Invalid case 2: Missing user password

	invalidUserJSON = `{
		"email": "` + USER_EMAIL + "_invalid_2" + `",
		"username": "` + USER_USERNAME + "_invalid_2" + `"
	}`

	reader = strings.NewReader(invalidUserJSON)                  // Convert string to reader
	request, err = client.Post(addUserUrl, CONTENT_TYPE, reader) // Create request with JSON body
	if err != nil {
		t.Error(err.Error())
	}
	if request.StatusCode != http.StatusInternalServerError {
		t.Errorf("Internal server error expected: %d", request.StatusCode) // HTTP request failed: Test failed
	}

	// Invalid case 3: Missing user username

	invalidUserJSON = `{
		"email": "` + USER_EMAIL + "_invalid_3" + `",
		"password": "` + USER_PWD + `"
	}`

	reader = strings.NewReader(invalidUserJSON)                  // Convert string to reader
	request, err = client.Post(addUserUrl, CONTENT_TYPE, reader) // Create request with JSON body
	if err != nil {
		t.Error(err.Error())
	}
	if request.StatusCode != http.StatusInternalServerError {
		t.Errorf("Internal server error expected: %d", request.StatusCode) // HTTP request failed: Test failed
	}
}

// loginUser is an API to add a new user
func TestLoginUser(t *testing.T) {
	fmt.Println("-- loginUser --")

	// Valid case 1

	userJSON := `{
		"password": "` + USER_PWD + `",
		"username": "` + USER_USERNAME + `"
	}`

	reader = strings.NewReader(userJSON)                            // Convert string to reader
	request, err := client.Post(loginUserUrl, CONTENT_TYPE, reader) // Create request with JSON body
	if err != nil {
		t.Error(err.Error())
	}
	if request.StatusCode != http.StatusOK {
		t.Errorf("Success expected: %d", request.StatusCode) // HTTP request failed: Test failed
	}

	// Display cookie
	// body, _ := ioutil.ReadAll(request.Body)
	// request.Body.Close()
	// fmt.Println(string(body))

	// Invalid case 1

	invalidUserJSON := `{
		"password": "` + INVALID_USER_PWD + `",
		"username": "` + INVALID_USER_USERNAME + `"
	}`

	reader = strings.NewReader(invalidUserJSON)                    // Convert string to reader
	request, err = client.Post(loginUserUrl, CONTENT_TYPE, reader) // Create request with JSON body
	if err != nil {
		t.Error(err.Error())
	}
	if request.StatusCode != http.StatusInternalServerError {
		t.Errorf("Internal server error expected: %d", request.StatusCode) // HTTP request failed: Test failed
	}
}

// checkUserLogged is an API to check if a user is logged
func TestCheckUserLogged(t *testing.T) {
	fmt.Println("-- checkUserLogged --")

	// Valid case 1

	request, err := client.Get(checkUserLoggedUrl) // Create request with JSON body
	if err != nil {
		t.Error(err.Error())
	}
	defer request.Body.Close()

	if request.StatusCode != http.StatusOK {
		t.Errorf("Success expected: %d", request.StatusCode) // HTTP request failed: Test failed
	}
}
