package models

import (
	"time"

	"gorm.io/gorm"
)

type PartyDate struct {
	gorm.Model
	Date      *time.Time
	PartyName string
	Party     Party
	Votes     []PartyDateVote
}

type PartyDateVote struct {
	gorm.Model
	Username    string
	PartyDateID int
	PartyDate   PartyDate
}
