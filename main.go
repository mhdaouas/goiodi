package main

import (
	"flag"
	jwt "github.com/dgrijalva/jwt-go"
	"gopkg.in/mgo.v2"
	"io/ioutil"
	"net/http"
)

// Read the key files before starting http handlers
func init() {
	signBytes, err := ioutil.ReadFile(PRIV_KEY_PATH)
	if err != nil {
		Log.Error(err.Error())
	}

	SignKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		Log.Error(err.Error())
	}

	verifyBytes, err := ioutil.ReadFile(PUB_KEY_PATH)
	if err != nil {
		Log.Error(err.Error())
	}

	VerifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		Log.Error(err.Error())
	}
}

func main() {

	var addr string
	var insecure bool

	// Catch user entered arguments
	flag.StringVar(&addr, HOST, PORT, "HTTP server address")
	flag.BoolVar(&insecure, "insecure", false, "Don't use SSL")
	flag.Parse()

	// Init logging system
	InitLogger()

	// Get DB session
	dbs, err := mgo.Dial(HOST)
	if err != nil {
		panic(err)
	}
	defer dbs.Close()

	// Switch the session to a monotonic behavior.
	dbs.SetMode(mgo.Monotonic, true)

	// Create the HTTP server router
	router := createRouter(dbs)

	// Listen to HTTP requests either in secure (TLS based) or insecure mode
	if insecure {
		Log.Notice("Listening on %s (HTTP)...\n", addr)
		err := http.ListenAndServe(addr, router)
		if err != nil {
			Log.Fatal(err)
		}
	} else {
		Log.Notice("Listening on %s (HTTPS)...\n", addr)
		// The certificate files can be created using this command in this
		// repository's root directory:
		//
		// go run $GOROOT/src/crypto/tls/generate_cert.go --host="localhost"
		// or go run /usr/lib/go/src/pkg/crypto/tls/generate_cert.go --host="localhost"
		//
		err := http.ListenAndServeTLS(addr, TLS_CERT_FILEPATH, TLS_KEY_FILEPATH, router)
		if err != nil {
			Log.Fatal(err)
		}
	}

}
