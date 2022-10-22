package models

import "gorm.io/gorm"

type SongPlaylist struct {
	gorm.Model
	Link      string
	PartyName string
	Party     Party
	Votes     []SongPlaylistVote
}

type SongPlaylistVote struct {
	gorm.Model
	Username       string
	SongPlaylistID int
	SongPlaylist   SongPlaylist
}
