package domain

import "gorm.io/gorm"

type Vehicle struct {
	gorm.Model
	MemberID    uint
	PlateNumber string `gorm:"index:plate_number,unique"`
	Brand       string
	ModelName   string
	Member      Member
}
