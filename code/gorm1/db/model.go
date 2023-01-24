package model

import (
	"gorm.io/gorm"
)

type Member struct {
	gorm.Model
	Name         	string
	Email        	string
	// Self-referential many-to-many relationship:
	Connections  	[]Member `gorm:"many2many:member_connections"`
	// One-to-many relationship:
	Notes			[]Note `gorm:"foreignKey:MemberID"`
}

type Note struct {
	gorm.Model
	// Has one relationship
	Author		Member `gorm:"foreignKey:ID"`
	Text		string
	MemberID	string	
}

