package models

type Error struct {
	Err        error
	StatusCode int
	Message    string
}
