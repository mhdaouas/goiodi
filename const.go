package main

// Constants
const (
	// location of the files used for signing and verification
	PRIV_KEY_PATH = "/etc/goiodi.rsa"     // openssl genrsa -out /etc/app.rsa keysize
	PUB_KEY_PATH  = "/etc/goiodi.rsa.pub" // openssl rsa -in /etc/app.rsa -pubout > /etc/app.rsa.pub
	AUTH_TOKEN    = "GoiodiAccessToken"

	// HTTP server related constants
	HOST              = "localhost"
	PORT              = ":8083"
	TLS_CERT_FILEPATH = "/etc/ssl/certs/goiodi.crt"
	TLS_KEY_FILEPATH  = "/etc/ssl/certs/goiodi.key"
	CONTENT_TYPE      = "application/json; charset=utf-8"

	// DB related constants
	DB_INIT_SCRIPT     = "initDB.js"
	GOIODI_DB          = "goiodi"
	USER_COLLECTION    = "users"
	WORD_COLLECTION    = "words"
	COMMENT_COLLECTION = "comments"
)
