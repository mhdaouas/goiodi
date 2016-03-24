package main

// Constants
const (
	// location of the files used for signing and verification
	PRIV_KEY_PATH = "goiodi.rsa"     // openssl genrsa -out /etc/app.rsa keysize
	PUB_KEY_PATH  = "goiodi.rsa.pub" // openssl rsa -in /etc/app.rsa -pubout > /etc/app.rsa.pub
	AUTH_TOKEN    = "GoiodiAccessToken"

	// HTTP server related constants
	HOST = "localhost"
	PORT = ":8083"

	// TLS certificate and key generated with:
	// sudo openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout goiodi.key
	// -out goiodi.crt
	TLS_CERT_FILEPATH = "goiodi.crt"
	TLS_KEY_FILEPATH  = "goiodi.key"

	// Content type
	CONTENT_TYPE = "application/json; charset=utf-8"

	// DB related constants
	DB_INIT_SCRIPT     = "initDB.js"
	GOIODI_DB          = "goiodi"
	USER_COLLECTION    = "users"
	WORD_COLLECTION    = "words"
	COMMENT_COLLECTION = "comments"
)
