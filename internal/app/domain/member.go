package domain

import "gorm.io/gorm"

type Member struct {
	gorm.Model
	FirstName   string
	LastName    string
	PhoneNumber string `gorm:"unique"`
	Email       string `gorm:"unique"`
	Vehicles    []Vehicle
	Faces       []Face
}
