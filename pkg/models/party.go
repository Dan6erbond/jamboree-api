package models

type Party struct {
	Name                   string `gorm:"primaryKey"`
	AdminCode              string
	Creator                string
	DateOptionsEnabled     bool
	DateVotingEnabled      bool
	LocationOptionsEnabled bool
	LocationVotingEnabled  bool
	Dates                  []PartyDate
	Locations              []PartyLocation
	Supplies               []Supply
	SongPlaylists          []SongPlaylist
}
