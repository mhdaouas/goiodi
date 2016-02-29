package main

// Dictionary word structure
type Word struct {
	Word         string `json:"word" bson:"word"`
	Definition   string `json:"definition" bson:"definition"`
	CreationTime int64  `json:"creation_time" bson:"creation_time"`
}

// Structure for the comments created by users for dictionary words
type Comment struct {
	Word         string `json:"word" bson:"word" db:"word"`
	Creator      string `json:"creator" bson:"creator"`
	Content      string `json:"content" bson:"content"`
	CreationTime int64  `json:"creation_time" bson:"creation_time"`
}
