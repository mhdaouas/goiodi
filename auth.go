package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"encoding/json"
	"errors"
	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/scrypt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

// Password encryption params
const (
	SALT_SIZE      = 32
	SALT_SECRET    = "fkdspfk@$%"
	SCRYPT_N       = 16384
	SCRYPT_R       = 8
	SCRYPT_P       = 1
	SCRYPT_KEY_LEN = 32
)

// Keys are held in global variables
var (
	VerifyKey *rsa.PublicKey
	SignKey   *rsa.PrivateKey
)

func GenerateSalt(secret []byte) ([]byte, error) {
	buf := make([]byte, SALT_SIZE, SALT_SIZE+sha1.Size)
	_, err := io.ReadFull(rand.Reader, buf)

	if err != nil {
		Log.Error("random read failed:", err.Error())
		return buf, err
	}

	hash := sha1.New()
	hash.Write(buf)
	hash.Write(secret)
	return hash.Sum(buf), err
}

// getPwdHashAndSalt generates a password hash and salt from a password
func getPwdHashAndSalt(password string) (string, string, error) {
	var salt, hash []byte
	salt, err := GenerateSalt([]byte(SALT_SECRET))
	if err != nil {
		Log.Error(err.Error())
		return "", "", err
	}

	hash, err = scrypt.Key([]byte(password), salt, SCRYPT_N, SCRYPT_R, SCRYPT_P, SCRYPT_KEY_LEN)
	if err != nil {
		Log.Error(err.Error())
		return "", "", err
	}

	saltStr := string(salt[:])
	hashStr := string(hash[:])
	return hashStr, saltStr, err
}

// getPwdHash generates a password hash from given password and salt
func getPwdHash(password, salt string) (string, error) {
	hash, err := scrypt.Key([]byte(password), []byte(salt), SCRYPT_N, SCRYPT_R, SCRYPT_P, SCRYPT_KEY_LEN)
	if err != nil {
		Log.Error(err.Error())
		return "", err
	}

	hashStr := string(hash[:])
	return hashStr, err
}

// addUser is an API to add a new user
func addUser(dbs *mgo.Session) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		type addedUser struct {
			Email    string `json:"email" bson:"email"`
			Username string `json:"username" bson:"username"`
			Password string `json:"password" bson:"password"`
		}
		var au addedUser

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

		err = json.Unmarshal(byteData, &au)
		if err != nil {
			Log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if au.Email != "" &&
			au.Username != "" &&
			au.Password != "" {

			var u User

			// Add the creation time of the word based on Epoch time
			now := time.Now()
			secs := now.Unix()
			u.CreationTime = secs

			// Get user e-mail
			u.Email = au.Email

			// Get user username
			u.Username = au.Username

			// Get user password hash and salt
			hash, salt, err := getPwdHashAndSalt(au.Password)
			if err != nil {
				Log.Error(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			u.PwdHash = hash
			u.PwdSalt = salt

			c := dbs.DB(GOIODI_DB).C(USER_COLLECTION)
			err = c.Insert(u)
			if err != nil {
				Log.Error(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			err := errors.New("The user has no related email/username/password")
			Log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// loginUser is an API to login a user with JWT JSON based authentication
func loginUser(dbs *mgo.Session) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type loggedUser struct {
			Username string `json:"username" bson:"username"`
			Password string `json:"password" bson:"password"`
		}
		var lu loggedUser
		var u User

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

		err = json.Unmarshal(byteData, &lu)
		if err != nil {
			Log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = dbs.DB(GOIODI_DB).C(USER_COLLECTION).
			Find(bson.M{
				"username": lu.Username,
			}).
			One(&u)
		if err != nil {
			Log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Check password validity
		luHash, err := getPwdHash(lu.Password, u.PwdSalt)
		if err != nil {
			Log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		Log.Debug("luhash: ", luHash)
		Log.Debug("uhash: ", u.PwdHash)
		if luHash != u.PwdHash {
			http.Error(w, "Invalid user password", http.StatusInternalServerError)
			return
		}

		// Update last log date
		// Get the last log-in time based on Epoch time
		now := time.Now()
		secs := now.Unix()
		luLastLoginTime := secs
		err = dbs.DB(GOIODI_DB).C(USER_COLLECTION).
			Update(
				bson.M{"username": lu.Username},
				bson.M{"$set": bson.M{"logged": true, "last_log_in": luLastLoginTime}},
			)
		if err != nil {
			Log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// create a signer for rsa 256
		t := jwt.New(jwt.GetSigningMethod("RS256"))

		// set our claims
		t.Claims["AccessToken"] = "level1"
		t.Claims["CustomUserInfo"] = struct {
			Name string
			Kind string
		}{u.Username, "human"}

		// Set the expiration time to 1 minute
		t.Claims["exp"] = time.Now().Add(time.Minute * 1).Unix()
		tokenString, err := t.SignedString(SignKey)
		if err != nil {
			err = errors.New("Sorry, error while signing Token!")
			Log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Set user cookie
		http.SetCookie(w, &http.Cookie{
			Name:       AUTH_TOKEN,
			Value:      tokenString,
			Path:       "/",
			RawExpires: "0",
		})
	}
}

// checkUserLogged checks if a user is logged by checking his cookie
func checkUserLogged(w http.ResponseWriter, r *http.Request) {

	// Check HTTP request method is the correct one
	if r.Method != "GET" {
		http.Error(w, "HTTP Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if we have a cookie authentication token
	tokenCookie, err := r.Cookie(AUTH_TOKEN)
	switch {
	case err == http.ErrNoCookie:
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	case err != nil:
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if the token is empty
	if tokenCookie.Value == "" {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Validate the token
	token, err := jwt.Parse(tokenCookie.Value, func(token *jwt.Token) (interface{}, error) {
		// Since we only use the one private key to sign the tokens,
		// we also only use its public counter part to verify
		return VerifyKey, nil
	})

	// Branch out into the possible error from signing
	switch err.(type) {

	case nil: // no error

		// The user token is invalid
		if !token.Valid {
			err = errors.New("Invalid authentication token")
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// The user has a valid token
		Log.Notice("Someone accessed resricted area! Token:%+v\n", token)

	// JWT token Validation error
	case *jwt.ValidationError:
		vErr := err.(*jwt.ValidationError)

		switch vErr.Errors {
		case jwt.ValidationErrorExpired:
			err = errors.New("Token Expired, get a new one.")
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return

		default:
			err = errors.New("Error while parsing JWT token!")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	// Unknown token validation error
	default:
		err = errors.New("Error while parsing JWT token!")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
