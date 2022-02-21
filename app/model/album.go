package model

type Album struct {
	Number int `json:"number"` 
	Year int `json:"year" validate:"required"`
	Album string `json:"album" validate:"required"`
	Artist string `json:"artist" validate:"required"`
	Genre string `json:"genre" validate:"required"`
	Subgenre string `json:"subgenre" validate:"required"`
}