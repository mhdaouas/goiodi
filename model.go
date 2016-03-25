package main

// Dictionary word
type Word struct {
	CreationTime int64     `json:"creation_time" bson:"creation_time"`
	Definition   string    `json:"definition" bson:"definition"`
	Word         string    `json:"word" bson:"word"`
	Comments     []Comment `json:"comments" bson:",omitempty"`
}

// Structure for the comments created by users for dictionary words
type Comment struct {
	Content      string `json:"content" bson:"content"`
	CreationTime int64  `json:"creation_time" bson:"creation_time"`
	Creator      string `json:"creator" bson:"creator"`
	Word         string `json:"word" bson:"word" db:"word"`
}

// Application user
type User struct {
	// Id       string `json:"id"`
	CreationTime int64  `json:"creation_time" bson:"creation_time"`
	Email        string `json:"email" bson:"email"`
	PwdHash      string `json:"pwd_hash" bson:"pwd_hash"`
	PwdSalt      string `json:"pwd_salt" bson:"pwd_salt"`
	LastLogTime  int64  `json:"last_login_time" bson:"last_login_time"`
	Logged       bool   `json:"logged" bson:"logged"`
	Username     string `json:"username" bson:"username"`
}
