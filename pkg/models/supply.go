package models

import "gorm.io/gorm"

type Supply struct {
	gorm.Model
	Name      string
	Quantity  int32
	Assignee  string
	IsUrgent  bool
	Emoji     string
	PartyName string
	Party     Party
}
