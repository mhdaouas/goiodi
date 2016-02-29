package main

import (
	"flag"
	"gopkg.in/mgo.v2"
	"net/http"
)

func main() {

	var addr string
	var insecure bool

	// Catch user entered arguments
	flag.StringVar(&addr, HOST, ":8083", "HTTP server address")
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
