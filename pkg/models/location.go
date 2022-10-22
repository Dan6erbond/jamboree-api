package models

import "gorm.io/gorm"

type PartyLocation struct {
	gorm.Model
	Location  string
	PartyName string
	Party     Party
	Votes     []PartyLocationVote
}

type PartyLocationVote struct {
	gorm.Model
	Username        string
	PartyLocationID int
	PartyLocation   PartyLocation
}
